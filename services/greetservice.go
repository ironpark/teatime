package services

type Greeter struct{}

func (g *Greeter) Greet(name string) string {
	return "Hello " + name + "!"
}
