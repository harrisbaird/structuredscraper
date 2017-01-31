package q_test

import (
	"strings"
	"testing"

	xmlpath "gopkg.in/xmlpath.v2"

	. "github.com/harrisbaird/flexiscraper/q"
	"github.com/nbio/st"
)

func TestReplace(t *testing.T) {
	tests := []struct {
		name     string
		template string
		value    []string
		result   []string
		wantErr  bool
	}{
		{"blank", "", []string{}, []string{}, false},
		{"valid", "hello %s", []string{"world"}, []string{"hello world"}, false},
		{"missing", "hello", []string{"world"}, []string{"hello%!(EXTRA string=world)"}, false},
		{"extra", "hello %s, %s", []string{"world"}, []string{"hello world, %!s(MISSING)"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Replace(tt.template)(nil, tt.value)
			st.Assert(t, err != nil, tt.wantErr)
			st.Assert(t, result, tt.result)
		})
	}
}

func TestRegexp(t *testing.T) {
	tests := []struct {
		name    string
		regex   string
		value   []string
		result  []string
		wantErr bool
	}{
		{"blank", "", []string{""}, []string{""}, false},
		{"invalid", "(.*", []string{"value"}, []string{"value"}, true},
		{"valid", "\\d+", []string{"Post 5"}, []string{"5"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Regexp(tt.regex)(nil, tt.value)
			st.Assert(t, err != nil, tt.wantErr)
			st.Assert(t, result, tt.result)
		})
	}
}

func TestXPath(t *testing.T) {
	r := strings.NewReader("<!DOCTYPE html><html><head><title>Hello world</title></head><body><h1>Test</h1></body></html>")
	node, err := xmlpath.ParseHTML(r)
	st.Assert(t, err, nil)

	tests := []struct {
		name    string
		exp     string
		result  []string
		wantErr bool
	}{
		{"blank", "", []string{}, true},
		{"valid", "//title", []string{"Hello world"}, false},
		{"invalid", "~~Test", []string{}, true},
		{"no match", "//invalid", []string{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := XPath(tt.exp)(node, []string{})
			st.Assert(t, err != nil, tt.wantErr)
			st.Assert(t, result, tt.result)
		})
	}
}
