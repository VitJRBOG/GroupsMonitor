# coding: utf8


import vk_api
import logger
import model.datamanager as datamanager


def update_sessions_list(sessions_list):
    sender = "VK authorization"
    PATH = datamanager.read_path()
    loads_json = datamanager.read_json(PATH + "bot_notificator/", "data")
    sessions_list["admin_session"] =\
        get_session(loads_json["admin_token"])
    message = "Session of admin has been succesfully created."
    logger.message_output(sender, message)
    sessions_list["bot_session"] =\
        get_session(loads_json["bot_token"])
    message = "Session of bot has been succesfully created."
    logger.message_output(sender, message)
    return sessions_list


def get_session(access_token):
    vk_session = vk_api.VkApi(token=access_token)
    vk_session._auth_token()

    return vk_session
