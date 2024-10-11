package store_test

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/goodsign/monday"
	"github.com/themulle/chronokeyaccess/internal/store"
	"github.com/themulle/chronokeyaccess/pkg/codemanager"
)

func TestLoadPersonalCodeCSV(t *testing.T) {
	// Setup test CSV file
	testFilePath := "test_personal_codes.csv"
	content := `Name,PinCode,Slot
Alice,1234,2024-2030,1
Bob,5678,*
Charlie,9101,Mo,2024
Peter,9101,"14:00-16:00 Mo,Mi",`
	
	if err := os.WriteFile(testFilePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test CSV file: %v", err)
	}
	defer os.Remove(testFilePath) // Clean up after the test

	tests := []struct {
		filePath   string
		expectErr  bool
		expected   codemanager.PersonalCodes
	}{
		{
			filePath: testFilePath,
			expectErr: false,
			expected: codemanager.PersonalCodes{
				{Name: "Alice", PinCode: 1234, CronString: "0 0 0 * * * 2024-2030", Duration: time.Hour*24-1},
				{Name: "Bob", PinCode: 5678, CronString: "0 0 0 * * * *", Duration: time.Hour*24-1},
				{Name: "Charlie", PinCode: 9101, CronString: "0 0 0 * 1 * 2024", Duration: time.Hour},
				{Name: "Peter", PinCode: 9101, CronString: "0 00 14 * * 1,3 *", Duration: 2*time.Hour},
			},
		},
		{
			filePath: "non_existing_file.csv",
			expectErr: true,
		},
		{
			filePath: "invalid_format.csv",
			expectErr: true,
		},
	}

	// Create an invalid CSV file for the test
	os.WriteFile("invalid_format.csv", []byte("TestInvalid,Format\n"), 0644)
	defer os.Remove("invalid_format.csv") // Clean up after the test

	for _, tt := range tests {
		t.Run(tt.filePath, func(t *testing.T) {
			codes, err := store.LoadPersonalCodeCSV(tt.filePath, monday.LocaleDeDE)

			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			fmt.Println(err)
			if !tt.expectErr && !comparePersonalCodes(codes, tt.expected) {
				t.Errorf("expected: %+v, got: %+v", tt.expected, codes)
			}
		})
	}
}

func TestWritePersonalCodeCSV(t *testing.T) {
	testFilePath := "output_personal_codes.csv"
	codes := codemanager.PersonalCodes{
		{Name: "Alice", PinCode: 1234, CronString: "Slot1", Duration: time.Hour},
		{Name: "Bob", PinCode: 5678, CronString: "Slot2",  Duration: time.Duration(2.5*float64(time.Hour))},
	}

	// Write the codes to a CSV file
	if err := store.WritePersonalCodeCSV(codes, testFilePath); err != nil {
		t.Fatalf("failed to write CSV: %v", err)
	}
	defer os.Remove(testFilePath) // Clean up after the test

	// Verify the content of the CSV file
	file, err := os.Open(testFilePath)
	if err != nil {
		t.Fatalf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV file: %v", err)
	}

	expectedHeader := []string{"Name", "PinCode", "Slot", "Dauer"}
	if len(records) != 3 || !compareSlices(records[0], expectedHeader) {
		t.Errorf("expected header to be %v, got %v", expectedHeader, records[0])
	}

	expectedRecords := [][]string{
		{"Alice", "1234", "Slot1", "1"},
		{"Bob", "5678", "Slot2", "2.5"},
	}

	for i, record := range expectedRecords {
		if !compareSlices(records[i+1], record) {
			t.Errorf("expected record %d to be %v, got %v", i, record, records[i+1])
		}
	}
}

// Helper function to compare PersonalCodes slices
func comparePersonalCodes(a, b codemanager.PersonalCodes) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Helper function to compare string slices
func compareSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
