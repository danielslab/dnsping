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
query_type: A
timeout in ms: 100
sleep in μs: 100000
Count: 5
quiet: false
------------------------------
sending packets...

MsgNumber       SendTime                            RTT(ms)         RCode      Answer              
0               05-01-2023 19:03:05.310986414       29              NOERROR    google.de.       300     IN      A       172.217.16.163 
1               05-01-2023 19:03:05.411098547       29              NOERROR    google.de.       300     IN      A       172.217.16.163 
2               05-01-2023 19:03:05.511626303       40              NOERROR    google.de.       300     IN      A       172.217.16.163 
3               05-01-2023 19:03:05.611708589       29              NOERROR    google.de.       300     IN      A       172.217.16.163 
4               05-01-2023 19:03:05.711793748       16              NOERROR    google.de.       300     IN      A       172.217.16.163 
------------------------------
send: 5 received: 5 loss: 0.00% min_rtt: 16ms avg_rtt: 29.00ms max_rtt: 40ms jitter: 4.38ms rtt_variance: 57.84ms^2
NOERROR:5 

rtt distribution:
------------------------------
0ms >= <= 5ms: 0.00% (count: 0)
5ms >= <= 10ms: 0.00% (count: 0)
10ms >= <= 20ms: 20.00% (count: 1)
20ms >= <= 30ms: 60.00% (count: 3)
30ms >= <= 40ms: 0.00% (count: 0)
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
IPv6:
```
dnsping -dnsserver "[fd20::]" -sleep 100000 -count 5 -quiet=false -timeout 1 -domain "google.de" -timeouts_only=true
dnsping Parameters:
------------------------------
dnsserver: [fd20::]
dnsport: 53
query_domain: google.de
query_type: A
timeout in ms: 1
sleep in μs: 100000
Count: 5
quiet: false
------------------------------
sending packets...

MsgNumber       SendTime                            RTT(ms)         RCode      Answer              
0               05-01-2023 19:04:03.915236370       read udp [fd00::1c99:a69e:42fa:6190]:40569->[fd20::]:53: i/o timeout 
3               05-01-2023 19:04:04.217037617       read udp [fd00::1c99:a69e:42fa:6190]:48291->[fd20::]:53: i/o timeout 
------------------------------
send: 5 received: 3 loss: 40.00% min_rtt: 0ms avg_rtt: 0.00ms max_rtt: 1ms jitter: 0.00ms rtt_variance: 0.22ms^2
NOERROR:3 

rtt distribution:
------------------------------
0ms >= <= 5ms: 100.00% (count: 3)
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
        count of messages to send (default 10)
  -dnsport string
        dnsport to sent requests (default "53")
  -dnsserver string
        dnsserver to sent requests (default "8.8.8.8")
  -domain string
        Request domain (default "google.de")
  -qtype string
        dns query type for request (default "A")
  -quiet
        displays only a summary every 10 seconds
  -sleep int
        time between querys in μs (default 100000)
  -timeout int
        dns-timeout in ms (default 1000)
  -timeouts_only
        displays only timeouts or paketloss
```
