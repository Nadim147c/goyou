package color

import (
	"fmt"
	"math"
)

// Hct represents a color in the HCT color space (Hue, Chroma, Tone).
// HCT provides a perceptually accurate color measurement system that can also
// accurately render what colors will appear as in different lighting
// environments.
type Hct struct {
	Hue    float64
	Chroma float64
	Tone   float64
}

// Ensure Color implements the color.Color interface
var _ digitalColor = (*Hct)(nil)

// From creates an HCT color from the provided hue, chroma, and tone values.
//
// hue: 0 <= hue < 360; invalid values are corrected.
// chroma: 0 <= chroma < ?; Chroma has a different maximum for any given hue and
// tone.
// tone: 0 <= tone <= 100; invalid values are corrected.
func NewHct(hue, chroma, tone float64) Hct {
	return solveToARGB(hue, chroma, tone).ToHct()
}

// ToInt returns the ARGB representation of this color.
func (h Hct) ToARGB() ARGB {
	return solveToARGB(h.Hue, h.Chroma, h.Tone)
}

// ToInt returns the ARGB representation of this color.
func (h Hct) ToXYZ() XYZ {
	return h.ToARGB().ToXYZ()
}

// ToInt returns the ARGB representation of this color.
func (h Hct) ToLab() Lab {
	return h.ToARGB().ToXYZ().ToLab()
}

// HctFromColor creates an HCT color from the provided ARGB integer.
func (h Hct) ToHct() Hct {
	return h
}

// HctFromColor creates an HCT color from the provided ARGB integer.
func (h Hct) ToCam() *Cam16 {
	return h.ToARGB().ToCam()
}

// ToInt returns the ARGB representation of this color.
func (h Hct) RGBA() (uint32, uint32, uint32, uint32) {
	return solveToARGB(h.Hue, h.Chroma, h.Tone).RGBA()
}

// String returns a string representation of the HCT color.
func (h Hct) String() string {
	return fmt.Sprintf("HCT(%.4f, %.4f, %.4f) %s", h.Hue, h.Chroma, h.Tone, h.ToARGB().AnsiBg("  "))
}

// Hash returns a uint64 hash representation of the HCT color.
// This is much more efficient than string-based hashing.
func (h Hct) Hash() uint64 {
	// Convert each float to bits and extract portions for the hash
	hueBits := math.Float64bits(h.Hue)
	chromaBits := math.Float64bits(h.Chroma)
	toneBits := math.Float64bits(h.Tone)

	// Create hash using FNV-1a inspired approach, but with direct bit operations
	// for better performance, combining all three components
	hash := uint64(14695981039346656037) // FNV offset basis

	// Mix in the hue bits
	hash ^= (hueBits & 0xFFFFFFFF)
	hash *= 1099511628211 // FNV prime

	// Mix in the chroma bits
	hash ^= (chromaBits & 0xFFFFFFFF)
	hash *= 1099511628211

	// Mix in the tone bits
	hash ^= (toneBits & 0xFFFFFFFF)
	hash *= 1099511628211

	return hash
}

// IsBlue determines if a hue is in the blue range.
func (h Hct) IsBlue() bool {
	return h.Hue >= 250 && h.Hue < 270
}

// IsYellow determines if a hue is in the yellow range.
func (h Hct) IsYellow() bool {
	return h.Hue >= 105 && h.Hue < 125
}

// IsCyan determines if a hue is in the cyan range.
func (h Hct) IsCyan() bool {
	return h.Hue >= 170 && h.Hue < 207
}

// IsBlue determines if a hue is in the blue range.
func IsBlue(hue float64) bool {
	return hue >= 250 && hue < 270
}

// IsYellow determines if a hue is in the yellow range.
func IsYellow(hue float64) bool {
	return hue >= 105 && hue < 125
}

// IsCyan determines if a hue is in the cyan range.
func IsCyan(hue float64) bool {
	return hue >= 170 && hue < 207
}

// InViewingConditions translates a color into different ViewingConditions.
//
// Colors change appearance. They look different with lights on versus off,
// the same color, as in hex code, on white looks different when on black.
// This is called color relativity, most famously explicated by Josef Albers
// in Interaction of Color.
//
// In color science, color appearance models can account for this and
// calculate the appearance of a color in different settings. HCT is based on
// CAM16, a color appearance model, and uses it to make these calculations.
//
// See ViewingConditions.Make for parameters affecting color appearance.
func (h *Hct) InViewingConditions(env *Environmnet) Hct {
	cam := h.ToARGB().ToCam()
	viewedInEnv := cam.Viewed(env)
	newCam := viewedInEnv.ToCam()
	return newCam.ToHct()
}
