package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
)

var (
	ErrExists    = errors.New("host already in the list")
	ErrNotExists = errors.New("host not in the list")
)

type HostList struct {
	Hosts []string
}

func (hl *HostList) search(host string) (bool, int) {
	sort.Strings(hl.Hosts)

	i := sort.SearchStrings(hl.Hosts, host)
	if i < len(hl.Hosts) {
		return true, i
	}
	return false, -1
}

func (hl *HostList) Add(host string) error {
	if found, _ := hl.search(host); found {
		return ErrExists
	}
	hl.Hosts = append(hl.Hosts, host)
	return nil
}

func (hl *HostList) Remove(host string) error {
	ok, i := hl.search(host)
	if !ok {
		return ErrNotExists
	}
	hl.Hosts = append(hl.Hosts[:i], hl.Hosts[i+1:]...)
	return nil
}

// Load obtains hosts from a hosts file
func (hl *HostList) Load(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		if errors.Is(err, ErrNotExists) {
			return nil
		}
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hl.Hosts = append(hl.Hosts, scanner.Text())
	}
	return nil
}

// Save saves hosts to a hosts file
func (hl *HostList) Save(fileName string) error {
	output := ""
	for _, h := range hl.Hosts {
		output += fmt.Sprintln(h)
	}
	return os.WriteFile(fileName, []byte(output), 0644)
}
