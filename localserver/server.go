package localserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

type FileHandler struct {
	staticPath string
	indexPath  string
}

func (f FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(f.staticPath, r.URL.Path)

	log.Println("ServeHTTP func")

	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		http.ServeFile(w, r, filepath.Join(f.staticPath, f.indexPath))
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(f.staticPath)).ServeHTTP(w, r)
}

type Server struct{}

func (s Server) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	fileHandler := FileHandler{staticPath: "build", indexPath: "index.html"}

	router.HandleFunc("/", fileHandler.ServeHTTP)
	router.HandleFunc("/home", homeHandler)
	router.HandleFunc("/login", loginHandler).Methods("GET")
	router.HandleFunc("/signup", signupHandler).Methods("GET")
	router.HandleFunc("/login", Signin).Methods("POST")
	router.HandleFunc("/signup", Signup).Methods("POST")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./build"))))

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	log.Println("Login Form")

	pageToLoad, _ := os.ReadFile("./build/login.html")

	render(w, r, parseTemplate(string(pageToLoad)))
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	log.Println("Login Form")

	pageToLoad, _ := os.ReadFile("./build/signup.html")

	render(w, r, parseTemplate(string(pageToLoad)))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	log.Println("Home page")

	pageToLoad, _ := os.ReadFile("./build/home.html")

	render(w, r, parseTemplate(string(pageToLoad)))
}

func render(w http.ResponseWriter, r *http.Request, tpl *template.Template) {
	buf := new(bytes.Buffer)

	if err := tpl.ExecuteTemplate(buf, "base", []byte{}); err != nil {
		fmt.Printf("\nRender Error: %v\n", err)
		return
	}
	w.Write(buf.Bytes())
}

func parseTemplate(content string) *template.Template {
	file, _ := os.ReadFile("./build/index.html")

	base := template.New("base")
	base.Funcs(template.FuncMap{
		"content": func() string { return content },
	}).Parse(string(file))

	return base
}
