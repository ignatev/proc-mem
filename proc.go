package main

import (
	 "os"
	 "strconv"
	// "syscall"
	"fmt"
	"log"
	"regexp"
	"bufio"
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
		result = re.FindStringSubmatch(line)[0]
		fmt.Println(result)
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
	fmt.Println(res)
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
			res[uintpid] = procVmRSS(pid)
		}
	}

	return res

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