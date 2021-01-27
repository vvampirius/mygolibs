package htmlinfo

import (
	"github.com/PuerkitoBio/goquery"
	"io"
)


func Parse(r io.Reader, url string) (HtmlInfo, error) {
	document, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		ErrorLogger.Println(url, err.Error())
		return HtmlInfo{}, err
	}

	htmlInfo := HtmlInfo{
		Url: url,
		Meta: make([]map[string]string, 0),
	}

	if head := document.Find(`head`); len(head.Nodes) == 1 {
		htmlInfo.Title = head.Find(`title`).Text()

		head.Find(`meta`).Each(func(_ int, s *goquery.Selection){
			meta := make(map[string]string)
			for _, attr := range s.Nodes[0].Attr {
				if attr.Key == `charset` { return }
				if attr.Key == `name` && attr.Val == `viewport` { return }
				if attr.Key == `data-tid` { continue }
				meta[attr.Key] = attr.Val
			}
			if Debug { DebugLogger.Println(meta) }
			htmlInfo.Meta = append(htmlInfo.Meta, meta)
		})
	}

	if htmlInfo.IsKinopoiskFilm() {
		htmlInfo.KinopoiskDescription = document.Find(`.film-details-block p`).Text()
		if Debug { DebugLogger.Println(htmlInfo.KinopoiskDescription) }
	}

	return htmlInfo, nil
}
