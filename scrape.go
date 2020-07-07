package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type pageName struct {
	ListedAs string `json:"listedas"`
	CodedAs  string `json:"codedas"`
}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {

	pages := []pageName{}

	// Instantiate default collector
	c := colly.NewCollector()

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	e.Request.Visit(e.Attr("href"))
	// })

	// c.OnHTML("body > div.o_content > div > div > div > br", func(e *colly.HTMLElement) {
	// 	fmt.Println("PAge and Date: ", e.Text)
	// })

	// On every a element which has href attribute call callback
	c.OnHTML("body > div.o_content > div > div > div > a:nth-child(2n+2)", func(e *colly.HTMLElement) {
		page := pageName{}
		page.ListedAs = e.Text
		p := e.Attr("href")
		page.CodedAs = strings.TrimPrefix(p, "/story/")
		// Print link
		pages = append(pages, page)
		fmt.Printf("Page Title: %q -> %s\n", e.Text, page.ListedAs)
		fmt.Printf("Page Code: %s\n", page.CodedAs)
	})

	c.OnXML("//h1", func(e *colly.XMLElement) {
		fmt.Println(e.Text)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://www.homestuck.com/log/story")
	serializeToJSON(pages)

}
func writeFile(file []byte) {
	if err := ioutil.WriteFile("output.json", file, 0644); err != nil {
		log.Fatalf("Unable to write file! %v", err)
	}
}
func serializeToJSON(n []pageName) {
	fmt.Println("Serializing Pages to JSON...")
	serialized, _ := json.Marshal(n)
	writeFile(serialized)
	fmt.Println("Successfully serialized the Homestuck to 'output.json'")
}
