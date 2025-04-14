package main

import "fmt"

type Animal interface {
	Speak()
}

type Cat struct{}
type Dog struct{}

func (c Cat) Speak() {
	fmt.Println("Meow")
}

func (d Dog) Speak() {
	fmt.Println("Woof")
}

func makeItSpeak(a Animal) {
	a.Speak()
}

func main() {
	c := Cat{}
	d := Dog{}

	c.Speak()
	makeItSpeak(c)

	d.Speak()
	makeItSpeak(d)
}
