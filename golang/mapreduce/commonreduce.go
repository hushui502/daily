package mapreduce

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

// doReduce does the job of a reduce worker: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.

// jobName: the name of the whole MapReduce job
// reduceTaskNumber: which reduce task this is
// nMap: the number of map tasks that were run ("M" in the paper)
func doReduce(jobName string, reduceTaskNumber int, nMap int, reduceF func(key string, values []string) string) {

	// step1: read map generator file ,same key merge put map[string][]string
	kvs := make(map[string][]string)

	for i := 0; i < nMap; i++ {
		fileName := reduceName(jobName, i, reduceTaskNumber)
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal("doReduce1: ", err)
		}

		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			err = dec.Decode(&kv)
			if err != nil {
				break
			}

			_, ok := kvs[kv.Key]
			if !ok {
				kvs[kv.Key] = []string{}
			}
			kvs[kv.Key] = append(kvs[kv.Key], kv.Value)
		}
		file.Close()
	}

	var keys []string
	for k := range kvs {
		keys = append(keys, k)
	}

	// step2: sort by keys
	sort.Strings(keys)

	// step3: create result file
	p := mergeName(jobName, reduceTaskNumber)
	file, err := os.Create(p)
	if err != nil {
		log.Fatal("doReduce2: ceate ", err)
	}
	enc := json.NewEncoder(file)

	// step4: call user reduce each key of kvs
	for _, k := range keys {
		res := reduceF(k, kvs[k])
		enc.Encode(KeyValue{k, res})
	}

	file.Close()
}
