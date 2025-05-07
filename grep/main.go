package main

import (
	"bufio"
	"container/list"
	"errors"
	"path/filepath"
	"sync"

	// "flag"
	"fmt"
	"io"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

const (
	programName = "grep"
)

var (
	ErrIsDirectory      = errors.New("is a directory")
	ErrPermissionDenied = errors.New("permission denied")
	ErrFileNotExist     = errors.New("no such file or directory")
	ErrInvalidFlags     = errors.New("invalid flags passed")
	ErrFileExists       = errors.New("file already exists")
)

type flagState struct {
	caseInsensitive bool
	invertMatch     bool
	recursive       bool
	output          string
	afterContext    int
	beforeContext   int
}

// var grepFlagState flagState

func errorHandler(filepath string, err error) {
	if len(filepath) > 0 {
		fmt.Fprintf(os.Stderr, "%s: %s: %3s\n", programName, filepath, err)
	} else {
		fmt.Fprintf(os.Stderr, "%s: %3s\n", programName, err)
	}
}

func SearchString(line, key string, CaseInsensitive, invertMatch bool) bool {
	if CaseInsensitive {
		line = strings.ToLower(line)
		key = strings.ToLower(key)
	}
	if strings.Contains(line, key) != invertMatch {
		return true
	}
	return false
}

func fileProcessor(reader io.Reader, key string, grepFlagState *flagState) ([]string, error) {
	var output []string
	scanner := bufio.NewScanner(reader)
	beforeQueue := list.New()
	matched := false
	var afterContextTracker int
	for scanner.Scan() {
		line := scanner.Text()
		if SearchString(line, key, grepFlagState.caseInsensitive, grepFlagState.invertMatch) {
			for e := beforeQueue.Front(); e != nil; e = e.Next() {
				output = append(output, e.Value.(string))
			}
			beforeQueue.Init()
			matched = true
			output = append(output, line)
			afterContextTracker = 0
		} else if matched {
			if afterContextTracker == grepFlagState.afterContext {
				matched = false
			} else if afterContextTracker < grepFlagState.afterContext {
				output = append(output, line)
				afterContextTracker++
			}
		}
		if grepFlagState.beforeContext > 0 && !matched {
			if beforeQueue.Len() == grepFlagState.beforeContext {
				beforeQueue.Remove(beforeQueue.Front())
			}
			beforeQueue.PushBack(line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return output, nil
}

func recursiveFileList(root string) []string {
	var filesList []string
	info, err := os.Stat(root)
	if err != nil {
		return []string{root}
	}
	if info.IsDir() {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				filesList = append(filesList, path)
			}
			return nil
		})
		if err != nil {
			return []string{root}
		}
	} else {
		return []string{root}
	}

	return filesList
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
	return fileProcessor(file, searchKey, grepFlagState)
}

func processStdin(reader io.Reader, key string, grepFlagstate *flagState, bufWriter io.Writer) error {
	scanner := bufio.NewScanner(reader)
	writer := bufio.NewWriter(bufWriter)
	defer writer.Flush()
	beforeQueue := list.New()
	matched := false
	var afterContextTracker int
	for scanner.Scan() {
		line := scanner.Text()
		if SearchString(line, key, grepFlagstate.caseInsensitive, grepFlagstate.invertMatch) {
			for e := beforeQueue.Front(); e != nil; e = e.Next() {
				writer.WriteString(e.Value.(string) + "\n")
			}
			beforeQueue.Init()
			matched = true
			writer.WriteString(line + "\n")
			afterContextTracker = 0
		} else if matched {
			if afterContextTracker < grepFlagstate.afterContext {
				writer.WriteString(line + "\n")
				afterContextTracker++
			} else if afterContextTracker == grepFlagstate.afterContext {
				matched = false
			}
		}
		if grepFlagstate.beforeContext > 0 && !matched {
			if beforeQueue.Len() == grepFlagstate.beforeContext {
				beforeQueue.Remove(beforeQueue.Front())
			}
			beforeQueue.PushBack(line)
		}
	}
	return scanner.Err()
}

func flagParser() (string, []string, flagState, error) {
	var fileList []string
	var searchKey string
	var grepFlagState flagState
	caseInsensitiveFlag := flag.BoolP("ignore-case", "i", false, "Ignore  case")
	invertMatchFlag := flag.BoolP("invert-match", "v", false, "Invert sense of matching, to select non-matching lines")
	recursiveFlag := flag.BoolP("recursive", "r", false, "like --directories=recurse")
	outputFlag := flag.StringP("output", "o", "", " grep [options...] [files....] -o [filename]")
	afterContextFlag := flag.IntP("after-context", "A", 0, "print NUM lines of trailing context")
	beforeContextFlag := flag.IntP("before-context", "B", 0, "print NUM lines of leading context")
	contextFlag := flag.IntP("context", "C", 0, "print NUM lines of output context")

	flag.Parse()
	if !flag.Parsed() {
		return searchKey, fileList, grepFlagState, ErrInvalidFlags
	}
	grepFlagState = flagState{
		caseInsensitive: *caseInsensitiveFlag,
		invertMatch:     *invertMatchFlag,
		output:          *outputFlag,
		recursive:       *recursiveFlag,
	}

	if *contextFlag > 0 {
		grepFlagState.afterContext = *contextFlag
		grepFlagState.beforeContext = *contextFlag
	} else {
		grepFlagState.afterContext = *afterContextFlag
		grepFlagState.beforeContext = *beforeContextFlag
	}

	searchKey = flag.Arg(0)
	if flag.NArg() > 1 {
		fileList = flag.Args()[1:]
	}
	if grepFlagState.recursive {
		var tempFileList []string
		if len(fileList) < 1 {
			fileList = append(fileList, ".")
		}
		for _, file := range fileList {
			tempFileList = append(tempFileList, recursiveFileList(file)...)
		}
		fileList = tempFileList
	}
	return searchKey, fileList, grepFlagState, nil
}

func printOnStdOut(filepath string, output []string) {
	for _, line := range output {
		if len(filepath) > 1 {
			fmt.Printf("%s: %s\n", filepath, line)
		} else {
			fmt.Println(line)
		}
	}
}

func WriteToFile(filepath string, output []string) error {
	_, err := os.Stat(filepath)
	if err == nil {
		return ErrFileExists
	}
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(strings.Join(output, "\n"))
	if err != nil {
		return err
	}
	return writer.Flush()
}

func main() {
	var osExitCode int
	searchKey, fileList, grepFlagState, err := flagParser()
	if err != nil {
		errorHandler("", err)
		osExitCode = 1
	} else {
		if len(fileList) < 1 {
			err := processStdin(os.Stdin, searchKey, &grepFlagState, os.Stdout)
			if err != nil {
				errorHandler("", err)
				osExitCode = 1
			}
		} else {
			var wg sync.WaitGroup
			type outputType struct {
				fileName string
				output   []string
			}
			outputChannel := make(chan outputType, len(fileList))
			// var mu sync.Mutex
			for _, filepath := range fileList {
				wg.Add(1)
				go func(filepath string, grepFlagState flagState) {
					defer wg.Done()
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
						outputChannel <- outputType{fileName: filepath, output: fileOutput}
					}
				}(filepath, grepFlagState)
			}
			wg.Wait()
			close(outputChannel)
			for msg := range outputChannel {
				if len(grepFlagState.output) > 0 {
					err := WriteToFile(grepFlagState.output, msg.output)
					if err != nil {
						errorHandler(msg.fileName, err)
						osExitCode = 1
					}
				} else {
					if len(fileList) > 1 && len(msg.output) > 0 {
						printOnStdOut(msg.fileName, msg.output)
						if grepFlagState.afterContext > 0 || grepFlagState.beforeContext > 0 {
							fmt.Println("--")
						}
					} else {
						printOnStdOut("", msg.output)
					}
				}
			}
			// wg.Wait()
		}
	}
	os.Exit(osExitCode)
}
