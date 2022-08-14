package main

//
//import (
//	"encoding/csv"
//	"fmt"
//	"log"
//	"os"
//	"path/filepath"
//)
//
//type Record struct { //it contains just a line of data
//	fields []string
//}
//
//func createRecordList(data [][]string) []Record {
//	var recordList []Record
//	for i, line := range data {
//		if i >= 0 {
//			var rec Record
//			rec.fields = line
//			recordList = append(recordList, rec)
//		}
//	}
//	return recordList
//}
//
//func creatPathList(pattern string) []string {
//	files, err := filepath.Glob(pattern)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return files
//}
//
//func readFromSource(pathList []string) []Record {
//	var lines []Record
//	for _, file := range pathList {
//		f, err := os.Open(file)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		csvReader := csv.NewReader(f)
//		data, err := csvReader.ReadAll()
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		lines = append(lines[:], createRecordList(data)...)
//		f.Close()
//	}
//	return lines
//}
//
//func main() {
//	pathList := creatPathList("C:/Users/EMINENT/Desktop/Tests/*.csv")
//	lines := readFromSource(pathList)
//
//	fmt.Printf("%+v\n", lines[10])
//}
//
