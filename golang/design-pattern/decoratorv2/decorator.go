package decoratorv2

type IDraw interface {
	Draw() string
}

type Square struct {}

func (s Square) Draw() string {
	return "this is a square"
}

type ColorSquare struct {
	color string
	square IDraw
}

func NewColorSquare(square IDraw, color string) ColorSquare {
	return ColorSquare{color:color, square:square}
}

func (c *ColorSquare) Draw() string {
	return c.square.Draw() + ", color is " + c.color
}
