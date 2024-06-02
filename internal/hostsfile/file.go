package hostsfile

import (
	"bufio"
	"fmt"
	"golang.org/x/exp/maps"
	"io"
	"os"
	"sort"
)

type records map[string]Record

type Record struct {
	Type    string
	Ip      string
	DNSName string
}

type HostsFile struct {
	Path    string
	records records
}

func (h *HostsFile) Get() ([]Record, error) {
	if h.records == nil {
		err := h.read()
		if err != nil {
			return nil, err
		}
	}

	res := maps.Values(h.records)
	sort.Slice(res, func(i, j int) bool {
		keyA := res[i].key()
		keyB := res[j].key()
		return keyA < keyB
	})
	return res, nil
}

func (h *HostsFile) Set(typ, dnsName, target string) {
	r := Record{
		Type:    typ,
		Ip:      target,
		DNSName: dnsName,
	}
	h.records[r.key()] = r
}

func (h *HostsFile) Remove(typ, dnsName string) {
	r := Record{
		Type:    typ,
		DNSName: dnsName,
	}
	delete(h.records, r.key())
}

func (r Record) key() string {
	return fmt.Sprintf("%:s%s", r.Type, r.DNSName)
}

func (h *HostsFile) Flush() error {
	f, err := os.Create(h.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	rec, err := h.Get()
	if err != nil {
		return err
	}

	for _, record := range rec {
		switch record.Type {
		case "A":
			_, err = io.WriteString(f, fmt.Sprintf("%s %s\n", record.Ip, record.DNSName))
			if err != nil {
				return err
			}
		case "TXT":
			_, err = io.WriteString(f, fmt.Sprintf("# TXT %s %s\n", record.DNSName, record.Ip))
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (h *HostsFile) read() error {
	h.records = make(map[string]Record)

	f, err := os.Open(h.Path)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		} else if len(line) >= 5 && line[:5] == "# TXT" {
			var dnsName, txt string
			_, err := fmt.Sscanf(line, "# TXT %s %s", &dnsName, &txt)
			if err != nil {
				return err
			}

			h.Set("TXT", dnsName, txt)
		} else if line[0] == '#' {
			continue
		} else {
			var ip, host string
			_, err := fmt.Sscanf(line, "%s %s", &ip, &host)
			if err != nil {
				return err
			}

			h.Set("A", host, ip)
		}

	}

	return nil
}