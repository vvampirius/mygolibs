package htmlinfo

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
)

var (
	kinoposkFilmIdRegexp = regexp.MustCompile(`://[w\.]+kinopoisk.ru/\w+/(\d+)`)
)

type HtmlInfo struct {
	Url string
	Title string
	Meta []map[string]string
	KinopoiskDescription string
}

func (htmlInfo *HtmlInfo) KinopoiskFilmId() string {
	match := kinoposkFilmIdRegexp.FindStringSubmatch(htmlInfo.Url)
	if len(match) != 2 { return "" }
	return match[1]
}

func (htmlInfo *HtmlInfo) IsKinopoiskFilm() bool {
	if htmlInfo.KinopoiskFilmId() != `` { return true }
	return false
}

func (htmlInfo *HtmlInfo) KinopoiskFilmRatingGif() string {
	filmId := htmlInfo.KinopoiskFilmId()
	if filmId == `` { return "" }
	return fmt.Sprintf("https://www.kinopoisk.ru/rating/%s.gif", filmId)
}

func (htmlInfo *HtmlInfo) GetTitle() string {
	if meta, ok := htmlInfo.findMeta(`property`, `title`); ok {
		if content, ok := meta[`content`]; ok && content != `` { return content }
	}
	return htmlInfo.Title
}

func (htmlInfo *HtmlInfo) findMeta(key, value string) (map[string]string, bool) {
	for _, meta := range htmlInfo.Meta {
		for k, v := range meta {
			if k == key && v == value { return meta, true }
		}
	}
	return map[string]string{}, false
}

func (htmlInfo *HtmlInfo) GetDescription() string {
	if htmlInfo.KinopoiskDescription != `` { return htmlInfo.KinopoiskDescription }
	if meta, ok := htmlInfo.findMeta(`property`, `og:description`); ok {
		if content, ok := meta[`content`]; ok && content != `` { return content }
	}
	if meta, ok := htmlInfo.findMeta(`name`, `description`); ok {
		if content, ok := meta[`content`]; ok && content != `` { return content }
	}
	return ""
}

func (htmlInfo *HtmlInfo) GetPosterUrl() string {
	if meta, ok := htmlInfo.findMeta(`property`, `og:image`); ok {
		if content, ok := meta[`content`]; ok && content != `` { return htmlInfo.fixPosterUrl(content) }
	}
	return ""
}

func (htmlInfo *HtmlInfo) fixPosterUrl(s string) string {
	return `https:`+s
}

func (htmlInfo *HtmlInfo) Save(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(*htmlInfo)
}

func (htmlInfo *HtmlInfo) Load(r io.Reader) error {
	newHtmlInfo := HtmlInfo{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&newHtmlInfo); err != nil { return err }
	*htmlInfo = newHtmlInfo
	return nil
}
