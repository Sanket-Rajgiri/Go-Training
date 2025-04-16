package wc

import (
	"bufio"
	"log"
	"os"
)

// const (
// 	lineDelimiter = '\n'
// )

func lineCount(filepath string) (int, error) {
	lineCounter := 0
	f, err := os.Open(filepath)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	defer f.Close()

	// reader := bufio.NewReader(f)
	// for {
	// 	_, err := reader.ReadString(lineDelimiter)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			lineCounter += 1
	// 			break
	// 		} else {
	// 			log.Print(err)
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
		log.Print(err)
		return 0, err
	}

	return lineCounter, nil
}

func byteCount(filepath string) (int, error) {
	byteCounter := 0
	f, err := os.Open(filepath)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		byteCounter++
	}
	if err := scanner.Err(); err != nil {
		log.Print(err)
		return 0, err
	}
	return byteCounter, nil
}

func wordCount(filepath string) (int, error) {
	wordCounter := 0
	f, err := os.Open(filepath)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		wordCounter++
	}
	if err := scanner.Err(); err != nil {
		log.Print(err)
		return 0, err
	}
	return wordCounter, nil
}
