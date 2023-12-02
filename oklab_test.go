// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 Hajime Hoshi

package oklab_test

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/oklab"
)

func TestPrecision(t *testing.T) {
	const num = 17

	rgbs := make([]color.RGBA, 0, num*num*num)

	for i := 0; i < num; i++ {
		for j := 0; j < num; j++ {
			for k := 0; k < num; k++ {
				rgbs = append(rgbs, color.RGBA{
					R: byte(0xff * float64(i) / (num - 1)),
					G: byte(0xff * float64(j) / (num - 1)),
					B: byte(0xff * float64(k) / (num - 1)),
					A: 0xff,
				})
			}
		}
	}

	for _, rgb := range rgbs {
		want := rgb
		got0 := color.RGBAModel.Convert(oklab.OklabModel.Convert(rgb)).(color.RGBA)
		if got0 != want {
			t.Errorf("got: %v, want: %v", got0, want)
		}
		got1 := color.RGBAModel.Convert(oklab.OklchModel.Convert(rgb)).(color.RGBA)
		if got1 != want {
			t.Errorf("got: %v, want: %v", got1, want)
		}
	}
}
