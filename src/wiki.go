package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"html/template"
)

const (
	PagePath string = "../res/pages/"
	PageViewRoute string = "/view/"
	PageEditRoute string = "/edit/"
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
	p, err := loadPage(title)
	if err != nil {
		http.ServeFile(w, r, PagePath + "NotFound")
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(PageEditRoute):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(PagePath + tmpl + ".html")
	t.Execute(w, p)
}

func getWorkingDir() string {
	path, _ := os.Getwd()
	return path + "/"
}

func main()  {
	http.HandleFunc(PageViewRoute, viewHandler)
	http.HandleFunc(PageEditRoute, editHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
