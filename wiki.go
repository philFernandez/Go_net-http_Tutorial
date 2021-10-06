package main

import (
	"fmt"
	"io/ioutil"
	"log"
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

func main() {
	p1 := &Page{Title: "Test Page", Body: []byte("This is a sample Page.")}
	err := p1.save()
	if err != nil {
		log.Fatalln("There was a problem writing the page:", err)
	}
	p2, err := loadPage("Test Page")
	if err != nil {
		log.Fatalln("There was a problem reading the page:", err)
	}
	fmt.Println(string(p2.Body))
}
