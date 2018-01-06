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
	
	url := "http://www.delfi.rs/knjige/knjige_1.html"
	
	doc, err := goquery.NewDocument(url)

	if err != nil {
		fmt.Printf("Unable to scrape %v\n", err)
		ps.count = 0
	}

	lasturl := doc.Find("div> #paginacija > table > tbody > tr > td > a").Last().AttrOr("href", "unknown")
	fmt.Println(url)
	fmt.Println(lasturl)
	if lasturl == "unknown" {

		ps.count = 0

	}

	re := regexp.MustCompile(`knjige_([0-9]+).html$`)
	last := re.FindStringSubmatch(lasturl)

	fmt.Println(lasturl)
	fmt.Println(last)

	// Extract last page number from url http://www.laguna.rs/s1_spisak_naslova_laguna.html
	ps.count, _ = strconv.Atoi(last[1])
	//fmt.Println(ps.count)
	
	
}

func initf() {

	pagination()
	//os.Exit(1)
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

func matched(s string, p string) (res bool) {

	re, _ := regexp.Compile(p)
	res = false

	if re.MatchString(s) == true {

		res = true

    }

	return res
}

func extractBook1(s string) (res map[string]string, err error) {
	//fmt.Println(s)
	return res, err
}

func extractBook(s string) (res map[string]string, err error) {

	res = make(map[string]string)

	doc, err := goquery.NewDocument(s)

	if err != nil {

		fmt.Println("Error:" + s)
		fmt.Println(err)
		return res, err
		
	}
	res["url"] = s

	res["title"] = doc.Find("#art_podaci > h1").First().Text()

	res["description"] = strings.TrimSpace(doc.Find("div > #art_opis > div").Text())

	res["author"] = doc.Find("div > #art_podaci > h2").First().Text()

	a1, _ := doc.Find("div > #art_korica > a").First().Attr("href")
	a2, _ := doc.Find("div > #art_korica > a > img").First().Attr("src")

	res["bigimg"] = "http://www.delfi.rs/" + a1
	res["smallimg"] = "http://www.delfi.rs/" + a2

	res["old_price"] = "0"
	res["new_price"] = "0"
	res["saving"] = "0"

	tmp := strings.TrimSpace(doc.Find("div > #art_stara_cena").First().Text())
	res["old_price"] = extract(tmp, "Cena:\\s([0-9\\.]+)\\sdin")

	tmp = doc.Find("div > #art_nova_cena").First().Text()
	res["new_price"] = extract(strings.TrimSpace(tmp), "([0-9\\.]+)\\sdin")

	tmp = doc.Find("div > #art_usteda").First().Text()
	res["saving"] = extract(strings.TrimSpace(tmp), "uštedite:\\s([0-9\\.]+)\\sdin")

	doc.Find("div > #art_podaci > h3").Each(func(i int, s *goquery.Selection) {

		txt := strings.TrimSpace(s.First().Text())

		if i == 0 {
			res["cats"] = extract(txt, "Žanrovi:\\s([\\pL,\\.\\s]+)$")
		}
		if i == 1 {
			
			res["publisher"] = extract(txt, "Izdavač:\\s([\\pL,\\.\\s]+)$")

		}

	})

	doc.Find("div > #art_podaci_div > div").Each(func(i int, s *goquery.Selection) {

		txt := strings.TrimSpace(s.First().Text())

		if matched(txt, "Pismo:\\s([\\pL,\\.\\s]+)$") {
			res["letter"] = strings.TrimSpace(extract(txt, "Pismo:\\s([\\pL,\\.\\s]+)$"))
		}
		if matched(txt, "Povez:\\s([\\pL,\\.\\s]+)$") {
			res["cover"] = strings.TrimSpace(extract(txt, "Povez:\\s([\\pL,\\.\\s]+)$"))
		}
		if matched(txt, "Broj strana:\\s([0-9\\.]+)") {
			res["number_of_pages"] = strings.TrimSpace(extract(txt, "Broj strana:\\s([0-9\\.]+)"))
		}
		if matched(txt, "Format:\\s([A-Za-z0-9,\\.x]+)(\\s)?[^cm]?") == true {
			//fmt.Println(txt)
			res["format"] = strings.TrimSpace(extract(txt, "Format:\\s([A-Za-z0-9,\\.x]+)(\\s)?[^cm]?"))
		}
		
	})

	return res, err

}

func CsvWrite(f *os.File, book map[string]string) {
	
	record := []string{
		book["url"],
		book["title"],
		book["author"],
		book["bigimg"],
		book["smallimg"],
		book["cats"],
		book["format"],
		book["cover"],
		book["number_of_pages"],
		book["letter"],
		book["published"],
		book["isbn"],
		book["publisher"],
		book["description"],
		book["old_price"],
		book["new_price"],
		book["saving"]}


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

	filename := "data/delfi/csv/page-" + strconv.Itoa(i) + ".csv"
	//fmt.Printf("%q", filename)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	t := "#table_spisa_artikala > table > tbody > tr > td > a"
	num := i + 1
	snum := strconv.Itoa(num)

	url := "http://www.delfi.rs/knjige/knjige_" + snum + ".html"
	
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

		fmt.Println(strconv.Itoa(j) + " http://www.delfi.rs/" + bl)

		book, _ = extractBook("http://www.delfi.rs/" + bl)

		CsvWrite(f, book)
		
		fmt.Printf("%v\n", book)
		
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