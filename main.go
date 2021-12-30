package main

import (
	"encoding/csv"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strings"
)

func writeCsv(data []string) {
	fileName := "data.csv"

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 777)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(data)
	if err != nil {
		log.Fatalln(err)
	}
}

func scapePageData(doc *goquery.Document) {
	doc.Find("ul.srp-results>li.s-item").Each(func (index int, item *goquery.Selection){
		a := item.Find("a.s-item__link")

		title := strings.TrimSpace(a.Text())
		url, _ := a.Attr("href")

		priceSpan := strings.TrimSpace(item.Find("span.s-item__price").Text())
		price := strings.Trim(priceSpan, " VND ")

		scrapedData := []string{title, price, url}

		writeCsv(scrapedData)
	})
}

func getHtml(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode >= 400 {
		log.Fatal("Status code of request != 2000")
	}

	return resp
}

func main () {
	url := "https://www.ebay.com/sch/i.html?_from=R40&_nkw=beatls+puzzle&_sacat=0&_ipg=200"

	var previousUrl string

	for {
		response := getHtml(url)
		defer response.Body.Close()

		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		scapePageData(doc)

		href, _:= doc.Find("nav.pagination>a.pagination__next").Attr("href")
		if href == previousUrl {
			break
		}else {
			url = href
			previousUrl = href
		}
	}
}

