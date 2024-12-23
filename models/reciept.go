package models

import (
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// Defining Item struct data type to store the values of the Reciept's Item values.
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Defining Reciept struct data type to store the reciept values as a whole
type Reciept struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// Map data structure to store the reciept in-memory and make to reserve the memory for it.
var AllReceipts = make(map[string]Reciept)

// Service which acts on the preassumed type of reciept to generate id and store it in the map.
func (reciept Reciept) SaveReciept() string {
	// Random unique is generated with the help of an external library, uuid and converted to string as a key.
	genratedId := uuid.New().String()
	// Genereated id and the reciept are stored in the map.
	AllReceipts[genratedId] = reciept
	return genratedId
}

// Service to calculate points for the reciept which accepts an id, based on which the calculation is done and returned.
func CalculatePoints(recieptId string) int64 {
	// Declaring a variable of type int64 to store the total points.
	var points int64
	// Declaring and Intialising start time and end time to calculate points based on the time 2:00pm and 4:00pm.
	startTime := time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)
	endTime := time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)
	// Retrieving and storing the reciept data from map based on the key(id).
	data := AllReceipts[recieptId]
	// Parsing the total from string to float and storing it.
	totalF, _ := strconv.ParseFloat(data.Total, 64)

	// If condition to check if the total is a round figure.
	// Using the mod function in the math library, for example math.Mod(5.0, 1) = 0.0, true and is a round figure,
	// and math.Mod(5.5,1) = 0.5, false because not a round figure.
	if math.Mod(totalF, 1) == 0 {
		points += 50
	}
	// If condition to check if the total is a multipe of 0.25.
	if int(totalF*100)%25 == 0 {
		points += 25
	}
	// Points are incremented for every 2 items in the list of items, calculated based on the length of items.
	points += int64(len(data.Items)/2) * 5
	// For loop to increment points by the product of price and 0.2 if the trimmed length of description is a multiple of 3,
	// and rounding it to the nearest integer before adding it to the points.
	for _, item := range data.Items {
		itemPrice, _ := strconv.ParseFloat(item.Price, 32)
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			points += int64(math.Ceil(itemPrice * 0.2))
		}
	}
	// Variable to store the parsed date string to a int type to check if the date is an odd and incrementing the points.
	convertedDay, _ := strconv.ParseInt(data.PurchaseDate[len(data.PurchaseDate)-2:], 10, 32)
	if convertedDay%2 != 0 {
		points += 6
	}
	// Declaring the format of the time which is 24 hrs.
	format := "15:04"
	// Variable to store the parsed string type of the date to date type.
	parsedTime, _ := time.Parse(format, data.PurchaseTime)
	// If condition to check if the purchase is made after 2:00 pm and before 4:00 pm.
	if parsedTime.After(startTime) && parsedTime.Before(endTime) {
		points += 10
	}
	// For loop to check and increment points for every alphanumeric characters in the retailer's name.
	for _, char := range data.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			points++
		}
	}
	return points
}
