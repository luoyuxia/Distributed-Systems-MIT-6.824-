package main
import "fmt"

type Bitcoin int
type Rectangle struct {
	Width float64
	Height float64
}

func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Shape interface {
	Area() float64
}

func checkArea(wanted float64, shape Shape)  error {
	fmt.Print(wanted == shape.Area())
	return nil
}

func (b Bitcoin) String() string  {
	return fmt.Sprintf("%d BTC", b)
}

func main()  {
	t := []int {1, 2, 3}
	fmt.Println(t)
	r := Rectangle{Width: 1, Height: 2}
	b := Bitcoin(12)
	fmt.Println(b)
	fmt.Println(r.Area())
}
