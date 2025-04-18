// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2020 Björn Ottosson
// SPDX-FileCopyrightText: 2023 Hajime Hoshi

// Package oklab provides color.Color implementations for Oklab and Oklch.
package oklab

import (
	"image/color"
	"math"
)

var _ color.Color = Oklab{}

// Oklab is a color in the Oklab color space.
type Oklab struct {
	L     float64
	A     float64
	B     float64
	Alpha float64
}

// RGBA implements color.Color.
func (c Oklab) RGBA() (uint32, uint32, uint32, uint32) {
	// https://bottosson.github.io/posts/oklab/#converting-from-linear-srgb-to-oklab

	l_ := c.L + 0.3963377774*c.A + 0.2158037573*c.B
	m_ := c.L - 0.1055613458*c.A - 0.0638541728*c.B
	s_ := c.L - 0.0894841775*c.A - 1.2914855480*c.B

	l := l_ * l_ * l_
	m := m_ * m_ * m_
	s := s_ * s_ * s_

	r_ := +4.0767416621*l - 3.3077115913*m + 0.2309699292*s
	g_ := -1.2684380046*l + 2.6097574011*m - 0.3413193965*s
	b_ := -0.0041960863*l - 0.7034186147*m + 1.7076147010*s

	r := toNonLinear(r_)
	g := toNonLinear(g_)
	b := toNonLinear(b_)
	a := c.Alpha

	if r < 0 {
		r = 0
	}
	if r > 1 {
		r = 1
	}
	if g < 0 {
		g = 0
	}
	if g > 1 {
		g = 1
	}
	if b < 0 {
		b = 0
	}
	if b > 1 {
		b = 1
	}
	if a < 0 {
		a = 0
	}
	if a > 1 {
		a = 1
	}

	return uint32(a * r * 0xffff), uint32(a * g * 0xffff), uint32(a * b * 0xffff), uint32(a * 0xffff)
}

func (c Oklab) oklch() Oklch {
	// https://www.w3.org/TR/css-color-4/#lab-to-lch
	if c.A == 0 && c.B == 0 {
		return Oklch{
			L:     c.L,
			C:     0,
			H:     math.NaN(),
			Alpha: c.Alpha,
		}
	}
	return Oklch{
		L:     c.L,
		C:     math.Hypot(c.A, c.B),
		H:     math.Atan2(c.B, c.A),
		Alpha: c.Alpha,
	}
}

var _ color.Color = Oklch{}

// Oklch is a color in the Oklch color space.
type Oklch struct {
	H     float64
	C     float64
	L     float64
	Alpha float64
}

// RGBA implements color.Color.
func (c Oklch) RGBA() (uint32, uint32, uint32, uint32) {
	return c.oklab().RGBA()
}

func (c Oklch) oklab() Oklab {
	// https://www.w3.org/TR/css-color-4/#lch-to-lab
	if math.IsNaN(c.H) {
		return Oklab{
			L:     c.L,
			A:     0,
			B:     0,
			Alpha: c.Alpha,
		}
	}
	return Oklab{
		L:     c.L,
		A:     c.C * math.Cos(c.H),
		B:     c.C * math.Sin(c.H),
		Alpha: c.Alpha,
	}
}

var (
	OklabModel color.Model = color.ModelFunc(oklabModelFunc)
	OklchModel color.Model = color.ModelFunc(oklchModelFunc)
)

func oklabModelFunc(clr color.Color) color.Color {
	if _, ok := clr.(Oklab); ok {
		return clr
	}
	if c, ok := clr.(Oklch); ok {
		return c.oklab()
	}

	r32, g32, b32, a32 := clr.RGBA()

	r := toLinear(float64(r32) / float64(a32))
	g := toLinear(float64(g32) / float64(a32))
	b := toLinear(float64(b32) / float64(a32))

	// https://bottosson.github.io/posts/oklab/#converting-from-linear-srgb-to-oklab

	l := 0.4122214708*r + 0.5363325363*g + 0.0514459929*b
	m := 0.2119034982*r + 0.6806995451*g + 0.1073969566*b
	s := 0.0883024619*r + 0.2817188376*g + 0.6299787005*b

	l_ := math.Cbrt(l)
	m_ := math.Cbrt(m)
	s_ := math.Cbrt(s)

	return Oklab{
		L:     0.2104542553*l_ + 0.7936177850*m_ - 0.0040720468*s_,
		A:     1.9779984951*l_ - 2.4285922050*m_ + 0.4505937099*s_,
		B:     0.0259040371*l_ + 0.7827717662*m_ - 0.8086757660*s_,
		Alpha: float64(a32) / 0xffff,
	}
}

func oklchModelFunc(clr color.Color) color.Color {
	if _, ok := clr.(Oklch); ok {
		return clr
	}
	return oklabModelFunc(clr).(Oklab).oklch()
}

func toNonLinear(x float64) float64 {
	// https://bottosson.github.io/posts/colorwrong/#what-can-we-do%3F
	if x >= 0.0031308 {
		return (1.055)*math.Pow(x, (1.0/2.4)) - 0.055
	}
	return 12.92 * x
}

func toLinear(x float64) float64 {
	// https://bottosson.github.io/posts/colorwrong/#what-can-we-do%3F
	if x >= 0.04045 {
		return math.Pow((x+0.055)/(1+0.055), 2.4)
	}
	return x / 12.92
}
