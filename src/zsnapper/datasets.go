package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/calmh/zfs"
)

var (
	datasetsMutex  sync.Mutex
	datasetsLoaded time.Time
	datasetsCache  []string
)

func datasets() []string {
	datasetsMutex.Lock()
	defer datasetsMutex.Unlock()

	if time.Since(datasetsLoaded) > time.Minute {
		datasetsCache = make([]string, 0, len(datasetsCache))

		dss, err := zfs.Filesystems("", 0)
		if err == nil {
			for _, ds := range dss {
				datasetsCache = append(datasetsCache, ds.Name)
			}
		} else if verbose {
			fmt.Println("Filesystems:", err)
		}

		dss, err = zfs.Volumes("", 0)
		if err == nil {
			for _, ds := range dss {
				datasetsCache = append(datasetsCache, ds.Name)
			}
		} else if verbose {
			fmt.Println("Volumes:", err)
		}
	}

	return datasetsCache
}

var (
	datasetMutexes      = make(map[string]*sync.Mutex)
	datasetMutexesMutex sync.Mutex
)

func mutexFor(dataset string) *sync.Mutex {
	datasetMutexesMutex.Lock()
	defer datasetMutexesMutex.Unlock()
	if m, ok := datasetMutexes[dataset]; ok {
		return m
	}
	m := new(sync.Mutex)
	datasetMutexes[dataset] = m
	return m
}
