package reports

import "github.com/AlexNa-Holdings/savva-reports/pdf"

func addSectionLegal(doc *pdf.Doc) {
	doc.NewSection(doc.T("legal_notice_title"))
	doc.MarkDownToPdf(doc.T("legal_notice"))
}
