package main

import (
	"iot-services/host-server/src/api"
	"iot-services/host-server/src/host_manager"
	"iot-services/host-server/src/mongo_manager"
	"os"
)

func main() {
	mongo_manager.MongoIni(os.Getenv("MONGO_DSN"))

	go host_manager.HostsProc("0.0.0.0:10000", 100000) //функция инициализации управления хостами
	//go kafka_consumer_manager.Kafka_consumer_service()
	api.Ini()
	/*
		for{
			runtime.GC()
			time.Sleep(1*time.Duration(time.Millisecond))
		}*/

	//service.Monitorini()
	//for {
	//	var temp string
	//	fmt.Fscan(os.Stdin, &temp)
	//	if temp == "g" {
	//		fmt.Println("GC")
	//		runtime.GC()
	//	}
	//}

}
