package main

import (
	"gotesttask/app"
	"gotesttask/conf"
	"gotesttask/queue"
	"log"
	"os"
)

func main(){
	conf, err := conf.ReadConfig(os.Args)
	if err != nil {
		log.Println("Failed to read configuration: ", err.Error())
	}

	queueImpl := queue.MakeImpl(conf.MessageInQueueNumber, conf.QueueNumber)
	
	err = app.Start(conf, queueImpl)
	if err != nil {
		log.Println("Error while running app: ", err.Error())
	}
	
	return
}