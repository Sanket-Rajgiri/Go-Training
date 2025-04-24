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
	programName                             = "wc"
)

var (
	ErrIsDirectory      = errors.New("is a directory")
	ErrPermissionDenied = errors.New("permission denied")
	ErrFileNotExist     = errors.New("no such file or directory")
	ErrMissingArguments = errors.New("missing arugments")
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

func cliOutput(wcflagState flagState, filepath string) (string, errorCode, error) {
	var outputString string
	if wcflagState.countLines {
		lines, err := lineCount(filepath)
		if err != nil {
			return " ", unexpectedErrorCode, err
		}
		outputString += fmt.Sprintf("%1d", lines)
	}
	if wcflagState.countWords {
		lines, err := wordCount(filepath)
		if err != nil {
			return " ", unexpectedErrorCode, err
		}
		outputString += fmt.Sprintf("%1d", lines)
	}
	if wcflagState.countBytes {
		lines, err := byteCount(filepath)
		if err != nil {
			return " ", unexpectedErrorCode, err

		}
		outputString += fmt.Sprintf("%1d", lines)
	}
	return outputString, 0, nil
}

func flagParser() flagState {
	lineFlag := flag.Bool("l", false, "count lines")
	wordFlag := flag.Bool("w", false, "count words")
	byteFlag := flag.Bool("c", false, "count bytes")
	flag.Parse()
	if !*lineFlag && !*wordFlag && !*byteFlag {
		return flagState{countLines: true, countWords: true, countBytes: true}
	}
	return flagState{countLines: *lineFlag, countWords: *wordFlag, countBytes: *byteFlag}
}

func checkFile(filepath string) (errorCode errorCode, err error) {
	fileinfo, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return fileNotExistErrorCode, ErrFileNotExist
		} else if os.IsPermission(err) {
			return filePermissionDeniedErrorCode, ErrPermissionDenied
		} else {
			return unexpectedErrorCode, err
		}
	}
	if fileinfo.IsDir() {
		return IsDirectoryErrorCode, ErrIsDirectory
	}
	return 0, nil
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "%s: %s\n", programName, ErrMissingArguments)
		os.Exit(int(unexpectedErrorCode))
	}
	flagState := flagParser()
	filepath := args[len(args)-1]
	exitCode, err := checkFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s: %3s\n", programName, filepath, err)
		os.Exit(int(exitCode))
	}
	output, exitCode, err := cliOutput(flagState, filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s: %3s\n", programName, filepath, err)
		os.Exit(int(exitCode))
	}
	fmt.Printf("%8s %s\n", output, filepath)
}
