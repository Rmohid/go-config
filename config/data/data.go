// Key value store for global configuration data, replace with whatever backend you like

package data

import (
	"sort"
	"sync"
)

var (
	data map[string]string
	mu   sync.Mutex
        d    KvData
)

type KvData interface {
	GetData() *map[string]string
	Delete(k string)
	Set(k, v string)
	Get(k string) string
	Exists(k string) bool
	Keys() []string
	Replace(newkv map[string]string)
}

type NaiveKv struct {
}

func init() {
	data = make(map[string]string)
        d = new(NaiveKv)
}
func New() KvData {
	return d
}
func SwapKvStore(newKvStore KvData) {
        d = newKvStore
}

func (d *NaiveKv) GetData() *map[string]string {
	return &data
}
func (d *NaiveKv) Delete(k string) {
	mu.Lock()
	defer mu.Unlock()
	delete(data, k)
}
func (d *NaiveKv) Set(k, v string) {
	mu.Lock()
	defer mu.Unlock()
	if k == "" {
		return
	}
	data[k] = v
}
func (d *NaiveKv) Get(k string) string {
	mu.Lock()
	defer mu.Unlock()
	return data[k]
}
func (d *NaiveKv) Exists(k string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, ok := data[k]
	return ok
}
func (d *NaiveKv) Keys() []string {
	mu.Lock()
	defer mu.Unlock()
	list := make([]string, 0, len(data))
	for k := range data {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}
func (d *NaiveKv) Replace(newkv map[string]string) {
	mu.Lock()
	defer mu.Unlock()
	// take old reference and garbage collect memory
	data = newkv
}
