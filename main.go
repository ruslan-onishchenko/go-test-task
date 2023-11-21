package main

import (
	"gotesttask/app"
	"gotesttask/conf"
	"gotesttask/queue"
	"log"
)

/*
	Запуск: ./gotesttask -p 9999 -mn 10 -qn 15 -t 3

	Параметры: ./gotesttask -h
	-mn - количество сообщений в очереди (default 10)
	-p  - порт (default "9999")
	-qn - количество очередей (default 10)
	-t  - тайм-аут (default 5)
*/
func main(){
	// парсим параметры
	conf, err := conf.ReadConfig()
	if err != nil {
		log.Println("Failed to read configuration: ", err.Error())
	}
	// Хранилище очередей реализовано отдельно от общего кода, так частично реализую гексагональную архитектуру
	queueImpl := queue.MakeImpl(conf.MessageInQueueNumber, conf.QueueNumber)
	// Сходу не продумал разделение логики и роутинг, потому у нас только один компонент queue.Impl,
	// который реализует интерфейс app.Queue
	err = app.Start(conf, queueImpl)
	if err != nil {
		log.Println("Error while running app: ", err.Error())
	}
	
	return
}