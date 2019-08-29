package main

import (
	"bytes"
	"net"
	"strings"
	"testing"
)

func Test4and6(t *testing.T) {
	exp := "216.58.215.78\n2a00:1450:401b:805::200e\n"
	check(t, exp, exp, "google.com")
}

func check(t *testing.T, source, exp, url string) {
	stdin, stdout, probe := mockedProbe(strings.Split(source, "\n"))
	stdin.Write([]byte(url))
	run(probe)
	res := stdout.Bytes()
	if exp != string(res) {
		t.Errorf("Incorrect result. Expected %q, got %q\n", exp, res)
	}
}

func TestEmpty(t *testing.T) {
	exp := ""
	check(t, exp, exp, "google.com")
}

func TestOnly4(t *testing.T) {
	resetFlags()
	*flagOnly4 = true
	source := "216.58.215.78\n2a00:1450:401b:805::200e\n"
	exp := "216.58.215.78\n"
	check(t, source, exp, "google.com")
}

func TestOnly6(t *testing.T) {
	resetFlags()
	*flagOnly6 = true
	source := "216.58.215.78\n2a00:1450:401b:805::200e\n"
	exp := "2a00:1450:401b:805::200e\n"
	check(t, source, exp, "google.com")
}

func TestFormat6(t *testing.T) {
	resetFlags()
	*flagFormat6 = true
	source := "216.58.215.78\n2a00:1450:401b:805::200e\n"
	exp := "216.58.215.78\n[2a00:1450:401b:805::200e]\n"
	check(t, source, exp, "google.com")
}

func TestWithHost(t *testing.T) {
	resetFlags()
	*flagWithHost = true
	source := "216.58.215.78\n2a00:1450:401b:805::200e\n"
	exp := "216.58.215.78 google.com\n2a00:1450:401b:805::200e google.com\n"
	check(t, source, exp, "google.com")
}

func TestFormat6WithHostAndSchema(t *testing.T) {
	resetFlags()
	*flagWithHost = true
	*flagFormat6 = true
	source := "216.58.215.78\n2a00:1450:401b:805::200e\n"
	exp := "216.58.215.78 google.com\n[2a00:1450:401b:805::200e] google.com\n"
	for _, pfix := range []string{"", "http://", "https://"} {
		check(t, source, exp, pfix+"google.com")
	}
}

func resetFlags() {
	*flagOnly4 = false
	*flagOnly6 = false
	*flagFormat6 = false
	*flagWithHost = false
}

func mockedProbe(ips []string) (*bytes.Buffer, *bytes.Buffer, Probe) {
	in := bytes.NewBuffer([]byte{})
	out := bytes.NewBuffer([]byte{})

	return in, out, Probe{
		lookupFunc: Result{ips}.mockLookup,
		reader:     in,
		writer:     out,
	}
}

type Result struct {
	ips []string
}

func (obj Result) mockLookup(url string) ([]net.IP, error) {
	if len(obj.ips) == 0 {
		return []net.IP{}, nil
	}
	res := []net.IP{}
	for _, ip := range obj.ips {
		if ip == "" {
			continue
		}
		res = append(res, net.ParseIP(ip))
	}
	return res, nil
}
