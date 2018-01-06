package main

import (
	"bitbucket.org/aleksandrah/htmlparser"
	//"bitbucket.org/aleksandrah/log2"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"time"
	"strings"
	"strconv"
	//"calendar"
)
type Html struct {
	Href   string
	Title  string
	Text   string
	Strong string
	Src    string
}
type Page struct {
	url    string
	token  string
	tokens []string
	html   Html
}
type Event struct {
	date string
	time string
	price string
	place string
}
type Movie struct {
	title string
	engtitle string
	country string
	category string
	duration string
	director string
	cast string
	description string
	synopsis string
	prizes string
	events []Event
}

var startfrom = 0
var limit = 115
var separator = ","
var encloser = "\""
var escaper = "\\"
var home = "http://www.fest.rs"
var first = "http://www.fest.rs/FEST-2016/102/Program.shtml/keyword=/product_categories_id=0/product_cinema_id=0/date=0/tagsid=0/limit=5/startfrom=0/sortby=vote/ascdesc=ASC/submit=next/location_in_text[D]=1"
var movie =  Movie{}
var movies =  []Movie{}

func GetPageLinks() (htmls []htmlparser.Html){

	htmls = htmlparser.Scrape(first, home, "#pages > li.number > a")
	fmt.Printf("[%v] Title: %+v\n", time.Now(), htmls)

	return htmls
}

func Details(i int, url string){

	//fmt.Printf("[%v] %+v\n", now, doc.Href)

	ds := htmlparser.ScrapeDocs(url)

	ds.Find("ul.box-P").Each(func(i int, s *goquery.Selection) {

		a := s.Find("li.body > ul > li.desc")

		data := a.Find("ul > li.data")
		projection := a.Find("ul > li.projection")

		title := data.Find("h2").Text()
		eng := data.Find("h3").Text()

		detaillink := s.Find("li.body > .links > a.view")
		detaillink.Attr("href")

		single := detaillink.First()
		href, _ := single.Attr("href")

		dl := htmlparser.ScrapeDocs(home + href)

		movie.synopsis = dl.Find("ul.box-P > li.body > ul > li.desc > ul > li.text > p").First().Text()
		movie.prizes = dl.Find("ul.box-P > li.body > ul > li.desc > ul > li.data > ul > li.v").Last().Text()
		/*country := data.Find("ul > li.v").First().Text()
		category := data.Find("ul > li.v > a").First().Text()
		duration := data.Find("ul > li.v").First().Text()*/


		country := ""
		category := ""
		duration := ""
		director := ""
		cast := ""

		skeys := data.Find("ul > li.k")
		svalues := data.Find("ul > li.v")

		keys := []string{}
		values := []string{}

		skeys.Each(func(k int, sk *goquery.Selection) {

			keys = append(keys, sk.Text())

		})
		svalues.Each(func(v int, sv *goquery.Selection) {

			values = append(values, sv.Text())

		})

		for key, value := range keys {

			if value == "Država:" {
				country = values[key]
			}else if value == "Program:" {
				category = values[key]
			}else if value == "Trajanje:" {
				duration = values[key]
			}else if value == "Režija:" {
				director = values[key]
			}else if value == "Uloge:" {
				cast = values[key]
			}
		}

		events := []Event{}

		price_time := ""
		location := ""

		projection.Find("p").Each(func(i1 int, s1 *goquery.Selection) {

			price_time = s1.Find("span").First().Text()

			res := strings.Split(price_time, "|")

			location = s1.Find("a").First().Text()

			event := Event{date: res[0], time: res[1], price: res[2], place: location}

			events = append(events, event)

		})


		movie.title = title
		movie.engtitle = eng
		movie.country = country
		movie.category = category
		movie.duration = duration
		movie.director = director
		movie.cast = cast
		movie.events = events

		movies = append(movies, movie)

	})



}

func print() {

	fmt.Printf("[%v] Title;TitleEng;Country;Category;Duration;Director;Cast;Synopsis;Prizes;Price;Date;Time;Location \n", time.Now(), movie.title)

	for _, movie := range movies {

		fmt.Printf("[%v] Title: %+v;", time.Now(), movie.title)
		fmt.Printf("[%v] Tile Eng: %+v;", time.Now(), movie.engtitle)
		fmt.Printf("[%v] Country: %+v;", time.Now(), movie.country)
		fmt.Printf("[%v] Category: %+v;", time.Now(), movie.category)
		fmt.Printf("[%v] Duration: %+v;", time.Now(), movie.duration)
		fmt.Printf("[%v] Director: %+v;", time.Now(), movie.director)
		fmt.Printf("[%v] Cast: %+v;", time.Now(), movie.cast)
		fmt.Printf("[%v] Synopsis: %+v;", time.Now(), movie.synopsis)
		fmt.Printf("[%v] Prizes: %+v;", time.Now(), movie.prizes)

		for _, ev := range movie.events{

			fmt.Printf("[%v] Price:%+v;", time.Now(), ev.price)
			fmt.Printf("[%v] Date:%+v;", time.Now(), ev.date)
			fmt.Printf("[%v] Time:%+v;", time.Now(), ev.time)
			fmt.Printf("[%v] Location: %+v;", time.Now(), ev.place)

		}
		fmt.Printf("[%v]\n", time.Now())
	}

}

func printCsv() {

	//fmt.Printf("Title;TitleEng;Country;Category;Duration;Director;Cast;Synopsis;Prizes;Price;Date;Time;Location\n")
	fmt.Printf("Subject,Start Date,Start Time,End Date,End Time,All Day,Description,Location,UID\n")

	for _, movie := range movies {

		fmt.Printf("\"%+v\";", movie.title)
		fmt.Printf("\"%+v\";", movie.engtitle)
		fmt.Printf("\"%+v\";", movie.country)
		fmt.Printf("\"%+v\";", movie.category)
		fmt.Printf("\"%+v\";", movie.duration)
		fmt.Printf("\"%+v\";", movie.director)
		fmt.Printf("\"%+v\";", movie.cast)
		fmt.Printf("\"%+v\";", movie.synopsis)
		fmt.Printf("\"%+v\";", movie.prizes)

		for _, ev := range movie.events{

			fmt.Printf("\"%+v\";", ev.price)
			fmt.Printf("\"%+v\";", ev.date)
			fmt.Printf("\"%+v\";", ev.time)
			fmt.Printf("\"%+v\";", ev.place)

		}
		fmt.Printf("\n")
	}

}

func printCsvCal() {

	fmt.Printf("Subject,Start Date,Start Time,End Date,End Time,All Day Event,Description,Location,Private\n")

	const dform = "2016-03-30"
	const tform = "22:00:00"

	for _, movie := range movies {

		for _, ev := range movie.events{

			event_start_date := strings.TrimSpace(ev.date) //02 Mar 2016
			event_start_time := strings.TrimSpace(ev.time)

			parts_dates := strings.Split(event_start_date, " ")

			d := "1"
			m := "Apr"

			if len(parts_dates) != 2 {
				fmt.Printf("Error: %q", parts_dates)
			}else{
				d = parts_dates[0]
				m = strings.TrimSpace(parts_dates[1])
			}

			r := strings.NewReplacer("'", "", ".", "") // 60 min
			sdur, _ := strconv.ParseInt(r.Replace(strings.TrimSpace(movie.duration)), 10, 64)
			d = r.Replace(d)

			t, _ := time.Parse(time.RFC822Z, d + " " + m + " 16 " + event_start_time + " +0100")

			t1 := time.Unix(t.Unix() + sdur*60, 0)

			fmt.Printf("%+v%v", strconv.Quote(movie.title), separator)
			//fmt.Printf("%q%v", t, separator)
			//fmt.Printf("%q%v", t1, separator)
			fmt.Printf("%q%v", t.Format("01/02/2006"), separator)
			fmt.Printf("%q%v", t.Format("15:04:05"), separator)
			fmt.Printf("%q%v", t1.Format("01/02/2006"), separator)
			fmt.Printf("%q%v", t1.Format("15:04:05"), separator)
			fmt.Printf("%v%v", "False", separator)
			fmt.Printf("%s%v", "", separator)
			fmt.Printf("%+v%v", strconv.Quote(strings.TrimSpace(ev.place)), separator)
			fmt.Printf("%+v", "True")
			fmt.Printf("\n")

		}
	}

}

func main() {


	for i:=0; i<limit; i+=5 {

		url := "http://www.fest.rs/FEST-2016/102/Program.shtml/keyword=/product_categories_id=0/product_cinema_id=0/date=0/tagsid=0/limit=5/startfrom=" + strconv.Itoa(i) + "/sortby=vote/ascdesc=ASC/submit=next/location_in_text[D]=1"

		Details(i, url)

	}

	//fmt.Printf("[%v] Number of movies: %+v\n", time.Now(), len(movies))
	printCsvCal()
	//var input string
	//fmt.Scanln(&input)

}
