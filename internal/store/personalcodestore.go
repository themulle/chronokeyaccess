package store

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/goodsign/monday"
	"github.com/themulle/chronokeyaccess/pkg/codemanager"
	"github.com/themulle/chronokeyaccess/pkg/cronslotstring"
)

var PersonalCodeCsvHeader []string = []string{"Name", "PinCode", "Slot"}

// Parse liest die CSV-Datei und gibt eine Liste von codemanager.PersonalCodes zurück
func LoadPersonalCodeCSV(filePath string, locale monday.Locale) (codemanager.PersonalCodes, error) {
	var combinedError error
	
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

	csg:=cronslotstring.NewCronSlotGenerator()
	csg.SetLocale(locale)

	var codes codemanager.PersonalCodes
	for line, record := range records {
		if len(record) < 3 {
			combinedError=errors.Join(combinedError, fmt.Errorf("invalid record in line %d: %v", line, record))
			continue //skip line
		}
		//trim spaces
		for i := range PersonalCodeCsvHeader {
			record[i]=strings.TrimSpace(record[i])
		}

		if slices.Equal(record[0:3],PersonalCodeCsvHeader) {
			continue  //skip header
		}

		//pin-Code check
		pinCode, err := strconv.ParseUint(record[1], 10, 32)
		if err != nil {
			combinedError=errors.Join(combinedError, fmt.Errorf("invalid PinCode in line %d: %v", line, record))
			continue  //skip line
		}

		cronString:=""
		var duration time.Duration
			cronString, duration, err = csg.ParseCronSlotString(strings.Join(record[2:]," "),)
			if err!=nil {
				combinedError=errors.Join(err)
				continue //skip line
			}

		
		codes = append(codes, codemanager.PersonalCode{
			Name:     record[0],
			PinCode:  uint(pinCode),
			CronString: cronString,
			Duration: duration,
		})
	}

	return codes, combinedError
}

func WritePersonalCodeCSV(codes codemanager.PersonalCodes, outputFilePath string) error {
	file, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("unable to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(PersonalCodeCsvHeader); err != nil {
		return fmt.Errorf("unable to write record to CSV file: %w", err)
	}

	for _, code := range codes {
		record := []string{
			code.Name,
			strconv.FormatUint(uint64(code.PinCode), 10),
			code.CronString,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("unable to write record to CSV file: %w", err)
		}
	}

	return nil
}
