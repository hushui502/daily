package pattern

import "fmt"

type Boy interface {
	name()
}

type OldBoy struct {

}

func (o *OldBoy) name() {
	fmt.Println("ole")
}

type YBoy struct {

}

func (y *YBoy) name() {
	fmt.Println("young")
}

type BoyFactory struct {
	
}

func (*BoyFactory) CreateBoy(like string) Boy {
	switch like {
	case "old":
		return &OldBoy{}
	case "y":
		return &YBoy{}
	default:

	}

	return nil
}
