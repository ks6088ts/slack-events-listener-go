package internal

import (
	"context"
	"errors"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// SheetClient denotes a Sheets API client
type SheetClient struct {
	service *sheets.Service
	sheetId string
}

// NewSheetClient create a Sheets API client
func NewSheetClient(filepath, sheetId string) (*SheetClient, error) {
	credential := option.WithCredentialsFile(filepath)
	service, err := sheets.NewService(context.TODO(), credential)
	if err != nil {
		return nil, errors.New("Failed to get service")
	}
	return &SheetClient{
		service: service,
		sheetId: sheetId,
	}, nil
}

// AppendValues appends arguments to a sheet
func (client *SheetClient) AppendValues(items ...string) error {
	// create ValueRange
	var row []interface{}
	var rows [][]interface{}
	var vr sheets.ValueRange
	for _, item := range items {
		row = append(row, item)
	}
	rows = append(rows, row)
	vr.Values = append(vr.Values, rows...)

	_, err := client.service.Spreadsheets.Values.Append(client.sheetId, "A1", &vr).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Do()
	if err != nil {
		return errors.New("Failed to insert rows")
	}
	return nil
}
