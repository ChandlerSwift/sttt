package main

//go:generate swagger generate spec

import (
	"html/template"
	"log"
	"net/http"

	"github.com/chandlerswift/sttt/internal/models"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type alert struct {
	Type    string
	Content template.HTML
}

type pageData struct {
	Alerts []alert
}

func main() {
	// Read config
	viper.SetConfigName("sttt")
	viper.AddConfigPath(".")
	viper.SetDefault("database_path", "sttt.db")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; that's okay.
		} else {
			log.Fatalf("Fatal error parsing config file: %v\n", err)
		}
	}

	// Set up database
	db, err := gorm.Open(sqlite.Open(viper.GetString("database_path")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{}, &models.Game{})

	r := mux.NewRouter()

	// Handle API routes
	api := r.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Serve index page on all unhandled routes
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	log.Println("Serving on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
