package internal

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []string `xml:"url"`
}

func GetHtmlBody(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	htmlContent := string(body)
	return htmlContent, nil
}

func ValidateUrls(mainUrl string, urls []string) []string {
	var cleanedUrls []string
	for _, url := range urls {
		if strings.HasPrefix(url, "/") {
			cleanedUrls = append(cleanedUrls, mainUrl+url)
		} else if strings.Contains(url, mainUrl) {
			cleanedUrls = append(cleanedUrls, url)
		}
	}
	return cleanedUrls
}

func BuildSiteMapWithDepth(depth int, mainUrl, url string, visitedLinks map[string]struct{}) map[string]struct{} {
	_, ok := visitedLinks[url]
	if ok {
		return visitedLinks
	}
	page, err := GetHtmlBody(url)
	if err != nil {
		log.Fatal(err)
	}
	r := strings.NewReader(page)
	links, err := ParseHtml(r)
	if err != nil {
		log.Fatal(err)
	}
	urls := ExtractUrls(links)
	urls = ValidateUrls(mainUrl, urls)
	visitedLinks[url] = struct{}{}
	if depth > 0 {
		depth--
	}
	if depth == 0 {
		return visitedLinks
	}
	for i := 0; i < len(urls); i++ {
		visitedLinks = BuildSiteMapWithDepth(depth, mainUrl, urls[i], visitedLinks)
	}
	return visitedLinks
}

func BuildXmlRepresentation(visitedLinks map[string]struct{}) string {
	xmlRepresentation := &UrlSet{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	for url, _ := range visitedLinks {
		xmlRepresentation.Urls = append(xmlRepresentation.Urls, url)
	}
	xmlcontent, _ := xml.MarshalIndent(xmlRepresentation, " ", "  ")
	return xml.Header + string(xmlcontent)
}
