package icheck

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// RequestValues is a collection of values that can be submitted along with a
// request that specifically allows for duplicate keys and encodes its entries
// in the same order that they were added.
type RequestValues struct {
	values []formValue
}

// Add adds a key/value tuple to the form.
func (f *RequestValues) Add(key, val string) {
	f.values = append(f.values, formValue{key, val})
}

// Encode encodes the values into “URL encoded” form ("bar=baz&foo=quux").
func (f *RequestValues) Encode() string {
	var buf bytes.Buffer
	for _, v := range f.values {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(v.Key))
		buf.WriteString("=")
		buf.WriteString(url.QueryEscape(v.Value))
	}
	return buf.String()
}

// Empty returns true if no parameters have been set.
func (f *RequestValues) Empty() bool {
	return len(f.values) == 0
}

// Set sets the first instance of a parameter for the given key to the given
// value. If no parameters exist with the key, a new one is added.
//
// Note that Set is O(n) and may be quite slow for a very large parameter list.
func (f *RequestValues) Set(key, val string) {
	for i, v := range f.values {
		if v.Key == key {
			f.values[i].Value = val
			return
		}
	}

	f.Add(key, val)
}

// Get retrieves the list of values for the given key.  If no values exist
// for the key, nil will be returned.
//
// Note that Get is O(n) and may be quite slow for a very large parameter list.
func (f *RequestValues) Get(key string) []string {
	var results []string
	for i, v := range f.values {
		if v.Key == key {
			results = append(results, f.values[i].Value)
		}
	}
	return results
}

// ToValues converts an instance of RequestValues into an instance of
// url.Values. This can be useful in cases where it's useful to make an
// unordered comparison of two sets of request values.
//
// Note that url.Values is incapable of representing certain Rack form types in
// a cohesive way. For example, an array of maps in Rack is encoded with a
// string like:
//
//     arr[][foo]=foo0&arr[][bar]=bar0&arr[][foo]=foo1&arr[][bar]=bar1
//
// Because url.Values is a map, values will be handled in a way that's grouped
// by their key instead of in the order they were added. Therefore the above
// may by encoded to something like (maps are unordered so the actual result is
// somewhat non-deterministic):
//
//     arr[][foo]=foo0&arr[][foo]=foo1&arr[][bar]=bar0&arr[][bar]=bar1
//
// And thus result in an incorrect request to Icheck.
func (f *RequestValues) ToValues() url.Values {
	values := url.Values{}
	for _, v := range f.values {
		values.Add(v.Key, v.Value)
	}
	return values
}

// A key/value tuple for use in the RequestValues type.
type formValue struct {
	Key   string
	Value string
}

// Params is the structure that contains the common properties
// of any *Params structure.
type Params struct {
	Exp []string
	// Headers may be used to provide extra header lines on the HTTP request.
	Headers http.Header
}

// ListParams is the structure that contains the common properties
// of any *ListParams structure.
type ListParams struct {
	Exp     []string
	Limit   int
	Skip    int
	Filters Filters
}

// Filters is a structure that contains a collection of filters for list-related APIs.
type Filters struct {
	f []*filter
}

// AppendTo adds the list of filters to the query string values.
func (f *Filters) AppendTo(values *RequestValues) {
	for _, v := range f.f {
		if len(v.Op) > 0 {
			values.Add(fmt.Sprintf("%v[%v]", v.Key, v.Op), v.Val)
		} else {
			values.Add(v.Key, v.Val)
		}
	}
}

// filter is the structure that contains a filter for list-related APIs.
// It ends up passing query string parameters in the format key[op]=value.
type filter struct {
	Key, Op, Val string
}

// AddFilter adds a new filter with a given key, op and value.
func (f *Filters) AddFilter(key, op, value string) {
	filter := &filter{Key: key, Op: op, Val: value}
	f.f = append(f.f, filter)
}

// AppendTo adds the common parameters to the query string values.
func (p *Params) AppendTo(body *RequestValues) {
	for _, v := range p.Exp {
		body.Add("expand[]", v)
	}
}

// Expand appends a new field to expand.
func (p *ListParams) Expand(f string) {
	p.Exp = append(p.Exp, f)
}

// AppendTo adds the common parameters to the query string values.
func (p *ListParams) AppendTo(body *RequestValues) {
	if len(p.Filters.f) > 0 {
		p.Filters.AppendTo(body)
	}

	if p.Skip > 0 {
		body.Add("skip", strconv.Itoa(p.Skip))
	}

	if p.Limit > 0 {
		if p.Limit > 100 {
			p.Limit = 100
		}

		body.Add("limit", strconv.Itoa(p.Limit))
	}

	for _, v := range p.Exp {
		body.Add("expand[]", v)
	}
}

// ToParams converts a ListParams to a Params by moving over any fields that
// have valid targets in the new type. This is useful because fields in
// Params can be injected directly into an http.Request while generally
// ListParams is only used to build a set of parameters.
func (p *ListParams) ToParams() *Params {
	return &Params{}
}
