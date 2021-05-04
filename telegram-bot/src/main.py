import asyncio
import logging
from threading import Thread
from bot.tg_bot import bot_ini
from kafka_manager.kafka_manager import consumer
from mongo_manager.mongo_manager import MongoManager
from typedef import typdef
from typedef.typdef import MONGO_DSN


logging.basicConfig(level=logging.INFO)


def kafka_consumer():
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    result = loop.run_until_complete(consumer())

if __name__ == '__main__':
    logging.info("tg_bot_service starting")
    while True:
        try:
            typdef.mongo_m = MongoManager(MONGO_DSN, 'users', 'tg_users', "dev_data", "150")
            break
        except:
            continue
    th = Thread(target=kafka_consumer)
    th.start()
    bot_ini()







