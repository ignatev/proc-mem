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

// // SetProcRoot sets the location of the proc filesystem.
// func SetProcRoot(root string) {
// 	procRoot = root
// }

// // walkProcPid walks over all numerical (PID) /proc entries, and sees if their
// // ./fd/* files are symlink to sockets. Returns a map from socket ID (inode)
// // to PID. Will return an error if /proc isn't there.
// func walkProcPid(namespaces *map[uint64]struct{}, connBuff *bytes.Buffer) (map[uint64]Proc, error) {
// 	fh, err := os.Open(procRoot)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dirNames, err := fh.Readdirnames(-1)
// 	fh.Close()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var (
// 		res  = map[uint64]Proc{}
// 		stat syscall.Stat_t
// 	)
// 	for _, dirName := range dirNames {
// 		pid, err := strconv.ParseUint(dirName, 10, 0)
// 		if err != nil {
// 			// Not a number, so not a PID subdir.
// 			continue
// 		}

// 		fdBase := procRoot + "/" + dirName + "/fd/"
// 		dfh, err := os.Open(fdBase)
// 		if err != nil {
// 			// Process is be gone by now, or we don't have access.
// 			continue
// 		}

// 		fdNames, err := dfh.Readdirnames(-1)
// 		dfh.Close()
// 		if err != nil {
// 			continue
// 		}

// 		// Read network namespace, and if we haven't seen it before,
// 		// read /proc/<pid>/net/tcp
// 		err = syscall.Lstat(fmt.Sprintf("%s/%d/ns/net", procRoot, pid), &stat)
// 		if err != nil {
// 			continue
// 		}

// 		if _, ok := (*namespaces)[stat.Ino]; !ok {
// 			(*namespaces)[stat.Ino] = struct{}{}
// 			readFile(fmt.Sprintf("%s/%d/net/tcp", procRoot, pid), connBuff)
// 			readFile(fmt.Sprintf("%s/%d/net/tcp6", procRoot, pid), connBuff)
// 		}

// 		var name string
// 		for _, fdName := range fdNames {
// 			// Direct use of syscall.Stat() to save garbage.
// 			err = syscall.Stat(fdBase+fdName, &stat)
// 			if err != nil {
// 				continue
// 			}

// 			// We want sockets only.
// 			if stat.Mode&syscall.S_IFMT != syscall.S_IFSOCK {
// 				continue
// 			}

// 			if name == "" {
// 				if name = procName(procRoot + "/" + dirName); name == "" {
// 					// Process might be gone by now
// 					break
// 				}
// 			}

// 			res[stat.Ino] = Proc{
// 				PID:  uint(pid),
// 				Name: name,
// 			}
// 		}
// 	}

// 	return res, nil
// }

// // procName does a pid->name lookup.
// func procName(base string) string {
// 	fh, err := os.Open(base + "/comm")
// 	if err != nil {
// 		return ""
// 	}

// 	name := make([]byte, 64)
// 	l, err := fh.Read(name)
// 	fh.Close()
// 	if err != nil {
// 		return ""
// 	}

// 	if l < 2 {
// 		return ""
// 	}

// 	// drop trailing "\n"
// 	return string(name[:l-1])
// }

// // readFile reads an arbitrary file into a buffer. It's a variable so it can
// // be overwritten for benchmarks. That's bad practice and we should change it
// // to be a dependency.
// var readFile = func(filename string, buf *bytes.Buffer) error {
// 	f, err := os.Open(filename)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = buf.ReadFrom(f)
// 	f.Close()
// 	return err
// }
