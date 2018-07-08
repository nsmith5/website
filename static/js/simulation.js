var n = 30;

var nodes = d3.range(n * n).map(function(i) {
  return {
    index: i,
    vx: 0,
    vy: 0
  };
});

var canvas = document.querySelector("canvas"),
    context = canvas.getContext("2d"),
    width = canvas.width,
    height = canvas.height;

var simulation = d3.forceSimulation(nodes)
    .velocityDecay(0.)
    .alphaDecay(0.0)
    .alpha(0.001)
    .force("collide", d3.forceCollide(4))
    // .force("y", d3.forceY())
    // .force("x", d3.forceX())
    .force("charge", d3.forceManyBody().strength(1))
    .on("tick", ticked);

function ticked() {
  context.clearRect(0, 0, width, height);
  context.save();
  context.translate(width / 2, height / 2);
  nodes.forEach(drawNode);
  context.restore();
}

function drawNode(d) {
  context.beginPath();
  context.moveTo(d.x + 3, d.y);
  context.arc(d.x, d.y, 3, 0, 2 * Math.PI);
  context.fillStyle = "steelblue";
  context.fill();
}