import logging
from asyncio import sleep
from aiokafka import AIOKafkaConsumer
from typedef.typdef import add_message_queue, INCOMING_TOPIC, KAFKA_DSN


async def consumer():
    while True:
        try:
            consumer = AIOKafkaConsumer(
                INCOMING_TOPIC,
                bootstrap_servers=KAFKA_DSN,
                group_id="my-group")
            await consumer.start()
        except:
            await sleep(10)
            continue
        try:
            # Consume messages
            async for msg in consumer:
                print("consumed: ", msg.topic, msg.partition, msg.offset,
                      msg.key, msg.value, msg.timestamp)
                add_message_queue(msg.value)
        except:
            logging.error("kafka_error")
            await sleep(10)
            continue


