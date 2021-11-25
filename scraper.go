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

	// return html element
	c.OnHTML(".e1nlzfl3", func(e *colly.HTMLElement) {
		// locate all information
		productName := e.ChildText(".css-1bjwylw")
		description := e.ChildText(".css-wfq7u")
		imageLink := e.ChildText(".css-79elbk img")
		price := e.ChildText(".css-o5uqvq")
		rating := e.ChildText(".css-153qjw7 span")
		merchantName := e.ChildText(".css-vbihp9 span")

		fmt.Println("imageLink")
		fmt.Println(imageLink)

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

	// scrapping page 10 times (get a 100 handphone)
	// for i := 1; i < 11; i++ {
	// 	fmt.Printf("Scrapping Page : %d\n", i)

	// 	c.Visit("https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=" + strconv.Itoa(i))
	// }
	c.Visit("https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=3")

}
