package main

import (
	rdc "code.google.com/p/tcgl/redis"
	"code.google.com/p/tcgl/applog"
	"fmt"
	"flag"
	"os"
	"io"
	"bufio"
	"text/template"
	"net/smtp"
	"strings"
	"path/filepath"
)

var server, from, to, pattern string
var enable_email, leftpush, rightpush bool
const dataPath = "/var/tmp/%v.txt"
const dataDir = "/var/tmp/"

func main() {
	setup()
	queues := getQueues()
	errs := make([]string, 0, len(queues))
	errs = append(errs, "Queue(s) not being processed:")

	for _, q := range queues {
		disk_entry := fromDisk(q)
		// fmt.Println("fromDisk:", disk_entry)
		redis_entry := fromRedis(q)
		// fmt.Println("fromRedis:", redis_entry)
		if disk_entry != redis_entry {
			if err := toDisk(q, redis_entry); err != nil {
				exit(3, "ERROR:", err.Error())
			}
		} else if redis_entry != `` {
			errs = append(errs, q)
		}
	}
	if len(errs) > 1 {
		exit(2, errs...)
	} else {
		exit(0)
	}
}

func getQueues() (qs []string) {
	set := map[string]bool{}
	rd := rdc.NewRedisDatabase(rdc.Configuration{})
	// resque* to optionally support custom namespaces
	for _, v := range rd.Command("keys", pattern).Values() {
		set[v.String()]=true
	}
	// this will pull old queues that have files left
	filepath.Walk(dataDir, func(p string, fi os.FileInfo, e error) error {
		if ! fi.IsDir() && fi.Mode() & os.ModeSymlink != os.ModeSymlink {
			p = filepath.Base(p)
			if match, _ := filepath.Match(pattern, p); match {
				e := filepath.Ext(p)
				p = strings.Replace(p, e, "", 1) // remove extension
				// fmt.Println(p, "match")
				set[p]=true
			}
		}
		return e
	})
	for k, _ := range set {
		qs = append(qs, k)
	}
	return
}

func setup() {
	// email flags
	flag.BoolVar(&enable_email, "e", false, "enable email message")
	flag.StringVar(&server, "s", "localhost:25", "smtp server")
	flag.StringVar(&from, "f", "", "From: address")
	flag.StringVar(&to, "t", "", "To: address")
	// queue flags
	flag.StringVar(&pattern, "p", "resque*:queue:*", "queue key pattern")
	flag.BoolVar(&leftpush, "l", false, "left push (test index -1)")
	flag.BoolVar(&rightpush, "r", true, "right push (test index 0)")

	flag.Usage = func() {
		fmt.Println("Usage: [options]", os.Args[0])
		fmt.Println("Options (for optional email message):")
		flag.PrintDefaults()
		fmt.Println("Returns:")
		fmt.Println("\t0 on success")
		fmt.Println("\t1 not used")
		fmt.Println("\t2 when queue is not being processed")
		fmt.Println("\t3 when there is an error with the check")
	}
	flag.Parse()
	if rightpush { leftpush = false }
	if leftpush { rightpush = false }

	// redis library can be chatty, shut it up
	devnull, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	applog.SetLogger(applog.NewStandardLogger(devnull))
}

// return top json blob from redis queue
func fromRedis(q string) (s string) {
	rd := rdc.NewRedisDatabase(rdc.Configuration{})
	switch rightpush {
	case true:  // rightpush
		s = rd.Command("lindex", q, 0).Value().String()
	case false: // leftpush
		s = rd.Command("lindex", q, -1).Value().String()
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
		if err != nil && err != io.EOF {
			exit(3, "ERROR: fromDisk;", err.Error())
		}
	}
	return
}

// store json blob from top of queue to disk for next run
func toDisk(q string, data string) (err error) {
	path := fmt.Sprintf(dataPath, q)
	linkCheck(path)
	if data == "" && exists(path) {
		os.Remove(path)
	} else {
		file, err := os.Create(path)
		if err == nil {
			defer file.Close()
			file.WriteString(data)
		}
	}
	return
}

// Exits with error if path is sym-link
func linkCheck(path string) {
	finfo, err := os.Lstat(path)
	if err == nil && finfo.Mode() & os.ModeSymlink == os.ModeSymlink {
		exit(3, "ERROR: Data file is sym-link;", path)
	}
}

func exists(path string) bool {
    _, err := os.Stat(path)
    if err == nil { return true }
    if os.IsNotExist(err) { return false }
    return false
}

// exit/sendmail
func exit(n int, msgs ...string) {
	output := strings.Join(msgs, " ")
	switch {
	case enable_email && n > 0:
		err := sendmail(output)
		if err != nil {
			fmt.Println(err)
		}
	case ! enable_email && len(output) > 0:
		fmt.Println(output)
	}
	os.Exit(n)
}

// email
var message string = `From: {{.From}}
To: {{.To}}
Subject: ({{.Host}}) {{.Msg}}

Hostname: {{.Host}}
Error: {{.Msg}}

To verify run queue-ok on ${{.Host}}.

	ssh {{.Host} queue-ok
`
type md struct {
	From, To, Host, Msg string
}

func sendmail(msg string) (err error) {
	c, err := smtp.Dial(server)
	if err != nil { return }
	c.Mail(from)
	c.Rcpt(to)
	wc, err := c.Data()
	if err != nil { return }
	defer wc.Close()
	mail, err := template.New("mail").Parse(message)
	if err != nil { return }
	hostname, _ := os.Hostname()
	err = mail.Execute(wc, &md{From: from, To: to, Host: hostname, Msg: msg})
	return
}

