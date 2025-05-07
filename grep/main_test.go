package main

import (
	"reflect"
	"strings"
	"testing"
)

func Test_OpenFile(t *testing.T) {
	type test struct {
		name          string
		filepath      string
		keyword       string
		want          []string
		grepFlagState flagState
		wantErr       bool
	}

	tests := []test{
		{
			name:          "SimpleTest1",
			filepath:      "testFiles/Simple.txt",
			keyword:       "test",
			want:          []string{"This is a test file."},
			grepFlagState: flagState{},
		},
		{
			name:     "Case Sensitive",
			filepath: "testFiles/case-sensitive.txt",
			keyword:  "Error",
			want:     []string{"Error: Something failed"},
		},
		{
			name:     "Case Insensitive",
			filepath: "testFiles/case-sensitive.txt",
			keyword:  "error",
			want: []string{
				"Error: Something failed",
				"error: case-insensitive test",
			},
			grepFlagState: flagState{caseInsensitive: true},
		},
		{
			name:          "Invert Match",
			filepath:      "testFiles/invert-match.txt",
			keyword:       "grapefruit",
			grepFlagState: flagState{invertMatch: true},
			want: []string{
				"apple",
				"banana",
				"pineapple",
				"mango",
			},
		},
		{
			name:     "Multiline",
			filepath: "testFiles/multiline.txt",
			keyword:  "line",
			want: []string{
				"This is line one",
				"This is line two",
				"This is line three",
				"This is line four",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := openFile(tt.filepath, tt.keyword, &tt.grepFlagState)
			if (err != nil) != tt.wantErr {
				t.Errorf("search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stdInput(t *testing.T) {
	type test struct {
		name          string
		inputString   string
		keyword       string
		grepFlagState flagState
		wantErr       bool
		want          []string
	}
	tests := []test{
		{
			name:          "SimpleTest",
			inputString:   "Hello World\nThis is a test file.\nGo is awesome.\nLet's grep some lines!\n",
			keyword:       "test",
			grepFlagState: flagState{caseInsensitive: false, invertMatch: false},
			wantErr:       false,
			want: []string{
				"This is a test file.",
			},
		},
		{
			name:          "case sensitive",
			inputString:   "Error: Something failed\nWarning: Low disk space\nInfo: All systems operational\nerror: case-insensitive test\n",
			keyword:       "Error",
			grepFlagState: flagState{caseInsensitive: false, invertMatch: false},
			wantErr:       false,
			want:          []string{"Error: Something failed"},
		},
		{
			name:          "case insensitive",
			inputString:   "Error: Something failed\nWarning: Low disk space\nInfo: All systems operational\nerror: case-insensitive test\n",
			keyword:       "error",
			grepFlagState: flagState{caseInsensitive: true, invertMatch: false},
			wantErr:       false,
			want: []string{
				"Error: Something failed",
				"error: case-insensitive test",
			},
		},
		{
			name:          "Invert Match",
			inputString:   "apple\nbanana\ngrapefruit\npineapple\nmango\n",
			keyword:       "grapefruit",
			grepFlagState: flagState{invertMatch: true},
			want: []string{
				"apple",
				"banana",
				"pineapple",
				"mango",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.inputString)
			got, err := fileProcessor(reader, tt.keyword, tt.grepFlagState.caseInsensitive, tt.grepFlagState.invertMatch)
			if (err != nil) != tt.wantErr {
				t.Errorf("search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("search() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_processStdin(t *testing.T) {
	type test struct {
		name          string
		inputString   string
		keyword       string
		grepFlagState flagState
		wantErr       bool
		want          []string
	}
	tests := []test{
		{
			name:          "Simple Match",
			inputString:   "Hello World\nThis is a test file.\nGo is awesome.\nLet's grep some lines!\n",
			keyword:       "test",
			grepFlagState: flagState{caseInsensitive: false, invertMatch: false},
			wantErr:       false,
			want: []string{
				"This is a test file.",
			},
		},
		{
			name:          "Case Sensitive Match",
			inputString:   "Error: Something failed\nWarning: Low disk space\nInfo: All systems operational\nerror: case-insensitive test\n",
			keyword:       "Error",
			grepFlagState: flagState{caseInsensitive: false, invertMatch: false},
			wantErr:       false,
			want:          []string{"Error: Something failed"},
		},
		{
			name:          "Case Insensitive Match",
			inputString:   "Error: Something failed\nWarning: Low disk space\nInfo: All systems operational\nerror: case-insensitive test\n",
			keyword:       "error",
			grepFlagState: flagState{caseInsensitive: true, invertMatch: false},
			wantErr:       false,
			want: []string{
				"Error: Something failed",
				"error: case-insensitive test",
			},
		},
		{
			name:          "Invert Match",
			inputString:   "apple\nbanana\ngrapefruit\npineapple\nmango\n",
			keyword:       "grapefruit",
			grepFlagState: flagState{invertMatch: true},
			wantErr:       false,
			want: []string{
				"apple",
				"banana",
				"pineapple",
				"mango",
			},
		},
		{
			name:          "No Match",
			inputString:   "apple\nbanana\ncherry\n",
			keyword:       "grapefruit",
			grepFlagState: flagState{caseInsensitive: false, invertMatch: false},
			wantErr:       false,
			want:          []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.inputString)
			writer := &strings.Builder{}

			err := processStdin(reader, tt.keyword, tt.grepFlagState, writer)
			if (err != nil) != tt.wantErr {
				t.Errorf("processStdin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := strings.Split(strings.TrimSpace(writer.String()), "\n")
			if len(tt.want) == 0 && writer.Len() == 0 {
				// Match
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processStdin() = %v, want %v", got, tt.want)
			}
		})
	}
}
