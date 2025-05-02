package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	programName = "grep"
)

var (
	ErrIsDirectory      = errors.New("is a directory")
	ErrPermissionDenied = errors.New("permission denied")
	ErrFileNotExist     = errors.New("no such file or directory")
	ErrInvalidFlags     = errors.New("invalid flags passed")
)

type flagState struct {
	caseInsensitive bool
	invertMatch     bool
	output          bool
}

// var grepFlagState flagState

func errorHandler(filepath string, err error) {
	if len(filepath) > 0 {
		fmt.Fprintf(os.Stderr, "%s: %s: %3s\n", programName, filepath, err)
	} else {
		fmt.Fprintf(os.Stderr, "%s: %3s\n", programName, err)
	}
}

func search(reader io.Reader, key string, caseInsensitive, invertMatch bool) ([]string, error) {
	var output []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if invertMatch {
			if caseInsensitive {
				if !strings.Contains(strings.ToLower(line), strings.ToLower(key)) {
					output = append(output, line)
				}
			} else {
				if !strings.Contains(line, key) {
					output = append(output, line)
				}
			}
		} else {
			if caseInsensitive {
				if strings.Contains(strings.ToLower(line), strings.ToLower(key)) {
					output = append(output, line)
				}
			} else {
				if strings.Contains(line, key) {
					output = append(output, line)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return output, nil
}

func openFile(filepath, searchKey string, grepFlagState *flagState) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotExist
		} else if os.IsPermission(err) {
			return nil, ErrPermissionDenied
		} else {
			if fileInfo, _ := os.Stat(filepath); fileInfo.IsDir() {
				return nil, ErrIsDirectory
			} else {
				return nil, err
			}
		}
	}
	defer file.Close()
	return search(file, searchKey, grepFlagState.caseInsensitive, grepFlagState.invertMatch)
}

func flagParser() (string, []string, flagState, error) {
	var fileList []string
	var searchKey string
	var grepFlagState flagState
	caseInsensitiveFlag := flag.Bool("i", false, "Ignore  case")
	invertMatchFlag := flag.Bool("v", false, "Invert sense of matching, to select non-matching lines")
	outputFlag := flag.Bool("o", false, " grep [options...] [files....] -o [filename]")
	flag.Parse()
	if !flag.Parsed() {
		return searchKey, fileList, grepFlagState, ErrInvalidFlags
	}
	grepFlagState = flagState{caseInsensitive: *caseInsensitiveFlag, output: *outputFlag, invertMatch: *invertMatchFlag}
	searchKey = flag.Arg(0)
	if flag.NArg() > 1 {
		fileList = flag.Args()[1:]
	}
	return searchKey, fileList, grepFlagState, nil
}

func main() {
	var osExitCode int
	searchKey, fileList, grepFlagState, err := flagParser()
	if err != nil {
		errorHandler("", err)
		osExitCode = 1
	} else {
		if len(fileList) < 1 {
			inputStream, err := io.ReadAll(os.Stdin)
			if err != nil {
				errorHandler("", err)
				osExitCode = 1
			} else {
				inputStreamReader := bytes.NewReader(inputStream)
				inputStreamReader.Seek(0, io.SeekStart)
				output, err := search(inputStreamReader, searchKey, grepFlagState.caseInsensitive, grepFlagState.invertMatch)
				if err != nil {
					errorHandler("", err)
					osExitCode = 1
				} else {
					fmt.Printf("%v\n", strings.Join(output, "\n"))
				}
			}
		} else {
			for _, filepath := range fileList {
				fileOutput, err := openFile(filepath, searchKey, &grepFlagState)
				// if err != nil {
				// 	errorHandler(filepath, err)
				// 	osExitCode = 1
				// }
				// fileOutput, err := search(bufio.NewReader(file), searchKey, *grepFlagState.caseInsensitive, *grepFlagState.invertMatch)
				if err != nil {
					errorHandler(filepath, err)
					osExitCode = 1
				} else {
					fmt.Printf("%v\n", strings.Join(fileOutput, "\n"))
				}

			}
		}
	}
	os.Exit(osExitCode)
}
