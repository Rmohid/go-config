package data

import (
	"log"
	"os"
	"testing"
)

var (
	Trace *log.Logger
)

func init() {
	// for Trace.Println(Keys())
	Trace = log.New(os.Stdout, "**: ", log.Lshortfile)
}

type testpair struct {
	value  string
	expect string
}
type testkv struct {
	key    string
	value  string
	expect string
}
type testbool struct {
	key    string
	value  string
	expect bool
}
type testkeys struct {
	inkeys  []string
	outkeys []string
	expect  []string
}

func resetData() {
	d.Clear()
}
func testEq(a, b []string) bool {
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
func loadSlice(a []string) {
	for i := range a {
		d.Set(a[i], a[i])
	}
}
func TestGet(t *testing.T) {
	resetData()
	var tests = []testpair{
		{"", ""},
		{"missing", ""},
		{"portInternal", ""},
	}
	for _, pair := range tests {
		v := d.Get(pair.value)
		if v != pair.expect {
			t.Error(
				"For", pair.value,
				"expected", pair.expect,
				"got", v,
			)
		}
	}
}
func TestSet(t *testing.T) {
	resetData()
	var tests = []testkv{
		{"", "", ""},
		{"", "abc", ""},
		{"foo", "bar", "bar"},
		{"foo", "", ""},
	}
	for _, pair := range tests {
		d.Set(pair.key, pair.value)
		v := d.Get(pair.key)
		if v != pair.expect {
			t.Error(
				"For", pair.key, pair.value,
				"expected", pair.expect,
				"got", v,
			)
		}
	}
}
func TestDelete(t *testing.T) {
	resetData()
	var key, value, expect string
	for {
		key, value, expect = "missing", "", ""
		d.Delete(key)
		value = d.Get(key)
		if value != expect {
			break
		}
		key, value, expect = "foo", "bar", ""
		d.Set(key, value)
		d.Delete(key)
		value = d.Get(key)
		if value != expect {
			break
		}
		key, value, expect = "fOo", "bar", ""
		d.Set(key, value)
		d.Delete(key)
		value = d.Get(key)
		if value != expect {
			break
		}

		return
	}
	t.Error(
		"For", key,
		"expected", expect,
		"got", value,
	)
}
func TestExists(t *testing.T) {
	resetData()
	var tests = []testbool{
		{"foo", "", true},
		{"foo", "bar", true},
		{"foo", "", true},
		{"missing", "", false},
	}
	for _, pair := range tests {
		if pair.expect {
			d.Set(pair.key, pair.value)
		}
		v := d.Exists(pair.key)
		if v != pair.expect {
			t.Error(
				"For", pair.key, pair.value,
				"expected", pair.expect,
				"got", v,
			)
		}
	}
}
func TestKeys(t *testing.T) {
	resetData()
	var tests = []testkeys{
		{inkeys: []string{}, outkeys: []string{}, expect: []string{}},
		{inkeys: []string{"foo", "bar"}, outkeys: []string{}, expect: []string{"bar", "foo"}},
		{inkeys: []string{"abar", "zbar"}, outkeys: []string{}, expect: []string{"abar", "bar", "foo", "zbar"}},
	}
	for _, set := range tests {
		loadSlice(set.inkeys)
		set.outkeys = d.Keys()
		if true != testEq(set.outkeys, set.expect) {
			t.Error(
				"For", set.inkeys,
				"expected", set.expect,
				"got", set.outkeys,
			)
		}
	}
}
