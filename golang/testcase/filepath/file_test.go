package filepath

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

// Just for learning how to use some method, so the test case will against the SRP

const (
	Separator = os.PathSeparator	// "/"
	ListSeparator = os.PathListSeparator	// ";"
)

func TestSeparator(t *testing.T) {
	s := "http://www.baidu.com/a/b"
	u, _ := url.Parse(s)

	path := u.Path
	fmt.Println(path)	// /a/b

	path = filepath.FromSlash(u.Path)	// if current os is windows, the path is  \a\b
	fmt.Println(path)

	if err := os.MkdirAll(path[1:], 0777); err != nil {
		fmt.Println(err)
	}

	// convert "\a\b" to "/a/b", if your server os is windows, you should use this to determine the right path
	path = filepath.ToSlash(path)		// /a/b
}

func TestBase(t *testing.T) {
	path := `a///b///c///d`
	path = filepath.FromSlash(path)
	fmt.Println(path)

	// 获取path最后一个分隔符之前的部分（不包含最后分隔符）
	d1 := filepath.Dir(path)
	fmt.Println(d1)		// a\b\c


	// 获取path最后一个分隔符之后的部分（不包含分隔符）
	f1 := filepath.Base(path)
	fmt.Println(f1)		// d

	// 以最后一个分隔符，分割两部分
	d2, f2 := filepath.Split(path)
	fmt.Println(d2, "  ", f2)		// a\\\b\\\c\\\    d

	// 获取文件的扩展名
	ext := filepath.Ext(path + ".jpg")
	fmt.Println(ext)	// .jpg
}

func TestRel(t *testing.T) {
	// 都是绝对路径
	s, err := filepath.Rel(`a/b/c`, `a/b/c/d/e`)
	fmt.Println(s, err)	// d\e <nil>

	// 都是相对路径
	s, err = filepath.Rel(`a/b/c`, `a/b/c/d/e`)
	fmt.Println(s, err) // d/e <nil>

	// 一个绝对一个相对
	s, err = filepath.Rel(`/a/b/c`, `a/b/c/d/e`)
	fmt.Println(s, err)
	//  Rel: can't make a/b/c/d/e relative to /a/b/c

	// 一个相对一个绝对
	s, err = filepath.Rel(`a/b/c`, `/a/b/c/d/e`)
	fmt.Println(s, err)
	//  Rel: can't make /a/b/c/d/e relative to a/b/c

	// 从 `a/b/c` 进入 `a/b/d/e`，只需要进入 `../d/e` 即可
	s, err = filepath.Rel(`a/b/c`, `a/b/d/e`)
	fmt.Println(s, err) // ../d/e <nil>
}

func TestJoin(t *testing.T) {
	// 将 elem 中的多个元素合并为一个路径，忽略空元素，清理多余字符。
	s := filepath.Join("a", "b", "", ":::", " ", `//c////d///`)
	fmt.Println(s) // a/b/:::/  /c/d
}

func TestClean(t *testing.T) {
	// 清除../ /// ./
	s := filepath.Clean("a/./b/:::/..// /c/..///d///")
	fmt.Println(s)	// a\b\ \d
}

func TestAbs(t *testing.T) {
	s1 := "a/b/c/d"
	fmt.Println(filepath.Abs(s1)) // D:\project\go\src\lib\test\filepath\a\b\c\d <nil>

	// 用`` symbol 主要是因为这样可以避免/的转换，节约时间，代码清晰
	s2 := `D:\project\go\src\lib\test\filepath\a\b\c\d`
	fmt.Println(filepath.IsAbs(s1))
	fmt.Println(filepath.IsAbs(s2))
}

func TestSplitList(t *testing.T) {
	path := `a/b/c:d/e/f:   g/h/i`
	s := filepath.SplitList(path)
	fmt.Printf("%q", s)  // ["a/b/c" "d/e/f" "   g/h/i"]
}

func TestVolumeName(t *testing.T) {
	path := `D:\project\go\src\lib\test\filepath\a\b\c\d`

	// 返回路径字符串中的卷名
	volumeName := filepath.VolumeName(path)
	fmt.Println(volumeName)		// "D:"
}

func TestEvalSymlinks(t *testing.T) {
	path := `C:\Users\hufan\AppData\Local\SourceTree\SourceTreelink.exe`

	// 返回链接（快捷方式）所指向的实际文件
	fmt.Println(filepath.EvalSymlinks(path))
}

func TestGlob(t *testing.T) {
	// 遍历出/usr子目录下的所有ab开头的项目(不区分大小写)
	list, err := filepath.Glob("/usr/*/[Bb][Aa]*")
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range list {
		fmt.Println(v)
	}
}

// 列出包含txt文件的目录，会跳过子目录
func findTxtDir(path string, info os.FileInfo, err error) error {
	ok, err := filepath.Match(`*.txt`, info.Name())
	if ok {
		fmt.Println(filepath.Dir(path), info.Name())
		// 这里skip就会导致当前的目录的子目录不会包含
		return filepath.SkipDir
	}

	return err
}


// 列出所有以 ab 开头的目录（全部，因为没有跳过任何项目）
func findabDir(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		ok, err := filepath.Match(`[aA][bB]*`, info.Name())
		if err != nil {
			return err
		}
		if ok {
			fmt.Println(path)
		}
	}

	return nil
}

func TestWalkFunc(t *testing.T) {
	// 列出含有 *.txt 文件的目录（不是全部，因为会跳过一些子目录）
	err := filepath.Walk(`/usr`, findTxtDir)
	fmt.Println(err)

	fmt.Println("==============================")

	// 列出所有以 ab 开头的目录（全部，因为没有跳过任何项目）
	err = filepath.Walk(`/usr`, findabDir)
	fmt.Println(err)
}