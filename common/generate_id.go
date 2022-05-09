/**
 * @Author: wenliangzhang
 * @Description:Snowflake算法生成ID
 * @File: generate_id
 * @Version: 1.0.0
 * @Date: 2022/4/18 8:42 PM
 */
package common

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

//Snowflake global var
var sequence int = 0
var lastTime int = -1

//every segment bit
var workerIdBits = 5
var datacenterIdBits = 5
var sequenceBits = 12

//every segment max number
var maxWorkerId int = -1 ^ (-1 << workerIdBits)
var maxDatacenterId int = -1 ^ (-1 << datacenterIdBits)
var maxSequence int = -1 ^ (-1 << sequenceBits)

//bit operation shift
var workerIdShift = sequenceBits
var datacenterShift = workerIdBits + sequenceBits
var timestampShift = datacenterIdBits + workerIdBits + sequenceBits

type Snowflake struct {
	datacenterId int
	workerId     int
	epoch        int
	mt           *sync.Mutex
}

var once sync.Once
var SFInstance *Snowflake

func GetSFInstance() *Snowflake {
	if SFInstance == nil {
		once.Do(func() {
			SFInstance, _ = NewSnowflake(0, 0, 1650211200000)
		})
	}
	return SFInstance
}

func NewSnowflake(datacenterId int, workerId int, epoch int) (*Snowflake, error) {

	if datacenterId > maxDatacenterId || datacenterId < 0 {
		return nil, errors.New(fmt.Sprintf("datacenterId cant be greater than %d or less than 0", maxDatacenterId))
	}
	if workerId > maxWorkerId || workerId < 0 {
		return nil, errors.New(fmt.Sprintf("workerId cant be greater than %d or less than 0", maxWorkerId))
	}
	if epoch > getCurrentTime() {
		return nil, errors.New(fmt.Sprintf("epoch time cant be after now"))
	}
	sf := Snowflake{datacenterId, workerId, epoch, new(sync.Mutex)}
	return &sf, nil
}

func (sf *Snowflake) GetUniqueId() int {
	sf.mt.Lock()
	defer sf.mt.Unlock()
	//get current time
	currentTime := getCurrentTime()
	//compute sequence
	if currentTime < lastTime {
		//occur clock back
		//panic or wait,wait is not the best way.can be optimized.
		currentTime = waitUntilNextTime(lastTime)
		sequence = 0
	} else if currentTime == lastTime {
		//at the same time(micro-second)
		sequence = (sequence + 1) & maxSequence
		if sequence == 0 { //overflow max num,wait next time
			currentTime = waitUntilNextTime(lastTime)
		}
	} else if currentTime > lastTime {
		//next time
		sequence = 0
		lastTime = currentTime
	}
	//generate id
	return (currentTime-sf.epoch)<<timestampShift | sf.datacenterId<<datacenterShift | sf.workerId<<workerIdShift | sequence
}

func waitUntilNextTime(lasttime int) int {
	currentTime := getCurrentTime()
	for currentTime <= lasttime {
		time.Sleep(1 * time.Second / 1000)
		//sleep micro second
		currentTime = getCurrentTime()
	}
	return currentTime
}

func getCurrentTime() int {
	return int(time.Now().UnixNano() / 1e6) //micro second
}
