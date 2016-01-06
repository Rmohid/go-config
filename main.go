// Template for command line application using JSON over http

package main

import (
	"fmt"
	"github.com/rmohid/go-template/config"
	"github.com/rmohid/go-template/dbg"
	"github.com/rmohid/go-template/webExternal"
	"os"
	"time"
)

var err error

func main() {

	// define all string based options
	var opts = [][]string{
		{"port", "localhost:7000", "external web port"},
		{"dbg.httpUrl", "localhost:7000"},
		{"dbg.verbosity", "0"},
	}

	if err = config.ParseArgs(opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	dbg.Log(2, config.Dump())
	dbg.Log(0,"listening on", config.Get("port"))

	go test()
	webExternal.Run()
}

func test() {
	for {
		time.Sleep(1 * time.Second)
		dbg.Log(2, "Debug log ", time.Now().Format(time.StampMilli))
	}
}
