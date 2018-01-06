package main

import (
	"encoding/csv"
	"fmt"
	//"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"math/rand"
	"os"
	//"regexp"
	//"bufio"
	"strconv"
	//"strings"
	"net/http"
	"time"
)

type Image struct {
	big   string
	small string
}

type Book struct {
	id              int
	url             string
	title           string
	description     string
	overview        map[string]string
	author          string
	publisher       string
	published_web   string
	isbn13          string
	published       string
	format          string
	number_of_pages string
	paperback       string
	letter          string
	price           string
	price_web       string
	category        Category
	subcategory     Category
	image           Image
}

type Category struct {
	id    int
	url   string
	name  string
	index int
	link  string
	//page  Page
	books []Book
}

func ImageDownload(id string, url string) {

	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}

	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create("images-books/" + id + ".jpg")
	if err != nil {
		log.Fatal(err)
	}
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	//fmt.Println(url)
	//fmt.Println("Success!")

	//amt := time.Duration(rand.Intn(250))
	//time.Sleep(time.Millisecond * amt)
}

/*func ImagesDownload(

}*/

func Books(p int) {

	filename := "scraped-books/" + strconv.Itoa(p) + "-books.csv"
	fmt.Println(filename)

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	r := csv.NewReader(file)
	r.LazyQuotes = true
	r.FieldsPerRecord = 20
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	tmp := [][]string{}

	for i, record := range records {

		j := len(record) - 1
		//fmt.Printlin(record[j])
		tmp[record[3]] = record[j]

		if i == 10 {
			//ImageDownload(record[3], record[j])
			go ImagesDownload(tmp)
		}
	}

	amt := time.Duration(rand.Intn(250))
	time.Sleep(time.Millisecond * amt)

	//return books

}

func main() {

	for i := 0; i < 28; i++ {

		go Books(i)

	}

	var input string
	fmt.Scanln(&input)

}
