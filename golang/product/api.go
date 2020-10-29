package product

type product interface {
	Store(name string, id int)
	Find(id int)
}



