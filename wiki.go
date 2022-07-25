package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

//this method will save Page's Body to a text file
func (p *Page) save() error { //takes a pointer to Page struct
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600) //WriteFile is stdlibrary function that writes byte slice to a file and returns error value, in case of success returns 'nil'
	//third paratemer means file should be created with read-write permissions for current user only
}

func loadPage(title string) (*Page, error) { //loadPage constructs file name from title parameter,reads file's contents into body and return a pointer to a Page literal
	filename := title + ".txt"
	body, err := os.ReadFile(filename) //returns []byte and error
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):] //extracting title from r.URL.Path(path component of requested URL), dropping "/view/" from start
	p, _ := loadPage(title)             //ignoring error value
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request) { //loads page,if it doesn't exist create an empty Page struct and displays an HTML form
	title := r.URL.Path[len("/edit/"):] //dropping starting "/edit/"
	p, err := loadPage(title)           //loading page
	if err != nil {                     //if there is error for example requested page doesn't exist
		p = &Page{Title: title} //creating a empty struct literal
	}
	fmt.Fprintf(w, "<h1>Editing %s</h1>"+
		"<form action=\"/save/%s\" method=\"POST\">"+
		"<textarea name=\"body\">%s</textarea><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form>",
		p.Title, p.Title, p.Body)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	//http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
