package main

import (
	"fmt"
	"log"
	"html/template"
	"net/http"
	"regexp"
	"io/ioutil"
	"github.com/gomarkdown/markdown"
)

const (
	myEmail = "navybluesilver@protonmail.ch"
	myFingerprint = "DE0F 14CE F6C2 819E 0ADC CF85 4153 56DD 6450 053C"
	port  = ":80"
)

var (
	templates = template.Must(template.ParseFiles("template/about.html", "template/donate.html", "template/article.html"))
	validPath = regexp.MustCompile("^/(article)/([a-zA-Z0-9]+)$")
)

type AboutPage struct {
	Title string
	Email string
	Fingerprint string
}

type DonatePage struct {
	Title string
}

type ArticlePage struct {
  Title string
  Article  template.HTML
}

//HTTP Handling
func aboutHandler(w http.ResponseWriter, r *http.Request) {
        p := &AboutPage{Title: "About", Email: myEmail, Fingerprint: myFingerprint }
				err := templates.ExecuteTemplate(w, "about.html", p)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}

func donateHandler(w http.ResponseWriter, r *http.Request) {
        p := &DonatePage{Title: "Donate" }
				err := templates.ExecuteTemplate(w, "donate.html", p)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}

func loadArticle(title string) (*ArticlePage, error) {
	file := fmt.Sprintf("article/%s.md", title)
	md, err := ioutil.ReadFile(file) // just pass the file name
	if err != nil {
			return nil, err
	}
	body := markdown.ToHTML(md, nil, nil)
	return &ArticlePage{Title: "Article", Article: template.HTML(body) }, nil
}

func articleHandler(w http.ResponseWriter, r *http.Request, title string) {
				p, err := loadArticle(title)
				err = templates.ExecuteTemplate(w, "article.html", p)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
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

func main() {
	http.HandleFunc("/", aboutHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/donate", donateHandler)
	http.HandleFunc("/article/", makeHandler(articleHandler))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("template"))))
	log.Fatal(http.ListenAndServe(port, nil))
}
