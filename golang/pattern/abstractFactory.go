package pattern

import "fmt"

type Man interface {
	name()
}

type AMan struct {

}

func (*AMan) name() {
	fmt.Println("a man")
}

type BMan struct {

}

func (*BMan) name() {
	fmt.Println("b man")
}

type Factory interface {
	CreateMan(like string) Man
}

type ChinaFactory struct {
	
}

func (*ChinaFactory) CreateMan(like string) Man {
	if like == "" {

	}
	return nil
}

type AMFactory struct {

}

func (*AMFactory) CreateMan(like string) Man {
	if like == "" {

	}

	return nil
}

type ManFactoryStore struct {
	factory Factory
}

func (store *ManFactoryStore) createMan(like string) Man {
	return store.factory.CreateMan(like)
}


