package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

type Lesson struct {
	Subject     string `json:"subject"`
	Type        string `json:"type"`
	TimeStart   string `json:"time_start"`
	TimeEnd     string `json:"time_end"`
	TeacherName string `json:"teacher_name"`
	Room        string `json:"room"`
	Building    string `json:"building"`
	Group       string `json:"group"`
}

type Day struct {
	DayNumber int      `json:"day_number"`
	Lessons   []Lesson `json:"lessons"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <input_file> <output_file>")
		return
	}

	input := os.Args[1]
	output := os.Args[2]

	yamlContent, err := ioutil.ReadFile(input)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	dayPattern := regexp.MustCompile(`day_number:\s*(\d+)`)
	lessonPattern := regexp.MustCompile(`- subject:\s*(.*)(?:\r?\n)\s*type:\s*(.*)(?:\r?\n)\s*time_start:\s*'(.*)'(?:\r?\n)\s*time_end:\s*'(.*)'(?:\r?\n)\s*teacher_name:\s*(.*)(?:\r?\n)\s*room:\s*'(.*)'(?:\r?\n)\s*building:\s*'(.*)'(?:\r?\n)\s*group:\s*(.*)(?:\r?\n)`)

	days := dayPattern.FindAllStringSubmatch(string(yamlContent), -1)
	lessons := lessonPattern.FindAllStringSubmatch(string(yamlContent), -1)
	var data []Day
	for i, day := range days {
		dayNumber, _ := strconv.Atoi(day[1])
		dayData := Day{DayNumber: dayNumber, Lessons: []Lesson{}}
		t := lessons
		if i == 0 {
			t = t[:2]
		} else {
			t = t[2:]
		}
		for _, lesson := range t {
			dayData.Lessons = append(dayData.Lessons, Lesson{
				Subject:     lesson[1],
				Type:        lesson[2],
				TimeStart:   lesson[3],
				TimeEnd:     lesson[4],
				TeacherName: lesson[5],
				Room:        lesson[6],
				Building:    lesson[7],
				Group:       lesson[8],
			})
		}
		data = append(data, dayData)
	}

	jsonOutput, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error generating JSON: %v\n", err)
		return
	}

	err = ioutil.WriteFile(output, jsonOutput, 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		return
	}
}
