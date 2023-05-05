package main

import (
	"flag"
	"fmt"
	"math"
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

func generateLabels(queue_channel chan string) {
	//alphabet^character_count generates 2481152873203736576 label combinations
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	character_count := 13 //13
	lenAlphabet := len(alphabet)
	numCombinations := int(math.Pow(float64(lenAlphabet), float64(character_count)))

	for i := 0; i < numCombinations; i++ {
		// Convert the index to a combination
		index := i
		combination := ""
		// Iterate over each character in the combination
		for j := 0; j < character_count; j++ {
			// Determine the rest of the division of the index by the length of the alphabet
			remainder := index % lenAlphabet

			//fmt.Printf("i: %d j: %d index: %d remainder=%d%%%d=%d newindex=%d/%d=%d %c\n", i, j, index, index, lenAlphabet, remainder, index, lenAlphabet, index/lenAlphabet, alphabet[remainder])

			// Divide the index by the length of the alphabet, so you know how many times the row has already been run through completely
			index = index / lenAlphabet
			// Add the corresponding character of the alphabet to the beginning of the combination
			combination = string(alphabet[remainder]) + combination
		}

		queue_channel <- combination
	}
	//close channel when all labels were genrated
	close(queue_channel)

}

func send_query(msgnumber int, dnsserver, dnsport, domain string, dnstype uint16, timeout uint16, quiet bool, statistic *statistic.Statistic, timeouts_only bool) {

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

	if quiet == false {

		if err == nil && timeouts_only == false {
			//no error

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

		}
		if err != nil {
			//error or timeout
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

func QPS_to_Time(qps int) int {
	if qps <= 0 || qps >= 1000000 {
		fmt.Println("QPS range has to be between 1 to 1000000")
		os.Exit(1)
	}
	time := 1000000 / qps
	return time
}

func main() {

	//Channel as Label-Queue
	queue_channel := make(chan string, 10000000)

	var dnsserver string
	var dnsport string
	var domain string
	var timeout int
	var count int
	var qps int
	var quiet bool
	var qtype string
	var timeouts_only bool
	var flame bool

	flag.StringVar(&dnsserver, "dnsserver", "8.8.8.8", "dnsserver to sent requests")
	flag.StringVar(&dnsport, "dnsport", "53", "dnsport to sent requests")
	flag.StringVar(&domain, "domain", "google.de", "Request domain")
	flag.IntVar(&timeout, "timeout", 1000, "dns-timeout in ms")
	flag.IntVar(&count, "count", 10, "count of messages to send. count 0 sets count unlimited")
	flag.IntVar(&qps, "qps", 5, "desired querys per second (1 to 1000000)")
	flag.BoolVar(&quiet, "quiet", false, "displays only a summary every 10 seconds")
	flag.BoolVar(&timeouts_only, "timeouts_only", false, "displays only timeouts or paketloss")
	flag.StringVar(&qtype, "qtype", "A", "dns query type for request")
	flag.BoolVar(&flame, "flame", false, "adds a 13 digit (aaaaaaaaaaaaa - zzzzzzzzzzzz) increasing subdomain in front of the domain for each query.")

	flag.Parse()

	dnstype := dns.StringToType[qtype]
	qps_time := time.Duration(QPS_to_Time(qps)) * time.Microsecond

	if count == 0 {
		//Set count to max integer when 0
		count = math.MaxInt
	}

	var waitGroup sync.WaitGroup

	fmt.Println("dnsping Parameters:")
	fmt.Println("------------------------------")
	fmt.Println("dnsserver: " + dnsserver)
	fmt.Println("dnsport: " + dnsport)
	fmt.Println("query_domain: " + domain)
	fmt.Println("flame: " + strconv.FormatBool(flame))
	fmt.Println("query_type: " + dns.TypeToString[dnstype])
	fmt.Println("timeout in ms: " + strconv.Itoa(timeout))
	fmt.Println("qps: " + strconv.Itoa(qps))
	fmt.Println("count: " + strconv.Itoa(count))
	fmt.Println("quiet: " + strconv.FormatBool(quiet))
	fmt.Println("------------------------------")
	fmt.Println("sending packets...")
	fmt.Println("")

	statistic := statistic.Init_Statistic()

	if quiet == false {
		// Print Head-line
		fmt.Printf("%-15s %-35s %-15s %-10s %-20s\n", "MsgNumber", "SendTime", "RTT(ms)", "RCode", "Answer snipped")
	} else {
		go func() {
			// Print Interim-Statstics
			old_send_counter := 0
			for {
				start_time := time.Now()
				time.Sleep(10 * time.Second)
				stop_time := time.Now()
				statistic.Print()
				old_send_counter = statistic.Print_QPS_on_Wire(start_time, stop_time, old_send_counter)
			}
		}()

	}

	//genarate random Labels if flame true
	if flame == true {
		go generateLabels(queue_channel)
	}

	start_time := time.Now()

	//capture ctr+c signal SIGINT
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		statistic.Print()
		statistic.Print_QPS_on_Wire(start_time, time.Now(), 0)
		statistic.RTT_Summary()
		os.Exit(1)
	}()

	for i := 1; i <= count; i++ {
		waitGroup.Add(1)
		//send querys parralel out if flame = false
		if flame == false {
			go func(i int) {
				send_query(i, dnsserver, dnsport, domain, dnstype, uint16(timeout), quiet, statistic, timeouts_only)
				waitGroup.Done()
			}(i)
		} else {
			subdomain_label, ok := <-queue_channel
			if ok {
				go func(i int) {
					new_domain := subdomain_label + "." + domain
					send_query(i, dnsserver, dnsport, new_domain, dnstype, uint16(timeout), quiet, statistic, timeouts_only)
					waitGroup.Done()
				}(i)
			} else {
				fmt.Println("received all Labels from label_queue")
				break
			}
		}
		time.Sleep(qps_time)
	}

	waitGroup.Wait()
	stop_time := time.Now()

	statistic.Print()
	statistic.Print_QPS_on_Wire(start_time, stop_time, 0)
	statistic.RTT_Summary()

}
