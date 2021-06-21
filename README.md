[Kalman filter in go](https://kalmango.herokuapp.com/).

# Run
`go run ./cmd/api`

Then go to `http://localhost:8080` in your browser.
The black circle is the ground truth location of the puck, and the blue circle is the point estimate of the position.
A heatmap visualizes the probability distribution of the estimated position, with white being higher and red being close to 0.
You can take measurements of the puck location by pressing the `Measure` button, which takes a noisy measurement of the position and velocity of the puck.
Each time the puck hits the edge of the screen, a control is applied to the puck to reverse the horizontal or vertical direction, keeping it bouncing within the screen.
The estimated position of the puck is not constrained to being on the screen.

![](figures/kf_demo.gif)