package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	fmt.Println("collect memory per process")
	procPid()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type proc struct {
	pid, mem, name string
}

func parseMem(line, pattern string) string {
	result := ""
	pattern = "^" + pattern + "\\s+(.+)"
	re := regexp.MustCompile(pattern)
	res := re.FindStringSubmatch(line)
	if len(res) != 0 {
		result = re.FindStringSubmatch(line)[1]
	}
	return result
}

func procVmRSS(pid string) proc {
	res := proc{}

	procstatus := "/proc/" + pid + "/status"
	f, err := os.Open(procstatus)
	check(err)
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		name := parseMem(sc.Text(), "Name:")
		if len(name) != 0 {
			res.name = name
		}
		mem := parseMem(sc.Text(), "VmRSS:")
		if len(mem) != 0 {
			res.mem = mem
		}
	}

	res.pid = pid
	return res
}

func procPid() map[uint64]proc {

	res := make(map[uint64]proc)
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
		uintpid, err := strconv.ParseUint(pid, 10, 0)
		if err != nil {
			continue
		} else {
			process := procVmRSS(pid)
			if len(process.mem) != 0 {
				res[uintpid] = process
			}
		}
	}
	fmt.Println(res)
	return res
}
