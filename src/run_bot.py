import datamanager
import starter
import logger
import core


sender = "Starter -> Starting"

PATH = datamanager.read_path(sender)

data_json = datamanager.read_json(sender, PATH, "data")

vk_admin_token = data_json["admin_token"]
vk_bot_token = data_json["bot_token"]

data_access_admin = {
    "token": vk_admin_token
}
data_access_bot = {
    "token": vk_bot_token
}

vk_admin_session = starter.autorization(sender, data_access_admin, "token")
vk_bot_session = starter.autorization(sender, data_access_bot, "token")

mess_for_log = "Program was started."
logger.message_output(sender, mess_for_log)

core.main(vk_admin_session, vk_bot_session)
