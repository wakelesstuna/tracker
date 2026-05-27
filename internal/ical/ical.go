package ical

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
)

var _ Interface = &Client{}

const PushUpCalenderUrl = "PUSH_UP_CALENDER_URL"

var calenderUrl = ""

func init() {
	calenderUrl = os.Getenv(PushUpCalenderUrl)
	if calenderUrl == "" {
		fmt.Printf("missing push up calendar url environment variable, set %s ", PushUpCalenderUrl)
	}
}

const (
	PushUpsKey = "pushups"
)

type Event struct {
	Title       string
	Description string
	Start       time.Time
	End         time.Time
}

type Interface interface {
	PushUps() ([]Event, error)
}
type Client struct {
}

func (p *Client) PushUps() ([]Event, error) {
	resp, err := http.Get(calenderUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	cal, err := ics.ParseCalendar(resp.Body)
	if err != nil {
		return nil, err
	}

	var events []Event
	for _, e := range cal.Events() {
		summary := e.GetProperty(ics.ComponentPropertySummary)
		desc := e.GetProperty(ics.ComponentPropertyDescription)
		startProp := e.GetProperty(ics.ComponentPropertyDtStart)

		if summary == nil || startProp == nil || desc == nil {
			continue
		}

		if strings.ToLower(summary.Value) != PushUpsKey {
			continue
		}

		start, err := parseEventTime(startProp)
		if err != nil {
			continue
		}

		events = append(events, Event{
			Title:       summary.Value,
			Description: desc.Value,
			Start:       start.Local(),
		})
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Start.Before(events[j].Start)
	})

	return events, nil
}
