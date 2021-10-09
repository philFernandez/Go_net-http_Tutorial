package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var (
	// Cache all templates so they don't have to be re-read every time they're loaded in the browser.
	templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))
	// Regex for path validation (restrict user to only these paths)
	validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
)

// Page is a data structure that represents a wiki page
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := "data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
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

func render(w http.ResponseWriter, templateFile string, p *Page) {
	// execute specified cached template using "p" as the object to be expanded in the template
	if err := templates.ExecuteTemplate(w, templateFile+".html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
	})
	http.HandleFunc("/view/", makeHandler(func(w http.ResponseWriter, r *http.Request, title string) {
		p, err := loadPage(title)
		if err != nil {
			http.Redirect(w, r, "/edit/"+title, http.StatusFound)
			return
		}
		render(w, "view", p)
	}))
	http.HandleFunc("/edit/", makeHandler(func(w http.ResponseWriter, r *http.Request, title string) {
		p, err := loadPage(title)
		if err != nil {
			p = &Page{Title: title}
		}
		render(w, "edit", p)
	}))
	http.HandleFunc("/save/", makeHandler(func(w http.ResponseWriter, r *http.Request, title string) {
		body := r.FormValue("body")
		p := &Page{Title: title, Body: []byte(body)}
		err := p.save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/view/"+title, http.StatusFound)
	}))
	const portNumber = ":8080"
	fmt.Println("Server Started on http://localhost" + portNumber)
	log.Fatal(http.ListenAndServe(portNumber, nil))
}
