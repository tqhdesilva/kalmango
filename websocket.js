// test some js stuff
// should work if you paste into console from localhost:8080
var socket = new WebSocket("ws://localhost:8080")
socket.onmessage = function (event) {
    console.log(event.data);
}

socket.send("data")