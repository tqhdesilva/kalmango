// https://github.com/lovasoa/multivariate-gaussian/blob/master/multivariate_gaussian.js
var n = numeric;

var sqrt2PI = Math.sqrt(Math.PI * 2);

/**
 * Represents a multivariate gaussian
* @param {{sigma: Array<Array<number>>, mu: Array<number>}} gaussian_parameters
**/
function Gaussian(parameters) {
    this.sigma = parameters.sigma;
    this.mu = parameters.mu;
    this.k = this.mu.length; // dimension
    try {
        var det = n.det(this.sigma);
        this._sinv = n.inv(this.sigma); // π ^ (-1)
        this._coeff = 1 / (Math.pow(sqrt2PI, this.k) * Math.sqrt(det));
        if ( !(isFinite(det) && det > 0 && isFinite(this._sinv[0][0]))) {
            throw new Error("Invalid matrix");
        }
    } catch(e) {
        this._sinv = n.rep([this.k, this.k], 0);
        this._coeff = 0;
    }
}

/**
 * Evaluates the density function of the gaussian at the given point
 */
Gaussian.prototype.density = function(x) {
    var delta = n.sub(x, this.mu); // 𝛿 = x - mu
    // Compute  Π = 𝛿T . Σ^(-1) . 𝛿
    var P = 0;
    for(var i=0; i<this.k; i++) {
        var sinv_line = this._sinv[i];
        var sum = 0;
        for(var j=0; j<this.k; j++) {
            sum += sinv_line[j] * delta[j];
        }
        P += delta[i] * sum
    }
    // Return: e^(-Π/2) / √|2.π.Σ|
    return this._coeff * Math.exp(P / -2);
};

// module.exports = Gaussian;