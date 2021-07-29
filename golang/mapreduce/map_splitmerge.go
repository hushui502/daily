package mapreduce

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

func (m *Master) merge() {
	debug("Merge phase")
	kvs := make(map[string]string)

	for i := 0; i < m.nReduce; i++ {
		p := mergeName(m.jobName, i)
		file, err := os.Open(p)
		if err != nil {
			log.Fatal("Merge: ", err)
		}
		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			err = dec.Decode(&kv)
			if err != nil {
				break
			}
			kvs[kv.Key] = kv.Value
		}
		file.Close()
	}

	var keys []string
	for k := range kvs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	file, err := os.Create("mrtmp." + m.jobName)
	if err != nil {
		log.Fatal("Merge: create ", err)
	}
	w := bufio.NewWriter(file)
	for _, k := range keys {
		fmt.Fprintf(w, "%s: %s\n", k, kvs[k])
	}

	w.Flush()
	file.Close()
}

// removeFile is a simple wrapper around os.Remove that logs errors.
func removeFile(n string) {
	err := os.Remove(n)
	if err != nil {
		log.Fatal("CleanupFiles ", err)
	}
}

// CleanupFiles removes all intermediate files produced by running mapreduce.
func (m *Master) CleanupFiles() {
	for i := range m.files {
		for j := 0; j < m.nReduce; j++ {
			removeFile(reduceName(m.jobName, i, j))
		}
	}
	for i := 0; i < m.nReduce; i++ {
		removeFile(mergeName(m.jobName, i))
	}
	removeFile("mrtmp." + m.jobName)
}