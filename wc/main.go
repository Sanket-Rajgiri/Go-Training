package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
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
	ErrInvalidFlags     = errors.New("invalid flags passed")
)

type flagState struct {
	countLines bool
	countWords bool
	countBytes bool
}

func countWithSplit(file []byte, split bufio.SplitFunc) (int, error) {
	// f, err := os.Open(filepath)
	// if err != nil {
	// 	return 0, err
	// }
	// defer f.Close()
	reader := bytes.NewReader(file)
	scanner := bufio.NewScanner(reader)
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

func lineCount(file []byte) (int, error) {
	return countWithSplit(file, bufio.ScanLines)
}

func wordCount(file []byte) (int, error) {
	return countWithSplit(file, bufio.ScanWords)
}

func byteCount(file []byte) (int, error) {
	return countWithSplit(file, bufio.ScanBytes)
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
	if !*lineFlag && !*wordFlag && !*byteFlag {
		wcflagState = flagState{countLines: true, countWords: true, countBytes: true}
	} else {
		wcflagState = flagState{countLines: *lineFlag, countWords: *wordFlag, countBytes: *byteFlag}
	}
	return wcflagState, fileList, 0, nil
}

func openFile(filepath string) ([]byte, errorCode, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fileNotExistErrorCode, ErrFileNotExist
		} else if os.IsPermission(err) {
			return nil, filePermissionDeniedErrorCode, ErrPermissionDenied
		} else {
			if fileInfo, _ := os.Stat(filepath); fileInfo.IsDir() {
				return nil, IsDirectoryErrorCode, ErrIsDirectory
			} else {
				return nil, unexpectedErrorCode, err
			}
		}
	}
	return file, 0, nil
}

func stdInputHandler() ([]byte, errorCode, error) {
	var file []byte
	file, err := io.ReadAll(os.Stdin)
	if err != nil {
		return file, unexpectedErrorCode, err
	}
	return file, 0, nil
}

func errorHandler(filepath string, err error) {
	if len(filepath) > 0 {
		fmt.Fprintf(os.Stderr, "%s: %s: %3s\n", programName, filepath, err)
	} else {
		fmt.Fprintf(os.Stderr, "%s: %3s\n", programName, err)
	}
}

func countGenerator(wcflagState flagState, file []byte) ([3]int, errorCode, error) {
	var output = [3]int{-1, -1, -1}
	if wcflagState.countLines {
		lines, err := lineCount(file)
		if err != nil {
			return output, unexpectedErrorCode, err
		}
		output[0] = lines
	}
	if wcflagState.countWords {
		words, err := wordCount(file)
		if err != nil {
			return output, unexpectedErrorCode, err
		}
		output[1] = words
	}
	if wcflagState.countBytes {
		bytes, err := byteCount(file)
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
	if len(fileList) == 0 {
		input, errorCode, err := stdInputHandler()
		if err != nil {
			errorHandler(" ", err)
			os.Exit(int(errorCode))
		}
		fileOutput, errorCode, err := countGenerator(wcFlagState, input)
		if err != nil {
			errorHandler(" ", err)
			os.Exit(int(errorCode))
		}
		generateCliOutput(wcFlagState, fileOutput, " ")
	}
	if len(fileList) > 0 {
		var totalCount = [3]int{0, 0, 0}
		var osExitCode int
		var wg sync.WaitGroup
		var mu sync.Mutex
		type outputType struct {
			fileName  string
			countList [3]int
		}
		outputsChannel := make(chan outputType, len(fileList)+1)

		for _, filepath := range fileList {
			wg.Add(1)
			go func(fp string, totalCount *[3]int) {
				defer wg.Done()
				file, _, err := openFile(filepath)
				if err != nil {
					errorHandler(filepath, err)
					osExitCode = 1
				}
				fileOutput, _, err := countGenerator(wcFlagState, file)
				if err != nil {
					errorHandler(filepath, err)
					osExitCode = 1
				}
				outputsChannel <- outputType{fileName: filepath, countList: fileOutput}
				for i := range fileOutput {
					if fileOutput[i] != -1 {
						mu.Lock()
						totalCount[i] += fileOutput[i]
						mu.Unlock()
					}
				}
			}(filepath, &totalCount)
		}
		wg.Wait()
		if len(fileList) > 1 {
			outputsChannel <- outputType{fileName: "total", countList: totalCount}
		}
		close(outputsChannel)
		for res := range outputsChannel {
			generateCliOutput(wcFlagState, res.countList, res.fileName)
		}

		os.Exit(osExitCode)
	}
}
