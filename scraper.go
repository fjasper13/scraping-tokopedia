package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

var test string

func main() {
	// MAKE CSV FILE
	file, err := os.Create("product.csv")
	check(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write column header
	header := []string{"Product Name", "Description", "Image Link", "Price", "Rating", "Merchant Name"}
	writer.Write(header)

	c := colly.NewCollector()
	detailProduct := c.Clone()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:57.0) Gecko/20100101 Firefox/57.0")
		fmt.Println("Visiting : ", r.URL)
	})

	// return html element
	c.OnHTML(".e1nlzfl3", func(e *colly.HTMLElement) {
		// locate all information
		productName := e.ChildText(".css-1bjwylw")
		description := e.ChildText(".css-wfq7u")
		imageLink := e.ChildText(".css-79elbk img")
		price := e.ChildText(".css-o5uqvq")
		rating := e.ChildText(".css-153qjw7 span")
		merchantName := e.ChildText(".css-vbihp9 span")

		productLink := e.ChildAttr("div.e1nlzfl3 > a", "href")
		productLink = e.Request.AbsoluteURL(productLink)
		detailProduct.Visit(productLink)

		//write information to csv file
		writer.Write([]string{
			productName,
			description,
			imageLink,
			price,
			rating,
			merchantName,
		})
	})

	//handle error on request
	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Status :", r.StatusCode)
		fmt.Println("Error :", e)
	})

	detailProduct.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:57.0) Gecko/20100101 Firefox/57.0")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		fmt.Println("Visiting : ", r.URL)
	})

	detailProduct.OnError(func(r *colly.Response, e error) {
		fmt.Println("Error To Get Detail Product")
		fmt.Println("Status :", r.StatusCode)
		fmt.Println("Error :", e)
	})

	detailProduct.OnHTML("div.css-41d95w", func(h *colly.HTMLElement) {
		test = h.ChildText(".css-xi606m")
		// test = h.ChildText(".css-xi606m h5")
		fmt.Println("test")
		fmt.Println(test)
	})

	// scrapping page 10 times (get a 100 handphone)
	// for i := 1; i < 11; i++ {
	// 	fmt.Printf("Scrapping Page : %d\n", i)

	// 	c.Visit("https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=" + strconv.Itoa(i))
	// }
	c.Visit("https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=1")

}
