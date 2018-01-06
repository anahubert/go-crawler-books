package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"regexp"
	"strconv"
	"strings"
	"math/rand"
	"encoding/csv"
	"github.com/PuerkitoBio/goquery"
)

type Pages struct {
	count int
}

var ps Pages

func pagination() {
	
	url := "http://www.laguna.rs/s1_spisak_naslova_laguna.html"
	
	doc, err := goquery.NewDocument(url)

	if err != nil {
		fmt.Printf("Unable to scrape %v\n", err)
		ps.count = 0
	}

	lasturl := doc.Find("div.paginacija > a.strelice").Last().AttrOr("href", "unknown")

	if lasturl == "unknown" {

		ps.count = 0

	}

	re := regexp.MustCompile(`s([0-9]+)_spisak_naslova_laguna\.html$`)
	last := re.FindStringSubmatch(lasturl)

	//fmt.Println(lasturl)
	//fmt.Println(last)

	// Extract last page number from url http://www.laguna.rs/s1_spisak_naslova_laguna.html
	ps.count, _ = strconv.Atoi(last[1])
	//fmt.Println(ps.count)
	
	
}

func initf() {

	pagination()
}

func extract(s string, p string) (res string) {
	
	//fmt.Printf("%v\n", p)
	re, _ := regexp.Compile(p)
	
	res = ""

	if re.MatchString(s) == true {

		a := re.FindStringSubmatch(s)
		res = a[1]

    }

	return res
}

func extractBook(s string) (res map[string]string, err error) {

	res = make(map[string]string)

	doc, err := goquery.NewDocument(s)

	if err != nil {

		fmt.Println(s)
		fmt.Println(err)
		return res, err
		
	}
	res["url"] = s

	res["title"] = doc.Find("div.podaci > h1").First().Text()

	res["description"] = strings.TrimSpace(doc.Find("div.sadrzaj > div.podaci_boks").Text())

	res["author"] = doc.Find("div.podaci > h2").First().Text()

	a1, _ := doc.Find("div.korica > a#single_image").First().Attr("href")
	a2, _ := doc.Find("div.korica > a#single_image > img#korica").First().Attr("src")

	res["bigimg"] = "http://www.laguna.rs/" + a1
	res["smallimg"] = "http://www.laguna.rs/" + a2

	res["cena"] = "0"
	res["usteda"] = "0"

	doc.Find("span.cena").Each(func(j1 int, sq7 *goquery.Selection) {

		e := extract(strings.TrimSpace(sq7.Text()), `([0-9.]+)din$`)
		if e != "" {

			res["cena"] = strings.TrimSpace(sq7.Text())

		}
		
	})

	doc.Find("span.usteda").Each(func(j2 int, sq8 *goquery.Selection) {

		e1 := extract(strings.TrimSpace(sq8.Text()), `([0-9.]+)din$`)

		if e1 != "" {

			res["usteda"] = strings.TrimSpace(sq8.Text())

		}
		
	})

	l := doc.Find("div.podaci_boks").First().Find("div").Length()

	doc.Find("div.podaci_boks").First().Find("div").Eq(l-2).First().Find("h3 > a").Each(func(j int, sq1 *goquery.Selection) {

		a := sq1.First().Text()
		res["cats" + strconv.Itoa(j)] = a

	})

    superMap := make(map[string]string)

	dtls := doc.Find("div.podaci_boks").Eq(1).First().Empty()
	
	k3 := 0

	dtls.Each(func(k2 int, sq2 *goquery.Selection) {
		
		details := strings.TrimSpace(sq2.Text())
		
		if details == "Format:" || details == "Povez:" || details == "Broj strana:" || details == "Pismo:" || details == "Godina izdanja:" || details == "ISBN:" {
			
			key := strings.ToLower(strings.Replace(details, ":", "", -1))
			superMap[key] = dtls.Get(k2+1).Data

			k3++
		}

	})

	for k, v := range superMap {
		res[k] = v
	}

	superMap = nil
	//fmt.Printf("%v\n", superMap)

	return res, err

}

func CsvWrite(f *os.File, book map[string]string) {
	//record := []string{strconv.Itoa(book.category.id), string(book.category.name), string(book.category.url)}
	record := []string{
		book["url"],
		book["title"],
		book["author"],
		book["bigimg"],
		book["smallimg"],
		book["cats0"],
		book["cats1"],
		book["cats2"],
		book["cats3"],
		book["cats4"],
		book["format"],
		book["povez"],
		book["broj strana"],
		book["pismo"],
		book["godina izdanja"],
		book["isbn"],
		book["description"],
		book["cena"],
		book["usteda"]}


	w := csv.NewWriter(f)
	if err := w.Write(record); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func list(i int) int {

	filename := "data/laguna/csv/page-" + strconv.Itoa(i) + ".csv"
	//fmt.Printf("%q", filename)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	t := "div.lista_jedna_knjiga > a"
	num := i + 1
	snum := strconv.Itoa(num)

	url := "http://www.laguna.rs/s" + snum + "_spisak_naslova_laguna.html"
	
	//fmt.Println(url)
	
	doc, err := goquery.NewDocument(url)

	if err != nil {

		fmt.Println("Num" + snum)
		fmt.Println(url)
		fmt.Println(err)
		return 0
		
	}

	book := make(map[string]string)

	doc.Find(t).Each(func(j int, s *goquery.Selection) {
	
		// Book link
		bl := s.First().AttrOr("href", "unknown")

		fmt.Println(strconv.Itoa(j) + " http://www.laguna.rs/" + bl)

		book, _ = extractBook("http://www.laguna.rs/" + bl)

		CsvWrite(f, book)
		
		//fmt.Printf("%v\n", book)
		
	})

	amt := time.Duration(rand.Intn(250))
	time.Sleep(time.Millisecond * amt)

	return 1
	
}

func main() {
	
	fmt.Println("START")

	initf() // open first url & count number of Pages
	
	// Iterate through each page & get list of books on each page
	for i:=0; i < ps.count; i++ {

		go list(i)
		
	}

	var input string
	fmt.Scanln(&input)

	fmt.Println("END")
	// Run in concurent mode and process each list
		// Open book on the list and get 
			/*  - Title
				- Url
			*	- Author
			*	- ISBN13
				- Format
				- NumberOfPages
				- Letter
				- Cover
				- Published
				- Category1
				- Category2
				- Category3
				- Category4
				- Category1Url
				- Category2Url
				- Category3Url
				- Category4Url
				- Image1Url
				- Image2Url
				- Description
				- Price
				- PriceDiscount
			*/
	
	
}