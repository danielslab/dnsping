# dnsping
ping for dns that does not wait for a response before sending the next packet.  
This program uses the dns library from https://github.com/miekg/dns.  
Originally built to measure the exact interruption time of the DNS service during a software upgrade.
  


## Features:
---
- send UDP queries, IPv4 and IPv6
- measuring the response times of dnsservers
- detailed response and loss statistics

## Examples:
---
IPv4:
```
dnsping -dnsserver 8.8.8.8 -sleep 100000 -count 5 -verbose=true -timeout 100 -domain "google.de"
dnsping Parameters:
------------------------------
dnsserver: 8.8.8.8
dnsport: 53
query_domain: google.de
query_type: A
timeout in ms: 100
sleep in μs: 100000
Count: 5
verbose: true
------------------------------
sending packetes...

MsgNumber       SendTime                            RTT(ms)         RCode      Answer              
0               04-25-2023 20:26:53.991498576       30              NOERROR    google.de.       300     IN      A       142.251.36.227 
1               04-25-2023 20:26:54.092500824       43              NOERROR    google.de.       300     IN      A       142.251.36.227 
2               04-25-2023 20:26:54.192738779       31              NOERROR    google.de.       300     IN      A       142.251.36.227 
3               04-25-2023 20:26:54.293130013       17              NOERROR    google.de.       300     IN      A       142.251.36.227 
4               04-25-2023 20:26:54.393958129       17              NOERROR    google.de.       300     IN      A       142.251.36.227 
------------------------------
send: 5 received: 5 loss: 0 min_rtt: 17 avg_rtt: 28 max_rtt: 43
NOERROR:5 
```
IPv6:
```
dnsping -dnsserver "[fd20::]" -sleep 100000 -count 5 -verbose=false -timeout 100 -domain "google.de"
dnsping Parameters:
------------------------------
dnsserver: [fd20::]
dnsport: 53
query_domain: google.de
query_type: A
timeout in ms: 100
sleep in μs: 100000
Count: 5
verbose: false
------------------------------
sending packetes...

------------------------------
send: 5 received: 5 loss: 0 min_rtt: 1 avg_rtt: 1 max_rtt: 1
NOERROR:5 

```

## Syntax
---
```
Usage of dnsping:
  -count int
        count of messages to send (default 10)
  -dnsport string
        dnsport to sent requests (default "53")
  -dnsserver string
        dnsserver to sent requests (default "8.8.8.8")
  -domain string
        Request domain (default "google.de")
  -sleep int
        time between querys in μs (default 100000)
  -timeout int
        dns-timeout in ms (default 1000)
  -type string
        dnstype for request (default "A")
  -verbose
        verbose output (default true)
```
