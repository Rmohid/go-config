// Handler for internet facing web interface

package webExternal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rmohid/go-config/config"
	"github.com/rmohid/go-config/dbg"
)

func Run() {
	serverExternal := http.NewServeMux()
	serverExternal.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(config.Get("port"), serverExternal))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	dbg.ErrLog(2, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
		dbg.ErrLog(1, "Form[%q] = %q\n", k, v)
	}
}
