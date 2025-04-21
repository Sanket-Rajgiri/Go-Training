package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

const (
	lineFlag      = "-l"
	wordFlag      = "-w"
	characterFlag = "-c"
	helpFlag      = "-h"
)

func lineCount(filepath string) (int, error) {
	lineCounter := 0
	f, _ := os.Open(filepath)
	defer f.Close()

	// reader := bufio.NewReader(f)
	// for {
	// 	_, err := reader.ReadString(lineDelimiter)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			lineCounter += 1
	// 			break
	// 		} else {
	//
	// 			return 0, err
	// 		}
	// 	}
	// 	lineCounter += 1
	// }

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lineCounter++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCounter, nil
}

func byteCount(filepath string) (int, error) {
	byteCounter := 0
	f, _ := os.Open(filepath)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		byteCounter++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return byteCounter, nil
}

func wordCount(filepath string) (int, error) {
	wordCounter := 0
	f, _ := os.Open(filepath)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		wordCounter++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return wordCounter, nil
}
func output(args []string) string {
	var output string
	filepath := args[len(args)-1]
	if slices.Contains(args, lineFlag) {
		lines, err := lineCount(filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s\n", args[0], err)
			os.Exit(125)
		}
		output += fmt.Sprintf("%8d", lines)
	}
	if slices.Contains(args, wordFlag) {
		words, err := wordCount(filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s\n", args[0], err)
			os.Exit(125)
		}
		output += fmt.Sprintf("%8d", words)
	}

	if slices.Contains(args, characterFlag) {
		bytes, err := byteCount(filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s\n", args[0], err)
			os.Exit(125)
		}
		output += fmt.Sprintf("%8d", bytes)
	}

	return output + " " + filepath
}

func main() {
	args := os.Args
	filepath := args[len(args)-1]
	fileinfo, err := os.Stat(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %s\n", args[0], err)
		fmt.Println(fileinfo.IsDir())
		if os.IsNotExist(err) {
			os.Exit(2)
		} else if os.IsPermission(err) {
			os.Exit(1)
		} else {
			os.Exit(125)
		}
	} else if fileinfo.IsDir() {
		fmt.Fprintf(os.Stderr, "%s read %s: Is a directory\n", args[0], filepath)
		os.Exit(21)
	}
	fmt.Fprintf(os.Stdout, "%s\n", output(args))
}
