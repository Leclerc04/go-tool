package svgutil

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/leclerc04/go-tool/agl/util/must"
)

// Manager controls a glyphMap stroring unicode:glyph pairs
type Manager struct {
	unitsPerEm float32
	glyphMap   map[string]glyph
}

// NewManager returns a new manager
func NewManager() *Manager {
	m := &Manager{
		glyphMap: make(map[string]glyph),
	}
	m.loadFont()
	return m
}

type glyphStruct struct {
	// attr
	D         string `xml:"d,attr"`
	Unicode   string `xml:"unicode,attr"`
	HorizAdvX string `xml:"horiz-adv-x,attr"`
	GlyphName string `xml:"glyph-name,attr"`
}

type fontFaceStruct struct {
	// attr
	UnitsPerEm string `xml:"units-per-em,attr"`
	FontFamily string `xml:"font-family,attr"`
}

type fontStruct struct {
	// attr
	HorizAdvX string `xml:"horiz-adv-x,attr"`
	// field
	FontFace fontFaceStruct `xml:"font-face"`
	Glyphs   []glyphStruct  `xml:"glyph"`
}

type svgStruct struct {
	Font fontStruct `xml:"defs>font"`
}

type glyph struct {
	d         string
	horizAdvX float32
}

// loadFont could load svg format font files, but now we hard-code the "Microsoft-Yahei.svg" here
func (m *Manager) loadFont() {
	var svg svgStruct
	var fontBytes []byte
	var err error

	fontBytes, err = Asset("pkg/util/svgutil/data/Microsoft-Yahei.svg")
	must.Must(err)

	err = xml.Unmarshal(fontBytes, &svg)
	must.Must(err)
	v, err := strconv.ParseFloat(svg.Font.FontFace.UnitsPerEm, 32)
	must.Must(err)
	m.unitsPerEm = float32(v)

	for _, g := range svg.Font.Glyphs {
		// <glyph> has no unicode symbol
		if g.Unicode == "" {
			continue
		}

		horizAdvX := g.HorizAdvX
		if g.HorizAdvX == "" {
			horizAdvX = svg.Font.HorizAdvX
		}
		v, err := strconv.ParseFloat(horizAdvX, 32)
		must.Must(err)
		m.glyphMap[g.Unicode] = glyph{
			d:         g.D,
			horizAdvX: float32(v),
		}
	}
}

func (m *Manager) textToPathWithSvgWidth(text string, svgWidth, size float32) ([]string, float32, float32) {
	// make sure there's no \n in text
	text = strings.Replace(text, "\n", " ", -1)

	line := text
	var previousIndex int
	var horizAdvY float32
	var linePaths []string
	for {
		var horizAdvX float32
		var paths []string
		horizAdvY += m.unitsPerEm

		if previousIndex == len(line) {
			break
		}
		line = line[previousIndex:]
		for i, c := range line {
			g := m.glyphMap[string(c)]
			paths = append(paths, fmt.Sprintf(`<path transform="translate(%f) rotate(180) scale(-1, 1)" d="%s" />`, horizAdvX, g.d))
			horizAdvX += g.horizAdvX
			if svgWidth-horizAdvX*size < m.unitsPerEm*size || i == len(line)-1 {
				previousIndex = i + 1
				break
			}
		}
		linePaths = append(linePaths, fmt.Sprintf(`<g transform="scale(%f) translate(0,%f)">
%s
</g>`, size, horizAdvY, strings.Join(paths, "\n")))
	}

	return linePaths, svgWidth, horizAdvY * size
}

func (m *Manager) textToPathWithoutSvgWidth(text string, size float32) ([]string, float32, float32) {
	lines := strings.Split(text, "\n")
	var horizAdvY float32
	var biggestX float32
	var linePaths []string
	for _, line := range lines {
		var horizAdvX float32
		var paths []string
		horizAdvY += m.unitsPerEm
		for _, c := range line {
			char := string(c)
			g := m.glyphMap[char]
			paths = append(paths, fmt.Sprintf(`<path transform="translate(%f) rotate(180) scale(-1, 1)" d="%s" />`, horizAdvX, g.d))
			horizAdvX += g.horizAdvX
			if biggestX < horizAdvX {
				biggestX = horizAdvX
			}
		}
		linePaths = append(linePaths, fmt.Sprintf(`<g transform="scale(%f) translate(0,%f)">
%s
</g>`, size, horizAdvY, strings.Join(paths, "\n")))
	}

	return linePaths, biggestX * size, horizAdvY * size
}

// TextToPath converts text to svg path, if svgWidth is specified, will auto wrap words; else, caller should make sure text is word-wrapped
func (m *Manager) TextToPath(text string, svgWidth, asize float32) string {
	var svgHeight float32
	var linePaths []string
	size := asize / m.unitsPerEm

	if svgWidth > 0 {
		linePaths, svgWidth, svgHeight = m.textToPathWithSvgWidth(text, svgWidth, size)
	} else {
		linePaths, svgWidth, svgHeight = m.textToPathWithoutSvgWidth(text, size)
	}

	return fmt.Sprintf(`<svg height="%fpx" width="%fpx">
%s
</svg>`, svgHeight, svgWidth, strings.Join(linePaths, "\n"))
}
