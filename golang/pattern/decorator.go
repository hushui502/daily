package pattern

import "fmt"

type Person interface {
	cost() int
	show()
}

type hufan struct {
	
}

func (h *hufan) show() {
	fmt.Println("hufan no wear")
}

func (h *hufan) cost() int {
	return 0
}

type clothesDecorator struct {
	person Person
}

func (c *clothesDecorator) cost() int {
	return 0
}

func (c *clothesDecorator) show() {

}

type Jacket struct {
	clothesDecorator
}

func (j *Jacket) cost() int {
	return j.person.cost() + 1
}

func (j *Jacket) show() {
	j.person.show()
	fmt.Println("hufan wear jacket")
}

func testwear() {
	hufan := &hufan{}

	jacket := &Jacket{}
	jacket.person = hufan
	jacket.show()


}



