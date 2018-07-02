package main

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	config "github.com/navybluesilver/config"
	"github.com/navybluesilver/lightning"
	"github.com/navybluesilver/futures"
	qrcode "github.com/skip2/go-qrcode"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	myEmail          = "navybluesilver@protonmail.ch"
	myFingerprint    = "DE0F 14CE F6C2 819E 0ADC CF85 4153 56DD 6450 053C"
	myBitcoinAddress = "3BwKJ23VEsWN9j678HE2KfU6dLEJCpHdJc"
	port             = ":443"
	defaultDonation  = 10000 //satoshis
)

var (
	templates        = template.Must(template.ParseFiles("template/about.html", "template/donate.html", "template/article.html", "template/disclaimer.html"))
	validArticlePath = regexp.MustCompile("^/(article)/([a-zA-Z0-9]+)$")
	certFile         = config.GetString("web.certFile")
	keyFile          = config.GetString("web.keyFile")
	fmap = template.FuncMap{
    "formatAsSatoshi": formatAsSatoshi,
  }
)


type PricingPage struct {
	Title string
}

type ArticlePage struct {
	Title   string
	Article template.HTML
}

type DonatePage struct {
	Title             string
	DonationAddress   string
	PaymentRequest    string
	PaymentRequestPNG string
	Donation          int64
}

func main() {
	//default
	http.HandleFunc("/", aboutHandler)

	//pages
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/article/", makeArticle(articleHandler))
	http.HandleFunc("/donate", donateHandler)
	http.HandleFunc("/disclaimer", disclaimerHandler)
	http.HandleFunc("/futures", futuresHandler)

	//invoice
	http.HandleFunc("/invoice/", makeInvoice(invoiceHandler))

	//files
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("template"))))

	//listen
	// redirect every http request to https
	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	log.Fatal(http.ListenAndServeTLS(port, certFile, keyFile, nil))
}

func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	host := strings.Split(req.Host, ":")[0]
	target := "https://" + host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		// see @andreiavrammsd comment: often 307 > 301
		http.StatusTemporaryRedirect)
}

//about
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadArticle("about")
	err = templates.ExecuteTemplate(w, "about.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//futures
func futuresHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("futures.html").Funcs(fmap).ParseFiles("template/futures.html"))
	err := t.ExecuteTemplate(w, "futures.html", futures.LoadFutures())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//pricing
func pricingHandler(w http.ResponseWriter, r *http.Request) {
	p := &PricingPage{Title: "Pricing"}
	err := templates.ExecuteTemplate(w, "pricing.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//article
func makeArticle(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validArticlePath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func articleHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadArticle(title)
	err = templates.ExecuteTemplate(w, "article.html", p)
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
	return &ArticlePage{Title: "", Article: template.HTML(body)}, nil
}

//donate
func donateHandler(w http.ResponseWriter, r *http.Request) {
	var paymentRequest string
	var donation int64

	// load a default invoice the first time
	if r.Method == "GET" {
		invoice, err := lightning.GetInvoice(defaultDonation, "intial")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		paymentRequest = invoice
		donation = defaultDonation
	} else {
		// load an invoice based on the specified donation
		r.ParseForm()
		satoshis, err := strconv.ParseInt(r.Form["satoshis"][0], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		invoice, err := lightning.GetInvoice(satoshis, "new donation")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		paymentRequest = invoice
		donation = satoshis
	}

	png := fmt.Sprintf("/invoice/%s", paymentRequest)
	p := &DonatePage{Title: "Donate", PaymentRequest: paymentRequest, DonationAddress: myBitcoinAddress, PaymentRequestPNG: png, Donation: donation}
	err := templates.ExecuteTemplate(w, "donate.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//disclaimer
func disclaimerHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadArticle("disclaimer")
	err = templates.ExecuteTemplate(w, "disclaimer.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//invoice
func makeInvoice(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		invoice := r.URL.Path[9:len(r.URL.Path)]
		fn(w, r, invoice)
	}
}

func invoiceHandler(w http.ResponseWriter, r *http.Request, payment_request string) {
	png, err := qrcode.Encode(payment_request, qrcode.Medium, 1500)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(png)))
	if _, err := w.Write(png); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func formatAsSatoshi(satoshi float64) (string, error) {
	if satoshi == 0 {
		return "", nil
	}
	return fmt.Sprintf("%.0f", satoshi), nil
}
