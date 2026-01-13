package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strconv"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/font/gofont/gobold"
	"github.com/golang/freetype/truetype"
)

// generateTrayIcon creates a tray icon with temperature text
func generateTrayIcon(temperature float64) ([]byte, error) {
	// Create a 64x64 image for the tray icon
	size := 64
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Fill with a gradient background (blue to purple)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// Calculate gradient
			ratio := float64(y) / float64(size)
			r := uint8(102 + (118-102)*ratio)
			g := uint8(126 + (75-126)*ratio)
			b := uint8(234 + (162-234)*ratio)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// Draw circular shape to make it look better
	center := size / 2
	radius := size / 2
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := x - center
			dy := y - center
			if dx*dx+dy*dy > radius*radius {
				img.Set(x, y, color.RGBA{0, 0, 0, 0}) // Transparent outside circle
			}
		}
	}

	// Prepare temperature text
	tempStr := strconv.Itoa(int(temperature)) + "°"
	
	// Draw text in the center
	point := fixed.Point26_6{
		X: fixed.I(size/2 - len(tempStr)*3),
		Y: fixed.I(size/2 + 7),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{255, 255, 255, 255}),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(tempStr)

	// Convert to PNG bytes
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// generateTrayIconWithWeather creates a tray icon based on weather data
func generateTrayIconWithWeather(weather *WeatherData) ([]byte, error) {
	// Create a 64x64 image
	size := 64
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Determine background color based on weather condition
	var bgColor color.RGBA
	switch weather.Condition {
	case "Sunny":
		bgColor = color.RGBA{255, 193, 7, 255} // Amber/Yellow
	case "Partly Cloudy":
		bgColor = color.RGBA{158, 158, 158, 255} // Gray
	case "Cloudy":
		bgColor = color.RGBA{117, 117, 117, 255} // Dark gray
	case "Rainy":
		bgColor = color.RGBA{33, 150, 243, 255} // Blue
	case "Stormy":
		bgColor = color.RGBA{63, 81, 181, 255} // Dark blue
	case "Snowy":
		bgColor = color.RGBA{224, 247, 250, 255} // Light cyan
	case "Foggy":
		bgColor = color.RGBA{189, 189, 189, 255} // Light gray
	default:
		bgColor = color.RGBA{102, 126, 234, 255} // Default purple
	}

	// Fill background
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Draw circular shape
	center := size / 2
	radius := size / 2
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := x - center
			dy := y - center
			if dx*dx+dy*dy > radius*radius {
				img.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}

	// Temperature text - just the number, no degree symbol for better visibility
	tempStr := strconv.Itoa(int(weather.Temperature))
	
	// Use a larger font
	ft, err := truetype.Parse(gobold.TTF)
	if err != nil {
		// Fallback to basic font if truetype fails
		return generateSimpleTrayIcon(weather, img)
	}

	// Create font face with much larger size (42pt)
	face := truetype.NewFace(ft, &truetype.Options{
		Size: 42,
		DPI:  72,
	})
	defer face.Close()

	// Measure text to center it
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{1, 1, 1, 255}),
		Face: face,
	}

	// Get text bounds for centering
	bounds, _ := d.BoundString(tempStr)
	textWidth := (bounds.Max.X - bounds.Min.X).Ceil()
	
	// Center the text (moved down 1 pixel)
	point := fixed.Point26_6{
		X: fixed.I((size - textWidth) / 2),
		Y: fixed.I(size/2 + 11),
	}

	d.Dot = point
	d.DrawString(tempStr)

	// Convert to PNG
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// generateSimpleTrayIcon is a fallback with large basic font
func generateSimpleTrayIcon(weather *WeatherData, img *image.RGBA) ([]byte, error) {
	size := 64
	tempStr := strconv.Itoa(int(weather.Temperature))
	
	// Draw large text manually using bigger basic font
	// Draw each character larger by drawing multiple times with offset
	startX := (size - len(tempStr)*20) / 2
	startY := size/2 + 5
	
	for i, ch := range tempStr {
		// Draw character multiple times to make it bold
		for dx := 0; dx < 3; dx++ {
			for dy := 0; dy < 3; dy++ {
				point := fixed.Point26_6{
					X: fixed.I(startX + i*20 + dx),
					Y: fixed.I(startY + dy),
				}
				d := &font.Drawer{
					Dst:  img,
					Src:  image.NewUniform(color.RGBA{255, 255, 255, 255}),
					Face: basicfont.Face7x13,
					Dot:  point,
				}
				d.DrawString(string(ch))
			}
		}
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	return buf.Bytes(), err
}
