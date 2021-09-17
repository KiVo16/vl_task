package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

const SampleUsersDataSize = 1000
const SampleRecordsDataSize = 500
const SampleAssignmentsPerUser = 20

type SampleType struct {
	User         User
	Record       Record
	RecordsCount int
	Err          error
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func (s Server) loadTestData() error {
	defer timeTrack(time.Now(), "factorial")

	fileNames, err := ioutil.ReadFile("sampleNames.json")
	if err != nil {
		return err
	}

	fileRecords, err := ioutil.ReadFile("sampleNames.json")
	if err != nil {
		return err
	}

	names := []string{}
	if err := json.Unmarshal(fileNames, &names); err != nil {
		return err
	}

	records := []string{}
	if err := json.Unmarshal(fileRecords, &records); err != nil {
		return err
	}

	//	c := make(chan SampleType)
	//c2 := make(chan SampleType)

	c := genNamesStream(names...)
	r := genRecordsStream(records...)
	c1 := createSampleUser(s, c)
	c2 := createSampleUser(s, c)
	c3 := createSampleUser(s, c)
	r1 := createSampleRecord(s, r)
	r2 := createSampleRecord(s, r)
	r3 := createSampleRecord(s, r)

	a1 := assign(s, c1, r1)
	a2 := assign(s, c2, r2)
	a3 := assign(s, c3, r3)

	count := 0
	for x := range merge(a1, a2, a3) {
		count++
		//fmt.Print("", x.RecordsCount)
		fmt.Println("created user: ", x.User.Name, x.User.ID, x.Record.Name, x.Record.ID)
	}

	log.Println(count)

	//log.Println(arr)

	return nil
}

func genNamesStream(names ...string) <-chan SampleType {
	out := make(chan SampleType)
	go func() {
		for i := 0; i < SampleUsersDataSize; i++ {
			name := names[randRange(0, len(names))]
			out <- SampleType{User: User{Name: name}}
		}
		close(out)
	}()
	return out
}

func genRecordsStream(records ...string) <-chan SampleType {
	out := make(chan SampleType)
	go func() {
		for i := 0; i < SampleRecordsDataSize; i++ {
			idx := randRange(0, len(records))
			record := records[idx]
			out <- SampleType{Record: Record{Name: record, Type: fmt.Sprintf("t%d", idx)}}
		}
		close(out)
	}()
	return out
}

func createSampleUser(s Server, in <-chan SampleType) <-chan SampleType {
	out := make(chan SampleType)
	go func() {
		for n := range in {
			u, err := s.createUser(n.User.Name)
			if err != nil {
				n.Err = err
				out <- n
				continue
			}
			n.User = *u
			out <- n
		}
		close(out)
	}()

	return out
}

func createSampleRecord(s Server, in <-chan SampleType) <-chan SampleType {
	out := make(chan SampleType)
	go func() {
		for n := range in {
			r, err := s.createRecord(n.Record.Name, n.Record.Type)
			if err != nil {
				n.Err = err
				out <- n
				continue
			}
			n.Record = *r
			out <- n
		}
		close(out)
	}()
	return out
}

func assign(s Server, c, r <-chan SampleType) <-chan SampleType {
	out := make(chan SampleType)
	temp := make(chan SampleType)

	go func() {
		for n := range c {
			temp <- n
		}
		close(temp)
	}()

	go func() {
		for n := range r {
			//arr := make([]SampleType, SampleAssignmentsPerUser)
			for i := 0; i < SampleAssignmentsPerUser; i++ {
				u := <-temp
				err := s.assignRecordToUser(u.User.ID, n.Record.ID)
				if err != nil {
					fmt.Println("assign err: ", err)
				}
				n.User = u.User
				n.RecordsCount += 1
				out <- n
				if n.RecordsCount < 20 {
					go func() {
						temp <- u
					}()
				}
			}
		}
		close(out)
		//close(temp)
	}()
	return out
}

func merge(cs ...<-chan SampleType) <-chan SampleType {
	var wg sync.WaitGroup
	out := make(chan SampleType)

	output := func(c <-chan SampleType) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
