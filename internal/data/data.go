package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type MigrationDiscussion struct {
	Title        string
	Body         string
	CategorySlug string
	Comments     []MigrationComment
}

type MigrationComment struct {
	Body           string
	AcceptedAnswer bool
}

func GetSODiscussions(dataFile string) []MigrationDiscussion {
	// Open the CSV file and getting the data
	data := getCSVFileData(dataFile)

	// Discussion holding slice
	discussions := []MigrationDiscussion{}

	// Iterate over the data and get the discussions
	for _, row := range data {
		if row[1] == "1" {
			comments := getSOComments(data, row[0], row[11])

			discussions = append(discussions, MigrationDiscussion{
				Title:        row[12],
				Body:         row[4],
				CategorySlug: "q-a",
				Comments:     comments,
			})
		}
	}
	return discussions
}

func getSOComments(data [][]string, postId string, answerId string) []MigrationComment {
	// Comment holding slice
	comments := []MigrationComment{}

	// Iterate over the data and get the comments
	for _, row := range data {
		if row[1] == "2" && row[10] == postId {
			acceptedAnswer := false
			if row[0] == answerId {
				acceptedAnswer = true
			}

			comments = append(comments, MigrationComment{
				Body:           row[4],
				AcceptedAnswer: acceptedAnswer,
			})
		}
	}
	return comments
}

func getCSVFileData(file string) [][]string {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, _ := csvReader.ReadAll()

	// Remove the header row
	if strings.ToLower(data[0][0]) == "id" {
		data = data[1:]
	}

	return data
}
