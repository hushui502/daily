package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"sync"
)

type File struct {
	Name  string
	Type  string
	Mtime string
}

type RebuildRequest struct {
	vol string
	url string
}

func get_files(url string) []File {
	var files []File
	dat, err := remote_get(url)
	if err != nil {
		fmt.Println("remote get file failed ", err)
		return files
	}
	err = json.Unmarshal([]byte(dat), &files)
	if err != nil {
		fmt.Println("unamrshal dat to files failed ", err)
	}

	return files
}

func rebuild(a *App, vol string, name string) bool {
	key, err := base64.StdEncoding.DecodeString(name)
	if err != nil {
		fmt.Println("base64 decode error", err)
		return false
	}

	kvolumes := key2volume(key, a.volumes, a.replaces, a.subvolumes)

	if !a.LockKey(key) {
		fmt.Println("lockKey issue", key)
		return false
	}
	defer a.UnlockKey(key)

	data, err := a.db.Get(key, nil)
	var rec Record
	if err != leveldb.ErrNotFound {
		rec = toRecord(data)
		rec.rvolumes = append(rec.rvolumes, vol)
	} else {
		rec = Record{[]string{vol}, NO, ""}
	}

	// sort by order in kvolumes
	pvalues := make([]string, 0)
	for _, v1 := range kvolumes {
		for _, v2 := range rec.rvolumes {
			if v1 == v2 {
				pvalues = append(pvalues, v1)
			}
		}
	}
	// insert the ones that are not there at the end
	for _, v2 := range rec.rvolumes {
		insert := true
		for _, v := range kvolumes {
			if v == v2 {
				insert = false
				break
			}
		}
		if insert {
			pvalues = append(pvalues, v2)
		}
	}

	if !a.PutRecord(key, Record{pvalues, NO, ""}) {
		fmt.Println("put error ", err)
		return false
	}

	fmt.Println(string(key), pvalues)

	return true
}

func valid(f File) bool {
	if len(f.Name) != 2 || f.Type != "directory" {
		return false
	}
	decoded, err := hex.DecodeString(f.Name)
	if err != nil {
		return false
	}
	if len(decoded) != 1 {
		return false
	}

	return true
}

func (a *App) Rebuild() {
	fmt.Println("rebuilding on ", a.volumes)

	// empty db
	iter := a.db.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		a.db.Delete(iter.Key(), nil)
	}

	var wg sync.WaitGroup
	reqs := make(chan RebuildRequest, 20000)

	for i := 0; i < 128; i++ {
		go func() {
			for req := range reqs {
				files := get_files(req.url)
				for _, f := range files {
					rebuild(a, req.vol, f.Name)
				}
				wg.Done()
			}
		}()
	}

	paras_volume := func(tvol string) {
		for _, i := range get_files(fmt.Sprintf("http://%s/", tvol)) {
			if valid(i) {
				for _, j := range get_files(fmt.Sprintf("http://%s/%s/", tvol, i.Name)) {
					if valid(j) {
						wg.Add(1)
						url := fmt.Sprintf("http://%s/%s/%s", tvol, i.Name, j.Name)
						reqs <- RebuildRequest{tvol, url}
					}
				}
			}
		}
	}

	for _, vol := range a.volumes {
		has_subvolumes := false
		for _, f := range get_files(fmt.Sprintf("http://%s/", vol)) {
			paras_volume(fmt.Sprintf("%s/%s", vol, f.Name))
			has_subvolumes = true
		}
		if !has_subvolumes {
			paras_volume(vol)
		}
	}

	close(reqs)
	wg.Wait()
}
