package pushups

import (
	"strconv"
	"strings"
	"time"

	"github.com/wakelesstuna/tracker/internal/ical"
)

type Interface interface {
	List() *PushupList
	TotalWeek(week int) int64
	TotalMonth(month int) int64
}

//	type Event struct {
//		Title       string
//		Description string
//		Start       time.Time
//		End         time.Time
//	}
//
//	type Interface interface {
//		PushUps() ([]Event, error)
//	}
type Pushups struct {
	Ical ical.Interface
}

type Pushup struct {
	Total int64
	Date  time.Time
}

type PushupList struct {
	Items []Pushup
}

func (p *Pushups) List() (*PushupList, error) {
	events, err := p.Ical.PushUps()
	if err != nil {
		return &PushupList{}, err
	}

	list := &PushupList{
		Items: make([]Pushup, 0, len(events)),
	}

	for _, e := range events {
		list.Items = append(list.Items, Pushup{
			Total: extractPushupCount(e),
			Date:  e.Start,
		})
	}

	return list, nil
}

// Total pushups for a given ISO week
func (p *Pushups) TotalWeek(week int) (int64, error) {
	list, err := p.List()
	if err != nil {
		return 0, err
	}
	now := time.Now()

	var total int64

	for _, item := range list.Items {
		year, w := item.Date.ISOWeek()
		nowYear, _ := now.ISOWeek()

		if w == week && year == nowYear {
			total += item.Total
		}
	}

	return total, nil
}

// Total pushups for a given month (1-12)
func (p *Pushups) TotalMonth(month int) (int64, error) {
	list, err := p.List()
	if err != nil {
		return 0, err
	}
	now := time.Now()

	var total int64

	for _, item := range list.Items {
		if int(item.Date.Month()) == month && item.Date.Year() == now.Year() {
			total += item.Total
		}
	}

	return total, nil
}

func extractPushupCount(e ical.Event) int64 {
	parts := strings.Split(e.Description, ",")

	var total int64

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		n, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			continue // skip invalid entries safely
		}

		total += n
	}

	return total
}
