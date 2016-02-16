// Template for command line application using JSON over http

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rmohid/go-config/config"
	"github.com/rmohid/go-config/dbg"
	"github.com/rmohid/go-config/webExternal"
)

var err error

func main() {

	// define all string based options
	var opts = [][]string{
		{"port", os.Getenv("HOSTNAME")+":7000", "external web port"},
		{"dbg.httpUrl", os.Getenv("HOSTNAME")+":7000"},
		{"dbg.verbosity", "0"},
	}

	if err = config.ParseArgs(opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	dbg.Log(2, config.Dump())
	dbg.Log(0, "listening on", config.Get("port"))

	go test()
	webExternal.Run()
}

func test() {
	for {
		time.Sleep(3 * time.Second)
		dbg.Log(4, "Debug log ", time.Now().Format(time.StampMilli))
	}
}
