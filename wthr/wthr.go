package main

import (
	"fmt"
	"localhost/weather"
	"os"
	"sort"
	"sync"
)

var (
	wg sync.WaitGroup
)

// we can delete the Report type and its methods
// if using the second way for solving this problem
type Report struct {
	Name        string
	Country     string
	Temp        float64
	Humidity    float64
	Description string
}

type ByName []Report

func (rs ByName) Len() int {
	return len(rs)
}

func (rs ByName) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs ByName) Less(i, j int) bool {
	return rs[i].Name < rs[j].Name
}

// get cmd line arguments and return a string slice
func collect_cities() []string {

	args := os.Args[1:]

	cities := []string{}

	cities = append(cities, args...)

	return cities
}

// many goroutines will call this function
// we want to create a critical section for sychronization
// if using the second solution
func city_weather(city string, ch chan *weather.Report) {

	defer wg.Done()

	weather.Temperature_of(city, ch)

}

func create_report(wr *weather.Report) *Report {

	r := new(Report)

	r.Name = wr.Name
	r.Country = wr.Country
	r.Temp = wr.Temp
	r.Humidity = wr.Humidity
	r.Description = wr.Description

	return r
}

func my_printer(rs []Report, n int) {

	for _, v := range rs {

		str := fmt.Sprintf("%s %s: %.2f deg C, %d%%, %s", v.Name, v.Country, v.Temp, int64(v.Humidity), v.Description)

		fmt.Println(str)

	}

	if n > 0 {
		fmt.Fprintf(os.Stderr, "Number of incorrect city names: %d\n", n)
	}

}

func main() {

	cities := collect_cities()

	rs := []Report{}

	var counter int

	wg.Add(len(cities))

	// From Albert
	// creating a go routine & immediately reading from the channel
	// associated with the go routine makes it sequential
	for _, v := range cities {

		ch := make(chan *weather.Report)

		go city_weather(v, ch)

		wr := <-ch

		var r *Report

		if len(wr.Name) > 0 {
			r = create_report(wr)
			rs = append(rs, *r)
		} else {
			counter++
		}

	}

	sort.Sort(ByName(rs))

	my_printer(rs, counter)

	// THE FOLLOWING CODE IS THE SECOND SOLUTION
	// We first collect all the cities from command line
	// cities is a string slice, so we can sort it
	// after sorting the slice, we can use channel or mutex for synchronization
	// the result then is printed out in weather.go

	// cities := collect_cities()

	// sort.Strings(cities)

	// ch := make(chan bool, 1)

	// lock := make(chan bool, 1)

	// for _, v := range cities {

	// 	// the main goroutine is blocked here until the New York goroutine unlocks
	// 	lock <- true

	// 	// the first city is New York City
	// 	go city_weather(v, ch, lock)

	// 	// the main goroutine sends the first message to the channel
	// 	ch <- true

	// }

	// close(ch)
	// close(lock)

	wg.Wait()

}
