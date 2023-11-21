package queue

import (
	"context"
	"fmt"
	"gotesttask/types"
	"sync"
)

// Impl - реализует интерфейс app.Queue.
// Хранилище очередей представляет из себя хеш-таблицу списков по строковому ключу(имя очереди).
// Очередь представляет из себя список, так как список позволяет без доп. затрат вставлять/удалять элементы,
// и у нас последовательные операции (пишем в начало списка, удаляем с конца списка)
type Impl struct{
	MaxMsgInQueue	int
	MaxQueues		int

	Storage map[string]*queue

	Mu *sync.Mutex
}

// MakeImpl - выделяем память хранилищу, инициализируем параметры
func MakeImpl(maxMsgInQueue, maxQueues int)*Impl{
	return &Impl{
		MaxMsgInQueue: maxMsgInQueue, 
		MaxQueues: maxQueues, 
		Storage: make(map[string]*queue),
		Mu: &sync.Mutex{}}
}

// Put - кладём сообщение в очередь 
func (i *Impl) Put(queueName string, msg types.MsgBody)error{
	if q, exist := i.Storage[queueName]; exist {// если очередь существует
		if q.Len >= i.MaxMsgInQueue{// если превысили количество сообщений в очереди, выходим с ошибкой
			return fmt.Errorf("Max message number in queue exceeded")
		}
		// i.Mu.Lock()
		q.put(msg.Message)
		// i.Mu.Unlock()
	} else { // иначе пытаемся создать новую очередб
		if len(i.Storage) >= i.MaxQueues{ // если превысили кол-во очередей, выходим с ошибкой
			return fmt.Errorf("Max queues number exceeded")
		}
		q = &queue{}

		// использование мютексов приводит к блокировке очереди в случае, если мы пытаемся
		// читать и писать в одну и ту же очередь.
		// этот недочёт решил (пока) оставить
		// i.Mu.Lock()
		q.put(msg.Message)
		// i.Mu.Unlock()

		i.Storage[queueName] = q //добавляем новую очередь в хранилище
	}
	
	return nil
}
// Get - достаём.
func (i *Impl) Get(ctx context.Context, queueName string)(types.MsgBody, error){
	var msg types.MsgBody
	for {
		select{
		case <-ctx.Done():// контекст с дедлайном: наступает тайм-аут - выходим из select
			return types.MsgBody{}, fmt.Errorf("Queue not exist")
		default : // пока не прилетит в ctx.Done(), пытаемся читать из очереди
			if q, exist := i.Storage[queueName]; exist {
				// сейчас мы контролируем, что читаем не из пустой очереди
				if q.Len > 0{ 
					// i.Mu.Lock()
					msg.Message = q.get()
					// i.Mu.Unlock()
					return msg, nil
				}
			}
		}
	}
}

// queue - структура одной очереди.
// В Tail хранится последний добавленный элемент,
// в Head - начало списка. С Head мы итерируемся, когда читаем из очереди
type queue struct{
	Tail 			*item
	Head 			*item

	Len				int
}

// item - элемент очереди
type item struct{
	Message string
	Prev *item
}

func (q *queue) put(msg string){
	if q.Tail == nil{// указатель нулевой, значит очередь пуста
		q.Tail = &item{Message: msg}
		q.Head = q.Tail // запоминаем начало очереди. Если в очереди один элемент, он же является хвостом
	} else {
		q.Tail.Prev = &item{Message: msg } // текущий элемент на хвост получает связь с новым 
		q.Tail = q.Tail.Prev // обновляем указатель на хвост
	}
	q.Len++
}

func (q *queue) get()string{
	if q.Len>0{ // иначе мы рискуем обратиться к полям Head, который ни на что не указывает
		msg := q.Head.Message // читаем

		q.Head = q.Head.Prev // итерируемся
		if q.Head == nil {
			// мы прочитали последний элемент очереди, 
			// Tail указывает на последний, уже прочитанный элемент, нужно обнулить
			q.Tail = nil 
		}

		q.Len--

		return msg
	}
	return ""
}
