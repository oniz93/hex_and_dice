package mapgen

import (
	opensimplex "github.com/ojrac/opensimplex-go"
)

// NoiseGenerator wraps OpenSimplex noise for terrain generation.
type NoiseGenerator struct {
	noise opensimplex.Noise
	scale float64
}

// NewNoiseGenerator creates a noise generator with the given seed and scale.
func NewNoiseGenerator(seed int64, scale float64) *NoiseGenerator {
	return &NoiseGenerator{
		noise: opensimplex.New(seed),
		scale: scale,
	}
}

// Eval2D returns a noise value in [0, 1] for the given coordinates.
// OpenSimplex returns [-1, 1], so we normalize to [0, 1].
func (ng *NoiseGenerator) Eval2D(x, y float64) float64 {
	raw := ng.noise.Eval2(x*ng.scale, y*ng.scale)
	return (raw + 1.0) / 2.0
}

// MultiOctave generates multi-octave noise for more natural terrain.
// Uses 3 octaves with decreasing amplitude and increasing frequency.
func (ng *NoiseGenerator) MultiOctave(x, y float64) float64 {
	value := 0.0
	amplitude := 1.0
	frequency := 1.0
	maxAmplitude := 0.0

	for i := 0; i < 3; i++ {
		value += ng.noise.Eval2(x*ng.scale*frequency, y*ng.scale*frequency) * amplitude
		maxAmplitude += amplitude
		amplitude *= 0.5
		frequency *= 2.0
	}

	// Normalize to [0, 1]
	return (value/maxAmplitude + 1.0) / 2.0
}
