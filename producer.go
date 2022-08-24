package main

import (
	"encoding/csv"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	p := Initialize()
	pathList := creatPathList2("C:/Users/EMINENT/Desktop/Tests/*.csv")

	DoItForMe(pathList, p)

	// Wait for all messages to be delivered
	p.Flush(15 * 1000)
	p.Close()
}

func Initialize() *kafka.Producer {
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

func Send(line []string, p *kafka.Producer, wg *sync.WaitGroup) {
	topic := "landmarks"
	wg.Done()

	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	for e := range p.Events() {
	//		switch ev := e.(type) {
	//		case *kafka.Message:
	//			if ev.TopicPartition.Error != nil {
	//				fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
	//			} else {
	//				//fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
	//					*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
	//				fmt.Printf("Produced: value = %s, timestamp = %v, topic = %s \n", string(ev.Value), ev.Timestamp, ev.TopicPartition)
	//			}
	//		}
	//	}
	//}()

	key := line[0]
	val := line[1]
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(val),
	}, nil)
	return
}

func creatPathList2(pattern string) []string {
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func DoItForMe(pathList []string, p *kafka.Producer) {
	//var lines []Record //
	t := time.Now()
	count := 0
	var wg sync.WaitGroup

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
			wg.Add(1)
			go Send(line, p, &wg)
		}

		count += len(data)

		f.Close()
	}

	wg.Wait()

	deltaT := time.Since(t).Seconds()
	fmt.Printf("handler took %f, count: %d, tps: %f \n", deltaT, count, float64(count)/deltaT)
	return
}
