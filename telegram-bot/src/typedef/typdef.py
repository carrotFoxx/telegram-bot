import os
from multiprocessing import Queue

MONGO_DSN = os.getenv('MONGO_DSN', 'mongodb://localhost:27017/')
KAFKA_DSN = os.getenv('KAFKA_DSN', 'localhost:30992')
IOT_SERVER_DSN = os.getenv('IOT_SERVER_DSN', 'http://localhost:8080')
INCOMING_TOPIC = os.getenv('INCOMING_TOPIC', 'msg_from_host')
API_TOKEN = os.getenv('API_TOKEN', '')

user_by_chat_id = dict()
device_by_device_id = dict()
users_by_device_id = dict()
sending_command_by_chat_id = dict()


mongo_m = None

message_queue = Queue()

def add_message_queue(msg):
    global message_queue
    message_queue.put(msg)

# def get_message_queue():
#     global message_queue
#     return message_queue.get()



DEFAULT_MESSAGE = 1
WARNING_MESSAGE = 2
ERROR_MESSAGE = 3
INFO_MESSAGE = 4
SERVICE_MESSAGE = 5

IMMERS_BOX = "immersion"

action_map_immer_box = {
    "OW_SEARCH": 1,
    "OW_READ_ADDRESES": 2,
    "RIG_1_ON": 3,
    "RIG_1_OFF": 4,
    "RIG_2_ON": 5,
    "RIG_2_OFF": 6,
    "RIG_RESET_ERROR": 7,
    "SET_UPPER_TEMP_LIMIT": 8,
    "SET_LOWER_TEMP_LIMIT": 9,
    "SET_TCP_SEND_INTERVAL": 10,
    "SET_AUTO_RIG_CONTROL": 11,
    "SET_POWER_COUNTER_PRESCALER_BY_CURRENT_VALUE": 12,
    "SET_AP": 13,
    "SET_HOST": 14,
    "SET_SERVER_URL": 15,
    "PUSH_POWER_BUTTON_A": 16,
    "PUSH_POWER_BUTTON_B": 17,
    "SET_POWER_COUNTER_PRESCALER_BY_PRESCALER_VALUE": 18,
}

HELP_MESSAGE_BY_IMMERSY_COMMAND = '"OW_SEARCH": 1,\r\n' \
                                  '"OW_READ_ADDRESES": 2, \r\n' \
                                  '"RIG_1_ON": 3, \r\n' \
                                  '"RIG_1_OFF": 4, \r\n' \
                                  '"RIG_2_ON": 5, \r\n' \
                                  '"RIG_2_OFF": 6, \r\n' \
                                  '"RIG_RESET_ERROR": 7, \r\n' \
                                  '"SET_UPPER_TEMP_LIMIT": 8, \r\n' \
                                  '"SET_LOWER_TEMP_LIMIT": 9, \r\n' \
                                  '"SET_TCP_SEND_INTERVAL": 10, \r\n' \
                                  '"SET_AUTO_RIG_CONTROL": 11, \r\n' \
                                  '"SET_POWER_COUNTER_PRESCALER_BY_CURRENT_VALUE": 12, \r\n' \
                                  '"SET_AP": 13, \r\n' \
                                  '"SET_HOST": 14, \r\n' \
                                  '"SET_SERVER_URL": 15, \r\n' \
                                  '"PUSH_POWER_BUTTON_A": 16, \r\n' \
                                  '"PUSH_POWER_BUTTON_B": 17, \r\n' \
                                  '"SET_POWER_COUNTER_PRESCALER_BY_PRESCALER_VALUE": 18,'

action_list_imer_box = ['RIG_1_ON',
                        'RIG_1_OFF',
                        'RIG_2_ON',
                        'RIG_2_OFF',
                        'RIG_RESET_ERROR',
                        'SET_UPPER_TEMP_LIMIT',
                        'SET_LOWER_TEMP_LIMIT',
                        'SET_TCP_SEND_INTERVAL',
                        'SET_AUTO_RIG_CONTROL',
                        'SET_POWER_COUNTER_PRESCALER',
                        'SET_AP',
                        'SET_HOST',
                        "SET_SERVER_URL",
                        "PUSH_POWER_BUTTON_A",
                        "PUSH_POWER_BUTTON_B",
                        ]

# define NONE 0
# define OW_SEARCH 1
# define OW_READ_ADDRESES 2
# define RIG_1_ON 3
# define RIG_1_OFF 4
# define RIG_2_ON 5
# define RIG_2_OFF 6
# define RIG_RESET_ERROR 7
# define SET_UPPER_TEMP_LIMIT 8
# define SET_LOWER_TEMP_LIMIT 9
# define SET_TCP_SEND_INTERVAL 10
# define SET_AUTO_RIG_CONTROL 11
# define SET_POWER_COUNTER_PRESCALER 12
# define SET_AP 13
# define SET_HOST 14
