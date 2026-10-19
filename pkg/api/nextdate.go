package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func nextDateHandler(res http.ResponseWriter, req *http.Request) {
	nowStr := req.FormValue("now")
	if nowStr == "" {
		nowStr = time.Now().Format(DateFormat)
	}
	now, err := time.Parse(DateFormat, nowStr)
	if err != nil {
		http.Error(res, "incorrect request: "+err.Error(), http.StatusBadRequest)
		return
	}
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")
	result, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(res, "incorrect request: "+err.Error(), http.StatusBadRequest)
		return
	}
	res.Write([]byte(result))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	dStartTime, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("incorrect start time. %w", err)
	}

	repeatParse := strings.Split(repeat, " ")
	periodType := repeatParse[0]

	var periodNumbers []int
	var periodNumbersSecondary []int

	switch len(repeatParse) {
	case 1:
		if periodType != "y" {
			return "", fmt.Errorf("incorrect repeat time. %s", periodType)
		}
	case 2:
		periodNumbers, err = separateNumbers(repeatParse[1])
		if err != nil {
			return "", fmt.Errorf("incorrect repeat time. %w", err)
		}
	case 3:
		periodNumbers, err = separateNumbers(repeatParse[1])
		if err != nil {
			return "", fmt.Errorf("incorrect repeat time. %w", err)
		}
		periodNumbersSecondary, err = separateNumbers(repeatParse[2])
		if err != nil {
			return "", fmt.Errorf("incorrect repeat time. %w", err)
		}
	default:
		return "", fmt.Errorf("incorrect repeat time. %v", repeatParse)
	}

	after := dStartTime.After(now)
	var newTime time.Time

	switch periodType {
	case "d":
		if periodNumbers[0] > 400 || periodNumbers[0] <= 0 {
			return "", fmt.Errorf("incorrect repeat time. value shoud be > 0 & <= 400. %d", periodNumbers[0])
		}
		if after {
			newTime = dStartTime.AddDate(0, 0, periodNumbers[0])
		} else {
			newTime = dStartTime
			for newTime.Before(now) {
				newTime = newTime.AddDate(0, 0, periodNumbers[0])
			}
		}
		return newTime.Format(DateFormat), nil
	case "m":
		for _, value := range periodNumbers {
			if value > 31 || value < -2 || value == 0 {
				return "", fmt.Errorf("incorrect repeat time. value shoud be in range [-2, 31], excluding 0. %d", value)
			}
		}
		for _, value := range periodNumbersSecondary {
			if value > 12 || value < 1 {
				return "", fmt.Errorf("incorrect repeat time. value shoud be in range [1...12]. %d", value)
			}
		}
		newTime = dStartTime
		if !after {
			newTime = now
		}
		found := false
		for !found {
			newTime = newTime.AddDate(0, 0, 1)
			for _, day := range periodNumbers {
				var lastDay int
				if day < 0 {
					lastDay = getLastDay(newTime).Day() + day + 1
				}
				if day == int(newTime.Day()) || lastDay == int(newTime.Day()) {
					if periodNumbersSecondary == nil {
						found = true
						break
					}
					for _, month := range periodNumbersSecondary {
						if month == int(newTime.Month()) {
							found = true
							break
						}
					}
				}
			}
		}
		return newTime.Format(DateFormat), nil
	case "y":
		if after {
			newTime = dStartTime.AddDate(1, 0, 0)
			return newTime.Format(DateFormat), nil
		}
		newTime = dStartTime
		for newTime.Before(now) {
			newTime = newTime.AddDate(1, 0, 0)
		}
		return newTime.Format(DateFormat), nil
	case "w":
		newTime = dStartTime
		if !after {
			newTime = now
		}
		for _, value := range periodNumbers {
			if value > 7 || value < 1 {
				return "", fmt.Errorf("incorrect repeat time. value shoud be in range [1...7]. %d", value)
			}
		}
		periodNumbers = normalizeWeekdays(periodNumbers)
		for {
			newTime = newTime.AddDate(0, 0, 1)
			for _, value := range periodNumbers {
				if value == int(newTime.Weekday()) {
					return newTime.Format(DateFormat), nil
				}
			}
		}
	default:
		return "", fmt.Errorf("unsupported format. %s", periodType)
	}
}

func separateNumbers(numbers string) ([]int, error) {
	separatedNumberStr := strings.Split(numbers, ",")
	periodNumbers := make([]int, 0, 10)
	for _, value := range separatedNumberStr {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			return make([]int, 0), fmt.Errorf("incorrect repeat time. %w", err)
		}
		periodNumbers = append(periodNumbers, valueInt)
	}
	return periodNumbers, nil
}

func normalizeWeekdays(days []int) []int {
	normalizedSlice := make([]int, 0, 10)
	for _, value := range days {
		if value == 7 {
			normalizedSlice = append(normalizedSlice, 0)
		} else {
			normalizedSlice = append(normalizedSlice, value)
		}
	}
	return normalizedSlice
}

func getLastDay(from time.Time) time.Time {
	return time.Date(from.Year(), from.Month()+1, 0, 0, 0, 0, 0, from.Location())
}
