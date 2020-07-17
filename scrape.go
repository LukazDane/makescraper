package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

type pageName struct {
	ListedAs string `json:"listedas"`
	CodedAs  string `json:"codedas"`
	Himg     string `json:"img"`
}

func main() {

	pages := []pageName{}

	// Instantiate default collector
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
		// Attach a debugger to the collector, prints found selector and inner info requiested
		// colly.Debugger(&debug.LogDebugger{}),
	)
	// Randomizes user agent to avoid being blocked by server for too many requests...again
	extensions.RandomUserAgent(c)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*http*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
	// Homestuck archive specific, checks each 2nd link on page
	c.OnHTML("body > div.o_content > div > div > div > a:nth-child(2n+2)", func(e *colly.HTMLElement) {
		page := pageName{}
		page.ListedAs = e.Text
		p := e.Attr("href")
		// return page url value to determine true page number
		page.CodedAs = p
		// page.URL = e.Request.URL
		e.Request.Visit(p)
		c.OnHTML("#content_container > div > img", func(e *colly.HTMLElement) {
			link := e.Attr("src")
			// fmt.Println(link)
			page.Himg = link

		})
		pages = append(pages, page)
		// fmt.Printf("Page Title: %q -> %s\n", e.Text, page.ListedAs)
		// fmt.Printf("Page Code: %s\n", page.CodedAs)
	})
	// when progam reaches end, return text to show that it is done
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
	// Site set to visit, explore making this an input field
	c.Visit("https://www.homestuck.com/log/story")
	c.Wait()
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
