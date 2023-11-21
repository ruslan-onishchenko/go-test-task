package app

import (
	"context"
	"encoding/json"
	"fmt"
	"gotesttask/conf"
	"gotesttask/types"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Queue - Интерфейс для менеджера очередей
type Queue interface{
	Put(queueName string, msg types.MsgBody)error
	Get(ctx context.Context, queueName string)(types.MsgBody, error)
}


// Start - здесь собираются все компоненты приложения, 
// создается экземпляр http сервера с хэндлером qHandler
// и происходит старт сервепа
func Start(conf conf.Configuration, queue Queue) error {
	addr:= fmt.Sprintf(":%d", conf.Port)
	h:=&qHandler{Queue: queue, Conf: conf} 
	if err := http.ListenAndServe(addr, h); err != nil {
		return fmt.Errorf("Listen and serve error: %s", err.Error())
	}
 
	return nil
}

// qHandler - обработчик запросов по нашей задаче.
// Т.к. у нас всего два вида запросов:
//
// PUT /queue/:queue ,
//
// GET /queue/:queue ,
//
// то нам достаточно различать запросы по методоам PUT и GET
type qHandler struct {
	Queue Queue
	Conf conf.Configuration
}

func(qh *qHandler)ServeHTTP(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case http.MethodPut:
		putKeeper(w, r, qh.Queue)
	case http.MethodGet:
		getKeeper(w, r, qh.Queue, qh.Conf.TimeOut)
	}
}


// putKeeper - читаем тело запроса и кладём в очередь
func putKeeper(w http.ResponseWriter, r *http.Request, q Queue){
	// узнаём имя очереди
	queueName, err:= getQueueName(r)
	if err != nil {
		writeAnswer(w, http.StatusBadRequest, err.Error())
		return
	}
	// читаем тело
	decoder := json.NewDecoder(r.Body)
	m:=types.MsgBody{}
	err = decoder.Decode(&m)
	if err != nil {
		writeAnswer(w, http.StatusBadRequest, err.Error())
		return
	}
	// кладём сообщение
	err = q.Put(queueName, m)
	if err!=nil{
		writeAnswer(w, http.StatusBadRequest, err.Error())
		return
	}
	// пишем что всё ок
	writeAnswer(w, http.StatusOK, "OK")
	return
}

// getQueueName - парсим только имя очереди из URL
func getQueueName(r *http.Request)(string, error){
	pathParts:= strings.Split(r.URL.Path,"/")
	if len(pathParts)<3{
		return "", fmt.Errorf("Incorrect path, it must be '/queue/:queue'")
	}
	return pathParts[2], nil
}

// getKeeper - достаём сообщение из очереди и отправляем в теле ответа
func getKeeper(w http.ResponseWriter, r *http.Request, q Queue, confTimeout int){
	// читаем имя очереди и таймаут
	queueName, timeout, err:= getQueueNameAndTimeOut(r)
	if err != nil {
		writeAnswer(w, http.StatusBadRequest, err.Error())
		return
	}
	// если таймут не указан, используем значение таймаута из параметров
	if timeout <1 {
		timeout = confTimeout
	}
	// Использую контекст с дедлайном вместо обычного канала, 
	// так как он позволяет сделать тайм-аут средствами пакета context легко и просто.
	// context.CancelFunc нам не нужен, так как канал Done закроется сам по истечении времени
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(timeout) * time.Second))
	// читаем
	m, err:=q.Get(ctx, queueName)
	if err!=nil{
		if err.Error() == "Queue not exist"{// если не найдено
			writeAnswer(w, http.StatusNotFound, err.Error())
			return
		}
		writeAnswer(w, http.StatusBadRequest, err.Error())
		return
	}

	writeAnswer(w, http.StatusOK, m)
	return 
}

// getQueueNameAndTimeOut - парсим имя очереди и таймаут из URL.
//
// Например: http://localhost/queue/pet?timeout=N
//
// pet - имя очереди
//
// timeout=N - таймаут
func getQueueNameAndTimeOut(r *http.Request)(string, int, error){
	pathParts:= strings.Split(r.URL.Path,"/")
	if len(pathParts)<3{
		return "", 0, fmt.Errorf("Incorrect path, it must be '/queue/:queue'")
	}
	timeout, _:= strconv.Atoi(r.URL.Query().Get("timeout"))
	return pathParts[2], timeout, nil
}

// writeAnswer - вспомогательная функция для записи кода ответа и тела
func writeAnswer(s http.ResponseWriter, statusCode int, body interface{}){
	serialized, _:= json.Marshal(body)
	s.WriteHeader(statusCode)
	s.Write(serialized)
}

