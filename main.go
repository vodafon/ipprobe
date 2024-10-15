package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/vodafon/swork"
)

var (
	flagProcs      = flag.Int("procs", 20, "concurrency")
	flagOnly4      = flag.Bool("ipv4", false, "only IPv4")
	flagOnly6      = flag.Bool("ipv6", false, "only IPv6")
	flagFormat6    = flag.Bool("format6", false, "URL friendly IPv6")
	flagWithHost   = flag.Bool("host", false, "print host too")
	flagOnlyPublic = flag.Bool("public", false, "print only public IPs")
)

type Probe struct {
	lookupFunc      func(string) ([]net.IP, error)
	reader          io.Reader
	writer          io.Writer
	privateIPBlocks []*net.IPNet
}

func main() {
	flag.Parse()
	if *flagProcs < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	probe := Probe{
		lookupFunc: net.LookupIP,
		reader:     os.Stdin,
		writer:     os.Stdout,
	}
	if *flagOnlyPublic {
		probe.privateIPBlocks = privateBlocks()
	}
	run(probe)
}

func (obj Probe) printIP(ip, host string) {
	if *flagWithHost {
		fmt.Fprintf(obj.writer, "%s %s\n", ip, host)
		return
	}
	fmt.Fprintf(obj.writer, "%s\n", ip)
}

func (obj Probe) Process(url string) {
	if strings.HasPrefix(url, "http") {
		for _, sch := range [2]string{"http://", "https://"} {
			url = strings.TrimLeft(url, sch)
		}
	}
	ips, err := obj.lookupFunc(url)
	if err != nil || len(ips) == 0 {
		return
	}
	for _, ip := range ips {
		if *flagOnlyPublic && obj.isPrivateIP(ip) {
			continue
		}
		isIPv6 := ip.To4() == nil
		if *flagOnly4 && isIPv6 {
			continue
		}
		if *flagOnly6 && !isIPv6 {
			continue
		}
		if *flagFormat6 && isIPv6 {
			obj.printIP(fmt.Sprintf("[%s]", ip), url)
			continue
		}
		obj.printIP(ip.String(), url)
	}
}

func run(probe Probe) {
	w := swork.NewWorkerGroup(*flagProcs, probe)

	sc := bufio.NewScanner(probe.reader)
	for sc.Scan() {
		w.StringC <- sc.Text()
	}

	close(w.StringC)

	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
	}

	w.Wait()
}

func privateBlocks() []*net.IPNet {
	var privateIPBlocks []*net.IPNet
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Fatalf("parse error on %q: %v", cidr, err)
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}
	return privateIPBlocks
}

func (obj Probe) isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	for _, block := range obj.privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}
