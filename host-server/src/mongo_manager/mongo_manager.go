package mongo_manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"iot-services/host-server/src/common"
	"log"
	"strconv"
	"sync"
	"time"
)

const IMMERS_BOX = "immers_box"

var mng *mongo.Client
var clientQerry *mongo.Client
var collectionQuerry *mongo.Collection

var mutex sync.Mutex
var docs []interface{}
var countMsg int = 0
var errorsMsg int = 0

type hostsStatusList_ map[int]int
type hosts_types_map_ map[int]string

var hostsStatusMap hostsStatusList_
var hosts_types_map hosts_types_map_

func AddHostValPackTimer(msg common.Common_msg) {
	//var doc_to_repository interface{}
	doc := bson.D{}
	//doc_to_repository := common.Immers_control_msg{}
	dev_type, _ := hosts_types_map[msg.Id]
	if dev_type == IMMERS_BOX {
		if msg.Type == 1 {
			doc_to_repository := common.Immers_control_msg{}
			err := json.Unmarshal([]byte(msg.Args), &doc_to_repository)
			if err != nil {
				log.Println(err)
			}
			//doc_to_repository.Rt = time.Now().UnixNano()
			//doc_to_repository.Id = host_id
			doc = bson.D{{"_id", time.Now().UnixNano()}, {"id", msg.Id}, {"type", msg.Type}, {"d", doc_to_repository}}
		} else {
			doc = bson.D{{"_id", time.Now().UnixNano()}, {"id", msg.Id}, {"type", msg.Type}, {"d", msg.Args}}
		}
	} else {
		return
	}

	mutex.Lock()
	docs = append(docs, doc)
	//docs = append(docs, doc_to_repository)
	mutex.Unlock()
}

func timer() {
	for {
		time.Sleep(100 * time.Duration(time.Millisecond))
		mutex.Lock()
		if len(docs) == 0 {
			mutex.Unlock()
			continue
		}

		collection := mng.Database("dev_data").Collection(strconv.FormatInt(int64(150), 10))
		insertManyResult, err := collection.InsertMany(context.TODO(), docs)
		if err != nil {
			//log.Fatal(err)
			log.Println(err)
		}
		_ = insertManyResult
		docs = docs[:0]
		mutex.Unlock()
	}
}

func CheckHostStatus(id int) (bool, error) { //функция обращается в бд devReg, коллекцию registredDevs

	type hostStat struct { //структура записи данных о статусе хстов в бд. Если меняется в БД, то меняем и тут
		Id     int
		Status string
	}

	var result hostStat
	filter := bson.D{{"id", id}}

	err := collectionQuerry.FindOne(context.TODO(), filter).Decode(&result)
	//	bs, err := collection.FindOne(context.TODO(), filter).DecodeBytes()
	if err != nil {
		fmt.Println(err)
	}
	err = errors.New("no host")
	if result.Status == "" {
		return false, err
	} else {
		if result.Status == "ok" {
			return true, nil
		} else {
			return false, nil
		}
	}
}

func MongoIni(uri string) {

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	client2, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	mng = client
	clientQerry = client2

	collection := clientQerry.Database("users").Collection("registred_devs") //todo переместить в переменные среды
	collectionQuerry = collection

	//exp_feelDBusers()
	GetUserStatusFromDBtoMap()
	go timer()
}

func GetUserStatusFromDBtoMap() { //загрузка статусов хостов для более быстрого доступа внутри программы
	hostsStatusMap = make(hostsStatusList_)
	hosts_types_map = make(hosts_types_map_)

	start := time.Now()

	type hostStat struct { //структура записи данных о статусе хстов в бд. Если меняется в БД, то меняем и тут
		Id     int
		Status string
		Type   string
	}
	fmt.Println("========")

	//var hostStats []*hostStat

	//cur, err := collectionQuerry.Find(context.TODO(), bson.D{{"status", "ok"}})//, findOptions)
	cur, err := collectionQuerry.Find(context.TODO(), bson.D{{}}) //, findOptions)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var elem hostStat
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
		}
		if elem.Status == "ok" {
			hostsStatusMap[elem.Id] = 0
		} else if elem.Status == "permDen" {
			hostsStatusMap[elem.Id] = 1
		} else {
			hostsStatusMap[elem.Id] = 2
		}
		hosts_types_map[elem.Id] = elem.Type
	}

	log.Println(len(hostsStatusMap), "host_manager status loaded from bd", time.Now().Sub(start))
}

func GetHostStatus(id int) int {
	if val, ok := hostsStatusMap[id]; ok {
		return val
	} else {
		return -1 //хост не найден
	}
}

func exp_feelDBusers() {
	for i := 0; i < 25000; i++ {
		insertResult, err := collectionQuerry.InsertOne(context.TODO(), bson.D{{"id", i}, {"status", "ok"}})
		if err != nil {
			//log.Fatal(err)
			log.Println(err)
		}
		log.Println(insertResult)

	}
}
