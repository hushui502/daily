package main

//
//import (
//	"encoding/json"
//	"fmt"
//	"time"
//)
//
//func main() {
//	nf, err := New(
//		WithCaptureTimeout(5 * time.Second),
//	)
//	if err != nil {
//		panic(err)
//	}
//
//	err = nf.Start()
//	if err != nil {
//		panic(err)
//	}
//	defer nf.Stop()
//
//	<-nf.Done()
//
//	var (
//		limit     = 5
//		recentSec = 3
//	)
//
//	rank, err := nf.GetProcessRank(limit, recentSec)
//	if err != nil {
//		panic(err)
//	}
//
//	bs, err := json.MarshalIndent(rank, "", "    ")
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(string(bs))
//}
//
