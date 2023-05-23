package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/journalctl"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

var streamLogs []*interfaces.StreamLog

type LogType string

func (d *GormDatabase) CreateLogAndReturn(body *interfaces.StreamLog) (*interfaces.StreamLog, error) {
	streamLog, err := d.CreateStreamLog(body) // add a log
	if err != nil {
		return nil, err
	}
	timeDuration := body.Duration + 1
	time.Sleep(time.Duration(timeDuration) * time.Second)
	streamLogData := d.GetStreamLog(streamLog) // get add
	d.DeleteStreamLog(streamLog)               // delete the log
	return streamLogData, nil                  // return the data
}

func (d *GormDatabase) GetStreamsLogs() []*interfaces.StreamLog {
	if streamLogs == nil {
		streamLogs = []*interfaces.StreamLog{}
	}
	return streamLogs
}

func (d *GormDatabase) GetStreamLog(uuid string) *interfaces.StreamLog {
	for _, _log := range streamLogs {
		if _log.UUID == uuid {
			return _log
		}
	}
	return nil
}

func (d *GormDatabase) CreateStreamLog(body *interfaces.StreamLog) (string, error) {
	body.UUID = nuuid.ShortUUID("log")
	body.Message = []string{}
	s := systemctl.New(false, 30)
	isRunning, status, err := s.IsRunning(body.Service)
	if !isRunning || err != nil {
		return status, errors.New(fmt.Sprintf("service not running %s", body.Service))
	}
	go createLogStream(body)
	return body.UUID, nil
}

func (d *GormDatabase) DeleteStreamLog(uuid string) bool {
	deleted := false
	for i, entry := range streamLogs {
		if entry.UUID == uuid {
			streamLogs = append(streamLogs[:i], streamLogs[i+1:]...)
			deleted = true
			break
		}
	}
	return deleted
}

func (d *GormDatabase) DeleteStreamLogs() {
	streamLogs = []*interfaces.StreamLog{}
}

func checkSubstrings(str string, subs ...string) (bool, int) {
	matches := 0
	isCompleteMatch := true
	for _, sub := range subs {
		if strings.Contains(str, sub) {
			matches += 1
		} else {
			isCompleteMatch = false
		}
	}
	return isCompleteMatch, matches
}

func createLogStream(body *interfaces.StreamLog) {
	log.Infof("starting log stream for service: %s for time: %d secounds", body.Service, body.Duration)
	entries, err := journalctl.NewJournalCTL().EntriesAfter(body.Service, "", "")
	for _, entry := range entries {
		lenKeyWordsFilter := len(body.KeyWordsFilter)
		if lenKeyWordsFilter > 0 {
			_, matches := checkSubstrings(entry.Message, body.KeyWordsFilter...)
			if matches == lenKeyWordsFilter {
				body.Message = append(body.Message, entry.Message)
			}
		} else {
			body.Message = append(body.Message, entry.Message)
		}
	}
	if err == nil {
		streamLogs = append(streamLogs, body)
	}
	time.Sleep(time.Duration(body.Duration) * time.Second)
	log.Infof("finished log stream for service: %s for time: %d secounds", body.Service, body.Duration)
}
