package conf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Configuration struct {
	Port					string	`json:"port"`
	QueueNumber				int		`json:"queue_number"`
	MessageInQueueNumber	int		`json:"message_in_queue_number"`
	TimeOut					int		`json:"timeout"`
}

func ReadConfig(args[]string) (Configuration, error){
	jsonFile, err := os.Open("./conf.json")
	if err != nil {
		return Configuration{}, fmt.Errorf("Config file opening error ", err)
	}
	defer jsonFile.Close()
	byteFileContent, err := io.ReadAll(jsonFile)
	if err != nil {
		return Configuration{}, fmt.Errorf("config file read error: ", err)
	}
	conf := Configuration{}
	err = json.Unmarshal(byteFileContent, &conf)
	if err != nil {
		return Configuration{}, fmt.Errorf("config file decoding error: ", err)
	}
	return conf, nil
	
	// return Configuration{
	// 	Port: "9999",
	// 	QueueNumber: 10,
	// 	MessageInQueueNumber: 20,
	// 	TimeOut: 5,
	// }, nil
}