package reports

import (
	"image"
	"time"

	"github.com/AlexNa-Holdings/savva-reports/data"
	"github.com/AlexNa-Holdings/savva-reports/pdf"
	"github.com/rs/zerolog/log"
)

type AuthorToShow struct {
	*data.User
	Posts []data.Post
}

const MAX_USERS_TO_SHOW = 5
const MAX_POSTS_FROM_AUTHOR_TO_SHOW = 3

func addSectionAuthors(doc *pdf.Doc, from, to time.Time) {

	if doc.Sponsored == nil {
		var err error
		doc.Sponsored, err = data.GetSponsoredBy(doc.UserAddress)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch sponsored by")
			return
		}
	}

	// Select authors to show
	authors := make([]AuthorToShow, 0)
	for _, s := range doc.Sponsored {

		user, err := data.GetUser(s.Author)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch user data")
			continue
		}

		posts, err := data.GetPostsByAuthor(s.Author, doc.UserAddress, from, to)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch posts by author")
			continue
		}

		if len(posts) == 0 {
			continue // skip this author
		}

		// Sort posts for the best

		if len(posts) > MAX_POSTS_FROM_AUTHOR_TO_SHOW {
			posts = posts[:MAX_POSTS_FROM_AUTHOR_TO_SHOW]
		}

		authors = append(authors, AuthorToShow{User: user, Posts: posts})
	}

	if len(authors) == 0 {
		return // skip the section
	}

	// Sort authors for the best

	if len(authors) > MAX_USERS_TO_SHOW {
		authors = authors[:MAX_USERS_TO_SHOW]
	}

	doc.NewSection(doc.T("authors.title"))

	doc.MarkDownToPdf(doc.T("authors.introduction"))
	doc.NewLine()

	for _, author := range authors {
		doc.AssureVertialSpace(230)
		doc.NewSubSection(author.BestName())

		for _, post := range author.Posts {
			doc.AssureVertialSpace(200)

			save_y := doc.GetY()

			doc.NewSubSubSection(post.GetTitle(doc.Locale))
			// doc.ImageFrom(post.ThumbnailImg, doc.GetX(), doc.GetY(), &gopdf.Rect{W: 160, H: 100})
			doc.DrawImageCover(post.ThumbnailImg, doc.GetX(), doc.GetY(), 160, 100)

			info := ""

			info += "*" + doc.T("posted") + "*: " + post.EffectiveTime.Format(time.RFC822) + "\n"
			info += "*" + doc.T("domain") + "*: " + post.Domain + "\n"

			doc.MarkDownToPdfEx(info, doc.GetX()+170, doc.GetY(),
				doc.PageWidth-doc.Margins.Right-doc.GetX()+170,
				100, false)

			content, err := post.GetContent(doc.Locale)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get content")
				continue
			}

			doc.SetY(save_y + 140)
			doc.NewLine()

			doc.GetImage = func(url string) (image.Image, error) {
				return post.GetImage(url)
			}
			doc.MarkDownToPdf(content)
			doc.GetImage = nil
		}
	}

}
