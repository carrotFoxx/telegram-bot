package common

type Common_msg struct { //общая структура сообщений
	Id   int    `json:"id"`
	Type int    `json:"t"`
	Args string `json:"args"`
}

type Immers_control_msg struct { //структура сообщения от контроллера рига
	//Rt int64  `bson:"_id" json:"id,omitempty"`
	Id    int
	Sens0 int16 //датчик температуры
	Sens1 int16
	Sens2 int16
	W1    uint16 //Ватт
	W2    uint16
	W3    uint16
	WH1   float32 //Ватт*час
	WH2   float32
	WH3   float32
	Rig1  uint8 // вкл/выкл фермы
	Rig2  uint8
}

type Command_msg struct { //входящее сообщение для отправления команды в хост
	Id   int    `json:"id"`
	Task int    `json:"task"`
	Arg  string `json:"arg"`
}

/*
1 - default_msg
2 - error_msg
3 - info_msg
4 - service_msg
*/

const DEFAULT_MESSAGE = 1
const ERROR_MESSAGE = 2
const WARNING_MESSAGE = 3
const INFO_MESSAGE = 4
const SERVICE_MESSAGE = 5
