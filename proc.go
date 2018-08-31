package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

func main() {
	fmt.Println("collect memory per process")
	procPid()

}

type proc struct {
	pid  uint
	mem  uint64
	name string
}

func procPid() map[uint]proc {

	res := make(map[uint]proc)
	pidDir, err := os.Open("/proc")
	if err != nil {
		fmt.Println("cannot open proc dir")
	}

	pids, err := pidDir.Readdirnames(-1)
	pidDir.Close()
	if err != nil {
		fmt.Println("error reading pids")
	}

	for _, pid := range pids {
		_, err := strconv.ParseUint(pid, 10, 0)
		if err != nil {
			continue
		} else {

			procname(pid)
		}
	}
	procmem(pids[333])

	return res

}

func parseShared(lines string) {
	re, err := regexp.Compile("^Shared.+")
	if err != nil {
		log.Println(err)
	}

	matches := re.FindAllString(lines, -1)
	//fmt.Print(lines)
	fmt.Println(matches)
}

func scanFile(filename string) error {
	begin := time.Now()

	total := 0 // count lines

	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		parseShared(sc.Text())
		total++
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return err
	}

	log.Printf("scan file, time_used: %v, lines=%v\n", time.Since(begin).Seconds(), total)

	return nil
}

func procname(pid string) {
	filename := "/proc/" + pid + "/comm"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	contents := buf.String()

	fmt.Print(pid, " ", contents)
}

func procmem(pid string) {
	filename := "/proc/" + pid + "/smaps"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	//	contents := buf.String()
	scanFile(filename)
	//fmt.Print(contents)
}