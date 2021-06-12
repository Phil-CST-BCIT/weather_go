package weather

import (
	"io/ioutil"
	"log"
)

// 1. read the api key from apikey.txt
// 2. create an url for the pass-in location, such as vancouver,ca
// 3. pass the newly created url to http.Get method to get the response
// 4. parse the response
// 5. compute the temperature
// 6.

func Apikey() string {
	b, err := ioutil.ReadFile("../apikey.txt")

	if err != nil {
		log.Fatal(err)
	}

	str := string(b)

	// fmt.Println(string(b))

	return str
}
