package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
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

	r := newRoom()

	//Root Handle.
	//Here we send a type templateHandler directly in without
	//creating a variabel/reference to it first. We can do that
	//by preceding it with &.
	//templateHandler have a serveHTTP method,
	//and then becomes a handler
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	//start the room.

	go r.run()

	//Start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
