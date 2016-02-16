// Handler for administrative web interface

package webInternal

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rmohid/go-config/config/data"
)

type webHandler struct {
	Path        string
	Handler     func(http.ResponseWriter, *http.Request)
	Description string
	Meta        string
}

var (
	d            = data.New()
	getHandlers  = []webHandler{}
	postHandlers = []webHandler{}
)

func Run() {
	getHandlers = []webHandler{
		webHandler{Path: "/", Handler: handler, Description: "Show usage"},
		webHandler{Path: "/key/", Handler: handleGetKey, Description: "Show a single key"},
		webHandler{Path: "/json", Handler: handleGetJson, Description: "Dump key/value store in JSON"},
		webHandler{Path: "/save", Handler: handleGetSave, Description: "Replace config file with current key/value store"},
		webHandler{Path: "/exit", Handler: handleGetExit, Description: "Exit program with 0"},
		webHandler{Path: "/clear", Handler: handleKvReset, Description: "Clear k/v store"},
	}
	postHandlers = []webHandler{
		webHandler{Path: "/", Handler: handlePostJson, Description: "Upload new key/value store in JSON", Meta: "application/json"},
	}

	serverInternal := http.NewServeMux()
	for _, h := range getHandlers {
		serverInternal.HandleFunc(h.Path, h.Handler)
	}
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
		fmt.Fprintln(w, err)
		return
	}
	if len(r.Form) > 0 {
		for k, v := range r.Form {
			d.Set(k, strings.Join(v, " "))
		}
	} else {
		handleGetUsage(w, r)
	}
}
func handleGetUsage(w http.ResponseWriter, r *http.Request) {
	getFmt := "\t%-10s %s\n"
	postFmt := "\t%-5s %-20s %s\n"
	fmt.Fprintln(w, "GET request handlers:")
	fmt.Fprintf(w, getFmt, "Path", "Description")
	for _, h := range getHandlers {
		fmt.Fprintf(w, getFmt, h.Path, h.Description)
	}
	fmt.Fprintln(w, "\nPOST request handlers:")
	fmt.Fprintf(w, postFmt, "Path", "Content-Type", "Description")
	for _, h := range postHandlers {
		fmt.Fprintf(w, postFmt, h.Path, h.Meta, h.Description)
	}
}
func handleGetJson(w http.ResponseWriter, r *http.Request) {
	dat, err := json.Marshal(d.GetData())
	if d.Get("config.readableJson") == "yes" {
		dat, err = json.MarshalIndent(d.GetData(), "", "  ")
	}
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprintf(w, "%s", dat)
}
func handlePostJson(w http.ResponseWriter, r *http.Request) {
	var newkv = make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newkv)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	d.Update(newkv)
}
func handleDelete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintln(w, err)
		return
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
		if d.Exists(r.URL.Path[i:]) {
			fmt.Fprintf(w, "%s\n", d.Get(r.URL.Path[i:]))
		} else {
			fmt.Fprintln(w, "Flag/Key list:\n")
			flag.CommandLine.SetOutput(w)
			flag.CommandLine.PrintDefaults()
		}
	}
}
func handleKvReset(w http.ResponseWriter, r *http.Request) {
	d.Clear()
}
func handleGetSave(w http.ResponseWriter, r *http.Request) {
	cfgFile := d.Get("config.file")
	jdata, err := json.Marshal(d.GetData())
	if d.Get("config.readableJson") == "yes" {
		jdata, err = json.MarshalIndent(d.GetData(), "", "  ")
	}
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	if err := ioutil.WriteFile(cfgFile, jdata, 0644); err != nil {
		fmt.Fprintln(w, err)
		log.Print(err)
		os.Exit(1)
	}
}
func handleGetExit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "exit")
	fmt.Println("User requested exit")
	os.Exit(0)
}
