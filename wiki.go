package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
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

//ParseFiles takes any number of string arguments that identify our template files and parses those files into templates named after base file name
var templates = template.Must(template.ParseFiles("template/edit.html", "template/view.html")) //Must panics when passed a error it is appropriate here if templates can't loaded only sensible thing to do is exit program

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	/*
		t, err := template.ParseFiles("template/" + tmpl + ".html") //ParseFiles will read edit.html and return a *template.Template
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) //error replies request with specified error message and HTTP code
			return
		}
		err = t.Execute(w, p) //t.Execute will execute template writes generated HTML to http.ResponseWriter
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	*/
	err := templates.ExecuteTemplate(w, tmpl+".html", p) //applies template associated with tmpl+".html" that has the given name to specified data object and writes the output to w
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	//	title := r.URL.Path[len("/view/"):] //extracting title from r.URL.Path(path component of requested URL), dropping "/view/" from start
	//	title, err := getTitle(w, r)
	p, err := loadPage(title) //ignoring error value
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) { //loads page,if it doesn't exist create an empty Page struct and displays an HTML form
	//	title := r.URL.Path[len("/edit/"):] //dropping starting "/edit/"
	//	title, err := getTitle(w, r)
	p, err := loadPage(title) //loading page
	if err != nil {           //if there is error for example requested page doesn't exist
		p = &Page{Title: title} //creating a empty struct literal
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	//	title := r.URL.Path[len("/save/"):]
	//	title, err := getTitle(w, r)
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

/*
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path) //returns first substring that matches regexp
	if m == nil {
		http.NotFound(w, r) //replies request with 404 and not found
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}
*/
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$") //expression for validation

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

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
