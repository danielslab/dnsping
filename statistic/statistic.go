package statistic

import (
	"fmt"
	"math"
	"strconv"
	"sync"
)

type Statistic struct {
	mutex            sync.RWMutex
	send_counter     int
	received_counter int
	rtt_slice        []int
	rcode_slice      map[string]int
}

func Init_Statistic() *Statistic {
	return &Statistic{
		send_counter:     0,
		received_counter: 0,
		rcode_slice:      make(map[string]int),
	}
}

func (statistic *Statistic) Increase_rcode_counter(rcode string) {
	statistic.mutex.Lock()
	defer statistic.mutex.Unlock()
	statistic.rcode_slice[rcode] = statistic.rcode_slice[rcode] + 1
}

func (statistic *Statistic) Increase_send_counter() {
	statistic.mutex.Lock()
	defer statistic.mutex.Unlock()
	statistic.send_counter += 1
}

func (statistic *Statistic) Increase_received_counter() {
	statistic.mutex.Lock()
	defer statistic.mutex.Unlock()
	statistic.received_counter += 1
}

func (statistic *Statistic) Append_rtt(rtt int) {
	statistic.mutex.Lock()
	defer statistic.mutex.Unlock()
	statistic.rtt_slice = append(statistic.rtt_slice, rtt)
}

func (statistic *Statistic) Print() {
	statistic.mutex.Lock()
	defer statistic.mutex.Unlock()

	loss := statistic.send_counter - statistic.received_counter
	sum := 0
	min_rtt := 0
	max_rtt := 0

	if len(statistic.rtt_slice) > 0 {
		min_rtt = statistic.rtt_slice[0]
		max_rtt = statistic.rtt_slice[0]
	}

	rcode := ""

	for i := 0; i < len(statistic.rtt_slice); i++ {
		if statistic.rtt_slice[i] < min_rtt {
			min_rtt = statistic.rtt_slice[i]
		}
		if statistic.rtt_slice[i] > max_rtt {
			max_rtt = statistic.rtt_slice[i]
		}
		// adding the values of array to the variable sum
		sum += (statistic.rtt_slice[i])
	}
	avg_rtt := (float64(sum)) / (float64(len(statistic.rtt_slice)))
	avg_rtt = math.Round(avg_rtt)

	fmt.Println("------------------------------")
	res := fmt.Sprintf("send: %d received: %d loss: %.2f%% min_rtt: %dms avg_rtt: %.2fms max_rtt: %dms jitter: %.2fms variance: %.2f", statistic.send_counter,
		statistic.received_counter, float64(loss)*100/float64(statistic.send_counter), min_rtt, avg_rtt, max_rtt, statistic.calculate_jitter(), statistic.calculate_variance())
	fmt.Println(res)

	for z, m := range statistic.rcode_slice {
		rcode = rcode + z + ":" + strconv.Itoa(m) + " "
	}
	fmt.Println(rcode)

}

func (statistic *Statistic) RTT_Summary() {

	fmt.Println("")
	fmt.Println("rtt distribution:")
	fmt.Println("------------------------------")

	for i := 0; i < 1000; {
		from := i
		if i >= 100 {
			i = i + 150
			statistic.calculate_rtt_stats(from, i)
		}
		if i < 100 && i >= 10 {
			i = i + 10
			statistic.calculate_rtt_stats(from, i)
		}
		if i < 10 {
			i = i + 5
			statistic.calculate_rtt_stats(from, i)
		}
		if i == 1000 {
			statistic.calculate_rtt_stats(i, 10000)
		}
	}
}

func (statistic *Statistic) calculate_rtt_stats(from int, to int) {
	statistic.mutex.Lock()
	defer statistic.mutex.Unlock()
	count := 0.0
	for i := 0; i < len(statistic.rtt_slice); i++ {
		if statistic.rtt_slice[i] >= from && statistic.rtt_slice[i] < to {
			count += 1
		}
	}

	value := count / float64(len(statistic.rtt_slice)) * 100
	//res := fmt.Sprintf("%dms to %dms: %.2f%% (count: %d)", from, to, value, int(count))
	res := fmt.Sprintf("%dms >= <= %dms: %.2f%% (count: %d)", from, to, value, int(count))
	fmt.Println(res)
}

func (statistic *Statistic) calculate_jitter() float64 {
	var latencyDiffs []float64
	var latencyDiff float64
	for i := 1; i < len(statistic.rtt_slice); i++ {
		//Calculate absolute diff
		latencyDiff = math.Abs(float64(statistic.rtt_slice[i] - statistic.rtt_slice[i-1]))
		latencyDiffs = append(latencyDiffs, latencyDiff)
	}

	var sumLatencyDiff float64
	//calculate sum of all elements in latencyDiffs
	for _, value := range latencyDiffs {
		sumLatencyDiff += value
	}

	avglatencyDiffs := sumLatencyDiff / float64(len(latencyDiffs))

	var jitter float64
	for _, value := range latencyDiffs {
		deviation := math.Abs(value - avglatencyDiffs)
		jitter += deviation
	}
	jitter = jitter / float64(len(latencyDiffs))

	return jitter

}

func (statistic *Statistic) calculate_variance() float64 {
	mean := 0.0
	for _, value := range statistic.rtt_slice {
		mean += float64(value)
	}
	mean /= float64(len(statistic.rtt_slice))

	variance := 0.0
	for _, value := range statistic.rtt_slice {
		variance += (float64(value) - mean) * (float64(value) - mean)
	}

	return variance / float64(len(statistic.rtt_slice))

}
