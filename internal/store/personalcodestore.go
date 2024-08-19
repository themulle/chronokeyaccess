package store

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/themulle/chronokeyaccess/pkg/codemanager"
)

// Parse liest die CSV-Datei und gibt eine Liste von codemanager.PersonalCodes zurück
func LoadPersonalCodeCSV(filePath string) (codemanager.PersonalCodes, error) {
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

	var codes codemanager.PersonalCodes
	for _, record := range records {
		if len(record) != 3 {
			return nil, fmt.Errorf("invalid record: %v", record)
		}

		pinCode, err := strconv.ParseUint(record[1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid PinCode in record: %v", record)
		}

		codes = append(codes, codemanager.PersonalCode{
			Name:     record[0],
			PinCode:  uint(pinCode),
			SlotName: record[2],
		})
	}

	return codes, nil
}

func WritePersonalCodeCSV(codes codemanager.PersonalCodes, outputFilePath string) error {
	file, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("unable to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, code := range codes {
		record := []string{
			code.Name,
			strconv.FormatUint(uint64(code.PinCode), 10),
			code.SlotName,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("unable to write record to CSV file: %w", err)
		}
	}

	return nil
}
