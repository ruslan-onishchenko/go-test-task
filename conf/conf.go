package conf

type Configuration struct {
	Port string
	QueueNumber int
	MessageInQueueNumber int
	TimeOut int
}

func ReadConfig(args[]string) (Configuration, error){
	return Configuration{
		Port: "9999",
		QueueNumber: 10,
		MessageInQueueNumber: 20,
		TimeOut: 5,
	}, nil
}