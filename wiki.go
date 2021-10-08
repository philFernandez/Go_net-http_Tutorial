package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// Page is a data structure that represents a wiki page
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func render(w http.ResponseWriter, templateFile string, p *Page) {
	t, _ := template.ParseFiles(templateFile + ".html")
	t.Execute(w, p)
}

func main() {
	http.HandleFunc("/view/", func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/view/"):]
		p, err := loadPage(title)
		if err != nil {
			http.Redirect(w, r, "/edit/"+title, http.StatusFound)
			return
		}
		render(w, "view", p)
	})
	http.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/edit/"):]
		p, err := loadPage(title)
		if err != nil {
			p = &Page{Title: title}
		}
		render(w, "edit", p)

	})
	http.HandleFunc("/save/", func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/save/"):]
		body := r.FormValue("body")
		p := &Page{Title: title, Body: []byte(body)}
		p.save()
		http.Redirect(w, r, "/view/"+title, http.StatusFound)
	})
	const portNumber = ":8080"
	fmt.Println("Server Started on http://localhost" + portNumber)
	log.Fatal(http.ListenAndServe(portNumber, nil))
}
