package pdf

import (
	_ "embed"
)

const (
	commonFontName   = "common"
	commonFontSize   = 15
	commonLineHeight = 10

	coverTitleFontName = "cover"

	titleFontName   = "title"
	titleFontSize   = 24
	titleLineHeight = 20
	titleLn         = 10

	subTitleFontName      = "subtitle"
	subTitleFontSize      = 22
	subTitleLineHeight    = 15
	subTitleMinSpaceAfter = 35
	subTitleTopMargin     = 10
	subtitleBotMargin     = 5

	monoFontName = "mono"
)

//go:embed fonts/Raleway/Raleway-Medium.ttf
var fontCoverTitle []byte

//go:embed fonts/Raleway/Raleway-ExtraBold.ttf
var fontCoverTitleBold []byte

//go:embed fonts/Roboto/Roboto-Regular.ttf
var fontCommon []byte

//go:embed fonts/Roboto/Roboto-Bold.ttf
var fontCommonBold []byte

//go:embed fonts/Roboto/Roboto-Italic.ttf
var fontCommonItalic []byte

//go:embed fonts/Roboto/Roboto-BoldItalic.ttf
var fontCommonBoldItalic []byte

//go:embed fonts/Montserrat/Montserrat-Bold.ttf
var fontTitle []byte

var fontSubTitle []byte = fontCommonBold

//go:embed fonts/JetBrainsMono/JetBrainsMonoNL-Regular.ttf
var fontMono []byte

func (c *Converter) setFonts() {
	// cover font
	c.pdf.AddUTF8FontFromBytes(coverTitleFontName, "", fontCoverTitle)
	c.pdf.AddUTF8FontFromBytes(coverTitleFontName, "B", fontCoverTitleBold)

	// common font
	c.pdf.AddUTF8FontFromBytes(commonFontName, "", fontCommon)
	c.pdf.AddUTF8FontFromBytes(commonFontName, "B", fontCommonBold)
	c.pdf.AddUTF8FontFromBytes(commonFontName, "I", fontCommonItalic)
	c.pdf.AddUTF8FontFromBytes(commonFontName, "BI", fontCommonBoldItalic)

	// title font
	c.pdf.AddUTF8FontFromBytes(titleFontName, "", fontTitle)

	// subtitle font
	c.pdf.AddUTF8FontFromBytes(subTitleFontName, "", fontSubTitle)

	// mono font for page numbers
	c.pdf.AddUTF8FontFromBytes(monoFontName, "", fontMono)
}
