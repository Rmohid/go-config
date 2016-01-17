// Key value web api for configuration data
// See github.com/rmohid/go-config for detailed description

package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/rmohid/go-config/config/data"
	"github.com/rmohid/go-config/config/webInternal"
)

type Option struct {
	Name, Default, Description string
	Value                      *string
}

const (
	NameIdx = iota
	DefaultIdx
	DescriptionIdx
)

var (
	wg       sync.WaitGroup
	d        = data.New()
	indexed  = make(map[string]*Option)
	argsIn   = make(chan Option)
	argsDone = make(chan bool)
)

func init() {
	// default options for config package
	opts := [][]string{
		{"config.port", "7100", "internal api web port"},
		{"config.file", os.Args[0] + ".json", "configuration file to use"},
		{"config.readableJson", "yes", "pretty print api json output"},
		{"config.enableFlagParse", "yes", "allow config to flag.Parse()"},
	}

	go argConsumer()
	PushArgs(opts)
}
func Delete(k string) {
	d.Delete(k)
}
func Set(k, v string) {
	d.Set(k, v)
}
func Get(k string) string {
	return d.Get(k)
}
func Exists(k string) bool {
	return d.Exists(k)
}
func Keys() []string {
	return d.Keys()
}
func Clear() {
	d.Clear()
}
func Dump() []string {
	var out []string
	for _, k := range Keys() {
		kv := fmt.Sprintf("%s=%s,", k, Get(k))
		out = append(out, kv)
	}
	return out
}
func PushArgs(inOpts [][]string) error {
	for i, _ := range inOpts {
		var o Option
		o.Name, o.Default = inOpts[i][NameIdx], inOpts[i][DefaultIdx]
		if len(inOpts[i]) > 2 {
			o.Description = inOpts[i][DescriptionIdx]
		}
		wg.Add(1)
		select {
		case <-argsDone:
			return nil
		case argsIn <- o:
		}
	}
	return nil
}
func ParseArgs(inOpts [][]string) error {
	PushArgs(inOpts)
	wg.Wait()
	close(argsDone)
	for _, elem := range indexed {
		elem.Value = flag.String(elem.Name, elem.Default, elem.Description)
	}
	// nothing is actally done until parse is called
	if Get("config.enableFlagParse") == "yes" {
		flag.Parse()
	}
	for _, elem := range indexed {
		d.Set(elem.Name, *elem.Value)
	}
	loadConfigFile()
	// Start the internal admin web interface
	if Get("dbg.verbosity") != "0" {
		fmt.Println("configuration on", "localhost:"+Get("config.port"))
	}
	if Get("config.port") != "" {
		go webInternal.Run()
	}
	return nil
}
func loadConfigFile() {
	var newkv = make(map[string]string)
	cfgFile := d.Get("config.file")
	if configJson, err := ioutil.ReadFile(cfgFile); err == nil {
		err := json.Unmarshal(configJson, &newkv)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for k, v := range newkv {
			d.Set(k, v)
		}
		return
	}
}
func argConsumer() {
	for {
		select {
		case <-argsDone:
			return
		case o := <-argsIn:
			wg.Done()
			if v, ok := indexed[o.Name]; ok == true {
				o.Description = (*v).Description
			}
			d.Set(o.Name, o.Default)
			indexed[o.Name] = &o
		}
	}
}
