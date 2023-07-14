package storage

import (
	"encoding/csv"
	"engine-db/entity"
	"io"
	"log"
	"os"
	"strings"
)

// read csv
func ReadCsv() []entity.User {
	// Open the CSV file
	filePath := "storage/data-user.csv"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// create a CSV reader
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	users := []entity.User{}
	num := 1

	// iterate over CSV rows
	for {
		// read 1 per 1 of rows
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		// skip row number 1 (using for header table)
		if num != 1 {
			// TODO: Fix parse Join Date

			role := strings.ToLower(row[15])
			if row[15] == "" {
				role = "employee"
			}

			users = append(users, entity.User{
				Name:    row[0],
				PrivyID: row[1],
				Email:   row[2],
				Status:  row[3],
				// JoinDate:               t,
				JobTitle:               row[5],
				Level:                  row[6],
				Directorate:            row[7],
				Division:               row[8],
				Homebase:               row[9],
				DirectLeader:           row[10],
				DirectLeaderJobTitle:   row[11],
				DirectLeaderEmployeeID: row[12],
				PICHrbp:                row[13],
				HrbpPrivyID:            row[14],
				Role:                   role,
			})
		}

		num++
	}

	return users
}
