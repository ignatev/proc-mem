package main

import (
	"fmt"
	"gopkg.in/alexcesaro/statsd.v2"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type metric struct {
	name  string
	value int
}

func collect(hostname string) []metric {
	resp, err := http.Get(hostname)
	check(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	metrics := strings.Split(string(body), "\n")
	collected := []metric{}
	for _, m := range metrics[:len(metrics)-1] { // pop last element -- empty string
		nv := strings.Split(m, " ")
		n := nv[0]
		v, e := strconv.Atoi(nv[1])
		check(e)
		collected = append(collected, metric{name: n, value: v})
		send(metric{name: n, value: v})
	}
	return collected
}

func send(m metric) {
	c, err := statsd.New(statsd.Address("127.0.0.1:8125"))
	check(err)
	defer c.Close()
	c.Gauge(m.name, m.value)
	c.Flush()
}

func collectParallel(hostname string) {
	for range time.Tick(time.Second * 3) {
		fmt.Println("collect metrics from " + hostname)
		collect("http://" + hostname + ":8080/metrics")
	}
}
