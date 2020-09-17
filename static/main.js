Vue.component('game-grid', {
  template: ''
})

let app = new Vue({
  el: '#app',
  data: {
    alerts: [
      {
        type: "alert-info",
        content: "<b>Info!</b> Here's some info."
      },
    ]
  },
  methods: {
    startGame: function() {
      alert("Starting game!")
    },
    joinGame: function() {
      alert("Joining game!")
    },
  }
});
