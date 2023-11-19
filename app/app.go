package app

import (
	"encoding/json"
	"fmt"
	"gotesttask/conf"
	"gotesttask/types"
	"net/http"
	"strings"
)



type Queue interface{
	Put(queueName string, msg types.MsgBody)error
	Get(queueName string)(types.MsgBody, error)
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

func(qh *qHandler)ServeHTTP(s http.ResponseWriter, r *http.Request){
	switch r.Method{
	case http.MethodPut:
		putKeeper(s, r, qh.Queue)
	case http.MethodGet:
		getKeeper(s, r, qh.Queue)
	}
}



func putKeeper(s http.ResponseWriter, r *http.Request, q Queue)error{
	queueName, err:= getQueueName(r)
	if err != nil {
		
	}

	decoder := json.NewDecoder(r.Body)
	m:=types.MsgBody{}
	err = decoder.Decode(&m)
	if err != nil{}
	err = q.Put(queueName, m)
	if err != nil{}

	return nil
}

func getQueueName(r *http.Request)(string, error){
	pathParts:= strings.Split(r.URL.Path,"/")
	if len(pathParts)<3{
		return "", fmt.Errorf("Incorrect path, it must be '/queue/:queue'")
	}
	return pathParts[2], nil
}

func getKeeper(s http.ResponseWriter, r *http.Request, q Queue)(types.MsgBody, error){
	queueName, err:= getQueueName(r)
	if err != nil {
		
	}
	m, err:=q.Get(queueName)

	a, err := json.Marshal(m)
	s.Write(a)
	return types.MsgBody{}, nil
}

