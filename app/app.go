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



type Queue interface{
	Put(queueName string, msg types.MsgBody)error
	Get(ctx context.Context, queueName string)(types.MsgBody, error)
}

func Start(conf conf.Configuration, queue Queue) error {
	addr:= fmt.Sprintf(":%s", conf.Port)
	h:=&qHandler{Queue: queue} 
	if err := http.ListenAndServe(addr, h); err != nil {
		return fmt.Errorf("Listen and serve error: %s", err.Error())
	}
 
	return nil
}

type qHandler struct {
	Queue Queue
}

func(qh *qHandler)ServeHTTP(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case http.MethodPut:
		putKeeper(w, r, qh.Queue)
	case http.MethodGet:
		getKeeper(w, r, qh.Queue)
	}
}



func putKeeper(w http.ResponseWriter, r *http.Request, q Queue)error{
	var serialized []byte

	queueName, err:= getQueueName(r)
	if err != nil {
		serialized, _ = json.Marshal(err.Error())
		writeAnswer(w, http.StatusBadRequest, serialized)
		return err
	}

	decoder := json.NewDecoder(r.Body)
	m:=types.MsgBody{}
	err = decoder.Decode(&m)
	if err != nil {
		serialized, _ = json.Marshal(err.Error())
		writeAnswer(w, http.StatusBadRequest, serialized)
		return err
	}
	
	err = q.Put(queueName, m)
	if err!=nil{
		serialized, _ = json.Marshal(err.Error())
		writeAnswer(w, http.StatusBadRequest, serialized)
		return err
	}
	
	serialized, err = json.Marshal(m)
	if err!=nil{
		serialized, _ = json.Marshal(err.Error())
		writeAnswer(w, http.StatusBadRequest, serialized)
		return err
	}

	writeAnswer(w, http.StatusOK, serialized)
	return nil
}

func getQueueName(r *http.Request)(string, error){
	pathParts:= strings.Split(r.URL.Path,"/")
	if len(pathParts)<3{
		return "", fmt.Errorf("Incorrect path, it must be '/queue/:queue'")
	}
	return pathParts[2], nil
}

func getKeeper(w http.ResponseWriter, r *http.Request, q Queue)(types.MsgBody, error){
	var serialized []byte
	queueName, timeout, err:= getQueueNameAndTimeOut(r)
	if err != nil {
		serialized, _ = json.Marshal(err.Error())
	}
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(timeout) * time.Second))
	
	m, err:=q.Get(ctx, queueName)
	if err!=nil{
		serialized, _ = json.Marshal(err.Error())
	}
	serialized, err = json.Marshal(m)
	writeAnswer(w, http.StatusOK, serialized)
	return types.MsgBody{}, nil
}

func getQueueNameAndTimeOut(r *http.Request)(string, int, error){
	pathParts:= strings.Split(r.URL.Path,"/")
	if len(pathParts)<3{
		return "", 0, fmt.Errorf("Incorrect path, it must be '/queue/:queue'")
	}
	timeout, _:= strconv.Atoi(r.URL.Query().Get("timeout"))
	return pathParts[2], timeout, nil
}

func writeAnswer(s http.ResponseWriter, statusCode int, body []byte){
	s.WriteHeader(statusCode)
	if body != nil{
		s.Write(body)
	}
}