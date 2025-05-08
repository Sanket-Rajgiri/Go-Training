package main

import (
	"bufio"
	"os"
	"reflect"
	"sort"
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
		wantErr       error
	}

	tests := []test{
		{
			name:     "Non-Existent File",
			filepath: "testFiles/sanket.txt",
			wantErr:  ErrFileNotExist,
		},
		{
			name:     "Folder Path",
			filepath: "testFiles/",
			wantErr:  ErrIsDirectory,
		},
		{
			name:     "Folder Path",
			filepath: "testFiles/do-not-read.txt",
			wantErr:  ErrPermissionDenied,
		},
		{
			name:          "SimpleTest1",
			filepath:      "testFiles/Simple.txt",
			keyword:       "test",
			want:          []string{"This is a test file."},
			grepFlagState: flagState{},
		},
		{
			name:     "Case Sensitive",
			filepath: "testFiles/case-insensitive.txt",
			keyword:  "Error",
			want:     []string{"Error: Something failed"},
		},
		{
			name:     "Case Insensitive",
			filepath: "testFiles/case-insensitive.txt",
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
		{
			name:          "After Context",
			filepath:      "testFiles/context.txt",
			keyword:       "random",
			grepFlagState: flagState{afterContext: 1, caseInsensitive: true},
			want: []string{
				"Line 2: Random text here",
				"Line 3: Another random line",
				"Line 4: This line contains the keyword match",
			},
		},
		{
			name:          "Before Context",
			filepath:      "testFiles/context.txt",
			keyword:       "random",
			grepFlagState: flagState{beforeContext: 1, caseInsensitive: true},
			want: []string{
				"Line 1: Introduction",
				"Line 2: Random text here",
				"Line 3: Another random line",
			},
		},
		{
			name:          "Before Context and After Context",
			filepath:      "testFiles/context.txt",
			keyword:       "random",
			grepFlagState: flagState{beforeContext: 1, afterContext: 1, caseInsensitive: true},
			want: []string{
				"Line 1: Introduction",
				"Line 2: Random text here",
				"Line 3: Another random line",
				"Line 4: This line contains the keyword match",
			},
		},
		{
			name:          "Before Context and After Context 2 ",
			filepath:      "testFiles/context2.txt",
			keyword:       "keyword",
			grepFlagState: flagState{caseInsensitive: true, beforeContext: 1, afterContext: 1},
			want: []string{
				"Line 2: Still nothing here",
				"Line 3: The keyword appears now",
				"Line 4: Follows the match",
				"Line 5: More irrelevant text",
				"Line 6: keyword is here again",
				"Line 7: Follow-up line",
			},
		},
		{
			name:          "Count",
			filepath:      "testFiles/Simple.txt",
			keyword:       "some",
			grepFlagState: flagState{caseInsensitive: true, count: true},
			want:          []string{"2"},
		},
		{
			name:          "Count with inverse",
			filepath:      "testFiles/case-insensitive.txt",
			keyword:       "some",
			grepFlagState: flagState{caseInsensitive: true, count: true, invertMatch: true},
			want:          []string{"3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := openFile(tt.filepath, tt.keyword, &tt.grepFlagState)
			if err != tt.wantErr {
				t.Errorf(" openFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("openFile(), got = %v, want %v", got, tt.want)
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
		{
			name:          "After Context",
			inputString:   "Line 1: This line has nothing\nLine 2: Still nothing here\nLine 3: The keyword appears now\nLine 4: Follows the match\nLine 5: More irrelevant text\nLine 6: keyword is here again\nLine 7: Follow-up line\nLine 8: The end of file\n",
			keyword:       "keyword",
			grepFlagState: flagState{caseInsensitive: true, afterContext: 1},
			want: []string{
				"Line 3: The keyword appears now",
				"Line 4: Follows the match",
				"Line 6: keyword is here again",
				"Line 7: Follow-up line",
			},
		},
		{
			name:          "Before Context",
			inputString:   "Line 1: This line has nothing\nLine 2: Still nothing here\nLine 3: The keyword appears now\nLine 4: Follows the match\nLine 5: More irrelevant text\nLine 6: keyword is here again\nLine 7: Follow-up line\nLine 8: The end of file\n",
			keyword:       "keyword",
			grepFlagState: flagState{caseInsensitive: true, beforeContext: 1},
			want: []string{
				"Line 2: Still nothing here",
				"Line 3: The keyword appears now",
				"Line 5: More irrelevant text",
				"Line 6: keyword is here again",
			},
		},
		{
			name:          "Before Context and After Context",
			inputString:   "Line 1: This line has nothing\nLine 2: Still nothing here\nLine 3: The keyword appears now\nLine 4: Follows the match\nLine 5: More irrelevant text\nLine 6: keyword is here again\nLine 7: Follow-up line\nLine 8: The end of file\n",
			keyword:       "keyword",
			grepFlagState: flagState{caseInsensitive: true, beforeContext: 1, afterContext: 1},
			want: []string{
				"Line 2: Still nothing here",
				"Line 3: The keyword appears now",
				"Line 4: Follows the match",
				"Line 5: More irrelevant text",
				"Line 6: keyword is here again",
				"Line 7: Follow-up line",
			},
		},
		{
			name:          "Count",
			inputString:   "Error: Something failed\nWarning: Low disk space\nInfo: All systems operational\nerror: case-insensitive test\n",
			keyword:       "error",
			grepFlagState: flagState{caseInsensitive: true, count: true},
			want:          []string{"2"},
		},
		{
			name:          "Count with inverse",
			inputString:   "apple\nbanana\ngrapefruit\npineapple\nmango\n",
			keyword:       "pineapple",
			grepFlagState: flagState{caseInsensitive: true, count: true, invertMatch: true},
			want:          []string{"4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.inputString)
			writer := &strings.Builder{}

			err := processStdin(reader, tt.keyword, &tt.grepFlagState, writer)
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
				t.Errorf("processStdin() , got= %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_recursiveFileList(t *testing.T) {
	type test struct {
		name      string
		inputPath string
		wantErr   error
		want      []string
	}
	tests := []test{
		{
			name:      "Path Exists",
			inputPath: "testFiles",
			want: []string{
				"testFiles/Simple.txt",
				"testFiles/context.txt",
				"testFiles/invert-match.txt",
				"testFiles/case-insensitive.txt",
				"testFiles/context2.txt",
				"testFiles/multiline.txt",
				"testFiles/do-not-read.txt",
			},
		},
		{
			name:      "Path Not Exists",
			inputPath: "sanket",
			want:      []string{"sanket"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := recursiveFileList(tt.inputPath)
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("recursiveFile(), got=%v, want=%v\n", got, tt.want)
			}
		})
	}
}

func Test_writeToFile(t *testing.T) {
	type test struct {
		filename    string
		inputString []string
		wantErr     error
	}
	// tests := []test{
	// 	{
	// 		name:        "File Create and Write",
	// 		filename:    "testFiles/sanket.txt",
	// 		inputString: []string{"line1", "line2", "line3"},
	// 	},
	// 	{
	// 		name:        "FileAlreadyExists",
	// 		filename:    "testFiles/sanket.txt",
	// 		wantErr:     ErrFileExists,
	// 		inputString: []string{},
	// 	},
	// }
	t.Run("File Create and Write", func(t *testing.T) {
		tt := test{
			filename:    "sanket.txt",
			inputString: []string{"line1", "line2", "line3"},
		}
		err := writeToFile(tt.filename, tt.inputString)
		defer os.Remove(tt.filename)
		if err != tt.wantErr {
			t.Errorf("writeToFile() gotErr=%v, expectedErr=%v", err, tt.wantErr)
		}
		file, osOpenErr := os.Open(tt.filename)
		if err != nil {
			t.Errorf("writeToFile() gotErr=%v", osOpenErr)
		}
		var fileReader []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fileReader = append(fileReader, scanner.Text())
		}
		if scannerErr := scanner.Err(); scannerErr != nil {
			t.Errorf("writeToFile() gotErr=%v", scannerErr)
		}
		sort.Strings(tt.inputString)
		sort.Strings(fileReader)
		if !reflect.DeepEqual(tt.inputString, fileReader) {
			t.Errorf("writeToFile() got=%v, want=%v", fileReader, tt.inputString)
		}

	})
	t.Run("File Already Exists", func(t *testing.T) {
		tt := test{
			filename: "testFiles/Simple.txt",
			// inputString: []string{"line1", "line2", "line3"},
			wantErr: ErrFileExists,
		}
		err := writeToFile(tt.filename, tt.inputString)
		if err != tt.wantErr {
			t.Errorf("writeToFile() gotErr=%v, expectedErr=%v", err, tt.wantErr)
		}
	})
}
