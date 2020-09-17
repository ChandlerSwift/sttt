/* Player 1 (Host, `"x"`):
 1. Connect to websocket server
 2. Request new game `{"cmd": "new_game"}`
 3. Receive new game code `{"response": "game_code", "game_code": "abc-def-ghi"}`

 3. (optional) `{"cmd": "get_state"}`
LOOP:
 4. Receive go-ahead for move `{"response": "state", "state": "...", "your_turn": true, "last_move": null}`
 5. Make move and send to server `{"cmd": "move", "location": "0 0"}` (location is `row column` format, e.g. `{4 6}` for the top right move in the center square )
 6. Server response:
     * Valid: Receive game state from server `{"response": "state", "state": "X O...", "your_turn": false}` (state is an 81-character string of `"x"`, `"o"`, and `" "`.) (proceed) 
     * Invalid: Receive error from server `{"response": "invalid", "state": "...", "your_turn": true}` (send new move)
 7. Repeat until win or loss: `{"response": "win|loss", "state": "..."}`
 */

// 1. Connect to websocket server
// TODO: autoreconnect on lost connection
const ws = new WebSocket(`ws://${window.location.hostname}:12075`);
ws.addEventListener("open", function (event) {
    console.log("Websocket Connection Successful!");
});

document.getElementById("new-game").onclick = function() {
    document.getElementById("new-game").disabled = true;
    document.getElementById("join-game").disabled = true;

    document.getElementById('game-info').innerHTML = "Your token is <b>X</b>";

    // 2. Request new game `{"cmd": "new_game"}`
    ws.send('{"cmd": "new_game"}');

    // TODO: loading...
}

document.getElementById("join-game").onclick = function() {
    document.getElementById("new-game").disabled = true;
    document.getElementById("join-game").disabled = true;

    document.getElementById('game-info').innerHTML = "Your token is <b>O</b>";

    ws.send(`{"cmd": "join_game", "game_code": "${prompt("Game Code:")}"}`);
}

function updateBoard(state, last_move, your_turn) {
    document.getElementById('game-input-submit').disabled = !your_turn;
    document.getElementById('game-input').disabled = !your_turn;
    // TODO: highlight cell
    // TODO: Implement dynamic method (canvas, perhaps?) that also supports drawing win lines
    let hl = [""]
    if (your_turn && last_move != null) {
    }
    document.getElementById("game-board").innerHTML = 
    `<div>${your_turn ? "Your turn" : "their turn" }<pre>
    ${state.charAt(0)}│${state.charAt(1)}│${state.charAt(2)} ┃ ${state.charAt(9)}│${state.charAt(10)}│${state.charAt(11)} ┃ ${state.charAt(18)}│${state.charAt(19)}│${state.charAt(20)}
    ─┼─┼─ ┃ ─┼─┼─ ┃ ─┼─┼─
    ${state.charAt(3)}│${state.charAt(4)}│${state.charAt(5)} ┃ ${state.charAt(12)}│${state.charAt(13)}│${state.charAt(14)} ┃ ${state.charAt(21)}│${state.charAt(22)}│${state.charAt(23)}
    ─┼─┼─ ┃ ─┼─┼─ ┃ ─┼─┼─
    ${state.charAt(6)}│${state.charAt(7)}│${state.charAt(8)} ┃ ${state.charAt(15)}│${state.charAt(16)}│${state.charAt(17)} ┃ ${state.charAt(24)}│${state.charAt(25)}│${state.charAt(26)}
    ━━━━━━╋━━━━━━━╋━━━━━━━
    ${state.charAt(27)}│${state.charAt(28)}│${state.charAt(29)} ┃ ${state.charAt(36)}│${state.charAt(37)}│${state.charAt(38)} ┃ ${state.charAt(45)}│${state.charAt(46)}│${state.charAt(47)}
    ─┼─┼─ ┃ ─┼─┼─ ┃ ─┼─┼─
    ${state.charAt(30)}│${state.charAt(31)}│${state.charAt(32)} ┃ ${state.charAt(39)}│${state.charAt(40)}│${state.charAt(41)} ┃ ${state.charAt(48)}│${state.charAt(49)}│${state.charAt(50)}
    ─┼─┼─ ┃ ─┼─┼─ ┃ ─┼─┼─
    ${state.charAt(33)}│${state.charAt(34)}│${state.charAt(35)} ┃ ${state.charAt(42)}│${state.charAt(43)}│${state.charAt(44)} ┃ ${state.charAt(51)}│${state.charAt(52)}│${state.charAt(53)}
    ━━━━━━╋━━━━━━━╋━━━━━━━
    ${state.charAt(54)}│${state.charAt(55)}│${state.charAt(56)} ┃ ${state.charAt(63)}│${state.charAt(64)}│${state.charAt(65)} ┃ ${state.charAt(72)}│${state.charAt(73)}│${state.charAt(74)}
    ─┼─┼─ ┃ ─┼─┼─ ┃ ─┼─┼─
    ${state.charAt(57)}│${state.charAt(58)}│${state.charAt(59)} ┃ ${state.charAt(66)}│${state.charAt(67)}│${state.charAt(68)} ┃ ${state.charAt(75)}│${state.charAt(76)}│${state.charAt(77)}
    ─┼─┼─ ┃ ─┼─┼─ ┃ ─┼─┼─
    ${state.charAt(60)}│${state.charAt(61)}│${state.charAt(62)} ┃ ${state.charAt(69)}│${state.charAt(70)}│${state.charAt(71)} ┃ ${state.charAt(78)}│${state.charAt(79)}│${state.charAt(80)}
    </pre>`;
}

document.getElementById('game-input-submit').addEventListener("click", function () {
    let move = document.getElementById('game-input').value;

    // check for valid input (i.e. '0 0' to '8 8')
    if (move.length != 3 || move[1] !== " " ||
            parseInt(move[0]) < 0 || parseInt(move[0]) > 8 ||
            parseInt(move[2]) < 0 || parseInt(move[2]) > 8) {
        alert('invalid move');
        return;
    }

    ws.send(`{"cmd": "move", "location": "${move}"}`);
    // TODO: click for input
    document.getElementById('game-input-submit').disabled = true;
});

ws.addEventListener("message", function (event) {
    console.log("WS RX: " + event.data)
    let data = JSON.parse(event.data);
    switch (data.response) {
        case "game_code":
            // TODO: save this somewhere
            alert(`Your game code is ${data.game_code}. Share this with your opponent to play!`);
            break;
        case "state":
            updateBoard(data.state, data.last_move, data.your_turn);
            if (data.your_turn) {
                if (data.last_move) {
                   alert(`Opponent played ${data.last_move}`);
                }
                // takeTurn();
            }
            break;
        case "invalid":
            alert("That move is not valid. Please try again.");
            // TODO takeTurn();
            break;
        case "win":
            alert("Game over. Congrats on your win!");
            break;
        case "loss":
            alert("Game over. You lose.");
            break;
        case "error":
            alert("ERROR: " + data.error);
            break;
        default:
            console.log("UNKNOWN RESPONSE: " + event.data);
    }
});
