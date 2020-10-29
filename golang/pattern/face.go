package pattern

import "fmt"

type CPU struct {

}

func (CPU) start() {
	fmt.Println("start cpu")
}

type Disk struct {

}

func (Disk) start() {
	fmt.Println("start disk")
}

type StartButton struct {
	
}

func (StartButton) start() {
	cpu := &CPU{}
	cpu.start()
	disk := &Disk{}
	disk.start()
}
