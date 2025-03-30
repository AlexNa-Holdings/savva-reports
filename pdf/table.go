package pdf

import (
	"github.com/rs/zerolog/log"
)

type Table struct {
	Cells     [][]string
	Header    []string
	W, H      int
	ColWidths []float64
	ColStyle  []Style
}

var HeaderStyle = &Style{
	FontName:  "DejaVuBold",
	FontSize:  12,
	FontColor: &Color{0, 0, 0},
	BGColor:   &SAVVA_COLOR,
	Align:     'C',
	Padding:   Padding{left: 5, top: 4, right: 5, bottom: 6},
}

var DefaultStyle = &Style{
	FontName:  "Arial",
	FontSize:  12,
	FontColor: &Color{0, 0, 0},
	Align:     'L',
}

func NewTable() *Table {
	t := &Table{
		W: 0,
		H: 0,
	}
	return t
}

func (t *Table) SetW(w int) {
	t.W = w

	if len(t.Cells) != w {
		t.ColWidths = make([]float64, w)
	}

	if len(t.ColStyle) != w {
		t.ColStyle = make([]Style, w)
		for i := 0; i < w; i++ {
			t.ColStyle[i] = *DefaultStyle
		}
	}
}

func (t *Table) SetHeader(s ...string) {
	if t.H == 0 && t.W == 0 {
		t.SetW(len(s))
	}

	if len(s) != t.W {
		log.Error().Msgf("Invalid number of columns: %d, expected: %d", len(s), t.W)
		return
	}
	t.Header = s
}

func (t *Table) AddRow(s ...string) {

	if t.H == 0 && t.W == 0 {
		t.SetW(len(s))
	}

	if len(s) != t.W {
		log.Error().Msgf("Invalid number of columns: %d, expected: %d", len(s), t.W)
		return
	}
	t.Cells = append(t.Cells, s)
	t.H++
}

func (doc *Doc) WriteTable(t *Table) {
	if t.W == 0 || t.H == 0 {
		log.Error().Msg("Invalid table: no columns or rows")
		return
	}

	n_zero_width_columns := 0
	total_width := 0.
	table_x := doc.margin_left
	for i := 0; i < t.W; i++ {
		total_width += t.ColWidths[i]
		if t.ColWidths[i] == 0 {
			n_zero_width_columns++
		}
	}

	if n_zero_width_columns > 0 {
		// use all page width for the table
		// and spread 0 width columns evenly
		for i := 0; i < t.W; i++ {
			if t.ColWidths[i] == 0 {
				t.ColWidths[i] = (doc.GetMarginWidth() - total_width) / float64(n_zero_width_columns)
			}
		}
		total_width = doc.GetMarginWidth()
	} else {
		table_x = doc.margin_left + (doc.GetMarginWidth()-total_width)/2
	}

	table_y := doc.GetY()

	header_height := doc.estimateHeaderHeight(t)

	if header_height > 0 {
		//Fiil the bg color
		if HeaderStyle.BGColor != nil {
			doc.SetFillColor(HeaderStyle.BGColor.R, HeaderStyle.BGColor.G, HeaderStyle.BGColor.B)
			doc.Rectangle(table_x, table_y-+header_height, table_x+total_width, table_y, "F", 0, 0)
		}

		x := table_x
		for i, text := range t.Header {
			doc.writeTextInWidth(text, x, table_y, t.ColWidths[i], HeaderStyle)
			x += t.ColWidths[i]
		}
		doc.SetY(doc.GetY() + header_height)
	}

}

func (doc *Doc) estimateHeaderHeight(t *Table) float64 {
	if t.Header == nil {
		return 0
	}
	header_height := 0.
	for i, text := range t.Header {
		header_height = max(header_height, doc.estimateTextHeight(text, t.ColWidths[i], HeaderStyle))
	}

	return header_height
}
