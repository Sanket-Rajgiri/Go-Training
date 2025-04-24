package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
)

type errorCode int

const (
	fileNotExistErrorCode         errorCode = 1
	filePermissionDeniedErrorCode errorCode = 2
	IsDirectoryErrorCode          errorCode = 21
	unexpectedErrorCode           errorCode = 125
	invalidFlagsErrorCode         errorCode = 22
	programName                             = "wc"
)

var (
	ErrIsDirectory      = errors.New("is a directory")
	ErrPermissionDenied = errors.New("permission denied")
	ErrFileNotExist     = errors.New("no such file or directory")
	ErrMissingArguments = errors.New("missing arugments")
	ErrInvalidFlags     = errors.New("invalid flags passed")
)

type flagState struct {
	countLines bool
	countWords bool
	countBytes bool
}

func countWithSplit(filepath string, split bufio.SplitFunc) (int, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(split)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return count, nil
}

func lineCount(filepath string) (int, error) {
	return countWithSplit(filepath, bufio.ScanLines)
}

func wordCount(filepath string) (int, error) {
	return countWithSplit(filepath, bufio.ScanWords)
}

func byteCount(filepath string) (int, error) {
	return countWithSplit(filepath, bufio.ScanBytes)
}

func flagParser() (flagState, []string, errorCode, error) {
	lineFlag := flag.Bool("l", false, "count lines")
	wordFlag := flag.Bool("w", false, "count words")
	byteFlag := flag.Bool("c", false, "count bytes")
	var wcflagState flagState
	var fileList []string
	flag.Parse()
	fileList = flag.Args()

	if !flag.Parsed() {
		return wcflagState, fileList, invalidFlagsErrorCode, ErrInvalidFlags
	}

	if len(fileList) < 1 {
		return wcflagState, fileList, invalidFlagsErrorCode, ErrMissingArguments
	}

	if !*lineFlag && !*wordFlag && !*byteFlag {
		wcflagState = flagState{countLines: true, countWords: true, countBytes: true}
	} else {
		wcflagState = flagState{countLines: *lineFlag, countWords: *wordFlag, countBytes: *byteFlag}
	}

	return wcflagState, fileList, 0, nil
}

func checkFile(filepath string) (errorCode errorCode, err error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return fileNotExistErrorCode, ErrFileNotExist
		} else if os.IsPermission(err) {
			return filePermissionDeniedErrorCode, ErrPermissionDenied
		} else {
			return unexpectedErrorCode, err
		}
	}
	if fileInfo.IsDir() {
		return IsDirectoryErrorCode, ErrIsDirectory
	}
	return 0, nil
}

func errorHandler(filepath string, err error) {
	if len(filepath) > 0 {
		fmt.Fprintf(os.Stderr, "%s: %s: %3s\n", programName, filepath, err)
	} else {
		fmt.Fprintf(os.Stderr, "%s: %3s\n", programName, err)
	}
}

func countGenerator(wcflagState flagState, filepath string) ([3]int, errorCode, error) {
	var output = [3]int{}
	exitCode, err := checkFile(filepath)
	if err != nil {
		return output, exitCode, err
	}
	if wcflagState.countLines {
		lines, err := lineCount(filepath)
		if err != nil {
			return output, unexpectedErrorCode, err
		}
		output[0] = lines
	}
	if wcflagState.countWords {
		words, err := wordCount(filepath)
		if err != nil {
			return output, unexpectedErrorCode, err
		}
		output[1] = words
	}
	if wcflagState.countBytes {
		bytes, err := byteCount(filepath)
		if err != nil {
			return output, unexpectedErrorCode, err

		}
		output[2] = bytes
	}
	return output, 0, nil
}

func generateCliOutput(wcFlagState flagState, fileOutput [3]int, filepath string) {
	var cliOutput string
	if wcFlagState.countLines {
		cliOutput += fmt.Sprintf("%8d", fileOutput[0])
	}
	if wcFlagState.countWords {
		cliOutput += fmt.Sprintf("%8d", fileOutput[1])
	}
	if wcFlagState.countBytes {
		cliOutput += fmt.Sprintf("%8d", fileOutput[2])
	}
	fmt.Printf("%8s %s\n", cliOutput, filepath)
}

func main() {
	wcFlagState, fileList, errorCode, err := flagParser()
	if err != nil {
		errorHandler(" ", err)
		os.Exit(int(errorCode))
	}
	if len(fileList) == 1 {
		fileOutput, errorcode, err := countGenerator(wcFlagState, fileList[0])
		if err != nil {
			errorHandler(fileList[0], err)
			os.Exit(int(errorcode))
		}
		generateCliOutput(wcFlagState, fileOutput, fileList[0])
	} else {
		var totalCount = [3]int{0, 0, 0}
		var osExitCode int
		for _, filepath := range fileList {
			fileOutput, _, err := countGenerator(wcFlagState, filepath)
			if err != nil {
				errorHandler(filepath, err)
				osExitCode = 1
			}
			generateCliOutput(wcFlagState, fileOutput, filepath)
			for i := range fileOutput {
				if fileOutput[i] != -1 {
					totalCount[i] += fileOutput[i]
				}
			}
		}
		generateCliOutput(wcFlagState, totalCount, "total")
		os.Exit(osExitCode)
	}
}
