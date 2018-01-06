package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Image struct {
	big   string
	small string
}
type Html struct {
	href   string
	title  string
	text   string
	strong string
	src    string
}
type Page struct {
	url    string
	token  string
	tokens []string
	html   Html
}

type Book struct {
	id              int
	url             string
	details         []Html
	page            Page
	prices          []Html
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
	page  Page
	books []Book
}

var (
	home       = Page{"http://www.knjizare-vulkan.rs", "", []string{}, Html{}}
	category   = Category{0, "", "", 0, "", Page{"", "#container_left > #categories_section > li > ul > li > a", []string{}, Html{}}, []Book{}}
	book       = Book{0, "", []Html{}, Page{"", "#content > ul.book_list > li.item > figure > a", []string{}, Html{}}, []Html{}, "", "", map[string]string{}, "", "", "", "", "", "", "", "", "", "", "", Category{}, Category{}, Image{}}
	pagin      = Page{"", "", []string{"#content > #pagination_holder > #pagination_holder_right > a.page_nav_link", "#content > #pagination_holder > #page_nav_right > a.page_nav_link"}, Html{}}
	details    = Page{"", "#content > ul.book_list > li.item > figure > a", []string{}, Html{}}
	categories = []Category{}
)

func Href(sel *goquery.Selection) (href string) {

	single := sel.First()
	href, _ = single.Attr("href")

	return href
}

func Title(sel *goquery.Selection) (title string) {

	single := sel.First()
	title = single.Text()

	return title
}

func Strong(sel *goquery.Selection) (result string) {

	single := sel.First()
	result = single.Find("strong").Text()

	return result
}

func Text(sel *goquery.Selection) (result string) {
	single := sel.First()
	result = single.Text()
	return result
}

func Src(sel *goquery.Selection) (result string) {
	single := sel.First()
	result, _ = single.Attr("src")

	return result
}

func Scrape(page Page) (result []Html) {

	var href string
	var title string
	var text string
	var strong string
	var src string

	doc, err := goquery.NewDocument(page.url)
	doc.Find(page.token).Each(func(i int, s *goquery.Selection) {

		href = Href(s)
		title = Title(s)
		text = Text(s)
		strong = Strong(s)
		src = Src(s)

		html := Html{href, title, text, strong, src}

		result = append(result, html)
	})

	//fmt.Printf("%+v\n", result)

	if err != nil {
		log.Fatal(err)
	}

	return result

}

func Last(page Page) (last int) {

	var err error = nil
	var result []Html

	pagin.token = pagin.tokens[0]
	pagin.url = page.url
	//fmt.Printf("%+v\n", pagin)
	result = Scrape(pagin)

	if len(result) == 0 {
		pagin.token = pagin.tokens[1]
		result = Scrape(pagin)
	}

	last = 1
	lent := len(result) - 1

	if lent > 0 {
		last, err = strconv.Atoi(result[lent].text)
	}

	//fmt.Printf("Lent: %v, Last: %v\n", lent, last)

	if err != nil {
		log.Fatal(err)
	}

	return last

}

func Books(p int, page Page, last int) /*(books []Book)*/ {

	filename := strconv.Itoa(p) + "-books.csv"
	//fmt.Printf("%q", filename)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	last += 1

	url := page.url

	for i := 1; i < last; i++ {

		s := strconv.Itoa(i)

		page.token = book.page.token
		page.url = url + "/" + s

		//fmt.Printf("i: %d, last: %d,  url: %s", i, last, page.url)

		htmls := Scrape(page)

		for _, html := range htmls {
			page.html = html

			book := PageBook(page)

			str := FormatBook(book)

			WriteString(f, str)

			//amt := time.Duration(rand.Intn(250))
			//time.Sleep(time.Millisecond * amt)
			//books = append(books, book)
		}
		//fmt.Printf("i: %d, last: %d,  url: %s\n", i, last, page.url)

	}

	amt := time.Duration(rand.Intn(250))
	time.Sleep(time.Millisecond * amt)

	//return books

}

func FormatBook(book Book) (str string) {

	str = fmt.Sprintf("%d,%q,%q,", book.category.id, book.category.name, book.category.url)
	str += fmt.Sprintf("%d,%q,%q,", book.subcategory.id, book.subcategory.name, book.subcategory.url)
	str += fmt.Sprintf("%d,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q",
		book.id, book.title, book.url, book.author, book.publisher, book.published_web, book.isbn13, book.published, book.format, book.number_of_pages, book.paperback, book.letter, book.price, book.price_web, book.description, book.image.small, book.image.big)
	str += fmt.Sprintf("\n")

	return str
}

func WriteString(f *os.File, text string) {

	if _, err := f.WriteString(text); err != nil {
		panic(err)
	}

}
func PrintBook(book Book) {

	fmt.Printf("%d,%q,%q,", book.category.id, book.category.name, book.category.url)
	fmt.Printf("%d,%q,%q,", book.subcategory.id, book.subcategory.name, book.subcategory.url)
	fmt.Printf("%d,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q,%q",
		book.id, book.title, book.url, book.author, book.publisher, book.published_web, book.isbn13, book.published, book.format, book.number_of_pages, book.paperback, book.letter, book.price, book.price_web, book.description, book.image.small, book.image.big)
	fmt.Printf("\n")
}

/*
:[{href: title:Autor: Ivan Kleut  text:Autor: Ivan Kleut  strong:Autor:} {href: title:Izdavač: GRAĐEVINSKA KNJIGA  text:Izdavač: GRAĐEVINSKA KNJIGA  strong:Izdavač:} {href: title:Na sajtu od: 24.12.2012. text:Na sajtu od: 24.12.2012. strong:Na sajtu od:} {href: title:ISBN: 9788639505271  text:ISBN: 9788639505271  strong:ISBN:} {href: title:Godina izdanja: 2007 text:Godina izdanja: 2007 strong:Godina izdanja:} {href: title:Format: 17x24  text:Format: 17x24  strong:Format:} {href: title:Broj strana: 328 text:Broj strana: 328 strong:Broj strana:} {href: title:Povez: Tvrd  text:Povez: Tvrd  strong:Povez:} {href: title:Pismo: Latinica  text:Pismo: Latinica  strong:Pismo:}]
*/
func BookDetails(htmls []Html, pattern string) (s string) {

	re := regexp.MustCompile(pattern)

	for _, h := range htmls {

		if re.FindString(h.text) != "" {

			a := re.Split(h.text, 2)
			return strings.TrimSpace(a[1])
		}

	}

	return "0"

}

func PageBook(page Page) (book Book) {

	//time.Sleep(time.Second * time.Duration(5))

	page.url = home.url + page.html.href
	page.token = "ul.book_rev_info > li"

	details := Scrape(page)
	//author,publisher,published_web,isbn13,published,format,number_of_pages,paperback,letter,price,price_web
	book.author = BookDetails(details, "Autor:")
	book.publisher = BookDetails(details, "Izdavač:")
	book.published_web = BookDetails(details, "Na sajtu od:")
	book.isbn13 = BookDetails(details, "ISBN:")
	book.published = BookDetails(details, "Godina izdanja:")
	book.format = BookDetails(details, "Format:")
	book.number_of_pages = BookDetails(details, "Broj strana:")
	book.paperback = BookDetails(details, "Povez:")
	book.letter = BookDetails(details, "Pismo:")

	page.token = "figure.book_rev > div.holder > h3"
	titles := Scrape(page)

	page.token = "#prices > div.book_rev_price"
	prices := Scrape(page)

	book.price = BookDetails(prices, "Cena:")
	book.price_web = BookDetails(prices, "Cena na sajtu:")

	page.token = "figure.book_rev > div.description"
	descriptions := Scrape(page)

	bookid := UrlID(page.url)
	book.id = bookid
	book.title = titles[0].text
	book.url = page.url
	book.description = descriptions[0].text

	page.token = "figure.book_rev > div.img > a"
	images_small := Scrape(page)

	image := Image{}

	if len(images_small) == 0 {
		image.small = ""
	} else {
		image.small = home.url + images_small[0].href
	}

	page.token = "figure.book_rev > div.img > a > img"
	images_big := Scrape(page)

	if len(images_big) != 0 {
		image.big = home.url + images_big[0].src

	} else {

		page.token = "figure.book_rev > div.img > img"
		images_big = Scrape(page)
		if len(images_big) != 0 {
			image.big = home.url + images_big[0].src
		} else {
			image.big = ""
		}
	}

	book.image = image

	page.token = "#breadcrumbs_list > li > h2 > a"
	categories := Scrape(page)

	page.token = "#breadcrumbs_list > li > a"
	subcategories := Scrape(page)

	r := strings.NewReplacer(">", "", "\n", "")

	cat := categories[0]
	subcat := subcategories[0]

	cat_name := r.Replace(cat.text)
	subcat_name := r.Replace(subcat.text)

	cat_link := home.url + cat.href
	subcat_link := home.url + subcat.href

	cat_id := UrlID(cat.href)
	subcat_id := UrlID(subcat.href)

	category1 := Category{}
	category2 := Category{}

	category1.id = cat_id
	category1.url = cat_link
	category1.name = cat_name

	category2.id = subcat_id
	category2.url = subcat_link
	category2.name = subcat_name

	book.category = category1
	book.subcategory = category2

	//fmt.Printf("Book: %+v\n", book)

	return book

}

func UrlID(url string) (id int) {
	uid := 0
	res := strings.Split(url, "/")
	size := len(res)
	last := size - 1

	if last > 3 {
		uid, _ = strconv.Atoi(res[last])
	}

	return uid

}

func main() {

	var last int = 1

	page := Page{}
	page.url = home.url
	page.token = category.page.token
	//var numCPU = runtime.GOMAXPROCS(0)
	//ch := make(chan []Book, numCPU)

	htmls := Scrape(page)

	//fmt.Printf("category_id,category_name,category_url,subcat_id,subcat_name,subcat_url,id,title,url,author,publisher,published_web,isbn13,published,format,number_of_pages,paperback,letter,price,price_web,description,image_small,image_big")
	//fmt.Printf("\n")

	for i, html := range htmls {

		if i != 1 && i != 4 && i != 6 {
			continue
		}
		page.url = home.url + html.href
		last = Last(page)

		//fmt.Printf("Book: %+v\n", page.url)
		//fmt.Printf("Last: %d - %d\n", i, last)

		//books := Books(page, last)

		go Books(i, page, last)
		//fmt.Printf("Books: %+v\n", books)
	}

	var input string
	fmt.Scanln(&input)

}
