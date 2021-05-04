import asyncio

import motor.motor_asyncio
from bson import ObjectId

from typedef.typdef import user_by_chat_id, device_by_device_id, users_by_device_id


class MongoManager:
    def __init__(self, dsn, user_bd, users, dev_bd, devs):
        self.async_client = motor.motor_asyncio.AsyncIOMotorClient(dsn)
        self.user_db = self.async_client[user_bd]
        self.dev_db = self.async_client[dev_bd]
        self.users_collection = self.user_db[users]
        self.dev_data_collection = self.dev_db[devs]
        loop = asyncio.get_event_loop()
        loop.run_until_complete(self.get_maps())

    async def get_maps(self):
        device_by_device_id.clear()
        users_by_device_id.clear()
        user_by_chat_id.clear()

        cursor = self.users_collection.find()
        while (await cursor.fetch_next):
            user = cursor.next_object()
            user_by_chat_id[user.get("chat_id")] = user
            for dev in user["devices"]:
                device_by_device_id[dev["id"]] = dev
                if users_by_device_id.get(dev["id"]) is None:
                    users_by_device_id[dev["id"]] = list()
                    users_by_device_id[dev["id"]].append(user)
                else:
                    users_by_device_id[dev["id"]].append(user)
            print(user)
        pass

    async def change_user_by_token(self, token, chat_id):
        try:
            await self.users_collection.find_one_and_update({"_id": ObjectId(token)}, {'$set': {'chat_id': chat_id}})
            await self.get_maps()
            return True
        except:
            return False

    async def get_devs_data(self, dev_id: int):
        cursor = self.dev_data_collection.find({"id": dev_id,"type":1}).sort("_id",-1).limit(1)
        while (await cursor.fetch_next):
            doc = cursor.next_object()
            return doc
        return None
