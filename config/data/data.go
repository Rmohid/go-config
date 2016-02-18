// Key value store for global configuration data, replace with whatever backend you like

package data

import (
	"sort"
	"sync"
)

type Kv struct {
	data        map[string]string
	d, shim     KvData
}
type KvList struct {
	list        map[string]Kv
}
var (
	kvlist      KvList
	mu          sync.Mutex
)

type KvData interface {
	GetData() *map[string]string
	Delete(k string)
	Set(k, v string)
	Get(k string) string
	Exists(k string) bool
	Keys() []string
	Clear()
	Update(newkv map[string]string)
}

func init() {
        kvlist = KvList{}
        kvlist.list[""] = Kv{}
}
func New() KvData {
	return shim
}
func SwapKvStore(newKvStore KvData) {
	d = newKvStore
}

type shimKv struct {
}

func (s *shimKv) GetData() *map[string]string    { return d.GetData() }
func (s *shimKv) Delete(k string)                { d.Delete(k) }
func (s *shimKv) Set(k, v string)                { d.Set(k, v) }
func (s *shimKv) Get(k string) string            { return d.Get(k) }
func (s *shimKv) Exists(k string) bool           { return d.Exists(k) }
func (s *shimKv) Keys() []string                 { return d.Keys() }
func (s *shimKv) Clear()                         { d.Clear() }
func (s *shimKv) Update(newkv map[string]string) { d.Update(newkv) }

type naiveKv struct {
}

func (d *naiveKv) GetData() *map[string]string {
	return &data
}
func (d *naiveKv) Delete(k string) {
	mu.Lock()
	defer mu.Unlock()
	delete(data, k)
}
func (d *naiveKv) Set(k, v string) {
	mu.Lock()
	defer mu.Unlock()
	if k == "" {
		return
	}
	data[k] = v
}
func (d *naiveKv) Get(k string) string {
	mu.Lock()
	defer mu.Unlock()
	return data[k]
}
func (d *naiveKv) Exists(k string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, ok := data[k]
	return ok
}
func (d *naiveKv) Keys() []string {
	mu.Lock()
	defer mu.Unlock()
	list := make([]string, 0, len(data))
	for k := range data {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}
func (d *naiveKv) Clear() {
	mu.Lock()
	defer mu.Unlock()
	newkv := make(map[string]string)
	// orphan old reference and garbage collect memory
	data = newkv
}
func (d *naiveKv) Update(newkv map[string]string) {
	mu.Lock()
	defer mu.Unlock()
	for k, v := range newkv {
		data[k] = v
	}
}
