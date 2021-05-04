package service

import (
	"IotHostServer/src/hosts"
	"IotHostServer/src/mongo_manager"
	"fmt"
	"time"
)

func Monitorini() {
	go printState()
	go updateUserStatusFromDBtoMap()
}

func printState() {
	for {
		fmt.Println("Hosts online ", host_manager.GetOnlineHostsCount())
		time.Sleep(5 * time.Duration(time.Second))
	}
}

func updateUserStatusFromDBtoMap() {
	for {
		time.Sleep(5 * time.Duration(time.Minute))
		mongo_manager.GetUserStatusFromDBtoMap()

	}
}
