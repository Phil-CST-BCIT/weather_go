package main

import (
	"fmt"
    "localhost/weather"
)

func main() {
    k := weather.Apikey()
	fmt.Println(k)
}
