package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/calmh/zfs"
	"github.com/robfig/cron"
)

var (
	cfgFile = "/opt/local/zsnapper/etc/zsnapper.yml"
	verbose = false
)

func main() {
	flag.StringVar(&cfgFile, "c", cfgFile, "Path to configuration file")
	flag.BoolVar(&verbose, "v", verbose, "Enables verbose output")
	flag.Parse()

	fd, err := os.Open(cfgFile)
	if err != nil {
		fatalln(err)
	}

	jobs, err := LoadJobs(fd)
	fd.Close()
	if err != nil {
		fatalln(err)
	}

	if len(jobs) == 0 {
		fatalln("No jobs to run")
	}

	j := cron.New()
	for _, job := range jobs {
		j.AddJob(job.Schedule, job)
		if verbose {
			fmt.Println("Added", job)
		}
	}
	j.Start()

	select {}
}

func timestamp() string {
	return time.Now().UTC().Format("20060102T150405Z")
}

func fatalln(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
}

type datasetList []*zfs.Dataset

func (l datasetList) Len() int {
	return len(l)
}
func (l datasetList) Swap(a, b int) {
	l[a], l[b] = l[b], l[a]
}
func (l datasetList) Less(a, b int) bool {
	return l[a].Name < l[b].Name
}
