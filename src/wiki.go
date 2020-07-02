package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

const (
	PagePath string = "../res/pages/"
	PageViewRoute string = "/view/"
	PageEditRoute string = "/edit/"
	PageSaveRoute string = "/save/"
)

type Page struct {
	Title string
	Body []byte
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

var templates = template.Must(template.ParseFiles(PagePath + "edit.html", PagePath + "view.html"))

func (p *Page) save() error {
	filename := PagePath + p.Title
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

func viewHandler(w http.ResponseWriter, r *http.Request, title string)  {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, PageEditRoute + title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		sendHttpError(w, err)
		return
	}
	http.Redirect(w, r, PageViewRoute + title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl + ".html", p)
	if err != nil {
		sendHttpError(w, err)
		return
	}
}

func sendHttpError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func getWorkingDir() string {
	path, _ := os.Getwd()
	return path + "/"
}

func main()  {
	http.HandleFunc(PageViewRoute, makeHandler(viewHandler))
	http.HandleFunc(PageEditRoute, makeHandler(editHandler))
	http.HandleFunc(PageSaveRoute, makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
