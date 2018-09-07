package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/metrics", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	perproclines := pidlines(perproc()) + totalmem(meminfo())
	io.WriteString(w, perproclines)
}

func pidlines(pids map[uint64]proc) string {
	result := ""
	for _, v := range pids {
		line := "process_memory{pid=\"" + v.pid + "\",proc_name=\"" + v.name + "\"} " + v.mem
		result = result + line + "\n"
	}
	return result
}

func totalmem(mem total) string {
	totalMem := "total_memory " + mem.total + "\n"
	freeMem := "free_memory " + mem.free + "\n"
	perc := "free_percentage " + fmt.Sprintf("%.2f", mem.percentage) + "\n"
	return totalMem + freeMem + perc
}
