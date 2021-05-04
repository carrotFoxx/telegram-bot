package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"iot-services/host-server/src/common"
	"iot-services/host-server/src/host_manager"
	"log"
	"net/http"
)

func Ini() {

	setRoutes()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func setRoutes() {
	http.HandleFunc("/command", send_command)
}

func send_command(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	msg := common.Command_msg{}
	err = json.Unmarshal(bodyBytes, &msg)
	if err != nil {
		log.Println(err)
		_, _ = fmt.Fprint(w, "{'status':'err'}") //temp
		fmt.Println(msg)
	} else {
		host_manager.Send_message_from_outside(msg)
		fmt.Println(msg)
		_, _ = fmt.Fprint(w, "{'status':'ok'}") //todo можно добавить более подробный ответ
	}

}
