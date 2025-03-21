package pdf

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/jung-kurt/gofpdf/v2"
)

const (
	marginLeft  = 15
	marginTop   = 15
	marginRight = 15

	paragraphLn     = 2
	paragraphIndent = "        "

	coverPageColorR = 200
	coverPageColorG = 0
	coverPageColorB = 0
)

type Converter struct {
	pdf *gofpdf.Fpdf

	w float64
	h float64

	marginLeft  float64
	marginTop   float64
	marginRight float64
	bold        bool
	italic      bool
	indent      bool
	alignment   string

	skipUnknownElems bool
	unknownElems     map[string]struct{}
}

func NewConverter(skip bool) *Converter {
	pdf := gofpdf.New("P", "mm", "A4", "")
	w, h := pdf.GetPageSize()
	return &Converter{
		pdf:              pdf,
		w:                w,
		h:                h,
		marginLeft:       marginLeft,
		marginTop:        marginTop,
		marginRight:      marginRight,
		skipUnknownElems: skip,
		unknownElems:     map[string]struct{}{},
	}
}

func (c *Converter) WritePDF(inPath string, outPath string) error {
	c.setFonts()
	c.pdf.SetMargins(c.marginLeft, c.marginTop, c.marginRight)

	f, err := os.Open(inPath)
	if err != nil {
		return fmt.Errorf("failed to open original book: %w", err)
	}
	defer f.Close()

	err = c.parse(f)
	if err != nil {
		return fmt.Errorf("faield to parse original book: %w", err)
	}

	for name := range c.unknownElems {
		fmt.Printf("unknown element: %s\n", name)
	}

	return c.pdf.OutputFileAndClose(outPath)
}

func (c *Converter) parse(source io.Reader) error {
	decoder := xml.NewDecoder(source)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read xml token: %w", err)
		}
		switch t := token.(type) {

		case xml.StartElement:
			switch t.Name.Local {

			case "p":
				c.indent = true

			case "strong":
				c.bold = true
				c.updateStyle()
			case "emphasis":
				c.italic = true
				c.updateStyle()

			case "title":
				var title Title
				err := decoder.DecodeElement(&title, &t)
				if err != nil {
					return fmt.Errorf("failed to decode title: %w", err)
				}
				c.writeTitle(&title)

			case "subtitle":
				var subTitle SubTitle
				err := decoder.DecodeElement(&subTitle, &t)
				if err != nil {
					return fmt.Errorf("failed to decode title: %w", err)
				}
				c.writeSubTitle(&subTitle)

			case "epigraph":
				c.pdf.SetLeftMargin(c.w / 3)
				c.alignment = "R"
				c.italic = true
				c.updateStyle()

			case "poem":
				c.alignment = "C"
				c.italic = true
				c.updateStyle()

			case "stanza":

			case "v":

			case "description":
				var descr Description
				if err := decoder.DecodeElement(&descr, &t); err != nil {
					return fmt.Errorf("failed to decode description: %w", err)
				}
				c.writeCoverPage(&descr.TitleInfo)
				c.pdf.AddPage()
				if err := c.parse(strings.NewReader(descr.TitleInfo.Annotation.Content)); err != nil {
					return fmt.Errorf("failed to parse annotation: %w", err)
				}
				c.pdf.SetFooterFunc(func() {
					c.pdf.SetFont(monoFontName, "", 12)
					c.pdf.SetY(-15)
					c.pdf.SetX(c.w / 4)
					c.pdf.CellFormat(c.w/2, 10, fmt.Sprintf("%d", c.pdf.PageNo()), "", 0, "C", false, 0, "")
				})

			case "FictionBook":
			case "body":
			case "section":
			default:
				c.unknownElems[t.Name.Local] = struct{}{}
				if c.skipUnknownElems {
					if err := decoder.DecodeElement(&struct{}{}, &t); err != nil {
						return fmt.Errorf("failed to decode unknown element: %w", err)
					}
				}
			}

		case xml.EndElement:
			switch t.Name.Local {

			case "p":
				// c.pdf.Write(commonLineHeight, "\n")
				c.pdf.Ln(commonLineHeight + paragraphLn)
			case "strong":
				c.bold = false
				c.updateStyle()
			case "emphasis":
				c.italic = false
				c.updateStyle()

			case "epigraph":
				c.pdf.SetLeftMargin(marginRight)
				c.pdf.Ln(commonLineHeight)
				c.alignment = "L"
				c.italic = false
				c.updateStyle()

			case "poem":
				c.alignment = "L"
				c.italic = false
				c.updateStyle()

			case "stanza":
				c.pdf.Ln(commonLineHeight)

			case "v":
				c.pdf.Write(commonLineHeight, "\n")

			case "body":
				return nil
			}

		case xml.CharData:
			s := string(t)
			if len(strings.TrimSpace(s)) == 0 {
				continue
			}
			if c.alignment == "" {
				if c.indent {
					s = paragraphIndent + s
					c.indent = false
				}
				c.pdf.Write(commonLineHeight, s)
			} else {
				c.pdf.WriteAligned(0, commonLineHeight, s, c.alignment)
			}
		}
	}

	return nil
}

func (c *Converter) writeCoverPage(info *TitleInfo) {
	c.pdf.AddPage()

	author := fmt.Sprintf("%s %s", info.Author.FirstName, info.Author.LastName)
	c.pdf.SetFont(coverTitleFontName, "", 30)
	c.pdf.SetY(50)
	c.pdf.MultiCell(0, 20, author, "", "C", false)

	title := strings.ToUpper(info.Title)
	c.pdf.SetFont(coverTitleFontName, "B", 32)
	c.pdf.Ln(30)
	c.pdf.MultiCell(0, 20, title, "", "C", false)

	c.pdf.SetDrawColor(coverPageColorR, coverPageColorG, coverPageColorB)
	c.pdf.SetLineWidth(10)
	c.pdf.Line(-10, 170, 220, 170)

	c.pdf.SetFont(commonFontName, "", 18)
	c.pdf.SetY(-50)
	c.pdf.MultiCell(0, 20, info.Date, "", "C", false)

	c.pdf.SetFont(commonFontName, "", commonFontSize)
}

func (c *Converter) writeTitle(title *Title) {
	c.pdf.AddPage()

	c.pdf.Bookmark(title.Pgs[0].Text, 0, c.pdf.GetY())

	c.pdf.SetFont(titleFontName, "", titleFontSize)
	for i := range title.Pgs {
		c.pdf.WriteAligned(0, titleLineHeight, title.Pgs[i].Text, "C")
		c.pdf.Ln(titleLineHeight)
	}
	c.pdf.Ln(titleLn)
	c.pdf.SetFont(commonFontName, "", commonFontSize)
}

func (c *Converter) writeSubTitle(subTitle *SubTitle) {
	c.pdf.SetFont(subTitleFontName, "", subTitleFontSize)
	c.pdf.Ln(subTitleTopMargin)

	lines := math.Ceil(c.pdf.GetStringWidth(subTitle.Text) / (c.w * 2 / 3))

	if c.pdf.GetY()+subTitleLineHeight*lines+subtitleBotMargin+subTitleMinSpaceAfter > c.h-c.marginTop {
		c.pdf.AddPage()
	}
	c.pdf.Bookmark(subTitle.Text, 1, c.pdf.GetY())
	c.pdf.SetX(c.w / 6)
	c.pdf.MultiCell(2*c.w/3, subTitleLineHeight, subTitle.Text, "", "C", false)

	c.pdf.Ln(subtitleBotMargin)
	c.pdf.SetFont(commonFontName, "", commonFontSize)

}

func (c *Converter) updateStyle() {
	if c.bold {
		if c.italic {
			c.pdf.SetFontStyle("BI")
		} else {
			c.pdf.SetFontStyle("B")
		}
	} else {
		if c.italic {
			c.pdf.SetFontStyle("I")
		} else {
			c.pdf.SetFontStyle("")
		}
	}
}
