package internal

import (
	"context"
	"errors"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func UpdateSheet(filepath, sheetId, text, timestamp string) error {
	credential := option.WithCredentialsFile(filepath)
	srv, err := sheets.NewService(context.TODO(), credential)
	if err != nil {
		log.Fatal(err)
		return errors.New("Failed to get service")
	}

	// create ValueRange
	var row []interface{}
	var rows [][]interface{}
	var vr sheets.ValueRange
	row = append(row, text)
	row = append(row, timestamp)
	rows = append(rows, row)
	vr.Values = append(vr.Values, rows...)

	_, err = srv.Spreadsheets.Values.Append(sheetId, "A1", &vr).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Do()
	if err != nil {
		return errors.New("Failed to insert rows")
	}
	return nil
}
