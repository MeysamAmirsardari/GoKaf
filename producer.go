package main

import (
	"encoding/csv"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"path/filepath"
)

func main() {
	p := Ininitialize()
	pathList := creatPathList2("C:/Users/EMINENT/Desktop/Tests/*.csv")

	DoItForMe(pathList, p)

	// Wait for all messages to be delivered
	p.Flush(15 * 1000)
	p.Close()
}

func Ininitialize() *kafka.Producer {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <config-file-path>\n",
			os.Args[0])
		os.Exit(1)
	}
	configFile := os.Args[1]
	conf := ReadConfig(configFile)

	p, err := kafka.NewProducer(&conf)

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	return p
}

func Send(line []string, p *kafka.Producer) {
	topic := "landmarks"

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	key := line[0]
	val := line[1]
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(val),
	}, nil)
}

func creatPathList2(pattern string) []string {
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func DoItForMe(pathList []string, p *kafka.Producer) {
	//var lines []Record
	for _, file := range pathList {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}

		csvReader := csv.NewReader(f)
		data, err := csvReader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		for _, line := range data {
			go Send(line, p)
		}

		f.Close()
	}
	return
}
