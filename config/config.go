// Key value web api for configuration data
// See github.com/rmohid/go-template for detailed description

package config

import (
	"flag"
	"fmt"
	"github.com/rmohid/go-template/config/data"
	"github.com/rmohid/go-template/config/webInternal"
	"sync"
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
	mu      sync.Mutex
	indexed map[string]*Option
	d       = data.New()
)

func init() {
	indexed = make(map[string]*Option)

	// default options for config package
	opts := [][]string{
		{"config.port", "localhost:7100", "internal api web port"},
		{"config.readableJson", "yes", "pretty print api json output"},
		{"config.enableFlagParse", "yes", "allow config to flag.Parse()"},
	}

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
func Replace(newkv map[string]string) {
	d.Replace(newkv)
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
	mu.Lock()
	defer mu.Unlock()
	for i, _ := range inOpts {
		var o Option
		if v, ok := indexed[inOpts[i][NameIdx]]; ok == true {
			o = *v
		}
		o.Name, o.Default = inOpts[i][NameIdx], inOpts[i][DefaultIdx]
		if len(inOpts[i]) > 2 {
			o.Description = inOpts[i][DescriptionIdx]
		}
		d.Set(o.Name, o.Default)
		indexed[o.Name] = &o
	}
	return nil
}
func ParseArgs(inOpts [][]string) error {

	PushArgs(inOpts)
	mu.Lock()
	defer mu.Unlock()
	for _, v := range indexed {
		elem := v
		elem.Value = flag.String(elem.Name, elem.Default, elem.Description)
	}
	// nothing is actally done until parse is called
	if Get("config.enableFlagParse") == "yes" {
		flag.Parse()
	}
	for _, elem := range indexed {
		d.Set(elem.Name, *elem.Value)
	}

	// Start the internal admin web interface
	if Get("dbg.verbosity") != "0" {
		fmt.Println("configuration on", Get("config.port"))
		go webInternal.Run()
	}
	return nil
}
