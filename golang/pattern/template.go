package pattern

import "fmt"

type Cooker interface {
	fire()
	cooke()
	outFire()
}

type CookMenu struct {
	
}

func (c *CookMenu) fire() {
	fmt.Println("fire")
}

func (c *CookMenu) outFire() {
	fmt.Println("out fire")
}

func (c *CookMenu) cooke() {

}

func doCook(cook Cooker) {
	cook.fire()
	cook.cooke()
	cook.outFire()
}

type Menu1 struct {
	CookMenu
}

func (*Menu1) cookie() {
	fmt.Println("1")
}

func testcookie() {
	menu1 := &Menu1{}
	doCook(menu1)
}
