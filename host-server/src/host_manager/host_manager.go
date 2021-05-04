package host_manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"iot-services/host-server/src/common"
	"iot-services/host-server/src/kafka_producer_manager"
	"iot-services/host-server/src/mongo_manager"
	"iot-services/host-server/src/tcpServer"
	"log"
	"net"
)

type host struct { //данные по соеденению с хостом
	connHandle net.Conn // хендл для связи
}

type hosts_online_ map[int]host    //мап данных по соеденению с хостом (реестр онлайн хостов)
var hosts_online_map hosts_online_ //создаем мап для хранения реестра клиентов онлайн

//---------------------------------------------------------------------------------------
//ВХОД В МОДУЛЬ
//функция обработки поступающих сообщений. Все сообщения со всех потоков группируются здесь.
//Это сделано для того, что бы мы могли иметь реестр подключенных хостов для работы с ними
func HostsProc(url string, buffer int) {

	hosts_online_map = make(hosts_online_)
	var tcpServ = tcpServer.TcpServer{url}

	//answers := make(chan tcpServer.MsgPakage)
	answers := make(chan tcpServer.MsgPakage, buffer)

	go tcpServ.Open(answers)

	for {
		hostIncommingMsg(<-answers) //многопоточность не допускается!!
	}
}

func hostIncommingMsg(packet tcpServer.MsgPakage) { // функция работы с входящим пакетом от всех потоков tcp сокетов
	msg, err := checkParseAndRegisterHost(packet)
	if err != nil {
		log.Println(err)
		return
	}
	//key := []byte{0x39, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x31, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36}                                                                                 //тестовое сообщение
	//packet.Conn.Write(Encrypt(key, "{\"id\":1729243175,\"task\":115,\"arg\":\"многопоточность не использовать в пакетном режиме по таймеру проверка, подключен ли хост,0123456789 012345678\"}")) //тестовое сообщение длиной для u->v.scan_buf_counter_chars == 140+64+32+16 && lenght_str == 128+64+32+16
	if msg.Type != common.DEFAULT_MESSAGE {
		b, _ := json.Marshal(msg)
		kafka_producer_manager.Send_kafka_msg(b)
	}

	//================= хост прошел проверки, можно дальше с ним работать
	mongo_manager.AddHostValPackTimer(msg)
}

func checkParseAndRegisterHost(packet tcpServer.MsgPakage) (common.Common_msg, error) {
	key := []byte{0x39, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x31, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36}
	parsedMsg := common.Common_msg{}
	if !packet.Connected { //проверка, подключен ли хост, или это сообщение об отключенном хосте
		deleteHostFromRegByConn(packet.Conn)
		err := errors.New("")
		return parsedMsg, err
	}

	DecryptMsg := DecryptFromString(key, packet.Txt)
	err := json.Unmarshal(DecryptMsg, &parsedMsg)
	if err != nil { //проверка, парсится ли сообщение
		//packet.Conn.Close()
		fmt.Println(err)
		fmt.Println(DecryptMsg)
		err := errors.New("can't parse command_msg:")
		return parsedMsg, err
	}

	if host_, ok := hosts_online_map[parsedMsg.Id]; !ok { //проверяем наличие хоста в локальном реестре хостонлайн hostsOnline
		status := mongo_manager.GetHostStatus(parsedMsg.Id)
		if status == -1 {
			err := errors.New(fmt.Sprintf("ancknown host %d", parsedMsg.Id))
			packet.Conn.Close()
			return parsedMsg, err
		} else if status == 1 {
			err := errors.New("host permission denied ")
			packet.Conn.Close()
			return parsedMsg, err
		} else if status != 0 {
			err := errors.New("host error: other..")
			packet.Conn.Close()
			return parsedMsg, err
		}
		//==============================если все проверки пройдены - добавление в реестр хостов-онлайн (ключ - id хоста, значение - структура с данными)

		h := host{}                        //создаем структуру данных о хосте
		h.connHandle = packet.Conn         //помещаем в нее сокет
		hosts_online_map[parsedMsg.Id] = h //помещаем структуру в слайс хостов онлайн
		_ = host_.connHandle               //заглушка
		//fmt.Println(hostsOnline)
	} else {
		if host_.connHandle != packet.Conn { //если есть в реестре - проверяем, совпадают ли соккеты,
			// иначе мы можем отправить сообщение не тому юзеру, если подключился хост с таким же id
			HostDisconnect(parsedMsg.Id)                                //отключаем старй хост (он также удаляется из map)
			log.Printf("dublicate host_manager! id %d: ", parsedMsg.Id) //логгируем это недоразумение
			h := host{}                                                 //создаем структуру данных о хосте
			h.connHandle = packet.Conn
			hosts_online_map[parsedMsg.Id] = h
			fmt.Println("change host", hosts_online_map)
		}
	}
	return parsedMsg, nil
}

//---------------------------------------------------------------------------------------
//функция удаления хоста из реестра и отключения от сокета
func HostDisconnect(id int) {
	host := hosts_online_map[id]
	socket := host.connHandle
	socket.Close()
	delete(hosts_online_map, id)
}

//---------------------------------------------------------------------------------------
//функция удаления хоста из реестра по номеру соккета
func deleteHostFromRegByConn(conn net.Conn) {
	for id, sock := range hosts_online_map {
		if sock.connHandle == conn {
			delete(hosts_online_map, id)
		}
	}

}

//---------------------------------------------------------------------------------------
//Функция возвращает количество хостов онлайн
func GetOnlineHostsCount() int {
	return len(hosts_online_map)
}

//Функция отправляет хосту сообщение
func Send_message_from_outside(msg common.Command_msg) {
	if _host, ok := hosts_online_map[msg.Id]; ok {
		out, _ := json.Marshal(msg)
		send_msg_to_host(_host, string(out))
	} else {
		fmt.Println("host_offline")
	}
}

func send_msg_to_host(_host host, msg string) {
	key := []byte{0x39, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x31, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36}
	_host.connHandle.Write(Encrypt(key, msg))
}

//тестовое сообщение
