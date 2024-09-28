package store

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/themulle/chronokeyaccess/pkg/accesslog"
	"github.com/themulle/chronokeyaccess/pkg/dateparser"
)

// Parse liest die CSV-Datei und gibt eine Liste von codemanager.PersonalCodes zurück
func LoadAccessLogCSV(filePath string) (accesslog.AccessLogs, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var logs accesslog.AccessLogs
	for _, record := range records {
		if len(record) < 2 {
			return nil, fmt.Errorf("invalid record: %v", record)
		}

		ts, err := dateparser.Parse(record[0])
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp in record: %v", record)
		}

		pinCode, err := strconv.ParseUint(record[1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid PinCode in record: %v", record)
		}

		var status string
		if len(record) >= 3 {
			status = record[2]
		}

		logs = append(logs, accesslog.AccessLog{
			Ts:      ts,
			PinCode: uint(pinCode),
			Status:  status,
		})
	}

	return logs, nil
}

func WriteAccessLogCSV(logs accesslog.AccessLogs, outputFilePath string) error {
	file, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("unable to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, accesslog := range logs {
		record := []string{
			accesslog.Ts.String(),
			strconv.FormatUint(uint64(accesslog.PinCode), 10),
			accesslog.Status,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("unable to write record to CSV file: %w", err)
		}
	}

	return nil
}
