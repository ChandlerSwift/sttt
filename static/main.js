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
      {
        type: "alert-success",
        content: "<b>Success!</b> Content was updated!."
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
