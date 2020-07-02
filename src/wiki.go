package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(PageSaveRoute):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		sendHttpError(w, err)
		return
	}
	http.Redirect(w, r, PageViewRoute + title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(PagePath + tmpl + ".html")
	if err != nil {
		sendHttpError(w, err)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		sendHttpError(w, err)
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
	http.HandleFunc(PageViewRoute, viewHandler)
	http.HandleFunc(PageEditRoute, editHandler)
	http.HandleFunc(PageSaveRoute, saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
