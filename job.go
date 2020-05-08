package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/calmh/zsnapper/zfs"
	"gopkg.in/yaml.v2"
)

type Job struct {
	Family    string   // "test"
	Datasets  []string // "data/foto"
	Schedule  string   // "0 */5 * * * *"
	Keep      int      // 12
	Recursive bool
}

func LoadJobs(r io.Reader) ([]Job, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var configs []Job
	if err := yaml.Unmarshal(bs, &configs); err != nil {
		return nil, err
	}

	return configs, nil
}

func (j Job) Run() {
	name := j.Family + "-" + timestamp()

	expandedSets := expandLines(j.Datasets)
	all := datasets()
	filteredDatasets := selectLines(expandedSets, all)

	for _, dataset := range filteredDatasets {
		if err := j.snapshot(dataset, name); err != nil {
			fmt.Println("Create snapshot:", err)
		}
		if err := j.clean(dataset); err != nil {
			fmt.Println("Clean snapshots:", err)
		}
	}
}

func (j Job) String() string {
	return fmt.Sprintf(`%v@%s, at "%s" (keep %d, recursive %v)`, j.Datasets, j.Family, j.Schedule, j.Keep, j.Recursive)
}

func (j Job) snapshot(dataset, name string) error {
	mutexFor(dataset).Lock()
	defer mutexFor(dataset).Unlock()

	ds, err := zfs.GetDataset(dataset)
	if err != nil {
		return err
	}

	ss, err := ds.Snapshot(name, j.Recursive)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Println("Created", ss.Name)
	}

	return nil
}

func (j Job) clean(dataset string) error {
	mutexFor(dataset).Lock()
	defer mutexFor(dataset).Unlock()

	ds, err := zfs.GetDataset(dataset)
	if err != nil {
		return err
	}

	snaps, err := ds.Snapshots()
	if err != nil {
		return err
	}

	prefix := dataset + "@" + j.Family + "-"
	var matching []*zfs.Dataset
	for _, snap := range snaps {
		if strings.HasPrefix(snap.Name, prefix) {
			matching = append(matching, snap)
		}
	}

	sort.Sort(datasetList(matching))

	if len(matching) <= j.Keep {
		return nil
	}

	var flags zfs.DestroyFlag
	if j.Recursive {
		flags |= zfs.DestroyRecursive
	}
	for _, snap := range matching[:len(matching)-j.Keep] {
		if err := snap.Destroy(flags); err != nil {
			return err
		}
		if verbose {
			fmt.Println("Destroyed", snap.Name)
		}
	}

	return nil
}
