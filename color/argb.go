package color

import (
	"errors"
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"
)

// Offset indiacates bit offset of Components in ARGB
const (
	blueOffset ARGB = iota * 8
	greenOffset
	redOffset
	alphaOffset
)

// Brightest is the max value of uint8 color
const Brightest = uint8(0xFF) // 255

// ARGB is an ARGB color packed into a uint32.
type ARGB uint32

// Ensure ARGB implements the color.Color interface
var _ digitalColor = (*ARGB)(nil)

// NewARGB creates a ARGB from individual 8-bit alpha, red, green, and blue
// components.
func NewARGB(a, r, g, b uint8) ARGB {
	return ARGB(a)<<alphaOffset |
		ARGB(r)<<redOffset |
		ARGB(g)<<greenOffset |
		ARGB(b)<<blueOffset
}

// ARGBFromInterface
func ARGBFromInterface(color color.Color) ARGB {
	r16, g16, b16, a16 := color.RGBA()

	// Convert from [0, 65535] to [0, 255]
	r8 := uint8(r16 >> 8)
	g8 := uint8(g16 >> 8)
	b8 := uint8(b16 >> 8)
	a8 := uint8(a16 >> 8)

	return NewARGB(a8, r8, g8, b8)
}

// Converts an L* value to an ARGB representation. lstar is L* in L*a*b*.
// returns ARGB representation of grayscale color with lightness matching L*
func ARGBFromLstar(lstar float64) ARGB {
	y := YFromLstar(lstar)
	component := Delinearized(y)
	return ARGBFromRGB(component, component, component)
}

// FromARGB creates a ARGB from xyz color space cordinates.
func ARGBFromXYZ(x, y, z float64) ARGB {
	return NewXYZ(x, y, z).ToARGB()
}

// FromARGB creates a ARGB from individual 8-bit red, green, and blue
// components.
func ARGBFromRGB(r, g, b uint8) ARGB {
	return NewARGB(0xFF, r, g, b)
}

// FromARGB creates a ARGB from individual 8-bit red, green, and blue
// components.
func ARGBFromLinRGB(r, g, b float64) ARGB {
	dr, dg, db := Delinearized3(r, g, b)
	return NewARGB(0xFF, dr, dg, db)
}

// ToCam16 convert ARGB Color to Cam16
func (c ARGB) ToCam() *Cam16 {
	return Cam16FromXyzInEnv(c.ToXYZ(), &DefaultEnviroment)
}

// ToHct convert ARGB Color to Hct
func (c ARGB) ToHct() Hct {
	cam := c.ToCam()
	return Hct{cam.Hue, cam.Chroma, c.LStar()}
}

// ToXYZ return XYZ color version for c.
func (c ARGB) ToXYZ() XYZ {
	r, g, b := c.Red(), c.Green(), c.Blue()

	// Convert RGB channel to linear color (0-1.0)
	lr, lg, lb := Linearized3(r, g, b)

	x, y, z := SRGB_TO_XYZ.MultiplyXYZ(lr, lg, lb).Values()
	return XYZ{x, y, z}
}

// ToLab convert Color to LabColor
func (c ARGB) ToLab() Lab {
	return c.ToXYZ().ToLab()
}

// ToLab convert Color to LabColor
func (c ARGB) ToARGB() ARGB {
	return c
}

func (c ARGB) Values() (uint8, uint8, uint8, uint8) {
	return c.Alpha(), c.Red(), c.Green(), c.Blue()
}

// RGBA implements the color.Color interface.
// It returns r, g, b, a values in the 0-65535 range.
func (c ARGB) RGBA() (uint32, uint32, uint32, uint32) {
	a, r, g, b := c.Values()
	// Convert from 8-bit to 16-bit by scaling: v * 0x101 == v * 257
	return uint32(r) * 0x101, uint32(g) * 0x101, uint32(b) * 0x101, uint32(a) * 0x101
}

// Lstart
func (c ARGB) LStar() float64 {
	r, g, b := c.Red(), c.Green(), c.Blue()
	// Convert RGB channel to linear color (0-1.0)
	lr, lg, lb := Linearized3(r, g, b)

	// Only calculate Y value of XYZ for LStar
	my1, my2, my3 := SRGB_TO_XYZ[1].Values()
	y := my1*lr + my2*lg + my3*lb
	return LstarFromY(y)
}

// AnsiFg wraps the given text with the ANSI escape sequence for the foreground color.
func (c ARGB) AnsiFg(text string) string {
	_, r, g, b := c.Values()
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, text)
}

// AnsiBg wraps the given text with the ANSI escape sequence for the background color.
func (c ARGB) AnsiBg(text string) string {
	_, r, g, b := c.Values()
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm%s\x1b[0m", r, g, b, text)
}

func (c ARGB) String() string {
	return c.HexRGB() + " " + c.AnsiBg("  ")
}

// TextMarshaler
func (c ARGB) MarshalText() ([]byte, error) {
	return []byte(c.HexRGBA()), nil
}

// TextUnmarshaler
func (c *ARGB) UnmarshalText(text []byte) error {
	argb, err := ARGBFromHex(string(text))
	if err != nil {
		return err
	}
	*c = argb
	return nil
}

// Alpha returns the 8-bit alpha component of the color.
func (c ARGB) Alpha() uint8 {
	return uint8((c >> alphaOffset) & 0xFF)
}

// Red returns the 8-bit red component of the color.
func (c ARGB) Red() uint8 {
	return uint8((c >> redOffset) & 0xFF)
}

// Green returns the 8-bit green component of the color.
func (c ARGB) Green() uint8 {
	return uint8((c >> greenOffset) & 0xFF)
}

// Blue returns the 8-bit blue component of the color.
func (c ARGB) Blue() uint8 {
	return uint8((c >> blueOffset) & 0xFF)
}

// HexARGB return #RRGGBB represetation of the color
func (c ARGB) HexRGB() string {
	return fmt.Sprintf("#%02X%02X%02X", c.Red(), c.Green(), c.Blue())
}

// HexARGB return #AARRGGBB represetation of the color
func (c ARGB) HexARGB() string {
	return fmt.Sprintf("#%02X%02X%02X%02X", c.Alpha(), c.Red(), c.Green(), c.Blue())
}

// HexRGBA return #RRGGBBAA represetation of the color
func (c ARGB) HexRGBA() string {
	return fmt.Sprintf("#%02X%02X%02X%02X", c.Red(), c.Green(), c.Blue(), c.Alpha())
}

// ARGBFromHexMust parses a hex color string and returns a Color.
// Supports formats: #RGB, #RGBA, #RRGGBB, #RRGGBBAA
func ARGBFromHexMust(hex string) ARGB {
	color, err := ARGBFromHex(hex)
	if err != nil {
		panic(err)
	}
	return color
}

// ARGBFromHex parses a hex color string and returns a Color.
// Supports formats: #RGB, #RGBA, #RRGGBB, #RRGGBBAA
func ARGBFromHex(hex string) (ARGB, error) {
	hex = strings.TrimPrefix(hex, "#")

	// Regex check if input is valid or not
	hexColorRegex := regexp.MustCompile(`^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	if !hexColorRegex.MatchString(hex) {
		return 0, errors.New("invalid hex color format")
	}

	// Expand shorthand formats to full 6/8-char form
	switch len(hex) {
	case 3: // #RGB → #RRGGBB
		hex = fmt.Sprintf("%c%c%c%c%c%c",
			hex[0], hex[0],
			hex[1], hex[1],
			hex[2], hex[2],
		)
	case 4: // #RGBA → #RRGGBBAA
		hex = fmt.Sprintf("%c%c%c%c%c%c%c%c",
			hex[0], hex[0],
			hex[1], hex[1],
			hex[2], hex[2],
			hex[3], hex[3],
		)
	}

	var r, g, b, a uint8 = 0, 0, 0, 0xFF

	switch len(hex) {
	case 6:
		val, err := strconv.ParseUint(hex, 16, 24)
		if err != nil {
			return 0, fmt.Errorf("invalid hex color: %w", err)
		}
		r = uint8(val >> 16)
		g = uint8((val >> 8) & 0xFF)
		b = uint8(val & 0xFF)

	case 8:
		val, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid hex color: %w", err)
		}
		r = uint8(val >> 24)
		g = uint8((val >> 16) & 0xFF)
		b = uint8((val >> 8) & 0xFF)
		a = uint8(val & 0xFF)
	}

	return NewARGB(a, r, g, b), nil
}
