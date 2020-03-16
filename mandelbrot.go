package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"os"
	"time"
)

func mandelbrotComplex(c complex128, maxIter int) int {
	z := c

	for i := 0; i <= maxIter; i++ {
		if cmplx.Abs(z) > 2 {
			return i
		}
		z = cmplx.Pow(z, 2) + c
	}
	return 0
}

func mandelbrot(c complex128, maxIter int) int {
	ca := real(c)
	cb := imag(c)
	za := ca
	zb := cb

	for i := 1; i <= maxIter; i++ {
		zas := za * za
		zbs := zb * zb

		zb = 2*za*zb + cb
		za = zas - zbs + ca

		if zas+zbs > 4 {
			return i
		}
	}
	return 0
}

func HsvToRgba(H, S, V float64) color.RGBA {
	// Taken and modified from https://github.com/lucasb-eyer/go-colorful/blob/master/colors.go
	if H == 0 {
		return color.RGBA{0, 0, 0, 255}
	}
	Hp := H / 60.0
	C := V * S
	X := C * (1.0 - math.Abs(math.Mod(Hp, 2.0)-1.0))

	m := V - C
	r, g, b := 0.0, 0.0, 0.0

	switch {
	case 0.0 <= Hp && Hp < 1.0:
		r = C
		g = X
	case 1.0 <= Hp && Hp < 2.0:
		r = X
		g = C
	case 2.0 <= Hp && Hp < 3.0:
		g = C
		b = X
	case 3.0 <= Hp && Hp < 4.0:
		g = X
		b = C
	case 4.0 <= Hp && Hp < 5.0:
		r = X
		b = C
	case 5.0 <= Hp && Hp < 6.0:
		r = C
		b = X
	}

	return color.RGBA{uint8((m + r) * 255), uint8((m + g) * 255), uint8((m + b) * 255), 255}
}

func main() {
	ResX := 1080
	ResY := 1080

	Pos := complex(0, 0)
	Scale := 2
	Iterations := 1000

	ResSmaller := func(a, b int) int { // min
		if a < b {
			return a
		}
		return b
	}(ResX, ResY)
	DistPerPixel := float64(Scale) / float64(ResSmaller/2)

	fmt.Println("Mandelbrot", Pos, Scale) //, ResSmaller, DistPerPixel)

	img := image.NewRGBA(image.Rect(0, 0, ResX, ResY))

	Left := float64(real(Pos)) + float64(-ResX/2)*DistPerPixel
	Top := float64(imag(Pos)) + float64(-ResY/2)*DistPerPixel

	timeBefore := time.Now()
	for x := 0; x < ResX; x++ {
		fmt.Println("x", x)
		Real := Left + float64(x)*DistPerPixel
		for y := 0; y < ResY; y++ {
			Imag := Top + float64(y)*DistPerPixel
			c := complex(Real, Imag)
			val := mandelbrot(c, Iterations)
			img.SetRGBA(x, y, HsvToRgba(float64(val), 1.0, 1.0))
		}
	}
	fmt.Println(time.Since(timeBefore))

	// File output
	file, err := os.Create("mandelbrot.jpg")
	if err != nil {
		log.Fatalf("failed create file: %s", err)
	}
	png.Encode(file, img)
}
