
var socket = new WebSocket("ws://localhost:8080/websocket")
var svg = d3.select("#graph").append("svg").attr("height", 500).attr("width", 500)
var circle = svg.append("circle").attr("cx", 250).attr("cy", 250).attr("fill", "black").attr("r", 10)
socket.addEventListener("message", function(event){
    var d = JSON.parse(event.data)
    circle.attr("cx", 50 * d.actual_position[0])
    circle.attr("cy", 50 * d.actual_position[1])
})
// socket.send("data")
