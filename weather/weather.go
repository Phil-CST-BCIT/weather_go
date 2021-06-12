package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// 1. read the api key from apikey.txt
// 2. create an url for the pass-in location, such as vancouver,ca
// 3. pass the newly created url to http.Get method to get the response
// 4. parse the response
// 5. compute the temperature

type Report struct {
	Name        string  `json:"name"`
	Country     string  `json:"country"`
	Temp        float64 `json:"temp"`
	Humidity    float64 `json:"humidity"`
	Description string  `json:"description"`
}

// read the api key from a txt file
// return the api key as a string
func apikey() string {
	b, err := ioutil.ReadFile("../apikey.txt")

	if err != nil {
		log.Fatal(err)
	}

	str := string(b)

	return str
}

// insert a location into the url query for a GET request
func url_for(location string) string {

	api_key := apikey()

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", location, api_key)

	return url
}

// retrieve the weather information from openwehthermap
// return the response body, which is a slice of bytes
func get_weather_info(url string) []byte {

	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, e := io.ReadAll(res.Body)

	if e != nil {
		log.Fatal(e)
	}

	return body

}

// obtain the desired data from the byte slice
func parse(b []byte) *Report {

	// we first declare an interface for storing the decoded json data
	var intf interface{}

	// encoding/json provides the capability to parse arbitrary JSON into an interface
	// parse the json and store it into intf
	err := json.Unmarshal(b, &intf)

	if err != nil {
		log.Fatal(err)
	}

	r := new(Report)

	parse_helper(intf, r)

	return r
}

// we have a go type interface, which means it can be anything
// we want to extract data from the interface and assign
// the data to the corresponding fields of Report
// we need to recursively walk through the data structure
func parse_helper(intf interface{}, r *Report) {

	// first, we check the type of the interface
	switch v := intf.(type) {

	// if it is of type array
	case []interface{}:

		// take each element in the array
		for _, u := range v {
			parse_helper(u, r)
		}

	// if it's of type map of string to anything
	case map[string]interface{}:

		if v["name"] != nil {

			// In order to assign the value to a field of Report
			// we need to check the concrete type of map v of string to something
			// if it is a map of string to string
			if str, ok := v["name"].(string); ok {
				r.Name = str
			}
		}

		if v["country"] != nil {

			if str, ok := v["country"].(string); ok {
				r.Country = str
			}
		}

		if v["humidity"] != nil {

			if fl, ok := v["humidity"].(float64); ok {
				r.Humidity = fl
			}
		}

		if v["temp"] != nil {

			if fl, ok := v["temp"].(float64); ok {
				r.Temp = fl
			}
		}

		if v["description"] != nil {

			if str, ok := v["description"].(string); ok {
				r.Description = str
			}
		}

		// get a vale from the map and use the value to find fields we need
		// until every value in the map is parsed
		for _, w := range v {
			parse_helper(w, r)
		}
	}
}

// for second solution
// func pretty_print(r *Report) {

// 	str := fmt.Sprintf("%s %s: %.2f deg C, %d%%, %s", r.Name, r.Country, r.Temp, int64(r.Humidity), r.Description)

// 	fmt.Println(str)

// }

func Temperature_of(location string, ch chan *Report) {

	url := url_for(location)

	b := get_weather_info(url)

	r := parse(b)

	ch <- r

	// pretty_print(r)

}
