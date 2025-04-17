package main

import (
	"testing"
)

type test struct {
	name     string
	filePath string
	want     int
}

func Test_lineCount(t *testing.T) {
	tests := []test{
		{
			name:     "Test1",
			filePath: "test1.txt",
			want:     10,
		},
		{
			name:     "Test1",
			filePath: "test2.txt",
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
			filePath: "test1.txt",
			want:     445,
		},
		{
			name:     "Test1",
			filePath: "test2.txt",
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
			filePath: "test1.txt",
			want:     78,
		},
		{
			name:     "Test1",
			filePath: "test2.txt",
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
