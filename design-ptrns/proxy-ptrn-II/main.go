package main

func main() {
	car := NewCarProxy(&Driver{12})
	car.Drive()
}
