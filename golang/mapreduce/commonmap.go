package mapreduce

import (
	"encoding/json"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
)

// doMap does the job of a map worker: it reads one of the input files
// (inFile), calls the user-defined map function (mapF) for that file's
// contents, and partitions the output into nReduce intermediate files.

// jobName: the name of the MapReduce job
// mapTaskNumber: which map task this is
// nReduce: the number of reduce task that will be run ("R" in the paper)
func doMap(jobName string, mapTaskNumber int, inFile string, nReduce int, mapF func(file string, contents string) []KeyValue) {
	// step1: read file
	contents, err := ioutil.ReadFile(inFile)
	if err != nil {
		log.Fatal("do map error for inFile ", err)
	}

	// step2: call user user-map method ,to get kv
	kvResult := mapF(inFile, string(contents))

	// ste3: use key of kv generator nReduce file ,partition
	tmpFiles := make([]*os.File, nReduce)
	encoders := make([]*json.Encoder, nReduce)

	for i := 0; i < nReduce; i++ {
		tmpFileName := reduceName(jobName, mapTaskNumber, i)
		tmpFiles[i], err = os.Create(tmpFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer tmpFiles[i].Close()

		encoders[i] = json.NewEncoder(tmpFiles[i])
	}

	for _, kv := range kvResult {
		hashKey := int(ihash(kv.Key)) % nReduce
		err := encoders[hashKey].Encode(&kv)
		if err != nil {
			log.Fatal("do map encoders ", err)
		}
	}
}

func ihash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))

	return h.Sum32()
}
