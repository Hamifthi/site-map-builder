package main

import (
	"flag"
	"fmt"
	"sitemapbuilder/internal"
)

func main() {
	url := flag.String("url", "https://www.calhoun.io/",
		"The root url to start the search from")
	depth := flag.Int("depth", -1,
		"The maximum depth to go through and build the site map. if not specified extract all of links")
	flag.Parse()
	visitedUrls := make(map[string]struct{})
	visitedUrls = internal.BuildSiteMapWithDepth(*depth, *url, *url, visitedUrls)
	xmlOutput := internal.BuildXmlRepresentation(visitedUrls)
	fmt.Println(xmlOutput)
}
