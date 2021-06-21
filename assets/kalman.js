var socket = new WebSocket(`${isHTTPS() ? 'wss' : 'ws'}://${window.location.host}/websocket`)
var svg = d3.select("#graph").append("svg").attr("height", 500).attr("width", 500)
var circle = svg.append("circle").attr("cx", 250).attr("cy", 250).attr("fill", "black").attr("r", 10)
var est_circle = svg.append("circle").attr("cx", 250).attr("cy", 250).attr("fill", "blue").attr("r", 10)


var i0 = d3.interpolateHsvLong(d3.hsv(120, 1, 0.65), d3.hsv(60, 1, 0.90)),
    i1 = d3.interpolateHsvLong(d3.hsv(60, 1, 0.90), d3.hsv(0, 0, 0.95)),
    interpolateTerrain = function (t) { return t < 0.5 ? i0(t * 2) : i1((t - 0.5) * 2); },
    color = d3.scaleSequential(interpolateTerrain).domain([90, 190]);

const scale = 50
var r = 10 * scale, c = 10 * scale;


function isHTTPS() {
    return window.location.protocol == "https:"
}

function scaler(x) {
    return x * scale
}

function getProbabilityDistribution(event_data) {
    dist_parameters = {
        sigma: event_data.estimated_covariance.map(function (x) { return x.map(scaler) }),
        mu: event_data.estimated_position.map(scaler)
    }
    var distribution = new Gaussian(dist_parameters);
    return distribution
}

function updatePDFDisplay(dist) {
    var data = []
    for (var i = 0; i <= r; i += scale) {
        for (var j = 0; j <= c; j += scale) {
            var p = - Math.log(dist.density([i, j]))
            data.push({
                x: i,
                y: j,
                value: p
            })
        }
    }
    heatmap.setData({
        min: 1,
        max: 1000,
        data: data
    })
}


var heatmap = h337.create({
    container: d3.select("#graph").node(),
    radius: 100,
})


function sendMeasureRequest() {
    socket.send("update")
    console.log("sent measurement request")
}

var counter = 0;

socket.addEventListener("message", function (event) {
    var d = JSON.parse(event.data)
    if (d["actual_position"] != null) {
        circle.transition().attr("cx", scale * d.actual_position[0]).attr("cy", scale * d.actual_position[1])
        est_circle.transition().attr("cx", scale * d.estimated_position[0]).attr("cy", scale * d.estimated_position[1])
        if (counter % 10 == 0) {
            var dist = getProbabilityDistribution(d)
            updatePDFDisplay(dist)
        }
        counter += 1;
    }
})
