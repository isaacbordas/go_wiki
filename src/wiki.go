package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	PagePath string = "../res/pages/"
	PageViewRoute string = "/view/"
)

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
	filename := p.Title
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := getWorkingDir() + PagePath + title
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request)  {
	title := r.URL.Path[len(PageViewRoute):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func getWorkingDir() string {
	path, _ := os.Getwd()
	return path + "/"
}

func main()  {
	http.HandleFunc(PageViewRoute, viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
