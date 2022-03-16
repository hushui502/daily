package hostfile

import (
	"io/ioutil"
	"strings"
)

func ReadHostsFile() ([]byte, error) {
	bs, err := ioutil.ReadFile(HostsPath)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func ParseHosts(hostsFileContent []byte, err error) (map[string][]string, error) {
	if err != nil {
		return nil, err
	}

	hostsMap := map[string][]string{}

LINE:
	for _, line := range strings.Split(strings.Trim(string(hostsFileContent), " \t\r\n"), "\n") {
		line = strings.Replace(strings.Trim(line, " \t"), "\t", " ", -1)
		if len(line) == 0 || line[0] == ';' || line[0] == '#' {
			continue
		}
		pieces := strings.SplitN(line, " ", 2)
		if len(pieces) > 1 && len(pieces[0]) > 0 {
			if names := strings.Fields(pieces[1]); len(names) > 0 {
				for _, name := range names {
					if strings.HasPrefix(name, "#") {
						continue LINE
					}
					hostsMap[pieces[0]] = append(hostsMap[pieces[0]], name)
				}
			}
		}
	}

	return hostsMap, nil
}

// ReverseLookup takes an IP address and returns a slice of matching hosts file
// entries.
func ReverseLookup(ip string) ([]string, error) {
	hostsMap, err := ParseHosts(ReadHostsFile())
	if err != nil {
		return nil, err
	}
	return hostsMap[ip], nil
}
