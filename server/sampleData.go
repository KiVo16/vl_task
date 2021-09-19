package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

const SampleUsersDataSize = 1000
const SampleRecordsDataSize = 500
const SampleAssignmentsPerUser = 20
const MaxRetries = 1

type SampleSummary struct {
	Err     error
	Retries int
}

type SampleUserChanType struct {
	User
	AssignErr     error
	AssignRetries int
	SampleSummary
}

func (s SampleUserChanType) GenerateSummary() string {
	return fmt.Sprintf("Type %v, retries: %d, error: %v\n, assignErr: %v, retriesAssign: %d", "User", s.Retries, s.Err, s.AssignErr, s.AssignRetries)
}

type SampleRecordChanType struct {
	Record
	SampleSummary
}

func (s SampleRecordChanType) GenerateSummary() string {
	return fmt.Sprintf("Type %v, retries: %d, error: %v\n", "Record", s.Retries, s.Err)
}

var lastInsertedUserId, lastInsertedRecordsId int = -1, -1
var userCount, recordsCount, assignErrorCount int64 = 0, 0, 0
var availableRecords = []int{}

var mu sync.Mutex

func (s Server) loadTestData() error {
	startTime := time.Now()

	fileNames, err := ioutil.ReadFile("./sampleData/sampleNames.json")
	if err != nil {
		return err
	}

	fileRecords, err := ioutil.ReadFile("./sampleData/sampleNames.json")
	if err != nil {
		return err
	}

	names := []string{}
	records := []string{}

	if err := json.Unmarshal(fileNames, &names); err != nil {
		return err
	}

	if err := json.Unmarshal(fileRecords, &records); err != nil {
		return err
	}

	lastUser := User{}
	lastRecord := Record{}

	s.db.Last(&lastUser)
	s.db.Last(&lastRecord)

	lastInsertedUserId = lastUser.ID
	lastInsertedRecordsId = lastRecord.ID

	bar := progressbar.NewOptions(SampleUsersDataSize+SampleRecordsDataSize,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("Generating sample data..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	var wg sync.WaitGroup

	c := genUsersStream(names...)
	r := genRecordsStream(records...)
	c1 := createSampleUser(s, c)
	c2 := createSampleUser(s, c)

	r1 := createSampleRecord(s, r)
	r2 := createSampleRecord(s, r)

	b1 := assignRecordToUser(s, c1)
	b2 := assignRecordToUser(s, c2)

	wg.Add(2)
	go func(in ...<-chan SampleUserChanType) {
		for user := range mergeUsers(in...) {
			//	fmt.Println(user.GenerateSummary())

			mu.Lock()
			if user.Err == nil {
				userCount += 1
			}

			if user.AssignErr != nil {
				assignErrorCount++
			}

			bar.Add(1)
			mu.Unlock()
		}
		wg.Done()
	}(b1, b2)
	go func(in ...<-chan SampleRecordChanType) {
		for record := range mergeRecords(in...) {
			//fmt.Println(record.GenerateSummary())

			if record.Err == nil {
				mu.Lock()

				recordsCount += 1
				bar.Add(1)
				mu.Unlock()
			}
		}
		wg.Done()
	}(r1, r2)

	wg.Wait()

	log.Printf("\nSummary: created users = %d, created records = %d, record assign error: %d, operation took: %v", userCount, recordsCount, assignErrorCount, time.Since(startTime))

	return nil
}

func genUsersStream(names ...string) chan SampleUserChanType {
	out := make(chan SampleUserChanType)
	go func() {
		for i := 0; i < SampleUsersDataSize; i++ {
			name := names[randRange(0, len(names)-1)]
			out <- SampleUserChanType{User: User{Name: name}}
		}
		close(out)
	}()
	return out
}

func genRecordsStream(records ...string) chan SampleRecordChanType {
	out := make(chan SampleRecordChanType)
	go func() {
		for i := 0; i < SampleRecordsDataSize; i++ {
			idx := randRange(0, len(records)-1)
			record := records[idx]
			out <- SampleRecordChanType{Record: Record{Name: record, Type: fmt.Sprintf("t%d", idx)}}
		}
		close(out)
	}()
	return out
}

func createSampleUser(s Server, in chan SampleUserChanType) chan SampleUserChanType {
	out := make(chan SampleUserChanType)
	go func() {
		for user := range in {
			u, err := s.createUser(user.Name)

			if err != nil {
				user.Err = err
				user.Retries += 1
				out <- user
				continue
			}

			user.User = *u
			user.Err = nil
			out <- user
		}
		close(out)
	}()

	return out
}

func createSampleRecord(s Server, in chan SampleRecordChanType) <-chan SampleRecordChanType {
	out := make(chan SampleRecordChanType)
	go func() {
		for n := range in {
			r, err := s.createRecord(n.Record.Name, n.Record.Type)
			if err != nil {
				n.Err = err
				n.Retries += 1
				out <- n
				continue
			}

			mu.Lock()
			availableRecords = append(availableRecords, r.ID)
			mu.Unlock()
			n.Record = *r
			out <- n
		}
		close(out)
	}()
	return out
}

func assignRecordToUser(s Server, c chan SampleUserChanType) <-chan SampleUserChanType {
	out := make(chan SampleUserChanType)

	go func() {
		for user := range c {
			batch := []UserRecord{}
			for i := 0; i < SampleAssignmentsPerUser; i++ {
				r := 0
				if len(availableRecords) > 0 {
					r = availableRecords[randRange(0, len(availableRecords)-1)]
				} else {
					r = randRange(0, lastInsertedRecordsId)
				}

				batch = append(batch, UserRecord{UserID: user.User.ID, RecordID: r})
			}

			if err := s.db.Create(&batch).Error; err != nil {
				user.AssignErr = errors.New("test")
				user.AssignRetries += 1
			}
			out <- user
		}
		close(out)
	}()
	return out
}

func mergeUsers(cs ...<-chan SampleUserChanType) <-chan SampleUserChanType {
	var wg sync.WaitGroup
	out := make(chan SampleUserChanType)

	output := func(c <-chan SampleUserChanType) {
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

func mergeRecords(cs ...<-chan SampleRecordChanType) <-chan SampleRecordChanType {
	var wg sync.WaitGroup
	out := make(chan SampleRecordChanType)

	output := func(c <-chan SampleRecordChanType) {
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
