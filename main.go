package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"edwardhorsey/url-shortener/base32"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Link struct {
	ID  int
	URL string
}

func main() {
	// Load environment variables from file.
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	dsn := os.Getenv("DSN")

	// Connect to database using DSN environment variable.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Create an API handler.
	handler := CreateNewHandler(db)

	// Start an HTTP API server.
	const addr = ":8080"
	log.Printf("successfully connected to database, starting HTTP server on %q", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}

type Handler struct {
	db *gorm.DB
}

func CreateNewHandler(db *gorm.DB) http.Handler {
	h := &Handler{db: db}

	r := mux.NewRouter()
	r.HandleFunc("/{code}/show", h.show).Methods(http.MethodGet)
	r.HandleFunc("/{code}", h.redirect).Methods(http.MethodGet)
	r.HandleFunc("/create", h.createLink).Methods(http.MethodPost)
	r.HandleFunc("/", h.index).Methods(http.MethodGet)

	return r
}

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {
	var link Link
	code := mux.Vars(r)["code"]
	decoded := base32.Decode(code)
	result := h.db.First(&link, decoded)
	if result.Error != nil {
		http.NotFound(w, r)
		return
	}

	redirectUrl := link.URL

	println("redirecting to " + redirectUrl)

	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func validateUrl(u string) bool {
	url, err := url.ParseRequestURI(u)

	return err == nil && url.Scheme != "" && url.Host != ""
}

func (h *Handler) createLink(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.Form.Get("url")
	url = strings.Trim(url, " ")

	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	if !validateUrl(url) {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	link := Link{URL: url}
	result := h.db.Create(&link)

	if result.Error != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	code := base32.Encode(link.ID)

	println(url)
	println("id", link.ID, "->", "code", code)

	shortUrl := r.Host + "/" + code

	template := template.Must(template.ParseFiles("templates/link.html"))
	template.Execute(w, map[string]string{"Link": shortUrl, "Url": url})
}

func (h *Handler) show(w http.ResponseWriter, r *http.Request) {
	var link Link
	code := mux.Vars(r)["code"]
	decoded := base32.Decode(code)

	result := h.db.First(&link, decoded)
	if result.Error != nil {
		http.NotFound(w, r)
		return
	}

	shortUrl := r.Host + "/" + code
	redirectUrl := link.URL

	template := template.Must(template.ParseFiles("templates/show.html"))
	template.Execute(w, map[string]string{"Link": shortUrl, "Url": redirectUrl})
}
