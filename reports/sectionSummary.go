package reports

import (
	"time"

	"github.com/AlexNa-Holdings/savva-reports/pdf"
)

func addSectionSummary(doc *pdf.Doc, from, to time.Time) {
	// Add a new section for the summary
	doc.NewSection(doc.T("section_summary"))

}
