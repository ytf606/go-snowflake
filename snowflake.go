package snowflake

import (
	"errors"
	"fmt"
	"time"
	"sync"
)

const (
	offsetTime = int64(1)
	serverBit = 7
	processBit = 8
	sequenceBit = 9
	maxServerId = -1 ^ (-1 << serverBit)
	maxProcessId = -1 ^ (-1 << processBit)
	maxSequenceId = -1 ^ (-1 << sequenceBit)
	offsetTimeShift = serverBit + processBit + sequenceBit
	serverShift = processBit + sequenceBit
	processShift = sequenceBit
	sequenceMask = maxSequenceId
)

type ProcessNode struct {
	sequence int64
	lastTime int64
	serverId int64
	processId int64
	offsetTime int64
	mutex sync.Mutex
}

func NewProcessWork(ServerId int64, ProcessId int64) (*ProcessNode, error) {
	if ServerId > maxServerId || ServerId < 0 {
		fmt.Sprintf("Server Id must be less %d and greater 0", maxServerId)
		return nil, errors.New("Server Id error")
	}
	if ProcessId > maxProcessId || ProcessId < 0 {
		fmt.Sprintf("Process Id must be less %d and greater 0", maxProcessId)
		return nil, errors.New("Process Id error")
	}
	processNode := &ProcessNode{}
	processNode.sequence = 0
	processNode.lastTime = -1
	processNode.serverId = ServerId
	processNode.processId = ProcessId
	processNode.offsetTime = offsetTime
	processNode.mutex = sync.Mutex{}
	return processNode, nil
}

func genTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func tilNextMillis(lastTime int64) int64 {
	time := genTime()
	for time <= lastTime {
		time = genTime()
	}
	return time
}

func (id *ProcessNode) Id() (int64, error) {
	id.mutex.Lock()
	defer id.mutex.Unlock()
	return id.create()
}

func (id *ProcessNode) create() (int64, error) {
	time := genTime()
	if time < id.lastTime {
		return 0, errors.New("local time backwards, please check it")	
	}
	if id.lastTime == time {
		id.sequence = (id.sequence + 1) & sequenceMask
		if id.sequence == 0 {
			time = tilNextMillis(id.lastTime)
		}
	} else {
		id.sequence = 0
	}
	id.lastTime = time
	return ((time - id.offsetTime) << offsetTimeShift) | (id.serverId << serverShift) | (id.processId << processShift) | id.sequence, nil
}
