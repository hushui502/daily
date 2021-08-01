package main

import (
	"fmt"
	"strings"
	"sync"
	"unicode"
)

/*
交替打印数字和字⺟
问题描述
使⽤两个  goroutine 交替打印序列，⼀个  goroutine 打印数字， 另外⼀
个  goroutine 打印字⺟， 最终效果如下：
解题思路
问题很简单，使⽤ channel 来控制打印的进度。使⽤两个 channel ，来分别控制数字和
字⺟的打印序列， 数字打印完成后通过 channel 通知字⺟打印, 字⺟打印完成后通知数
字打印，然后周⽽复始的⼯作。
源码参考
1 12AB34CD56EF78GH910IJ1112KL1314MN1516OP1718QR1920ST2122UV2324WX2526YZ2728
*/
func printA(number, letter chan bool, wg *sync.WaitGroup) {
	// number
	go func() {
		i := 1
		for {
			select {
			case <-number:
				fmt.Print(i)
				i++
				fmt.Print(i)
				i++
				letter <- true
			default:
			}
		}
	}()

	wg.Add(1)
	// letter
	go func(wg *sync.WaitGroup) {
		str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

		i := 0
		for {
			select {
			case <-letter:
				if i >= len(str)-1 {
					wg.Done()
					return
				}
				fmt.Print(string(str[i]))
				i++
				if i >= len(str) {
					i = 0
				}
				fmt.Print(string(str[i]))
				i++
				number <- true
			default:
			}
		}
	}(wg)

	number <- true
	wg.Wait()
}

/*
判断字符串中字符是否全都不同
问题描述
请实现⼀个算法，确定⼀个字符串的所有字符【是否全都不同】。这⾥我们要求【不允
许使⽤额外的存储结构】。 给定⼀个string，请返回⼀个bool值,true代表所有字符全都
不同，false代表存在相同的字符。 保证字符串中的字符为【ASCII字符】。字符串的⻓
度⼩于等于【3000】。
解题思路
这⾥有⼏个重点，第⼀个是 ASCII字符 ， ASCII字符 字符⼀共有256个，其中128个是常
⽤字符，可以在键盘上输⼊。128之后的是键盘上⽆法找到的。
然后是全部不同，也就是字符串中的字符没有重复的，再次，不准使⽤额外的储存结
构，且字符串⼩于等于3000。
如果允许其他额外储存结构，这个题⽬很好做。如果不允许的话，可以使⽤golang内置
的⽅式实现。
*/
func isUniqueString(s string) bool {
	if strings.Count(s, "") > 3000 {
		return false
	}

	for _, v := range s {
		if v > 127 {
			return false
		}
		if strings.Count(s, string(v)) > 1 {
			return false
		}
	}

	return true
}

func isUniqueString2(s string) bool {
	if strings.Count(s, "") > 3000 {
		return false
	}

	for i, v := range s {
		if v > 127 {
			return false
		}
		if strings.LastIndex(s, string(v)) != i {
			return false
		}
	}

	return true
}

/*
翻转字符串
问题描述
请实现⼀个算法，在不使⽤【额外数据结构和储存空间】的情况下，翻转⼀个给定的字
符串(可以使⽤单个过程变量)。
给定⼀个string，请返回⼀个string，为翻转后的字符串。保证字符串的⻓度⼩于等于
5000。
解题思路
翻转字符串其实是将⼀个字符串以中间字符为轴，前后翻转，即将str[len]赋值给str[0],
将str[0] 赋值 str[len]。
源码参考
*/
func reverseString(s string) (string, bool) {
	runes := []rune(s)

	if len(runes) > 5000 {
		return s, false
	}

	for i := 0; i < len(runes)/2; i++ {
		runes[i], runes[len(runes)-1-i] = runes[len(runes)-1-i], runes[i]
	}

	return string(runes), true
}

/*
给定两个字符串，请编写程序，确定其中⼀个字符串的字符重新排列后，能否变成另⼀
个字符串。 这⾥规定【⼤⼩写为不同字符】，且考虑字符串重点空格。给定⼀个string
s1和⼀个string s2，请返回⼀个bool，代表两串是否重新排列后可相同。 保证两串的
⻓度都⼩于等于5000。
解题思路
⾸先要保证字符串⻓度⼩于5000。之后只需要⼀次循环遍历s1中的字符在s2是否都存
在即可。
源码参考
*/

func isRegroup(s1, s2 string) bool {
	sl1 := len([]rune(s1))
	sl2 := len([]rune(s2))

	if sl1 > 5000 || sl2 > 5000 || sl1 != sl2 {
		return false
	}

	for _, v := range s1 {
		if strings.Count(s1, string(v)) != strings.Count(s2, string(v)) {
			return false
		}
	}

	return true
}

/*
问题描述
请编写⼀个⽅法，将字符串中的空格全部替换为“%20”。 假定该字符串有⾜够的空间存
放新增的字符，并且知道字符串的真实⻓度(⼩于等于1000)，同时保证字符串由【⼤⼩
写的英⽂字⺟组成】。 给定⼀个string为原始的串，返回替换后的string。
解题思路
两个问题，第⼀个是只能是英⽂字⺟，第⼆个是替换空格。
*/
func replaceBlank(s string) (string, bool) {
	if len([]rune(s)) > 5000 {
		return s, false
	}
	for _, v := range s {
		if string(v) != " " && unicode.IsLetter(v) == false {
			return s, false
		}
	}

	return strings.Replace(s, " ", "%20", -1), true
}

func reverseString2(s string) string {
	runes := []rune(s)

	for i := 0; i < len(runes)/2; i++ {
		runes[i], runes[len(runes)-i-1] = runes[len(runes)-i-1], runes[i]
	}

	return string(runes)
}

/*
给定⼀个字符串，找到它的第⼀个不重复的字符，并返回它的索引。如果不存在，则返
回 -1 。 案例:
s = "leetcode" 返回 0.
s = "loveleetcode", 返回 2
*/
func firstUniqueChar(s string) int {
	var arr [26]int

	for i, k := range s {
		arr[k-'a'] = i
	}

	for i, k := range s {
		if arr[k-'a'] == i {
			return i
		}
	}

	return -1
}

func isPalindrome(s string) bool {
	if s == "" {
		return false
	}
	s = strings.ToLower(s)

	if len(s) == 2 {
		return s[0] == s[1]
	}

	left, right := 0, len(s)-1

	for left < right {
		if !(s[left] >= 'a' && s[left] <= 'z') || (s[left] >= '0' && s[left] <= '9') {
			left++
			continue
		}
		if !(s[right] >= 'a' && s[right] <= 'z') || (s[right] >= '0' && s[right] <= '9') {
			right--
			continue
		}
		if s[left] != s[right] {
			return false
		}

		left++
		right--
	}

	return true
}

func maxSlideWindow(nums []int, k int) []int {
	len := len(nums)
	res := make([]int, 0)
	index := 0

	for index < len {
		maxNum := nums[index]

		if index > len-k {
			break
		}

		for i := index + 1; i < index+k; i++ {
			if nums[i] > maxNum {
				maxNum = nums[i]
			}
		}

		res = append(res, maxNum)
		index++
	}

	return res
}

func bubbleSort(arr []int) []int {
	if len(arr) == 0 {
		return []int{}
	}

	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}

	return arr
}

func selectSort(arr []int) []int {
	l := len(arr)
	if l == 0 {
		return arr
	}
	for i := 0; i < l; i++ {
		min := i
		for j := i + 1; j < l; j++ {
			if arr[j] < arr[min] {
				min = j
			}
		}
		arr[i], arr[min] = arr[min], arr[i]
	}

	return arr
}
