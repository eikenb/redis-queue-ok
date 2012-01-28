/*
Insert some entries into redis to run tests with.
It currently just uses json blobs hardcoded below.
*/

package main

import (
	"tideland-rdc.googlecode.com/hg"
	"fmt"
	"flag"
)

func main() {
	blobs := []string{
		`{ "firstname":"John", "lastname":"Eikenberry" }`,
	}
	rd := rdc.NewRedisDatabase(rdc.Configuration{})

	for _, i := range argsAsQueues() {
		for _, j := range blobs {
			if !rd.Command("rpush", i, j).IsOK() {
				fmt.Printf("Push error; %v in %v\n", j, i)
			}
		}
	}
}

func argsAsQueues() []string {
	flag.Parse()
	queues := make([]string, 0, flag.NArg())
	for i := 0; i < flag.NArg(); i++ {
		queues = append(queues, "resque:queue:"+flag.Arg(i))
	}
	return queues
}
