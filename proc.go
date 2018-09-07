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
//	procPid()
	meminfo()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type proc struct {
	pid, mem, name string
}

func parseMem(line, pattern string) string {	//
	result := ""								//	extracting only name and memory size
	pattern = "^" + pattern + "\\s+(\\S+)"		//	from strings:
	re := regexp.MustCompile(pattern)			//
	res := re.FindStringSubmatch(line)			//	`Name: 	Telegram` -> Telegram
	if len(res) != 0 {							//	`VmRSS:	18888234 kB` -> 18888234
		result = re.FindStringSubmatch(line)[1] //
	}
	return result
}

func vmrss(pid string) proc {
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

type total struct {
	total, free string
	percentage float64
}

func meminfo() total {
	result := total{}

	meminfo := "/proc/meminfo"
	f, err := os.Open(meminfo)
	check(err)
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		total := parseMem(sc.Text(), "MemTotal:")
		if len(total) != 0 {
			result.total = total
		}
		free := parseMem(sc.Text(), "MemAvailable:")
		if len(free) != 0 {
			result.free = free
		}
	}
	totalf, err := strconv.ParseFloat(result.total, 32)
	check(err)
	freef, err := strconv.ParseFloat(result.free, 32)
	check(err)

	result.percentage = freef / totalf * 100

	fmt.Println(result)
	return result
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
			process := vmrss(pid)
			if len(process.mem) != 0 {
				res[uintpid] = process
				fmt.Println(process)
			}
		}
	}
	return res
}
