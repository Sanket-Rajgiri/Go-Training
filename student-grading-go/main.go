package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Grade string

const (
	A             Grade   = "A"
	B             Grade   = "B"
	C             Grade   = "C"
	F             Grade   = "F"
	numberOfTests float32 = 4.0
)

type student struct {
	firstName, lastName, university                string
	test1Score, test2Score, test3Score, test4Score int
}

type studentStat struct {
	student
	finalScore float32
	grade      Grade
}

func parseCSV(filePath string) []student {
	students := make([]student, 0)
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error reading csv : ", err)

	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	if _, err := csvReader.Read(); err != nil {
		log.Fatal(err)
	}

	for {
		data, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		testScores := make([]int, int(numberOfTests))
		for i, v := range data[3:] {
			tempScore, err := strconv.Atoi(v)
			if err != nil {
				log.Fatal(err)
			}
			testScores[i] = tempScore
		}
		tempStudent := student{
			firstName:  data[0],
			lastName:   data[1],
			university: data[2],
			test1Score: testScores[0],
			test2Score: testScores[1],
			test3Score: testScores[2],
			test4Score: testScores[3],
		}
		students = append(students, tempStudent)
	}
	return students
}

func calculateGrade(students []student) []studentStat {
	finalStudentStat := []studentStat{}
	for i := range students {
		var totalScore, finalScore float32
		totalScore = float32(students[i].test1Score + students[i].test2Score + students[i].test3Score + students[i].test4Score)
		finalScore = totalScore / numberOfTests
		var grade Grade
		if finalScore >= 70 {
			grade = A
		} else if finalScore >= 50 {
			grade = B
		} else if finalScore >= 35 {
			grade = C
		} else {
			grade = F
		}
		finalStudentStat = append(finalStudentStat, studentStat{students[i], finalScore, grade})

	}
	return finalStudentStat
}

func findOverallTopper(gradedStudents []studentStat) studentStat {
	var highScoringStudentstat studentStat
	for _, studentStat := range gradedStudents {
		if studentStat.finalScore > highScoringStudentstat.finalScore {
			highScoringStudentstat = studentStat
		}
	}
	return highScoringStudentstat
}

func findTopperPerUniversity(gs []studentStat) map[string]studentStat {
	topperPerUniversity := make(map[string]studentStat)
	for _, studentStat := range gs {
		if topperPerUniversity[studentStat.university].finalScore < studentStat.finalScore {
			topperPerUniversity[studentStat.university] = studentStat
		}
	}
	return topperPerUniversity
}

func main() {
	fmt.Println(findOverallTopper(calculateGrade(parseCSV("./grades.csv"))))
	fmt.Println(findTopperPerUniversity(calculateGrade(parseCSV("./grades.csv"))))
}
