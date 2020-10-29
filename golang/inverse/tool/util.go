package tool

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"os"
	"time"
)

func isFileOrDir(filename string, decideDir bool) bool {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false
	}
	isDir := fileInfo.IsDir()
	if decideDir {
		return isDir
	}

	return !isDir
}

func checkFileIsExist(filepath string) bool {
	exist := true
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		exist = false
	}

	return exist
}

func GenRandom(start int, end int, count int) []int {
	if end < start || (end-start) < count {
		return nil
	}

	nums := make([]int, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		num := r.Intn(end-start)+start

		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}

func MD5Uri(uri string) string {
	ctx := md5.New()
	ctx.Write([]byte(uri))

	return hex.EncodeToString(ctx.Sum(nil))
}