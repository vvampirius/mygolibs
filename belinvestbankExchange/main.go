package belinvestbankExchange

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"golang.org/x/net/html"
)

var (
	LOGGER = log.New(os.Stderr, `belinvestbankExchange#`, log.Lshortfile)
	URL = `https://www.belinvestbank.by/exchange-rates/courses-tab-cashless`
	ErrNotFound = errors.New(`HTML element not found`)
)

type Currency struct {
	Id string
	Nominal float64
	Buy float64
	Sell float64
}

func parseAndSetFloat(src string, dst *float64) error {
	src = strings.ReplaceAll(src, "\n", ``)
	src = strings.ReplaceAll(src, ` `, ``)
	f, err := strconv.ParseFloat(src, 64)
	if err != nil { LOGGER.Println(err.Error()) }
	*dst = f
	return nil
}

func isInSlice(ss []string, s string) bool {
	for _, v := range ss {
		if v == s { return true }
	}
	return false
}

func getNodeAttributeValue(attributes []html.Attribute, key string) string {
	for _, attribute := range attributes {
		if attribute.Key == key { return attribute.Val }
	}
	return ""
}

func getElementNodes(node *html.Node, name string) []*html.Node {
	nodes := make([]*html.Node, 0)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.ElementNode { continue }
		if child.Data != name { continue }
		nodes = append(nodes, child)
	}
	return nodes
}

func getDiv(node *html.Node, id string) *html.Node {
	if node.Type == html.ElementNode && node.Data == `div` &&
		getNodeAttributeValue(node.Attr, `id`) == id { return node }
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if div := getDiv(child, id); div != nil { return div }
	}
	return nil
}

func getTbody(node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == `tbody` { return node }
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if div := getTbody(child); div != nil { return div }
	}
	return nil
}

func getCurrencies(node *html.Node) (map[string]Currency, error) {
	currencies := make(map[string]Currency)
	for _, tr := range getElementNodes(node, `tr`) {
		currency := Currency{}
		for _, td := range getElementNodes(tr, `td`) {
			classes := strings.Split(getNodeAttributeValue(td.Attr, `class`), ` `)
			if isInSlice(classes, `courses-table__td_nominal`) {
				if err := parseAndSetFloat(td.FirstChild.Data, &currency.Nominal); err != nil { return nil, err }
			}
			if isInSlice(classes, `courses-table__td_buy`) {
				if err := parseAndSetFloat(td.FirstChild.Data, &currency.Buy); err != nil { return nil, err}
			}
			if isInSlice(classes, `courses-table__td_sell`) {
				if err := parseAndSetFloat(td.FirstChild.Data, &currency.Sell); err != nil { return nil, err }
			}
			if isInSlice(classes, `courses-table__td_iso`) {
				currency.Id = strings.ReplaceAll(td.FirstChild.Data, "\n", ``)
				currency.Id = strings.ReplaceAll(currency.Id, ` `, ``)
			}
		}
		currencies[currency.Id] = currency
	}
	return currencies, nil
}

// use nil instead reader to make http request
func Get(r io.Reader) (map[string]Currency, error) {
	if r == nil {
		response, err := http.Get(URL)
		if err != nil {
			LOGGER.Println(err.Error())
			return nil, err
		}
		r = response.Body
		defer response.Body.Close()
	}

	document, err := html.Parse(r)
	if err != nil {
		LOGGER.Println(err.Error())
		return nil, err
	}

	div := getDiv(document, `courses-tab-cashless-content`)
	if div == nil {
		LOGGER.Println(ErrNotFound.Error())
		return nil, err
	}

	tbody := getTbody(div)
	if tbody == nil {
		LOGGER.Println(ErrNotFound.Error())
		return nil, err
	}

	return getCurrencies(tbody)
}