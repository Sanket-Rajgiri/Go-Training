package main

import (
	"testing"
)

type test struct {
	name              string
	filePath          string
	want              int
	expectedError     error
	expectedErrorCode errorCode
}

func Test_checkFile(t *testing.T) {
	tests := []test{
		{
			name:              "test1",
			filePath:          "testFiles/test1.txt",
			expectedError:     nil,
			expectedErrorCode: 0,
		},
		{
			name:              "test2",
			filePath:          "testFiles/test3.txt",
			expectedError:     ErrFileNotExist,
			expectedErrorCode: 1,
		},
		{
			name:              "test3",
			filePath:          "testFiles",
			expectedError:     ErrIsDirectory,
			expectedErrorCode: 21,
		},
	}
	for _, test := range tests {
		errorCode, fileError := checkFile(test.filePath)
		if test.expectedError != fileError {
			t.Errorf("Expected error %v but got %v", test.expectedError, fileError)
		}
		if test.expectedErrorCode != errorCode {
			t.Errorf("Expected errorCode %v  but got %v", test.expectedErrorCode, errorCode)
		}
	}
}
func Test_lineCount(t *testing.T) {
	tests := []test{
		{
			name:     "Test1",
			filePath: "test/test1.txt",
			want:     10,
		},
		{
			name:     "Test1",
			filePath: "test/test2.txt",
			want:     0,
		},
	}

	for _, test := range tests {
		got, err := lineCount(test.filePath)
		if err != nil {
			t.Errorf("Got error while running lineCount : %v", err)
		}
		if got != test.want {
			t.Errorf("Expected %v but got %v", test.want, got)
		}
	}
}

func Test_byteCount(t *testing.T) {
	tests := []test{
		{
			name:     "Test1",
			filePath: "test/test1.txt",
			want:     445,
		},
		{
			name:     "Test1",
			filePath: "test/test2.txt",
			want:     0,
		},
	}
	for _, test := range tests {
		got, err := byteCount(test.filePath)
		if err != nil {
			t.Errorf("Got error while running byteCount : %v", err)
		}
		if got != test.want {
			t.Errorf("Expected %v but got %v", test.want, got)
		}
	}
}

func Test_wordCount(t *testing.T) {
	tests := []test{
		{
			name:     "Test1",
			filePath: "test/test1.txt",
			want:     78,
		},
		{
			name:     "Test1",
			filePath: "test/test2.txt",
			want:     0,
		},
	}
	for _, test := range tests {
		got, err := wordCount(test.filePath)
		if err != nil {
			t.Errorf("Got error while running wordCount : %v", err)
		}
		if got != test.want {
			t.Errorf("Expected %v but got %v", test.want, got)
		}
	}
}
