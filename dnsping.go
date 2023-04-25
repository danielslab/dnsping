package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/danielslab/dnsping/statistic"
	"github.com/miekg/dns"
)

func send_query(msgnumber int, dnsserver, dnsport, domain string, dnstype uint16, timeout uint16, verbose bool, statistic *statistic.Statistic) {

	//Create DNS-Query-Nessage
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{domain + ".", dnstype, dns.ClassINET}

	//Create dnsclient with timeout
	c := new(dns.Client)
	c.Dialer = &net.Dialer{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}

	//Save sendtime
	dt := time.Now().Format("01-02-2006 15:04:05.000000000")
	//send dns message
	in, rtt, err := c.Exchange(m1, dnsserver+":"+dnsport)

	if verbose == true {

		if err == nil {
			//fmt.Println("Kein Fehler:")
			//fmt.Println("MsgNu:"+strconv.Itoa(msgnumber), "SendTime:"+dt, in, rtt, err)

			// if RCode NoERROR
			if in.Rcode == 0 {
				if len(in.Answer) > 0 {
					fmt.Printf("%-15s %-35s %-15d %-10s %-20s \n", strconv.Itoa(msgnumber), dt, rtt.Milliseconds(), dns.RcodeToString[in.Rcode], in.Answer[0].String())
				} else {
					fmt.Printf("%-15s %-35s %-15d %-10s %-20s \n", strconv.Itoa(msgnumber), dt, rtt.Milliseconds(), dns.RcodeToString[in.Rcode], "no Answer to decode")
				}

			} else {
				fmt.Printf("%-15s %-35s %-15d %-10s \n", strconv.Itoa(msgnumber), dt, rtt.Milliseconds(), dns.RcodeToString[in.Rcode])
			}

		} else {
			//fmt.Println("Fehler:")
			//fmt.Println("MsgNu:"+strconv.Itoa(msgnumber), "SendTime:"+dt, in, rtt, err)
			fmt.Printf("%-15s %-35s %-20s \n", strconv.Itoa(msgnumber), dt, err)
		}

	}

	// increse send counter
	statistic.Increase_send_counter()

	if err == nil {
		// increse received counter
		statistic.Increase_received_counter()
		statistic.Append_rtt(int(rtt.Milliseconds()))
		statistic.Increase_rcode_counter(dns.RcodeToString[in.Rcode])
	}

}

func main() {

	var dnsserver string
	var dnsport string
	var domain string
	var timeout int
	var count int
	var sleep int
	var verbose bool
	var dtype string

	flag.StringVar(&dnsserver, "dnsserver", "8.8.8.8", "dnsserver to sent requests")
	flag.StringVar(&dnsport, "dnsport", "53", "dnsport to sent requests")
	flag.StringVar(&domain, "domain", "google.de", "Request domain")
	flag.IntVar(&timeout, "timeout", 1000, "dns-timeout in ms")
	flag.IntVar(&count, "count", 10, "count of messages to send")
	flag.IntVar(&sleep, "sleep", 100000, "time between querys in μs")
	flag.BoolVar(&verbose, "verbose", true, "verbose output")
	flag.StringVar(&dtype, "type", "A", "dnstype for request")

	dnstype := dns.StringToType[dtype]

	flag.Parse()

	var waitGroup sync.WaitGroup

	//dnsserver := "[fd20::]"
	//dnsserver := "172.16.0.1"
	//dnsport := "53"
	//domain := "fritz.box"
	//dnstype := dns.TypeA
	//timeout := 20 //in ms
	//count := 10
	//sleep := 100
	//verbose := true

	fmt.Println("dnsping Parameters:")
	fmt.Println("------------------------------")
	fmt.Println("dnsserver: " + dnsserver)
	fmt.Println("dnsport: " + dnsport)
	fmt.Println("query_domain: " + domain)
	fmt.Println("query_type: " + dns.TypeToString[dnstype])
	fmt.Println("timeout in ms: " + strconv.Itoa(timeout))
	fmt.Println("sleep in μs: " + strconv.Itoa(sleep))
	fmt.Println("Count: " + strconv.Itoa(count))
	fmt.Println("verbose: " + strconv.FormatBool(verbose))
	fmt.Println("------------------------------")
	fmt.Println("sending packetes...")
	fmt.Println("")

	statistic := statistic.Init_Statistic()

	if verbose == true {
		// Print Head-line
		fmt.Printf("%-15s %-35s %-15s %-10s %-20s\n", "MsgNumber", "SendTime", "RTT(ms)", "RCode", "Answer")
	} else {
		go func() {
			for {
				// Print Interim-Statstics
				time.Sleep(10 * time.Second)
				statistic.Print()
			}
		}()

	}

	//capture ctr+c signal SIGINT
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		statistic.Print()
		os.Exit(1)
	}()

	//send querys parralel out
	for i := 0; i < count; i++ {
		time.Sleep(time.Duration(sleep) * time.Microsecond)
		waitGroup.Add(1)
		go func(i int) {
			send_query(i, dnsserver, dnsport, domain, dnstype, uint16(timeout), verbose, statistic)
			waitGroup.Done()
		}(i)

	}

	waitGroup.Wait()

	statistic.Print()

}
