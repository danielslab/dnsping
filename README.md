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
dnsping -dnsserver 8.8.8.8 -sleep 100000 -count 5 -timeout 100 -domain "google.de"
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
sending packetes...

MsgNumber       SendTime                            RTT(ms)         RCode      Answer              
0               04-30-2023 20:18:53.428613339       30              NOERROR    google.de.       300     IN      A       172.217.16.163 
1               04-30-2023 20:18:53.529026533       17              NOERROR    google.de.       300     IN      A       172.217.16.163 
2               04-30-2023 20:18:53.629649322       18              NOERROR    google.de.       300     IN      A       172.217.16.163 
3               04-30-2023 20:18:53.730197407       30              NOERROR    google.de.       300     IN      A       172.217.16.163 
4               04-30-2023 20:18:53.830456475       30              NOERROR    google.de.       300     IN      A       172.217.16.163 
------------------------------
send: 5 received: 5 loss: 0.00% min_rtt: 17ms avg_rtt: 25.00ms max_rtt: 30ms
NOERROR:5 

rtt distribution:
------------------------------
0ms to 5ms: 0.00%
5ms to 10ms: 0.00%
10ms to 20ms: 40.00%
20ms to 30ms: 0.00%
30ms to 40ms: 60.00%
40ms to 50ms: 0.00%
50ms to 60ms: 0.00%
60ms to 70ms: 0.00%
70ms to 80ms: 0.00%
80ms to 90ms: 0.00%
90ms to 100ms: 0.00%
100ms to 250ms: 0.00%
250ms to 400ms: 0.00%
400ms to 550ms: 0.00%
550ms to 700ms: 0.00%
700ms to 850ms: 0.00%
850ms to 1000ms: 0.00%
1000ms to 10000ms: 0.00%
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
sending packetes...

MsgNumber       SendTime                            RTT(ms)         RCode      Answer              
0               04-30-2023 20:16:39.719218229       read udp [fd00::a735:a204:59cf:5f4a]:34155->[fd20::]:53: i/o timeout 
2               04-30-2023 20:16:39.920010956       read udp [fd00::a735:a204:59cf:5f4a]:45043->[fd20::]:53: i/o timeout 
------------------------------
send: 5 received: 3 loss: 40.00% min_rtt: 0ms avg_rtt: 0.00ms max_rtt: 0ms
NOERROR:3 

rtt distribution:
------------------------------
0ms to 5ms: 100.00%
5ms to 10ms: 0.00%
10ms to 20ms: 0.00%
20ms to 30ms: 0.00%
30ms to 40ms: 0.00%
40ms to 50ms: 0.00%
50ms to 60ms: 0.00%
60ms to 70ms: 0.00%
70ms to 80ms: 0.00%
80ms to 90ms: 0.00%
90ms to 100ms: 0.00%
100ms to 250ms: 0.00%
250ms to 400ms: 0.00%
400ms to 550ms: 0.00%
550ms to 700ms: 0.00%
700ms to 850ms: 0.00%
850ms to 1000ms: 0.00%
1000ms to 10000ms: 0.00%

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
  -quiet
        displays only a summary every 10 seconds
  -sleep int
        time between querys in μs (default 100000)
  -timeout int
        dns-timeout in ms (default 1000)
  -timeouts_only
        displays only timeouts or paketloss
  -type string
        dnstype for request (default "A")
exit status 2
```
