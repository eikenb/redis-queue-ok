package main

import (
	"tideland-rdc.googlecode.com/hg"
	"fmt"
	"flag"
	"os"
	"bufio"
)

const dataPath = "/var/tmp/%v.txt"
const resque = "resque:queue:"

func main() {
	setUsage()
	queues := argsAsQueues()
	errs := make([]string, 0, len(queues))

	for _, q := range queues {
		disk_entry := fromDisk(q)
		// fmt.Println("fromDisk:", disk_entry)
		redis_entry := fromRedis(q)
		// fmt.Println("fromRedis:", redis_entry)
		if disk_entry != redis_entry {
			if err := toDisk(q, redis_entry); err != nil {
				fmt.Println("ERROR:", err)
				os.Exit(3)
			}
		} else if redis_entry != `` {
			errs = append(errs, q)
		}
	}
	if len(errs) > 0 {
		fmt.Println("Queue(s) not being processed:", errs)
		os.Exit(2)
	} else {
		fmt.Println("Queue(s) OK:", queues)
	}
}

// process arguments as queue names
func argsAsQueues() []string {
	flag.Parse()
	queues := make([]string, 0, flag.NArg())
	for i := 0; i < flag.NArg(); i++ {
		queues = append(queues, resque+flag.Arg(i))
	}
	return queues
}

// --help output
func setUsage() {
	flag.Usage = func() {
		fmt.Println("Usage:", os.Args[0], "QUEUE [...]")
	}
}

// return top json blob from redis queue
func fromRedis(q string) (s string) {
	rd := rdc.NewRedisDatabase(rdc.Configuration{})
	switch l := rd.Command("llen", q).ValueAsInt(); {
	case l > 0:
		rsb := rd.Command("lindex", q, 0)
		s = rsb.Value().String()
	}
	return
}

// return stored queue json blob from disk
func fromDisk(q string) (str string) {
	path := fmt.Sprintf(dataPath, q)
	linkCheck(path)
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		reader := bufio.NewReader(file)
		str, err = reader.ReadString('\n')
		if err != nil && err != os.EOF {
			fmt.Println("ERROR: fromDisk;", err)
			os.Exit(3)
		}
	}
	return
}

// store json blob from top of queue to disk for next run
func toDisk(q string, j string) (err os.Error) {
	path := fmt.Sprintf(dataPath, q)
	linkCheck(path)
	file, err := os.Create(path)
	if err == nil {
		defer file.Close()
		file.WriteString(j)
	}
	return
}

// Exits with error if path is sym-link
func linkCheck(path string) {
	finfo, err := os.Lstat(path)
	if err == nil && finfo.IsSymlink() {
		fmt.Println("ERROR: Data file is sym-link;", path)
		os.Exit(3)
	}
}
