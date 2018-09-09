package main

import (
	"fmt"
	"gopkg.in/alexcesaro/statsd.v2"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"os"
	"log"
	"bufio"
)

type metric struct {
	name  string
	value int
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	hosts := inventory("hosts")
	for _, host := range hosts {
		go collectParallel(host)
	}
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func inventory(path string) []string {
	var res []string
	if _, err := os.Stat(path); err == nil {
		f, err := os.Open(path)
		check(err)
		defer f.Close()

		sc := bufio.NewScanner(f)
		for sc.Scan() {
			res = append(res, sc.Text())
		}
	}

	return res
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
		fmt.Println("scraping metrics from " + hostname)
		collect("http://" + hostname + ":8080/metrics")
	}
}
