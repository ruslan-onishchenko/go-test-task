package queue

import (
	"context"
	"fmt"
	"gotesttask/types"
	"sync"
)
type Impl struct{
	MaxMsgInQueue	int
	MaxQueues		int

	Storage map[string]*queue

	Mu *sync.Mutex
}

//
func MakeImpl(maxMsgInQueue, maxQueues int)*Impl{
	return &Impl{
		MaxMsgInQueue: maxMsgInQueue, 
		MaxQueues: maxQueues, 
		Storage: make(map[string]*queue),
		Mu: &sync.Mutex{}}
}

func (i *Impl) Put(queueName string, msg types.MsgBody)error{
	
	if len(i.Storage) >= i.MaxQueues{
		return fmt.Errorf("Max queues number exceeded")
	}
	if q, exist := i.Storage[queueName]; exist {
		if q.Len >= i.MaxMsgInQueue{
			return fmt.Errorf("Max message number in queue exceeded")
		}
		// i.Mu.Lock()
		q.put(msg.Message)
		// i.Mu.Unlock()
	} else {
		q = &queue{}

		// i.Mu.Lock()
		q.put(msg.Message)
		// i.Mu.Unlock()

		i.Storage[queueName] = q
	}
	
	return nil
}

func (i *Impl) Get(ctx context.Context, queueName string)(types.MsgBody, error){
	var msg types.MsgBody

	for {
		select{
		case <-ctx.Done():
			return types.MsgBody{}, fmt.Errorf("Queue not exist")
		default :
			if q, exist := i.Storage[queueName]; exist {
				if q.Len > 0{
					// i.Mu.Lock()
					msg.Message = q.get()
					// i.Mu.Unlock()
					return msg, nil
				}
			}
		}
	}

	// return msg, nil
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
		// q.Head.Prev = q.Head
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
	if q.Tail == nil {
		q.Head = nil
	}

	q.Len--

	return msg
}
