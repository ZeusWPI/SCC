package view

import (
	"image"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
)

// ImageToString converts an image to a string
// If either widht or height is 0 then the aspect ratio is kept
func ImageToString(img image.Image, width, height int) string {
	if width == 0 || height == 0 {
		return imageToString(imaging.Resize(img, width, height, imaging.Lanczos))
	}

	imgW := imaging.Resize(img, width, 0, imaging.Lanczos)
	if imgW.Bounds().Dy() <= height {
		return imageToString(imgW)
	}

	return imageToString(imaging.Resize(img, 0, height, imaging.Lanczos))
}

func imageToString(img image.Image) string {
	b := img.Bounds()
	imageWidth := b.Max.X
	h := b.Max.Y
	str := strings.Builder{}

	for heightCounter := 0; heightCounter < h; heightCounter += 2 {
		for x := imageWidth; x < img.Bounds().Dx(); x += 2 {
			str.WriteString(" ")
		}

		for x := 0; x < imageWidth; x++ {
			c1, _ := colorful.MakeColor(img.At(x, heightCounter))
			color1 := lipgloss.Color(c1.Hex())
			c2, _ := colorful.MakeColor(img.At(x, heightCounter+1))
			color2 := lipgloss.Color(c2.Hex())

			style := lipgloss.NewStyle().Foreground(color1)
			// Prevent a dark line at the bottom for specific heights
			if heightCounter != h-1 || heightCounter%2 != 0 {
				style = style.Background(color2)
			}
			str.WriteString(style.Render("▀"))
		}

		str.WriteString("\n")
	}

	return str.String()
}

// GetOuterWidth returns the outer border size of a lipgloss Style
func GetOuterWidth(style lipgloss.Style) int {
	return style.GetHorizontalFrameSize() + style.GetHorizontalPadding()
}

// GetWidth returns the inner width of a lipgloss Style
func GetWidth(style lipgloss.Style) int {
	return style.GetWidth() - GetOuterWidth(style)
}

// GetOuterHeight returns the outer border size of a lipgloss Style
func GetOuterHeight(style lipgloss.Style) int {
	return style.GetVerticalFrameSize() + style.GetVerticalPadding()
}

// GetHeight returns the inner width of a lipgloss Style
func GetHeight(style lipgloss.Style) int {
	return style.GetHeight() - GetOuterHeight(style)
}
