# dnsping
ping for dns that does not wait for responses before sending the next packet.  
Originally built to measure the exact interruption time of the dns service during a software upgrade.  
This program uses the dns library from https://github.com/miekg/dns.  


## Features:
---
- sends udp and tcp queries
- measuring the response times of dnsservers
- detailed response and loss statistics

## Examples:
---
IPv4:
```
dnsping -dnsserver 8.8.8.8 -count 5 -timeout 100 -domain "google.de"
dnsping Parameters:
------------------------------
src-addr: 
dnsserver: 8.8.8.8
dnsport: 53
query_domain: google.de
flame: false
query_type: A
timeout in ms: 100
qps: 5
count: 5
quiet: false (interim-stats-timer: 10)
tcp: false
timeouts_only: false
------------------------------
sending packets...

MsgNumber       SendTime                            RTT(ms)         RCode      Answer snipped      
1               05-07-2023 12:08:39.994507474       21              NOERROR    google.de.       269     IN      A       172.217.16.163 
2               05-07-2023 12:08:40.195365301       17              NOERROR    google.de.       284     IN      A       172.217.16.163 
3               05-07-2023 12:08:40.395440607       12              NOERROR    google.de.       229     IN      A       172.217.16.163 
4               05-07-2023 12:08:40.596195298       17              NOERROR    google.de.       230     IN      A       172.217.16.163 
5               05-07-2023 12:08:40.796211690       16              NOERROR    google.de.       230     IN      A       172.217.16.163 
------------------------------
send: 5 received: 5 loss: 0.00% (count: 0) min_rtt: 12ms avg_rtt: 17.00ms max_rtt: 21ms jitter: 1.38ms rtt_variance: 8.24ms^2
RCodes: NOERROR:5 
tx pps: 4.987186
rx pps: 4.987186

rtt distribution:
------------------------------
rtt < 5ms: 0.00% (count: 0)
rtt < 10ms: 0.00% (count: 0)
rtt < 20ms: 80.00% (count: 4)
rtt < 30ms: 20.00% (count: 1)
rtt < 40ms: 0.00% (count: 0)
rtt < 50ms: 0.00% (count: 0)
rtt < 60ms: 0.00% (count: 0)
rtt < 70ms: 0.00% (count: 0)
rtt < 80ms: 0.00% (count: 0)
rtt < 90ms: 0.00% (count: 0)
rtt < 100ms: 0.00% (count: 0)
rtt < 250ms: 0.00% (count: 0)
rtt < 400ms: 0.00% (count: 0)
rtt < 550ms: 0.00% (count: 0)
rtt < 700ms: 0.00% (count: 0)
rtt < 850ms: 0.00% (count: 0)
rtt < 1000ms: 0.00% (count: 0)
rtt < 10000ms: 0.00% (count: 0)
```
Example with ipv6 dnssserver where only packet losses/timeouts are displayed :
```
dnsping -dnsserver "[fd20::]" -qps 10 -count 5 -quiet=false -timeout 1 -domain "google.de" -timeouts_only=true
dnsping Parameters:
------------------------------
src-addr: 
dnsserver: [fd20::]
dnsport: 53
query_domain: google.de
flame: false
query_type: A
timeout in ms: 1
qps: 10
count: 5
quiet: false (interim-stats-timer: 10)
tcp: false
timeouts_only: false
------------------------------
sending packets...

MsgNumber       SendTime                            RTT(ms)         RCode      Answer snipped      
2               05-07-2023 12:09:58.606020319       read udp [fd00::f4f2:c231:942:e7e0]:40763->[fd20::]:53: i/o timeout 
4               05-07-2023 12:09:58.807051536       read udp [fd00::f4f2:c231:942:e7e0]:38936->[fd20::]:53: i/o timeout 
------------------------------
send: 5 received: 3 loss: 40.00% (count: 2) min_rtt: 0ms avg_rtt: 0.00ms max_rtt: 1ms jitter: 0.50ms rtt_variance: 0.22ms^2
RCodes: NOERROR:3 
tx pps: 5.000000
rx pps: 3.000000

rtt distribution:
------------------------------
rtt < 5ms: 100.00% (count: 3)
rtt < 10ms: 0.00% (count: 0)
rtt < 20ms: 0.00% (count: 0)
rtt < 30ms: 0.00% (count: 0)
rtt < 40ms: 0.00% (count: 0)
rtt < 50ms: 0.00% (count: 0)
rtt < 60ms: 0.00% (count: 0)
rtt < 70ms: 0.00% (count: 0)
rtt < 80ms: 0.00% (count: 0)
rtt < 90ms: 0.00% (count: 0)
rtt < 100ms: 0.00% (count: 0)
rtt < 250ms: 0.00% (count: 0)
rtt < 400ms: 0.00% (count: 0)
rtt < 550ms: 0.00% (count: 0)
rtt < 700ms: 0.00% (count: 0)
rtt < 850ms: 0.00% (count: 0)
rtt < 1000ms: 0.00% (count: 0)
rtt < 10000ms: 0.00% (count: 0)
```

## Syntax
---
```
Usage of dnsping:
  -count int
        count of messages to send. count 0 sets count to max (default 5)
  -dnsport string
        dst-port to send requests (default "53")
  -dnsserver string
        dnsserver to sent requests (default "8.8.8.8")
  -domain string
        request domain (default "google.de")
  -flame
        adds a 13 digit (aaaaaaaaaaaaa - zzzzzzzzzzzz) increasing subdomain in front of the domain for each query.
  -interim int
        time between interim-stats for quiet-mode (default 10)
  -qps int
        desired querys per second (1 to 1000000) (default 5)
  -qtype string
        dns query type for request (default "A")
  -quiet
        displays only a interim-stats every 10 seconds
  -src string
        local address to send requests
  -tcp
        send tcp querys instead of udp
  -timeout int
        dns-timeout in ms (default 1000)
  -timeouts_only
        displays only timeouts or paketloss
  -version
        print version of dnsping
```

## statically linked build
---
```
git clone https://github.com/danielslab/dnsping
cd dnsping
go mod tidy
go build -ldflags "-linkmode 'external' -extldflags '-static'"
```
