package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func getPages(url string) int {
	res, _ := http.Get(url)

	if res.StatusCode != 200{
		log.Fatalf("status code: %d %s", res.StatusCode, res.Status)
	}

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	data := doc.Find(".vacancysearch-xs-header-text").Text()
	intdata, _ := strconv.Atoi(strings.Split(data, " ")[1])

	pages := intdata / 50

	fmt.Print(pages)
	return pages
}

func parseUrl(url string, cont *widget.Entry){
	res, _ := http.Get(url)

	if res.StatusCode != 200{
		log.Fatalf("status code: %d %s", res.StatusCode, res.Status)
	}

	doc, _ := goquery.NewDocumentFromReader(res.Body)

	text := cont.Text
	doc.Find(".g-user-content").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find("a").Text()
		vacurl, _:= s.Find("a").Attr("href")
		fmt.Printf("Review %d: %s\n", i, title)
		if title != "" {
			text += title + "\n"
			text += vacurl + "\n\n"
		}
		cont.SetText(text)
	})
}

func parseHHru(pyCont *widget.Entry, goCont *widget.Entry) {
	urlPy := fmt.Sprintf("https://penza.hh.ru/search/vacancy?schedule=remote&clusters=" +
		"true&ored_clusters=true&enable_snippets=true&salary=&text=Python+junior&page=")
	urlGo := fmt.Sprintf("https://penza.hh.ru/search/vacancy?schedule=remote&clusters=" +
		"true&ored_clusters=true&enable_snippets=true&salary=&text=Golang+junior&page=")

	pyPages := getPages(urlPy + strconv.Itoa(0))
	goPages := getPages(urlGo + strconv.Itoa(0))
	for i := 0; i <= pyPages; i++{
		fmt.Print(urlPy + strconv.Itoa(i))
		parseUrl(urlPy + strconv.Itoa(i), pyCont)
	}
	for i := 0; i <= goPages; i++{
		fmt.Print(urlGo + strconv.Itoa(i))
		parseUrl(urlGo + strconv.Itoa(i), goCont)
	}

}

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("Main")
	w.Resize(fyne.NewSize(600, 600))
	w.SetFixedSize(true)


	pyEntry := widget.NewMultiLineEntry()
	pyEntry.Resize(fyne.NewSize(250, 500))
	goEntry := widget.NewMultiLineEntry()
	goEntry.Resize(fyne.NewSize(250, 500))

	grid := container.New(layout.NewGridLayout(2))
	label1 := widget.NewLabelWithStyle("Python", fyne.TextAlignCenter, fyne.TextStyle{})
	label2 := widget.NewLabelWithStyle("Golang", fyne.TextAlignCenter, fyne.TextStyle{})
	grid.Add(label1)
	grid.Add(label2)

	split := container.NewHSplit(pyEntry, goEntry)
	top := container.NewVBox(widget.NewLabelWithStyle("Golang and Python vacancy",
		fyne.TextAlignCenter, fyne.TextStyle{Bold:true}), grid)
	parsebutton := widget.NewButton("Start", func() {
		parseHHru(pyEntry, goEntry)})

	content := container.New(layout.NewBorderLayout(top, parsebutton, nil, nil), top, split, parsebutton )

	w.SetContent(content)

	w.ShowAndRun()
}