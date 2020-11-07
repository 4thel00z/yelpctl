package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	path    = flag.String("path", "", "path to yelp dataset, must not be empty")
	bboxRaw = flag.String("bbox", "", "comma separated list of 4 floats with the bounding box coordinates, in this format:\n(lat_min,lat_max,lng_min,lng_max)\nMust not be empty.")
)

type Business struct {
	BusinessID  string     `json:"business_id"`
	Name        string     `json:"name"`
	Address     string     `json:"address"`
	City        string     `json:"city"`
	State       string     `json:"state"`
	PostalCode  string     `json:"postal_code"`
	Latitude    float64    `json:"latitude"`
	Longitude   float64    `json:"longitude"`
	Stars       float64    `json:"stars"`
	ReviewCount int        `json:"review_count"`
	IsOpen      int        `json:"is_open"`
	Attributes  Attributes `json:"attributes"`
	Categories  string     `json:"categories"`
	Hours       Hours      `json:"hours"`
}

type Attributes struct {
	BusinessAcceptsCreditCards string `json:"BusinessAcceptsCreditCards"`
	BikeParking                string `json:"BikeParking"`
	GoodForKids                string `json:"GoodForKids"`
	BusinessParking            string `json:"BusinessParking"`
	ByAppointmentOnly          string `json:"ByAppointmentOnly"`
	RestaurantsPriceRange2     string `json:"RestaurantsPriceRange2"`
}

type Hours struct {
	Monday    string `json:"Monday"`
	Tuesday   string `json:"Tuesday"`
	Wednesday string `json:"Wednesday"`
	Thursday  string `json:"Thursday"`
	Friday    string `json:"Friday"`
	Saturday  string `json:"Saturday"`
	Sunday    string `json:"Sunday"`
}

// lat_min, lat_max, lng_min, lng_max
type Boundingbox []float64

func (bb Boundingbox) Contains(lat, lng float64) bool {
	return bb[0] <= lat && lat <= bb[1] && bb[2] <= lng && lng <= bb[3]
}

func FromStrings(raw []string) (result Boundingbox, err error) {
	if len(raw) != 4 {
		err = errors.New("len of bounding box not 4")
		return
	}
	result = make(Boundingbox, 4)
	for i, val := range raw {
		var parsed float64
		parsed, err = strconv.ParseFloat(val, 64)
		if err != nil {
			return
		}
		result[i] = parsed
	}
	return
}

func cry(f string, vals ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stderr, f, vals...)
}

func main() {

	flag.Parse()

	if *path == "" {
		cry("\n%s\n", "path must not be empty")
		flag.PrintDefaults()
		os.Exit(1)
	}

	bbox, err := FromStrings(strings.Split(*bboxRaw, ","))

	if err != nil {
		cry("\n%s\n", "could not convert bbox")
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(*path)
	if err != nil {
		cry("\n%s\n", err.Error())
		flag.PrintDefaults()
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var business Business
		content := scanner.Bytes()
		err := json.Unmarshal(content, &business)
		if err != nil {
			cry("\n%s\n", err.Error())
			flag.PrintDefaults()
			os.Exit(1)
		}

		if bbox.Contains(business.Latitude, business.Longitude) {
			fmt.Println(string(content))
		}

	}

	if err := scanner.Err(); err != nil {
		cry("\n%s\n", err.Error())
		flag.PrintDefaults()
		os.Exit(1)
	}
}
