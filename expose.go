package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var hostname, _ = os.Hostname()

func main() {
	http.HandleFunc("/metrics", handler)
	http.HandleFunc("/statsd", getStats)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	perproclines := pidlines(perproc()) + totalmem(meminfo())
	io.WriteString(w, perproclines)
}

func getStats(w http.ResponseWriter, r *http.Request) {
	collect("http://localhost:8080/metrics")
}

func pidlines(pids map[uint64]proc) string {
	result := ""
	for _, v := range pids {
		line := "process_memory{pid=\"" + v.pid + "\",proc_name=\"" + v.name + "\",hostname=\"" + hostname + "\"} " + v.mem
		result = result + line + "\n"
	}
	return result
}

func totalmem(mem total) string {
	totalMem := "total_memory{hostname=\"" + hostname + "\"} " + mem.total + "\n"
	freeMem := "free_memory{hostname=\"" + hostname + "\"} " + mem.free + "\n"
	perc := "free_percentage{hostname=\"" + hostname + "\"} " + fmt.Sprintf("%.0f", mem.percentage) + "\n"
	return totalMem + freeMem + perc
}
