package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/postmannen/bp-learning/chapter1/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
)

//templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Load and execute the templates

	t.once.Do(func() {
		t.templ = template.Must(
			//filepath.join will create the complete path ./templates/<filename>
			template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.templ.Execute(w, r)
}

func main() {

	var addr = flag.String("addr", ":8080", "The addr of the  application.")
	flag.Parse() // parse the flags

	// setup gomniauth
	gomniauth.SetSecurityKey("PUT YOUR AUTH KEY HERE")
	gomniauth.WithProviders(
		//facebook.New("key", "secret",
		//	"http://localhost:8080/auth/callback/facebook"),
		//github.New("key", "secret",
		//	"http://localhost:8080/auth/callback/github"),
		google.New("1008756175538-30lreagdvr41c2molmvtv57tf79jv5o7.apps.googleusercontent.com",
			"l4agwoZIaDEHkVTQEr_YA8X5",
			"http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	//Root Handle.
	//Here we send a type templateHandler directly in without
	//creating a variabel/reference to it first. We can do that
	//by preceding it with &.
	//templateHandler have a serveHTTP method,
	//and then becomes a handler
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	//start the room.

	go r.run()

	//Start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
