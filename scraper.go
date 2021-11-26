package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

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

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:57.0) Gecko/20100101 Firefox/57.0")
		fmt.Println("Visiting : ", r.URL)
	})

	// counter for items
	counter := 0
	// counter to get merchant name
	counter2 := 1

	// return html element
	c.OnHTML(".e1nlzfl3", func(e *colly.HTMLElement) {
		// locate all information
		productName := e.ChildText(".css-1bjwylw")
		description := e.ChildText(".css-wfq7u")
		imageLink := e.ChildText(".css-16vw0vn > img")
		price := e.ChildText(".css-o5uqvq")

		// count the star
		var rating int
		e.ForEach(".css-177n1u3", func(index int, j *colly.HTMLElement) {
			rating++
		})

		// get merchant name only without merchant location
		var merchantName string
		e.ForEach(".css-vbihp9 > .css-1kr22w3", func(index int, k *colly.HTMLElement) {
			if counter2%2 == 0 {
				merchantName = k.Text
			}
			counter2++
		})

		//write information to csv file
		writer.Write([]string{
			productName,
			description,
			imageLink,
			price,
			strconv.Itoa(rating),
			merchantName,
		})
		counter++
	})

	//handle error on request
	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Status :", r.StatusCode)
		fmt.Println("Error :", e)
	})

	// loop to scrap 100 items
	for i := 1; counter < 100; i++ {
		fmt.Printf("Scrapping Page : %d\n", i)

		c.Visit("https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=" + strconv.Itoa(i))
	}
	// c.Visit("https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=1")

}
