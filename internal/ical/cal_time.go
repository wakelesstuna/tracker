package ical

import (
	"time"

	ics "github.com/arran4/golang-ical"
)

func parseEventTime(prop *ics.IANAProperty) (time.Time, error) {
	value := prop.Value

	tzid := ""
	if len(prop.ICalParameters["TZID"]) > 0 {
		tzid = prop.ICalParameters["TZID"][0]
	}

	return parseICalTime(value, tzid)
}

func parseICalTime(value string, tzid string) (time.Time, error) {
	layout := "20060102T150000"

	if tzid == "" {
		return time.Parse(layout, value)
	}

	loc, err := time.LoadLocation(tzid)
	if err != nil {
		return time.Time{}, err
	}

	return time.ParseInLocation(layout, value, loc)
}
