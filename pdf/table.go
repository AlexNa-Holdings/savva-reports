package pdf

import (
	"github.com/rs/zerolog/log"
)

type Table struct {
	Cells     [][]string
	Header    []string
	W, H      int
	ColWidths []float64
	ColAllign []uint8 // 'L', 'C', 'R'
	colStyle  []Style
}

var HeaderStyle = &Style{
	FontName:  "Arial",
	FontSize:  12,
	FontColor: Color{0, 0, 0}}

func NewTable() *Table {
	t := &Table{
		W: 0,
		H: 0,
	}
	return t
}

func (t *Table) SetW(w int) {
	t.W = w
	t.ColWidths = make([]float64, w)
	t.ColAllign = make([]uint8, w)
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
		total_width = doc.GetMarginWidth()
		for i := 0; i < t.W; i++ {
			if t.ColWidths[i] == 0 {
				t.ColWidths[i] = total_width / float64(t.W)
			}
		}
	} else {
		table_x = doc.margin_left + (doc.GetMarginWidth()-total_width)/2
	}

	header_height := doc.estimateHeaderHeight(t)

	if header_height > 0 {
		// write header
		for i, text := range t.Header {
			doc.TextCentered(text, table_x+t.ColWidths[i]/2, doc.cy+header_height)
			table_x += t.ColWidths[i]
		}
		doc.cy += header_height
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
