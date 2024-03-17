package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
	"unicode"
)

type Duration struct {
	Years    float64
	Months   float64
	Weeks    float64
	Days     float64
	Hours    float64
	Minutes  float64
	Seconds  float64
	Negative bool
}

const (
	parsingPeriod = iota
	parsingTime

	hoursPerDay   = 24
	hoursPerWeek  = hoursPerDay * 7
	hoursPerMonth = hoursPerYear / 12
	hoursPerYear  = hoursPerDay * 365

	nsPerSecond = 1000000000
	nsPerMinute = nsPerSecond * 60
	nsPerHour   = nsPerMinute * 60
	nsPerDay    = nsPerHour * hoursPerDay
	nsPerWeek   = nsPerHour * hoursPerWeek
	nsPerMonth  = nsPerHour * hoursPerMonth
	nsPerYear   = nsPerHour * hoursPerYear
)

var (
	// ErrUnexpectedInput is returned when an input in the duration string does not match expectations
	ErrUnexpectedInput = errors.New("unexpected input")
)

func ParseDurationString(d string) (int64, error) {
	state := parsingPeriod
	duration := &Duration{}
	num := ""
	var err error

	for _, char := range d {
		switch char {
		case 'P':
			state = parsingPeriod
		case 'T':
			state = parsingTime
		case 'Y':
			if state != parsingPeriod {
				return 0, ErrUnexpectedInput
			}

			duration.Years, err = strconv.ParseFloat(num, 64)
			if err != nil {
				return 0, err
			}
			num = ""
		case 'M':
			if state == parsingPeriod {
				duration.Months, err = strconv.ParseFloat(num, 64)
				if err != nil {
					return 0, err
				}
				num = ""
			} else if state == parsingTime {
				duration.Minutes, err = strconv.ParseFloat(num, 64)
				if err != nil {
					return 0, err
				}
				num = ""
			}
		case 'W':
			if state != parsingPeriod {
				return 0, ErrUnexpectedInput
			}

			duration.Weeks, err = strconv.ParseFloat(num, 64)
			if err != nil {
				return 0, err
			}
			num = ""
		case 'D':
			if state != parsingPeriod {
				return 0, ErrUnexpectedInput
			}

			duration.Days, err = strconv.ParseFloat(num, 64)
			if err != nil {
				return 0, err
			}
			num = ""
		case 'H':
			if state != parsingTime {
				return 0, ErrUnexpectedInput
			}

			duration.Hours, err = strconv.ParseFloat(num, 64)
			if err != nil {
				return 0, err
			}
			num = ""
		case 'S':
			if state != parsingTime {
				return 0, ErrUnexpectedInput
			}

			duration.Seconds, err = strconv.ParseFloat(num, 64)
			if err != nil {
				return 0, err
			}
			num = ""
		default:
			if unicode.IsNumber(char) || char == '.' {
				num += string(char)
				continue
			}

			fmt.Printf("Unknown character: %v\n", char)
			return 0, ErrUnexpectedInput
		}
	}

	timeduration := int64(duration.ToTimeDuration().Seconds())

	return timeduration, nil
}

func (duration *Duration) ToTimeDuration() time.Duration {
	var timeDuration time.Duration

	// zero checks are here to avoid unnecessary math operations, on a durations such as `PT5M`
	if duration.Years != 0 {
		timeDuration += time.Duration(math.Round(duration.Years * nsPerYear))
	}
	if duration.Months != 0 {
		timeDuration += time.Duration(math.Round(duration.Months * nsPerMonth))
	}
	if duration.Weeks != 0 {
		timeDuration += time.Duration(math.Round(duration.Weeks * nsPerWeek))
	}
	if duration.Days != 0 {
		timeDuration += time.Duration(math.Round(duration.Days * nsPerDay))
	}
	if duration.Hours != 0 {
		timeDuration += time.Duration(math.Round(duration.Hours * nsPerHour))
	}
	if duration.Minutes != 0 {
		timeDuration += time.Duration(math.Round(duration.Minutes * nsPerMinute))
	}
	if duration.Seconds != 0 {
		timeDuration += time.Duration(math.Round(duration.Seconds * nsPerSecond))
	}
	if duration.Negative {
		timeDuration = -timeDuration
	}

	return timeDuration
}
