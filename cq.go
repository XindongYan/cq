package cqueue

import (
	"bufio"
	"encoding/json"
	"os"
)

type CircularQueue struct {
	Length int
	Values []interface{}
	_flag  int
	isLoad bool
}

func NewCircularQueue(length int, load bool) *CircularQueue {
	cq := &CircularQueue{
		Length: length,
		Values: make([]interface{}, length),
		_flag:  0,
		isLoad: load,
	}
	if cq.isLoad {
		cq.loadToMemory()
	}
	return cq
}

func (cq *CircularQueue) Push(values ...interface{}) {
	for i := range values {
		cq._flag = cq._flag % cq.Length
		cq.Values[cq._flag] = values[i]
		cq._flag++
		if cq.isLoad {
			cq.write()
		}
	}
}

func (cq *CircularQueue) Output(length int32, reverse bool) (output []interface{}) {
	needLength := cq.realLen()
	if length < needLength {
		needLength = length
	}
	if needLength == 0 {
		return
	}
	if reverse {
		output = make([]interface{}, needLength)
	}
	_flag := cq._flag
	for i := int32(0); i < needLength; i++ {
		_flag -= 1
		if reverse && output != nil {
			output[needLength-1-i] = cq.Values[_flag]
		} else {
			output = append(output, cq.Values[_flag])
		}
		if _flag == 0 {
			_flag = cq.Length
		}
	}
	return
}

func (cq *CircularQueue) realLen() (len int32) {
	for _, v := range cq.Values {
		if v != nil {
			len++
		}
	}
	return
}

func (cq *CircularQueue) loadToMemory() {
	cq.Push(cq.read()...)
}

func (cq *CircularQueue) read() (data []interface{}) {
	b, _ := os.ReadFile("./circular.json")
	_ = json.Unmarshal(b, &data)
	return
}

func (cq *CircularQueue) write() {
	f, err := os.Create("./circular.json")
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	b, e := json.Marshal(cq.Values)
	if e != nil {
		return
	}
	_, _ = w.Write(b)
	_ = w.Flush()
}
