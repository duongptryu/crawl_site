package main

import (
	"encoding/csv"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strings"
)

func writeCsv(data [][]string) {
	file, err := os.Create("data.csv")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Println(err)
		}
	}
}

func scapePageData(doc *goquery.Document) [][]string {
	var result [][]string
	doc.Find("ul.srp-results>li.s-item").Each(func(index int, item *goquery.Selection) {
		a := item.Find("a.s-item__link")

		title := strings.TrimSpace(a.Text())
		url, _ := a.Attr("href")

		priceSpan := strings.TrimSpace(item.Find("span.s-item__price").Text())
		price := strings.Trim(priceSpan, " VND ")

		scrapedData := []string{title, price, url}

		result = append(result, scrapedData)
	})

	return result
}

func getHtml(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode >= 400 {
		log.Fatal("Status code of request != 200")
	}

	return resp
}

func main() {
	url := "https://www.ebay.com/sch/i.html?_from=R40&_nkw=beatls+puzzle&_sacat=0&_ipg=200"

	var previousUrl string

	var totalResult [][]string

	for {
		response := getHtml(url)
		defer response.Body.Close()

		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		result := scapePageData(doc)
		totalResult = append(totalResult, result...)

		href, _ := doc.Find("nav.pagination>a.pagination__next").Attr("href")
		if href == "" {
			break
		}
		if href == previousUrl {
			break
		} else {
			url = href
			previousUrl = href
		}
	}

	writeCsv(totalResult)
}
