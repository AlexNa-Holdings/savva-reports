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
	FontColor: &Color{0xf7, 0xf7, 0xf7},
	BGColor:   &SAVVA_DARK_COLOR,
	Align:     'C',
	Padding:   Padding{left: 5, top: 4, right: 5, bottom: 6},
}

var DefaultStyle = &Style{
	FontName:  "Arial",
	FontSize:  12,
	FontColor: &Color{0, 0, 0},
	Padding:   Padding{left: 5, top: 4, right: 0, bottom: 6},
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
	y := table_y
	last_row_y := y

	header_height := doc.estimateHeaderHeight(t)
	first_row_height := doc.estimateRowHeight(t, 0)

	if doc.GetY()+header_height+first_row_height > doc.pageHeight-doc.margin_bottom {
		doc.NextPage()
		table_y = doc.GetY()
		y = table_y
	}

	if header_height > 0 { // Write header
		if HeaderStyle.BGColor != nil {
			doc.SetFillColor(HeaderStyle.BGColor.R, HeaderStyle.BGColor.G, HeaderStyle.BGColor.B)
			doc.SetStrokeColor(HeaderStyle.BGColor.R, HeaderStyle.BGColor.G, HeaderStyle.BGColor.B)
			doc.SetLineWidth(0.5)
			doc.Rectangle(table_x, y-header_height, table_x+total_width, y, "DF", 0, 0)
		}

		x := table_x
		for i, text := range t.Header {
			doc.writeTextInWidth(text, x, y, t.ColWidths[i], HeaderStyle)
			x += t.ColWidths[i]
		}
		y += header_height
	}

	for i, row := range t.Cells {
		row_height := doc.estimateRowHeight(t, i)
		if y+row_height > doc.pageHeight-doc.margin_bottom {
			// may be add later ..continue to next page

			// Draw the border around the table
			doc.SetStrokeColor(HeaderStyle.BGColor.R, HeaderStyle.BGColor.G, HeaderStyle.BGColor.B)
			doc.SetLineWidth(0.5)
			// vertical lines
			var v_x = table_x
			for i := 0; i < t.W; i++ {
				doc.Line(v_x, table_y, v_x, last_row_y)
				v_x += t.ColWidths[i]
			}
			doc.Line(v_x, table_y, v_x, last_row_y)                        // right
			doc.Line(table_x, last_row_y, table_x+total_width, last_row_y) // bottom

			doc.NextPage()
			table_y = doc.GetY()
			y = table_y
			if header_height > 0 { // Write header
				if HeaderStyle.BGColor != nil {
					doc.SetFillColor(HeaderStyle.BGColor.R, HeaderStyle.BGColor.G, HeaderStyle.BGColor.B)
					doc.SetStrokeColor(HeaderStyle.BGColor.R, HeaderStyle.BGColor.G, HeaderStyle.BGColor.B)
					doc.SetLineWidth(0.5)
					doc.Rectangle(table_x, y-header_height, table_x+total_width, y, "DF", 0, 0)
				}

				x := table_x
				for i, text := range t.Header {
					doc.writeTextInWidth(text, x, y, t.ColWidths[i], HeaderStyle)
					x += t.ColWidths[i]
				}
				y += header_height
			}
		}

		x := table_x
		for j, text := range row {
			doc.writeTextInWidth(text, x, y, t.ColWidths[j], &t.ColStyle[j])
			x += t.ColWidths[j]
		}

		last_row_y = y
		y += row_height
	}

	// Draw the border around the table
	doc.SetStrokeColor(HeaderStyle.BGColor.R, HeaderStyle.BGColor.G, HeaderStyle.BGColor.B)
	doc.SetLineWidth(0.5)
	// vertical lines
	var v_x = table_x
	for i := 0; i < t.W; i++ {
		doc.Line(v_x, table_y, v_x, last_row_y)
		v_x += t.ColWidths[i]
	}
	doc.Line(v_x, table_y, v_x, last_row_y)                        // right
	doc.Line(table_x, last_row_y, table_x+total_width, last_row_y) // bottom

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

func (doc *Doc) estimateRowHeight(t *Table, row int) float64 {
	if row >= len(t.Cells) {
		return 0
	}

	row_height := 0.
	for i, text := range t.Cells[row] {
		row_height = max(row_height, doc.estimateTextHeight(text, t.ColWidths[i], &t.ColStyle[i]))
	}

	return row_height
}
