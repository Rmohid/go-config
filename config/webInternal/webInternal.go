// Handler for administrative web interface

package webInternal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rmohid/go-config/config/data"
)

var d = data.New()

func Run() {
	serverInternal := http.NewServeMux()
	serverInternal.HandleFunc("/", handler)
	serverInternal.HandleFunc("/key/", handleGetKey)
	serverInternal.HandleFunc("/json", handleGetJson)
	serverInternal.HandleFunc("/kv/reset", handleKvReset)
	log.Fatal("webInternal.Run(): ", http.ListenAndServe(d.Get("config.port"), serverInternal))
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGet(w, r)
	case "POST":
		switch strings.Join(r.Header["Content-Type"], "") {
		case "application/json":
			handlePostJson(w, r)
		default:
		}
	case "DELETE":
		handleDelete(w, r)

	}
}
func handleGet(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	if len(r.Form) > 0 {
		for k, v := range r.Form {
			d.Set(k, strings.Join(v, " "))
		}
	} else {
		configKeys := d.Keys()
		for _, k := range configKeys {
			fmt.Fprintf(w, "Config[%q] = %q\n", k, d.Get(k))
		}
	}
}
func handleGetJson(w http.ResponseWriter, r *http.Request) {
	dat, err := json.Marshal(d.GetData())
	if d.Get("config.readableJson") == "yes" {
		dat, err = json.MarshalIndent(d.GetData(), "", "  ")
	}
	if err != nil {
		log.Print(err)
	}
	fmt.Fprintf(w, "%s", dat)
}
func handlePostJson(w http.ResponseWriter, r *http.Request) {
	var newkv = make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newkv)
	if err != nil {
		log.Print(err)
	}
	d.Update(newkv)
}
func handleDelete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	if len(r.Form) > 0 {
		for k, _ := range r.Form {
			d.Delete(k)
		}
	}
}
func handleGetKey(w http.ResponseWriter, r *http.Request) {
	var i = strings.LastIndex(r.URL.Path, "/key/") + len("/key/")
	if i > 0 {
		fmt.Fprintf(w, "%s\n", d.Get(r.URL.Path[i:]))
	}
}
func handleKvReset(w http.ResponseWriter, r *http.Request) {
	d.Clear()
}
