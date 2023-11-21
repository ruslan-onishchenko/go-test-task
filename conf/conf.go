package conf

import (
	"flag"
)

type Configuration struct {
	Port					int		`json:"port"`
	QueueNumber				int		`json:"queue_number"`
	MessageInQueueNumber	int		`json:"message_in_queue_number"`
	TimeOut					int		`json:"timeout"`
}

// ReadConfig - парсинг параметров командной строки.
// Дефолт - ./gotesttask -p 9999 -mn 10 -qn 10 -t 5
func ReadConfig() (Configuration, error){
	port:=					flag.Int("p", 9999, "порт")
	queueNumber:=			flag.Int("qn", 10, "количество очередей")
	messageInQueueNumber:=	flag.Int("mn", 10, "количество сообщений в очереди")
	timeOut:= 				flag.Int("t", 5, "тайм-аут")
	flag.Parse()
	return Configuration{
		Port:					*port,
		QueueNumber:			*queueNumber,
		MessageInQueueNumber:	*messageInQueueNumber,
		TimeOut: 				*timeOut,
	}, nil
}