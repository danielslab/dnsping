# dnsping
ping for dns that does not wait for a response before sending the next packet.  
Originally built to measure the exact interruption time of the DNS service during a software upgrade.  
This program uses the dns library from https://github.com/miekg/dns.  

  


## Features:
---
- send UDP queries, IPv4 and IPv6
- measuring the response times of dnsservers
- detailed response and loss statistics

## Examples:
---
IPv4:
```
dnsping -dnsserver 8.8.8.8 -count 5 -timeout 100 -domain "google.de"
dnsping Parameters:
------------------------------
dnsserver: 8.8.8.8
dnsport: 53
query_domain: google.de
flame: false
query_type: A
timeout in ms: 100
qps: 5
count: 5
quiet: false
------------------------------
sending packets...

MsgNumber       SendTime                            RTT(ms)         RCode      Answer snipped      
1               05-05-2023 14:48:52.600647146       43              NOERROR    google.de.       300     IN      A       142.251.36.227 
2               05-05-2023 14:48:52.801061651       31              NOERROR    google.de.       300     IN      A       142.251.36.227 
3               05-05-2023 14:48:53.002170262       29              NOERROR    google.de.       300     IN      A       142.251.36.227 
4               05-05-2023 14:48:53.202679493       31              NOERROR    google.de.       300     IN      A       142.251.36.227 
5               05-05-2023 14:48:53.403414168       18              NOERROR    google.de.       300     IN      A       142.251.36.227 
------------------------------
send: 5 received: 5 loss: 0.00% min_rtt: 18ms avg_rtt: 30.00ms max_rtt: 43ms jitter: 5.25ms rtt_variance: 63.04ms^2
RCodes: NOERROR:5 
tx pps: 4.981294
rx pps: 4.891380

rtt distribution:
------------------------------
0ms >= <= 5ms: 0.00% (count: 0)
5ms >= <= 10ms: 0.00% (count: 0)
10ms >= <= 20ms: 20.00% (count: 1)
20ms >= <= 30ms: 20.00% (count: 1)
30ms >= <= 40ms: 40.00% (count: 2)
40ms >= <= 50ms: 20.00% (count: 1)
50ms >= <= 60ms: 0.00% (count: 0)
60ms >= <= 70ms: 0.00% (count: 0)
70ms >= <= 80ms: 0.00% (count: 0)
80ms >= <= 90ms: 0.00% (count: 0)
90ms >= <= 100ms: 0.00% (count: 0)
100ms >= <= 250ms: 0.00% (count: 0)
250ms >= <= 400ms: 0.00% (count: 0)
400ms >= <= 550ms: 0.00% (count: 0)
550ms >= <= 700ms: 0.00% (count: 0)
700ms >= <= 850ms: 0.00% (count: 0)
850ms >= <= 1000ms: 0.00% (count: 0)
1000ms >= <= 10000ms: 0.00% (count: 0)
```
Example with ipv6 dnssserver where only packet losses/timeouts are displayed :
```
dnsping -dnsserver "[fd20::]" -qps 10 -count 5 -quiet=false -timeout 10 -domain "google.de" -timeouts_only=true
dnsping Parameters:
------------------------------
dnsserver: [fd20::]
dnsport: 53
query_domain: google.de
flame: false
query_type: A
timeout in ms: 10
qps: 10
count: 5
quiet: false
------------------------------
sending packets...

MsgNumber       SendTime                            RTT(ms)         RCode      Answer snipped      
1               05-05-2023 14:47:36.649159629       read udp [fd00::778d:700f:c72a:a365]:57738->[fd20::]:53: i/o timeout 
------------------------------
send: 5 received: 4 loss: 20.00% min_rtt: 1ms avg_rtt: 1.00ms max_rtt: 1ms jitter: 0.00ms rtt_variance: 0.00ms^2
RCodes: NOERROR:4 
tx pps: 9.929743
rx pps: 7.919308

rtt distribution:
------------------------------
0ms >= <= 5ms: 100.00% (count: 4)
5ms >= <= 10ms: 0.00% (count: 0)
10ms >= <= 20ms: 0.00% (count: 0)
20ms >= <= 30ms: 0.00% (count: 0)
30ms >= <= 40ms: 0.00% (count: 0)
40ms >= <= 50ms: 0.00% (count: 0)
50ms >= <= 60ms: 0.00% (count: 0)
60ms >= <= 70ms: 0.00% (count: 0)
70ms >= <= 80ms: 0.00% (count: 0)
80ms >= <= 90ms: 0.00% (count: 0)
90ms >= <= 100ms: 0.00% (count: 0)
100ms >= <= 250ms: 0.00% (count: 0)
250ms >= <= 400ms: 0.00% (count: 0)
400ms >= <= 550ms: 0.00% (count: 0)
550ms >= <= 700ms: 0.00% (count: 0)
700ms >= <= 850ms: 0.00% (count: 0)
850ms >= <= 1000ms: 0.00% (count: 0)
1000ms >= <= 10000ms: 0.00% (count: 0)

```

## Syntax
---
```
Usage of dnsping:
  -count int
        count of messages to send. count 0 sets count unlimited (default 10)
  -dnsport string
        dnsport to sent requests (default "53")
  -dnsserver string
        dnsserver to sent requests (default "8.8.8.8")
  -domain string
        Request domain (default "google.de")
  -flame
        adds a 13 digit (aaaaaaaaaaaaa - zzzzzzzzzzzz) increasing subdomain in front of the domain for each query.
  -qps int
        desired querys per second (1 to 1000000) (default 5)
  -qtype string
        dns query type for request (default "A")
  -quiet
        displays only a summary every 10 seconds
  -timeout int
        dns-timeout in ms (default 1000)
  -timeouts_only
        displays only timeouts or paketloss
```
