package main

import "github.com/skip2/go-qrcode"

func main() {
	qrcode.WriteFile("https://www.baidu.com", qrcode.Medium, 256, "D:\\project\\go\\src\\awesomeProject2\\qrcode\\qr.png")
}
