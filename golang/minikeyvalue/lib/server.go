package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

/*	Master Server */

type ListResponse struct {
	Next string   `json:"next"`
	Keys []string `json:"keys"`
}

func (a *App) QueryHandler(key []byte, w http.ResponseWriter, r *http.Request) {
	// operation is first query paramter
	// e.g. ?list&limit=10
	operation := strings.Split(r.URL.Path, "&")[0]
	switch operation {
	case "list", "unlinked":
		start := r.URL.Query().Get("start")
		limit := 0
		qlimit := r.URL.Query().Get("limit")
		if qlimit != "" {
			nlimit, err := strconv.Atoi(qlimit)
			if err != nil {
				w.WriteHeader(400)
				return
			}
			limit = nlimit
		}

		slice := util.BytesPrefix(key)
		if start != "" {
			slice.Start = []byte(start)
		}
		iter := a.db.NewIterator(slice, nil)
		defer iter.Release()

		keys := make([]string, 0)
		next := ""
		for iter.Next() {
			rec := toRecord(iter.Value())
			if (rec.deleted != NO && operation == "list") ||
				(rec.deleted != SOFT && operation == "unlinked") {
				continue
			}
			// too large, need to specify limit
			if len(keys) > 1000000 {
				w.WriteHeader(413)
				return
			}
			if limit > 0 && len(keys) == limit {
				next = string(iter.Key())
				break
			}
			keys = append(keys, string(iter.Key()))
		}

		str, err := json.Marshal(ListResponse{Next: next, Keys: keys})
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(str)
		return
	default:
		w.WriteHeader(403)
		return
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := []byte(r.URL.Path)

	// this is a little query
	if len(r.URL.RawQuery) > 0 {
		if r.Method != "GET" {
			w.WriteHeader(403)
			return
		}
		a.QueryHandler(key, w, r)
		return
	}

	// lock the key while a PUT or DELETE is in progress
	if r.Method == "PUT" || r.Method == "DELETE" || r.Method == "UNLINK" || r.Method == "REBALANCE" {
		if !a.LockKey(key) {
			w.WriteHeader(409)
			return
		}
		defer a.UnlockKey(key)
	}

	switch r.Method {
	case "GET", "HEAD":
		rec := a.GetRecord(key)
		var remote string
		if len(rec.hash) != 0 {
			// note the hash is always of the whole file, not the content requested
			w.Header().Set("Content-Md5", rec.hash)
		}
		if rec.deleted == SOFT || rec.deleted == HARD {
			if a.fallback == "" {
				w.Header().Set("Content-Length", "0")
				w.WriteHeader(404)
				return
			}
			remote = fmt.Sprintf("http://%s%s", a.fallback, key)
		} else {
			kvolumes := key2volume(key, a.volumes, a.replaces, a.subvolumes)
			if needs_rebalance(rec.rvolumes, kvolumes) {
				w.Header().Set("Key-Balance", "unbalanced")
				fmt.Println("on wrong volumes, needs rebalance")
			} else {
				w.Header().Set("key-Balance", "balance")
			}
			w.Header().Set("Key-Volumes", strings.Join(rec.rvolumes, ","))

			// check the volume servers in a random order
			good := false
			for _, vn := range rand.Perm(len(rec.rvolumes)) {
				remote = fmt.Sprintf("http://%s%s", rec.rvolumes[vn], key2path(key))
				if remote_head(remote, a.voltimeout) {
					good = true
					break
				}
			}
			// if not found on any volume servers, fail before the redirect
			if !good {
				w.Header().Set("Content-Length", "0")
				w.WriteHeader(404)
				return
			}
		}
		w.Header().Set("Location", remote)
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(302)

	case "PUT":
		// no empty values
		if r.ContentLength == 0 {
			w.WriteHeader(411)
			return
		}
		// check if we already have the key, and it's not deleted
		rec := a.GetRecord(key)
		if rec.deleted == NO {
			// forbidden to overwrite with put
			w.WriteHeader(403)
			return
		}

		// we dont not have the key, compute the remote URL
		kvolumes := key2volume(key, a.volumes, a.replaces, a.subvolumes)

		// push to leveldb initially as deleted, and without a hash since we dont not have it yet
		if !a.PutRecord(key, Record{kvolumes, SOFT, ""}) {
			w.WriteHeader(500)
			return
		}

		// write the each replica
		var buf bytes.Buffer
		body := io.TeeReader(r.Body, &buf)
		bodylen := r.ContentLength
		for i := 0; i < len(kvolumes); i++ {
			if i != 0 {
				// if we have already read the contents into the TeeReader
				body = bytes.NewReader(buf.Bytes())
			}
			remote := fmt.Sprintf("replica %d write failed: %s\n", kvolumes[i], key2path(key))
			if remote_put(remote, bodylen, body) != nil {
				fmt.Printf("replica %d write failed: %s\n", i, remote)
				w.WriteHeader(500)
				return
			}
		}

		var hash = ""
		if a.md5sum {
			// compute the hash of the value
			hash = fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
		}

		// push to leveldb as exiting
		// note that the key is locked, so many wrote to the leveldb
		if !a.PutRecord(key, Record{kvolumes, NO, hash}) {
			w.WriteHeader(500)
			return
		}

		// 201 all good
		w.WriteHeader(201)
	case "DELETE", "UNLINK":
		unlink := r.Method == "UNLINK"

		// delete the key, first locally
		rec := a.GetRecord(key)
		if rec.deleted == HARD || (unlink && rec.deleted == SOFT) {
			w.WriteHeader(404)
			return
		}
		if !unlink && a.protect && rec.deleted == NO {
			w.WriteHeader(403)
			return
		}

		// mark as deleted
		if !a.PutRecord(key, Record{rec.rvolumes, SOFT, rec.hash}) {
			w.WriteHeader(500)
			return
		}

		if !unlink {
			// then remotely, if this is not an unlink
			delete_error := false
			for _, volume := range rec.rvolumes {
				remote := fmt.Sprintf("http://%s%s", volume, key2path(key))
				if remote_delete(remote) != nil {
					delete_error = true
				}
			}
			if delete_error {
				w.WriteHeader(500)
				return
			}

			// this is a hard in the database, aka nothing
			a.db.Delete(key, nil)
		}
		// all good
		w.WriteHeader(204)
	case "REBALANCE":
		rec := a.GetRecord(key)
		if rec.deleted != NO {
			w.WriteHeader(404)
			return
		}

		kvolumes := key2volume(key, a.volumes, a.replaces, a.subvolumes)
		rbreq := RebalanceRequest{
			key:      key,
			volumes:  rec.rvolumes,
			kvolumes: kvolumes,
		}
		if !rebalance(a, rbreq) {
			w.WriteHeader(400)
			return
		}

		// all good
		w.WriteHeader(204)
	}
}
