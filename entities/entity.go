package entities

import (
	"regexp"
	"strings"
)

var BaseURL = ""
var MqttBroker = ""
var MqttWSSBroker = ""

func localize(lang, th string, en string) string {
	switch {
	case lang == "th" && th != "-":
		return th
	case lang == "en" && en != "-":
		return en
	default:
		if th != "-" {
			return th
		} else {
			return en
		}
	}
}
func ELocalize(lang, th string, en string) string {
	switch {
	case lang == "th" && th != "-":
		return th
	case lang == "en" && en != "-":
		return en
	default:
		if th != "-" {
			return th
		} else {
			return en
		}
	}
}

// NewGeoLocationPoint new geo location point
func NewGeoLocationPoint(lat, lng float64) *GeoLocation {
	return &GeoLocation{
		Type:        "Point",
		Coordinates: []float64{lng, lat},
	}
}

// GeoLocation geolocation
type GeoLocation struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

func IsPhoneNumber(phoneNumber string) bool {
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

	return re.MatchString(phoneNumber)
}

func GetPhoneNumbersFromRawPhoneNumbers(rawPhoneNumber string) *[]string {
	removedAllSpacePhoneNumber := strings.ReplaceAll(rawPhoneNumber, " ", "")
	slicedPhoneNumbers := strings.Split(removedAllSpacePhoneNumber, ",")
	numberOfPhoneNumbers := len(slicedPhoneNumbers)
	if numberOfPhoneNumbers > 0 {
		return &slicedPhoneNumbers
	}

	return nil
}
