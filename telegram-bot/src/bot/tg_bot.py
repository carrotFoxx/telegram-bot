import asyncio
import datetime
import json
import logging
import time
import typing
from asyncio import sleep

from aiogram import Bot, Dispatcher, executor, types
from aiogram.contrib.fsm_storage.memory import MemoryStorage
from aiogram.contrib.middlewares.logging import LoggingMiddleware
from aiogram.utils.callback_data import CallbackData
from aiogram.utils.exceptions import MessageNotModified

from http_.http_ import set_request, send_command_by_device_id
from typedef import typdef
from typedef.typdef import user_by_chat_id, device_by_device_id, action_list_imer_box, IMMERS_BOX, \
    users_by_device_id, message_queue, ERROR_MESSAGE, INFO_MESSAGE, sending_command_by_chat_id, \
    HELP_MESSAGE_BY_IMMERSY_COMMAND, WARNING_MESSAGE, API_TOKEN

logging.basicConfig(level=logging.INFO)
log = logging.getLogger('broadcast')



bot = Bot(token=API_TOKEN, parse_mode=types.ParseMode.HTML)
storage = MemoryStorage()
dp = Dispatcher(bot, storage=storage)
dp.middleware.setup(LoggingMiddleware())

posts_cb = CallbackData('post', 'id', 'action')  # post:<id>:<action>
def get_main_keyboard(chat_id) -> types.InlineKeyboardMarkup:
    devices = get_devices_by_chat_id(chat_id)
    markup = types.InlineKeyboardMarkup()
    for device in devices:
        markup.add(
            types.InlineKeyboardButton(
                device["name"],
                callback_data=posts_cb.new(id=device["id"], action='view')),
        )
    return markup


async def format_post(dev_id: int) -> (str, types.InlineKeyboardMarkup):
    text = await get_sensor_data(dev_id)
    markup = types.InlineKeyboardMarkup()
    rows = get_buttons(dev_id)
    for row in rows:
        markup.row(*row)

    markup.add(types.InlineKeyboardButton('обновить', callback_data=posts_cb.new(id=dev_id, action='refresh')))
    markup.add(types.InlineKeyboardButton('<< назад', callback_data=posts_cb.new(id=dev_id, action='list')))
    return text, markup


def get_buttons(dev_id):
    rows = list()
    if device_by_device_id.get(dev_id) is not None:
        if device_by_device_id[dev_id]["type"] == "immersion":
            rows.append([
                types.InlineKeyboardButton('риг А вкл', callback_data=posts_cb.new(id=dev_id, action='RIG_1_ON')),
                types.InlineKeyboardButton('риг А выкл', callback_data=posts_cb.new(id=dev_id, action='RIG_1_OFF')),
                types.InlineKeyboardButton('риг А кнопка',
                                           callback_data=posts_cb.new(id=dev_id, action='PUSH_POWER_BUTTON_A'))
            ])
            rows.append([
                types.InlineKeyboardButton('риг Б вкл', callback_data=posts_cb.new(id=dev_id, action='RIG_2_ON')),
                types.InlineKeyboardButton('риг Б выкл', callback_data=posts_cb.new(id=dev_id, action='RIG_2_OFF')),
                types.InlineKeyboardButton('риг Б кнопка',
                                           callback_data=posts_cb.new(id=dev_id, action='PUSH_POWER_BUTTON_B'))
            ])
            rows.append([
                types.InlineKeyboardButton('послать команду', callback_data=posts_cb.new(id=dev_id, action='BUTTON_SEND_COMMAND')),
                types.InlineKeyboardButton('помощь', callback_data=posts_cb.new(id=dev_id, action='BUTTON_HELP_BY_COMMANS')),

            ])
            return rows


@dp.message_handler(commands='start')
async def cmd_start(query: types.CallbackQuery):
    if user_by_chat_id.get(query.chat.id) is None:
        await query.answer("Enter your token pls")
    else:
        keyboard = types.ReplyKeyboardMarkup(resize_keyboard=True)
        button_1 = types.KeyboardButton(text="Мои устройства")
        keyboard.add(button_1)
        text = "Hello " + user_by_chat_id.get(query.chat.id)['username'] + "!\r\nHere your devices:";
        await query.answer(text, reply_markup=get_main_keyboard(query.chat.id))
        await query.answer(reply_markup=keyboard)


@dp.callback_query_handler(posts_cb.filter(action='list'))
async def query_show_list(query: types.CallbackQuery):
    await query.message.edit_text('Устройства', reply_markup=get_main_keyboard(query.message.chat.id))


@dp.callback_query_handler(posts_cb.filter(action='refresh'))
async def query_refresh_data(query: types.CallbackQuery, callback_data: typing.Dict[str, str]):
    text, markup = await format_post(int(callback_data['id']))
    await query.message.edit_text(text, reply_markup=markup)


@dp.callback_query_handler(posts_cb.filter(action='view'))
async def query_view(query: types.CallbackQuery, callback_data: typing.Dict[str, str]):
    dev_id = int(callback_data['id'])
    dev = device_by_device_id.get(dev_id)
    if not dev:
        return await query.answer('Unknown device!')

    text, markup = await format_post(dev_id)
    await query.message.edit_text(text, reply_markup=markup)


@dp.callback_query_handler(posts_cb.filter(action=action_list_imer_box))
async def query_push_button(query: types.CallbackQuery, callback_data: typing.Dict[str, str]):
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    await set_request(callback_data)

@dp.callback_query_handler(posts_cb.filter(action=["BUTTON_SEND_COMMAND"]))
async def query_push_button(query: types.CallbackQuery, callback_data: typing.Dict[str, str]):
    sending_command_by_chat_id[query.message.chat.id] = callback_data['id']
    await query.message.answer("enter your command:")

@dp.callback_query_handler(posts_cb.filter(action=["BUTTON_HELP_BY_COMMANS"]))
async def query_push_button(query: types.CallbackQuery, callback_data: typing.Dict[str, str]):
    await query.message.answer(HELP_MESSAGE_BY_IMMERSY_COMMAND)


@dp.errors_handler(exception=MessageNotModified)
async def message_not_modified_handler(update, error):
    return True  # errors_handler must return True if error was handled correctly


@dp.message_handler()
async def echo(query: types.CallbackQuery):
    if user_by_chat_id.get(query.chat.id) is None:
        if await typdef.mongo_m.change_user_by_token(query.text, query.message.chat.id):
            await query.message.answer("Welcome!", reply_markup=get_main_keyboard(query.message.chat.id))
        else:
            await query.message.answer("Enter your token pls")
    else:
        if query.text is not None and query.text.startswith(':'):
            if sending_command_by_chat_id.get(query.chat.id) is not None:
                await send_command_by_device_id(sending_command_by_chat_id.get(query.chat.id), query.text)

        else:
            text = "Hello " + user_by_chat_id.get(query.chat.id)['username'] + "!\r\nHere your devices:";
            await query.answer(text, reply_markup=get_main_keyboard(query.chat.id))
    try:
        del sending_command_by_chat_id[query.chat.id]
    except:
        pass


async def refresher():
    while True:
        await typdef.mongo_m.get_maps()
        print("refesh_data_from_db")
        await sleep(60)


async def scheduler():
    while True:
        while message_queue.qsize() != 0:
            try:
                dev_msg = json.loads(message_queue.get())
            except:
                pass
            # todo здесь можно классифицировать и модфицировать сообщение перед отправкой
            #🟢🔴⚠ℹ‼
            msg = '' #⚠
            if dev_msg.get('t') == ERROR_MESSAGE:
                msg += '‼ '
            if dev_msg.get('t') == WARNING_MESSAGE:
                msg += '⚠ '
            if dev_msg.get('t') == INFO_MESSAGE:
                msg += 'ℹ '

            msg += dev_msg.get('args').replace('{', '').replace('}', '')
            print(msg)
            if dev_msg.get('id') is not None:
                try:
                    users = users_by_device_id.get(dev_msg['id'])
                    if users is not None:
                        for user in users:
                            dev_name = ''
                            try:
                                if user["notification"] >= dev_msg['t']:
                                    for device in user['devices']:
                                        if device['id'] == dev_msg['id']:
                                            dev_name = device['name']
                            except:
                                pass
                            await send_tg_message(user['chat_id'], '<b>'+dev_name +'</b>: ' + msg)
                except Exception as e:
                    logging.error(e)
        await sleep(0.1)


async def on_startup(_):
    asyncio.create_task(scheduler())
    asyncio.create_task(refresher())


def bot_ini():
    executor.start_polling(dp, skip_updates=True, on_startup=on_startup)
    pass


async def set_new_user(token, chat_id):
    await typdef.mongo_m.change_user_by_token(token, chat_id)


def get_devices_by_chat_id(chat_id):
    user = typdef.user_by_chat_id.get(chat_id)
    return user["devices"]


async def get_sensor_data(dev_id):
    dev_data = await typdef.mongo_m.get_devs_data(dev_id)
    if dev_data is not None:
        if device_by_device_id.get(dev_id).get('type') == IMMERS_BOX:  # для иммерсионной утсановки
            last_msg_time = time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(dev_data["_id"] / 1000000000))
            back_time = (datetime.datetime.now() - datetime.datetime.fromtimestamp(dev_data["_id"] / 1000000000)).seconds
            temp_inside_tank = dev_data["d"]["sens0"]
            online_indicator = '🟢' if back_time<120 else '🔴'
            temp_incoming_water = dev_data["d"]["sens1"]
            temp_outcomming_water = dev_data["d"]["sens2"]
            rig_power = dev_data["d"]["w1"]
            rig_enrgy = dev_data["d"]["wh1"]
            rigA_on = dev_data["d"]["rig1"]
            rigB_on = dev_data["d"]["rig2"]
        switch = {1: 'вкл', 0: 'выкл'}
        answ = "Время последнего показания: <b>" + last_msg_time + "</b>  "+ online_indicator+"\r\n" \
               + "Температура ванны: <b>" + str(temp_inside_tank / 10) + "</b> грС\r\n" \
               + "Температура входящего теплоносителя: <b>" + str(temp_incoming_water / 10) + "</b> грС\r\n" \
               + "Температура исходящего теплоносителя: <b>" + str(temp_outcomming_water / 10) + "</b> грС\r\n" \
               + "Моментальная мощность установки: <b>" + str(rig_power) + "</b> Ватт\r\n" \
               + "Количество потребленной энергии: <b>" + str(round(rig_enrgy / 1000, 3)) + "</b> ВтЧ\r\n" \
               + "Питание рига А: <b>" + switch[rigA_on] + "</b>\r\n" \
               + "Питание рига Б: <b>" + switch[rigB_on] + "</b>"
        return answ
    else:
        return "Нет данных 🔴"


async def send_tg_message(chat_id: int, text: str, disable_notification: bool = False) -> bool:
    """
    Safe messages sender
    :param chat_id:
    :param text:
    :param disable_notification:
    :return:
    """
    try:
        await bot.send_message(chat_id, text, disable_notification=disable_notification)
        # await bot.send_message(chat_id, text)
    # except exceptions.BotBlocked:
    #     log.error(f"Target [ID:{chat_id}]: blocked by user")
    # except exceptions.ChatNotFound:
    #     log.error(f"Target [ID:{chat_id}]: invalid user ID")
    # except exceptions.RetryAfter as e:
    #     log.error(f"Target [ID:{chat_id}]: Flood limit is exceeded. Sleep {e.timeout} seconds.")
    #     await asyncio.sleep(e.timeout)
    #     return await send_kafka_message(chat_id, text)  # Recursive call
    # except exceptions.UserDeactivated:
    #     log.error(f"Target [ID:{chat_id}]: user is deactivated")
    # except exceptions.TelegramAPIError:
    #     log.exception(f"Target [ID:{chat_id}]: failed")
    except Exception as e:
        pass
    else:
        log.info(f"Target [ID:{chat_id}]: success")
        return True
    return False
