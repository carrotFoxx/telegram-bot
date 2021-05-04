# IotHostServer
IotHostServer

ENV VARS:

    MONGO_DSN - mongodb://localhost:27017
    KAFKA_DSN - localhost:30992
    OUTGOING_TOPIC - msg_from_host
    INCOMING_TOPIC - msg_to_host
    GROUP_ID - group

Для создания образа:
docker build -t iotserver .

Далее этот образ будет использоваться при запуске telebot

    MONGO_DSN=mongodb://localhost:27017;KAFKA_DSN=localhost:30992;OUTGOING_TOPIC=msg_from_host;INCOMING_TOPIC=msg_to_host;GROUP_ID=group