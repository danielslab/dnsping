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

			// Divide the index by the length of the alphabet, so you know how many times the row has already been run through completely
			index = index / lenAlphabet
			// Add the corresponding character of the alphabet to the beginning of the combination
			combination = string(alphabet[remainder]) + combination
		}

		queue_channel <- combination
	}
	//close channel when all labels were generated
	close(queue_channel)

}

func send_query(msgnumber int, cfg *dnspingConfig, domain *string, statistic *statistic.Statistic) {

	//Create DNS-Query-Message
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{*domain + ".", cfg.dnstype, dns.ClassINET}

	//create dnsclient with timeout and protocol, default is udp
	c := new(dns.Client)

	if cfg.tcp == true {
		c.Net = "tcp"
		laddr := net.TCPAddr{
			IP:   net.ParseIP(cfg.source),
			Zone: "",
		}
		c.Dialer = &net.Dialer{
			Timeout:   time.Duration(cfg.timeout) * time.Millisecond,
			LocalAddr: &laddr,
		}
	} else {
		c.Net = "udp"
		laddr := net.UDPAddr{
			IP:   net.ParseIP(cfg.source),
			Zone: "",
		}
		c.Dialer = &net.Dialer{
			Timeout:   time.Duration(cfg.timeout) * time.Millisecond,
			LocalAddr: &laddr,
		}
	}

	//Save sendtime
	send_time := time.Now().Format("01-02-2006 15:04:05.000000000")

	// increase send counter
	statistic.Increase_send_counter()

	//send dns message
	in, rtt, err := c.Exchange(m1, cfg.dnsserver+":"+cfg.dnsport)

	if cfg.quiet == false {

		if err == nil && cfg.timeouts_only == false {
			//no error

			// if RCode NoERROR
			if in.Rcode == 0 {
				if len(in.Answer) > 0 {
					fmt.Printf("%-15s %-35s %-15d %-10s %-20s \n", strconv.Itoa(msgnumber), send_time, rtt.Milliseconds(), dns.RcodeToString[in.Rcode], in.Answer[0].String())
				} else {
					fmt.Printf("%-15s %-35s %-15d %-10s %-20s \n", strconv.Itoa(msgnumber), send_time, rtt.Milliseconds(), dns.RcodeToString[in.Rcode], "no Answer to decode")
				}

			} else {
				fmt.Printf("%-15s %-35s %-15d %-10s \n", strconv.Itoa(msgnumber), send_time, rtt.Milliseconds(), dns.RcodeToString[in.Rcode])
			}

		}
		if err != nil {
			//error or timeout
			fmt.Printf("%-15s %-35s %-20s \n", strconv.Itoa(msgnumber), send_time, err)
		}

	}

	if err == nil {
		// increase received counter
		statistic.Increase_received_counter()
		statistic.Append_rtt(int(rtt.Milliseconds()))
		statistic.Increase_rcode_counter(dns.RcodeToString[in.Rcode])
	}

}

func QPS_to_Time(qps int) int {
	if qps <= 0 || qps > 1000000 {
		fmt.Println("QPS range has to be between 1 to 1000000")
		os.Exit(1)
	}
	time := 1000000 / qps
	return time
}

type dnspingConfig struct {
	dnsserver     string
	dnsport       string
	domain        string
	timeout       int
	count         int
	qps           int
	quiet         bool
	qtype         string
	dnstype       uint16
	timeouts_only bool
	flame         bool
	tcp           bool
	source        string
	interim       int
}

func (cfg *dnspingConfig) Print() {
	fmt.Println("dnsping Parameters:")
	fmt.Println("------------------------------")
	fmt.Println("src-addr: " + cfg.source)
	fmt.Println("dnsserver: " + cfg.dnsserver)
	fmt.Println("dnsport: " + cfg.dnsport)
	fmt.Println("query_domain: " + cfg.domain)
	fmt.Println("flame: " + strconv.FormatBool(cfg.flame))
	fmt.Println("query_type: " + dns.TypeToString[cfg.dnstype])
	fmt.Println("timeout in ms: " + strconv.Itoa(cfg.timeout))
	fmt.Println("qps: " + strconv.Itoa(cfg.qps))
	fmt.Println("count: " + strconv.Itoa(cfg.count))
	fmt.Println("quiet: "+strconv.FormatBool(cfg.quiet), "(interim-stats-timer: "+strconv.Itoa(cfg.interim)+")")
	fmt.Println("tcp: " + strconv.FormatBool(cfg.tcp))
	fmt.Println("timeouts_only: " + strconv.FormatBool(cfg.timeouts_only))
	fmt.Println("------------------------------")
	fmt.Println("sending packets...")
	fmt.Println("")
}

func main() {

	//Channel as Label-Queue
	queue_channel := make(chan string, 10000000)

	//init type dnspingConfig
	cfg := dnspingConfig{}

	flag.StringVar(&cfg.dnsserver, "dnsserver", "8.8.8.8", "dnsserver to sent requests")
	flag.StringVar(&cfg.dnsport, "dnsport", "53", "dst-port to send requests")
	flag.StringVar(&cfg.domain, "domain", "google.de", "request domain")
	flag.IntVar(&cfg.timeout, "timeout", 1000, "dns-timeout in ms")
	flag.IntVar(&cfg.count, "count", 5, "count of messages to send. count 0 sets count to max")
	flag.IntVar(&cfg.qps, "qps", 5, "desired querys per second (1 to 1000000)")
	flag.BoolVar(&cfg.quiet, "quiet", false, "displays only a interim-stats every 10 seconds")
	flag.BoolVar(&cfg.timeouts_only, "timeouts_only", false, "displays only timeouts or paketloss")
	flag.StringVar(&cfg.qtype, "qtype", "A", "dns query type for request")
	flag.BoolVar(&cfg.flame, "flame", false, "adds a 13 digit (aaaaaaaaaaaaa - zzzzzzzzzzzz) increasing subdomain in front of the domain for each query.")
	flag.BoolVar(&cfg.tcp, "tcp", false, "send tcp querys instead of udp")
	flag.StringVar(&cfg.source, "src", "", "local address to sent requests")
	flag.IntVar(&cfg.interim, "interim", 10, "time between interim-stats for quiet-mode")

	flag.Parse()

	cfg.dnstype = dns.StringToType[cfg.qtype]
	qps_time := time.Duration(QPS_to_Time(cfg.qps)) * time.Microsecond

	if cfg.interim <= 0 || cfg.count < 0 || cfg.timeout <= 0 {
		fmt.Println("wrong input values")
		os.Exit(1)
	}

	if cfg.count == 0 {
		//Set count to max integer when 0
		cfg.count = math.MaxInt
	}

	var waitGroup sync.WaitGroup

	//Print dnsping settings
	cfg.Print()

	statistic := statistic.Init_Statistic()

	if cfg.quiet == false {
		// Print Head-line
		fmt.Printf("%-15s %-35s %-15s %-10s %-20s\n", "MsgNumber", "SendTime", "RTT(ms)", "RCode", "Answer snipped")
	} else {
		go func(interim int) {
			// Print Interim-Statstics
			old_rx_counter := 0
			old_tx_counter := 0
			for {
				start_time := time.Now()
				time.Sleep(time.Duration(interim) * time.Second)
				stop_time := time.Now()
				statistic.Print_Summary()
				old_tx_counter = statistic.Print_tx_pps_on_Wire(start_time, stop_time, old_tx_counter)
				old_rx_counter = statistic.Print_rx_pps_on_Wire(start_time, stop_time, old_rx_counter)
			}
		}(cfg.interim)
	}

	//generate random labels if flame true
	if cfg.flame == true {
		go generateLabels(queue_channel)
	}

	start_time := time.Now()

	//capture ctr+c signal SIGINT
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		statistic.Print_Summary()
		statistic.Print_tx_pps_on_Wire(start_time, time.Now(), 0)
		statistic.Print_rx_pps_on_Wire(start_time, time.Now(), 0)
		statistic.RTT_Summary()
		os.Exit(1)
	}()

	for msgnumber := 1; msgnumber <= cfg.count; msgnumber++ {

		waitGroup.Add(1)
		//send querys parralel out if flame = false
		if cfg.flame == false {
			go func(msgnumber int) {
				send_query(msgnumber, &cfg, &cfg.domain, statistic)
				waitGroup.Done()
			}(msgnumber)
		} else {
			subdomain_label, ok := <-queue_channel
			if ok {
				go func(msgnumber int) {
					new_domain := subdomain_label + "." + cfg.domain
					send_query(msgnumber, &cfg, &new_domain, statistic)
					waitGroup.Done()
				}(msgnumber)
			} else {
				fmt.Println("received all Labels from label_queue")
				break
			}
		}
		time.Sleep(qps_time)
	}

	send_stop_time := time.Now()
	waitGroup.Wait()
	recv_stop_time := time.Now()

	statistic.Print_Summary()
	statistic.Print_tx_pps_on_Wire(start_time, send_stop_time, 0)
	statistic.Print_rx_pps_on_Wire(start_time, recv_stop_time, 0)
	statistic.RTT_Summary()

}
