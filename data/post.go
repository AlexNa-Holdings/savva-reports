package data

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"image"
	"regexp"
	"strings"
	"time"

	"github.com/AlexNa-Holdings/savva-reports/cmn"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type SavvaContent_v2_0_Locale_Chapter struct {
	Title    string `yaml:"title"`     //
	Data     string `yaml:"data"`      // embedded content data
	DataPath string `yaml:"data_path"` // path to the content data
}

type SavvaContent_v2_0_Locale struct {
	Title       string                             `yaml:"title"`        //
	TextPreview string                             `yaml:"text_preview"` // short preview to show in the lists
	Categories  []string                           `yaml:"categories"`   //
	Tags        []string                           `yaml:"tags"`         //
	Data        string                             `yaml:"data"`         // embedded content data
	DataPath    string                             `yaml:"data_path"`    // path to the content data
	Chapters    []SavvaContent_v2_0_Locale_Chapter `yaml:"chapters"`
	// NSFW        bool                               `yaml:"nsfw"`
}

type Receipient_v2_0 struct {
	Pass     string `yaml:"pass"`
	Modifier string `yaml:"modifier"`
}

type SavvaContent_v2_0 struct {
	SpecVersion string `yaml:"savva_spec_version"` // "2.0" Savva Content Specification Version
	MimeType    string `yaml:"mime_type"`          //
	Data        string `yaml:"data"`               // embedded content data
	DataPath    string `yaml:"data_path"`          // path to the content data
	Thumbnail   string `yaml:"thumbnail"`          // embedded thumbnail data
	RootCID     string `yaml:"root_savva_cid"`     // root Savva CID
	ParentCID   string `yaml:"parent_savva_cid"`   // parent Savva CID
	NSFT        bool   `yaml:"nsfw"`               // Not Safe For Work
	Fundraiser  int64  `yaml:"fundraiser"`         // Fundraiser ID
	Encryption  struct {
		Type              string                     `yaml:"type"`
		KeyExchangeAlg    string                     `yaml:"key_exchange_alg"`
		KeyExchangePubKey string                     `yaml:"key_exchange_pub_key"`
		Receipients       map[string]Receipient_v2_0 `yaml:"receipients"`
	} `yaml:"encryption"`

	Locales map[string]SavvaContent_v2_0_Locale `yaml:"locales"`
}

type Post struct {
	SavvaCid      string
	ShortCid      sql.NullString
	AuthorAddr    string
	PosterAddr    sql.NullString
	Domain        string
	Guid          string
	Ipfs          string
	TimeStamp     time.Time
	EffectiveTime time.Time
	TotalChilds   int
	c_v2_0        SavvaContent_v2_0
	ThumbnailImg  image.Image
}

func GetPostsByAuthor(address, sponsor string, from, to time.Time) ([]Post, error) {

	r := make([]Post, 0)

	rows, err := cmn.C.DB.Query(`
		SELECT savva_cid, short_cid, author_addr, poster_addr, domain, guid, ipfs, time_stamp, effective_time, total_childs
		FROM savva_content
		WHERE
		savva_content.content_type = 'post' 
		AND EXISTS ( SELECT 1 FROM clubs_members cm WHERE cm.domain = savva_content.domain AND cm.author_addr = savva_content.author_addr AND cm.member_addr = $2 )
		AND author_addr = $1 AND effective_time BETWEEN $3 AND $4
	`, address, sponsor, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.SavvaCid, &p.ShortCid, &p.AuthorAddr, &p.PosterAddr, &p.Domain, &p.Guid, &p.Ipfs, &p.TimeStamp, &p.EffectiveTime, &p.TotalChilds); err != nil {
			return nil, err
		}

		err := LoadPostInfo(&p)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to load post %s", p.SavvaCid)
			continue
		}

		err = p.LoadThumbnail()
		if err != nil {
			log.Error().Err(err).Msgf("Failed to load thumbnail for post %s", p.SavvaCid)
			continue
		}

		r = append(r, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return r, nil
}

func LoadPostInfo(post *Post) error {
	fileContent := cmn.C.IPFS(post.Ipfs + "/info.yaml")
	if fileContent == nil {
		return fmt.Errorf("failed to load post %s", post.SavvaCid)
	}

	// determine the spec verison
	// pattern := `"savva_spec_version"\s*:\s*"([^"]+)"`
	// if isDir {
	pattern := `savva_spec_version:\s*['"]?(\d+\.\d+)['"]?`
	// }

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(string(fileContent))

	if len(matches) == 0 {
		log.Error().Msgf("Cannot determine SavvaSpecVersion. CID:(%s)", post.Ipfs)
		return fmt.Errorf("Cannot determine SavvaSpecVersion. CID:(%s)", post.Ipfs)
	}
	version := matches[1]

	log.Trace().Msgf("SavvaSpecVersion: %s", version)

	switch {
	case version == "2.0":
		err := yaml.Unmarshal(fileContent, &post.c_v2_0)

		log.Trace().Msgf("spec version 2.0: %s", post.c_v2_0.SpecVersion)

		if err != nil {
			log.Error().Msgf("Cannot unmarshal info.yaml from IPFS. CID:(%s) Error:(%v)", post.Ipfs, err)
			return err
		}

		return nil
	default:
		return errors.New("Unknown SavvaSpecVersion: " + version)
	}
}

func (p *Post) LoadThumbnail() error {
	if p.c_v2_0.Thumbnail == "" {
		user, err := GetUser(p.AuthorAddr)
		if err != nil {
			log.Error().Msgf("Failed to load user %s", p.AuthorAddr)
			return fmt.Errorf("failed to load user %s", p.AuthorAddr)
		}
		p.ThumbnailImg = user.AvatarImg
		return nil
	}

	dp := strings.TrimSpace(p.c_v2_0.Thumbnail)

	if !strings.HasPrefix(dp, "/") {
		dp = "/" + dp
	}

	thumbnail := cmn.C.IPFS(p.Ipfs + dp)
	if thumbnail == nil {
		log.Error().Msgf("Failed to load thumbnail for post %s", p.SavvaCid)
		return fmt.Errorf("failed to load thumbnail for post %s", p.SavvaCid)
	}

	var err error
	p.ThumbnailImg, _, err = image.Decode(bytes.NewReader(thumbnail))
	if err != nil {
		log.Error().Msgf("Failed to decode thumbnail for post %s", p.SavvaCid)
		return fmt.Errorf("failed to decode thumbnail for post %s", p.SavvaCid)
	}

	return nil
}

func (p *Post) GetLocale(locale string) (SavvaContent_v2_0_Locale, bool) {
	if p.c_v2_0.SpecVersion == "2.0" {
		l, ok := p.c_v2_0.Locales[locale]
		if ok {
			return l, true
		}

		l, ok = p.c_v2_0.Locales["en"]
		if ok {
			return l, true
		}

		// return the first locale
		for _, l := range p.c_v2_0.Locales {
			return l, true
		}

		return l, false
	}
	return SavvaContent_v2_0_Locale{}, false
}

func (p *Post) GetTitle(locale string) string {
	if p.c_v2_0.SpecVersion == "2.0" {

		l, ok := p.GetLocale(locale)
		if ok {
			return l.Title
		}
	}
	return ""
}

func (p *Post) GetContent(locale string) (string, error) {
	if p.c_v2_0.SpecVersion == "2.0" {
		l, ok := p.GetLocale(locale)
		if ok {

			if l.Data == "" {

				dp := strings.TrimSpace(l.DataPath)

				if !strings.HasPrefix(dp, "/") {
					dp = "/" + dp
				}

				content := cmn.C.IPFS(p.Ipfs + dp)
				if content == nil {
					return "", fmt.Errorf("failed to load content for post %s", p.SavvaCid)
				}
				l.Data = string(content)
			}
			return l.Data, nil
		}
	}
	return "", fmt.Errorf("failed to get content for post %s", p.SavvaCid)
}

func (p *Post) GetImage(url string) (image.Image, error) {
	if p.c_v2_0.SpecVersion == "2.0" {
		l, ok := p.GetLocale("en")
		if ok {
			if l.Data == "" {

				if !strings.HasPrefix(url, "/") {
					url = "/" + url
				}

				content := cmn.C.IPFS(p.Ipfs + url)
				if content == nil {
					return nil, fmt.Errorf("failed to load content for post %s", p.SavvaCid)
				}

				img, _, err := image.Decode(bytes.NewReader(content))
				if err != nil {
					return nil, fmt.Errorf("failed to decode image for post %s", p.SavvaCid)
				}
				return img, nil
			}
		}
	}
	return nil, fmt.Errorf("failed to get image for post %s", p.SavvaCid)
}
