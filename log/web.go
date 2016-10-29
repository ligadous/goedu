package main

//Versao Site Static - Desenvolvimento

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	LISTEN = ":8080"
)

//Log Seletivo

var accesslog *log.Logger

func init() {
	logconfig()
	go Loop()
}

func Loop() {
	var i int
	for {
		log.Printf("Loop %d", i)

		time.Sleep(time.Second)
		i++
	}

}

func logconfig() {

	//Discard first
	log.SetOutput(ioutil.Discard)

	//if isverbose {
	//Activate for debug propose
	log.SetOutput(os.Stderr)
	//}
	log.SetPrefix("Debug] ")

	//log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	//Access log
	accesslog = log.New(os.Stderr, "Access] ", log.Ldate|log.Ltime|log.Lshortfile)

}

func debugon(w http.ResponseWriter, r *http.Request) {
	log.SetOutput(os.Stderr)

	fmt.Fprint(w, "Debug On")
}

func debugoff(w http.ResponseWriter, r *http.Request) {
	log.SetOutput(ioutil.Discard)

	fmt.Fprint(w, "Debug Off")
}

//First level - Log Handlers
func logHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accesslog.Printf("%v: %v", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

//Dinamic
func logHandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return logHandler(next)
}

//Statics
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Wolrd!")
}

//Substituir /php/load_page.php por http.Handle
func main() {
	http.Handle("/debugon", logHandlerFunc(debugon))
	http.Handle("/debugoff", logHandlerFunc(debugoff))
	http.Handle("/", logHandlerFunc(helloHandler))
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(LISTEN, nil))
}
