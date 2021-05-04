import aiohttp

from typedef.typdef import action_map_immer_box, IMMERS_BOX, device_by_device_id, IOT_SERVER_DSN





async def set_request(callback_data):
    payload = dict()
    #dev_data = await typdef.mongo_m.get_devs_data(int(callback_data.get('id')))
    dev_data = device_by_device_id.get(int(callback_data.get('id')))
    if dev_data is not None:
        if dev_data["type"] == IMMERS_BOX: #формирование сообщения для эммирсионной установки
            payload['id'] = int(callback_data.get('id'))
            payload['task'] = action_map_immer_box[callback_data["action"]]
            payload['arg'] = ""

        await send_payload(payload)
    # print(callback_data)
    # print(payload)


async def send_command_by_device_id(device_id: int, command: str): #todo добавить права использования комманд
    #todo интерпретировать команды
    command = command[1:]
    command = command.split(' ')
    task = int(command[0])
    payload = dict()
    dev_data = device_by_device_id.get(int(device_id))
    if dev_data["type"] == IMMERS_BOX: #формирование сообщения для эммирсионной установки
        payload['id'] = int(device_id)
        payload['task'] = task
        if len(command)>1:
            payload['arg'] = command[1]
        else:
            payload['arg'] = ''
    await send_payload(payload)


async def send_payload(payload):
    async with aiohttp.ClientSession() as session:
        headers = {
            'Content-Type': 'application/json'
        }
        print(str(payload).replace("'", '"'))
        async with session.put(IOT_SERVER_DSN + '/command', data=str(payload).replace("'", '"'),
                               headers=headers) as resp:
            # print(resp.status)
            print(await resp.text())
