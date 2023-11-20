package queue

import (
	"fmt"
	"gotesttask/types"
)
type Impl struct{
	MaxMsgInQueue	int
	MaxQueues		int

	Storage map[string]*queue
}

//
func MakeImpl(maxMsgInQueue, maxQueues int)*Impl{
	return &Impl{MaxMsgInQueue: maxMsgInQueue, MaxQueues: maxQueues, Storage: make(map[string]*queue)}
}

func (i *Impl) Put(queueName string, msg types.MsgBody)error{
	if len(i.Storage) >= i.MaxQueues{
		return fmt.Errorf("Max queues number exceeded")
	}
	if q, exist := i.Storage[queueName]; exist {
		if q.Len >= i.MaxMsgInQueue{
			return fmt.Errorf("Max message number in queue exceeded")
		}
		q.put(msg.Message)
	} else {
		q = &queue{}
		q.put(msg.Message)
		i.Storage[queueName] = q
	}
	return nil
}

func (i *Impl) Get(queueName string)(types.MsgBody, error){
	var msg types.MsgBody

	if q, exist := i.Storage[queueName]; exist {
		if q.Len <= 0{
			return types.MsgBody{}, fmt.Errorf("Queue is empty")
		}
		msg.Message = q.get()
	} else {
		return types.MsgBody{}, fmt.Errorf("Queue not exist")
	}

	return msg, nil
}

//
type queue struct{
	Head 			*item
	Tail 			*item

	Len				int
}

type item struct{
	Message string
	Next *item
	Prev *item
}

func (q *queue) put(msg string){
	if q.Head == nil{
		q.Head = &item{Message: msg}
		q.Tail = q.Head
	} else {
		new := &item{Message: msg, Next: q.Head}
		q.Head.Prev = new
		q.Head = new
	}
	q.Len++
}

func (q *queue) get()string{
	msg := q.Tail.Message

	q.Tail = q.Tail.Prev
	q.Tail.Next = nil

	q.Len--

	return msg
}
