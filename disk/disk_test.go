package disk

import (
	"os"
	"testing"
)

func TestJsonFile(t *testing.T) {
	// test data
	type testStruct struct {
		A int
		B string
		C []int
	}
	testData := testStruct{
		A: 1,
		B: "test",
		C: []int{1, 2, 3},
	}
	// write test data to file
	err := WriteJsonFile[testStruct](testData, "../testdata/test.json", 0644)
	if err != nil {
		t.Error("Error writing file: ", err)
	}
	// read test data from file
	var readData testStruct
	err = ReadJsonFile[testStruct]("../testdata/test.json", &readData)
	if err != nil {
		t.Error("Error reading file: ", err)
	}
	// compare test data to read data
	if testData.A != readData.A {
		t.Errorf("TestJsonFile A = %d; want %d", readData.A, testData.A)
	}
	if testData.B != readData.B {
		t.Errorf("TestJsonFile B = %s; want %s", readData.B, testData.B)
	}
	if len(testData.C) != len(readData.C) {
		t.Errorf("TestJsonFile C = %d; want %d", len(readData.C), len(testData.C))
	}
	for i := range testData.C {
		if testData.C[i] != readData.C[i] {
			t.Errorf("TestJsonFile C[%d] = %d; want %d", i, readData.C[i], testData.C[i])
		}
	}

	// clean up
	err = os.Remove("../testdata/test.json")
	if err != nil {
		t.Error("Error removing file: ", err)
	}
}
