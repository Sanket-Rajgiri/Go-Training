package main

import (
	"testing"
)

type test struct {
	name              string
	filePath          string
	lineCount         int
	wordCount         int
	byteCount         int
	expectedError     error
	expectedErrorCode errorCode
}

var tests = []test{
	{
		name:              "test1",
		filePath:          "testFiles/test1.txt",
		expectedError:     nil,
		lineCount:         10,
		wordCount:         78,
		byteCount:         445,
		expectedErrorCode: 0,
	},
	{
		name:              "test2",
		filePath:          "testFiles/test2.txt",
		expectedError:     nil,
		lineCount:         0,
		wordCount:         0,
		byteCount:         0,
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
	{
		lineCount: 4032,
		wordCount: 24302,
		byteCount: 145435,
		filePath:  "testFiles/shakespeare-db/All's Well That Ends Well.txt",
	},
	{
		lineCount: 5058,
		wordCount: 26894,
		byteCount: 169548,
		filePath:  "testFiles/shakespeare-db/Antony and Cleopatra.txt",
	},
	{
		lineCount: 3640,
		wordCount: 22765,
		byteCount: 134602,
		filePath:  "testFiles/shakespeare-db/As You Like It.txt",
	},
	{
		lineCount: 2673,
		wordCount: 16129,
		byteCount: 96686,
		filePath:  "testFiles/shakespeare-db/Comedy of Errors.txt",
	},
	{
		lineCount: 5132,
		wordCount: 29147,
		byteCount: 180992,
		filePath:  "testFiles/shakespeare-db/Coriolanus.txt",
	},
	{
		lineCount: 4842,
		wordCount: 28775,
		byteCount: 177616,
		filePath:  "testFiles/shakespeare-db/Cymbeline.txt",
	},
	{
		lineCount: 5403,
		wordCount: 32062,
		byteCount: 196392,
		filePath:  "testFiles/shakespeare-db/Hamlet.txt",
	},
	{
		lineCount: 3999,
		wordCount: 25943,
		byteCount: 155693,
		filePath:  "testFiles/shakespeare-db/Henry IV, part 1.txt",
	},
	{
		lineCount: 4285,
		wordCount: 27747,
		byteCount: 168615,
		filePath:  "testFiles/shakespeare-db/Henry IV, part 2.txt",
	},
	{
		lineCount: 4154,
		wordCount: 27424,
		byteCount: 166207,
		filePath:  "testFiles/shakespeare-db/Henry V.txt",
	},
	{
		lineCount: 3674,
		wordCount: 22811,
		byteCount: 142855,
		filePath:  "testFiles/shakespeare-db/Henry VI, part 1.txt",
	},
	{
		lineCount: 4157,
		wordCount: 26729,
		byteCount: 163176,
		filePath:  "testFiles/shakespeare-db/Henry VI, part 2.txt",
	},
	{
		lineCount: 3983,
		wordCount: 25874,
		byteCount: 157821,
		filePath:  "testFiles/shakespeare-db/Henry VI, part 3.txt",
	},
	{
		lineCount: 4188,
		wordCount: 25889,
		byteCount: 159232,
		filePath:  "testFiles/shakespeare-db/Henry VIII.txt",
	},
	{
		lineCount: 3587,
		wordCount: 20787,
		byteCount: 126454,
		filePath:  "testFiles/shakespeare-db/Julius Caesar.txt",
	},
	{
		lineCount: 3328,
		wordCount: 21689,
		byteCount: 131639,
		filePath:  "testFiles/shakespeare-db/King John.txt",
	},
	{
		lineCount: 4846,
		wordCount: 27642,
		byteCount: 169824,
		filePath:  "testFiles/shakespeare-db/King Lear.txt",
	},
	{
		lineCount: 4045,
		wordCount: 22850,
		byteCount: 140390,
		filePath:  "testFiles/shakespeare-db/Love's Labour's Lost.txt",
	},
	{
		lineCount: 3251,
		wordCount: 18164,
		byteCount: 113189,
		filePath:  "testFiles/shakespeare-db/Macbeth.txt",
	},
	{
		lineCount: 3907,
		wordCount: 23071,
		byteCount: 140306,
		filePath:  "testFiles/shakespeare-db/Measure for Measure.txt",
	},
	{
		lineCount: 3448,
		wordCount: 22131,
		byteCount: 131655,
		filePath:  "testFiles/shakespeare-db/Merchant of Venice.txt",
	},
	{
		lineCount: 3862,
		wordCount: 23569,
		byteCount: 140721,
		filePath:  "testFiles/shakespeare-db/Merry Wives of Windsor.txt",
	},
	{
		lineCount: 2809,
		wordCount: 17074,
		byteCount: 104100,
		filePath:  "testFiles/shakespeare-db/Midsummer Night's Dream.txt",
	},
	{
		lineCount: 3692,
		wordCount: 22473,
		byteCount: 132701,
		filePath:  "testFiles/shakespeare-db/Much Ado About Nothing.txt",
	},
	{
		lineCount: 4950,
		wordCount: 27784,
		byteCount: 168957,
		filePath:  "testFiles/shakespeare-db/Othello.txt",
	},
	{
		lineCount: 3301,
		wordCount: 19481,
		byteCount: 119588,
		filePath:  "testFiles/shakespeare-db/Pericles.txt",
	},
	{
		lineCount: 3515,
		wordCount: 23837,
		byteCount: 144558,
		filePath:  "testFiles/shakespeare-db/Richard II.txt",
	},
	{
		lineCount: 5052,
		wordCount: 31300,
		byteCount: 192148,
		filePath:  "testFiles/shakespeare-db/Richard III.txt",
	},
	{
		lineCount: 4164,
		wordCount: 25712,
		byteCount: 154968,
		filePath:  "testFiles/shakespeare-db/Romeo and Juliet.txt",
	},
	{
		lineCount: 3712,
		wordCount: 22021,
		byteCount: 133380,
		filePath:  "testFiles/shakespeare-db/Taming of the Shrew.txt",
	},
	{
		lineCount: 3066,
		wordCount: 17345,
		byteCount: 107428,
		filePath:  "testFiles/shakespeare-db/The Tempest.txt",
	},
	{
		lineCount: 3444,
		wordCount: 19571,
		byteCount: 121471,
		filePath:  "testFiles/shakespeare-db/Timon of Athens.txt",
	},
	{
		lineCount: 3324,
		wordCount: 21655,
		byteCount: 133080,
		filePath:  "testFiles/shakespeare-db/Titus Andronicus.txt",
	},
	{
		lineCount: 4860,
		wordCount: 27440,
		byteCount: 171091,
		filePath:  "testFiles/shakespeare-db/Troiles and Cressida.txt",
	},
	{
		lineCount: 3577,
		wordCount: 21362,
		byteCount: 125638,
		filePath:  "testFiles/shakespeare-db/Twelfth Night.txt",
	},
	{
		lineCount: 3220,
		wordCount: 18199,
		byteCount: 109549,
		filePath:  "testFiles/shakespeare-db/Two Gentlemen of Verona.txt",
	},
}

func Test_openFile(t *testing.T) {
	for _, test := range tests {
		_, errorCode, fileError := openFile(test.filePath)
		if test.expectedError != fileError {
			t.Errorf("Expected error %v but got %v", test.expectedError, fileError)
		}
		if test.expectedErrorCode != errorCode {
			t.Errorf("Expected errorCode %v  but got %v", test.expectedErrorCode, errorCode)
		}
	}
}

func Test_lineCount(t *testing.T) {
	for _, test := range tests {
		file, _, _ := openFile(test.filePath)
		got, err := lineCount(&file)
		if err != nil {
			t.Errorf("Got error while running lineCount : %v", err)
		}
		if got != test.lineCount {
			t.Errorf("Expected %v but got %v", test.lineCount, got)
		}
	}
}
func Test_byteCount(t *testing.T) {
	for _, test := range tests {
		file, _, _ := openFile(test.filePath)
		got, err := byteCount(&file)
		if err != nil {
			t.Errorf("Got error while running byteCount : %v", err)
		}
		if got != test.byteCount {
			t.Errorf("Expected %v but got %v", test.byteCount, got)
		}
	}
}

func Test_wordCount(t *testing.T) {
	for _, test := range tests {
		file, _, _ := openFile(test.filePath)
		got, err := wordCount(&file)
		if err != nil {
			t.Errorf("Got error while running wordCount : %v", err)
		}
		if got != test.wordCount {
			t.Errorf("Expected %v but got %v", test.wordCount, got)
		}
	}
}

// func BechMark(b tes)
