package main

import "fmt"

type Person interface {
	WithName(name string) Person
	Name() string
	WithFavoriteColors(favoriteColors []string) Person
	NumFavoriteColors() int
	FavoriteColorAt(i int) string
	WithFavoriteColorAt(i int, favoriteColor string) Person
	AtVersion(version int) Person
}

type person struct {
	name           string
	favoriteColors []string
	history []person
}


func (p person) WithName(name string) Person {
	p.name = name
	return p.nextState()
}

func (p person) Name() string {
	return p.name
}

func (p person) WithFavoriteColors(favoriteColors []string) Person {
	defer func() {
		p.nextState()
	}()

	p.favoriteColors = favoriteColors
	return p
}

func (p person) FavoriteColors() []string {
	return p.favoriteColors
}

// 为了避免切片的公用底层，改位每次返回一个新的切片
//func updateFavoriteColors(p person) person {
//	return p.WithFavoriteColors(append([]string{"red"}, p.FavoriteColors()[1:]...))
//}

func (p person) NumFavoriteColors() int {
	return len(p.favoriteColors)
}

func (p person) FavoriteColorAt(i int) string {
	return p.favoriteColors[i]
}

func (p *person) nextState() Person {
	fmt.Printf("nextstate: %#+v\n", p)
	p.history = append(p.history, *p)
	return p
}

func (p person) AtVersion(version int) Person {
	return p.history[version]
}


// 直接避免返回切片
func (p person) WithFavoriteColorAt(i int, favoriteColor string) Person {
	p.favoriteColors = append(p.favoriteColors[:i],
		append([]string{favoriteColor}, p.favoriteColors[i+1:]...)...)

	return p
}

func updateFavoriteColors(p Person) Person {
	return p.WithFavoriteColorAt(0, "red")
}

func NewPerson() Person {
	return person{}.
		WithName("no name")
}


func main() {
	me := NewPerson().
		WithName("Elliot").
		WithFavoriteColors([]string{"black", "blue"})

	// We discard the result, but it will be put into the history.
	updateFavoriteColors(me)

	fmt.Printf("%s\n", me.AtVersion(0).Name())
	fmt.Printf("%s\n", me.AtVersion(1).Name())
}