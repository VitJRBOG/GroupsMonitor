# coding: utf-8


import os
import logger
import datamanager
import vk_api
# import core


class Start():
    def path_checking(self, sender):
        sender += " -> File \"path.txt\" checking"

        if os.path.exists("path.txt") is False:
            file_text = open("path.txt", "w")
            file_text.write("")
            file_text.close()

            mess_for_log = "Was created file \"path.txt\"."
            logger.message_output(sender, mess_for_log)

    def log_file_checking(self, sender, PATH):
        sender += " -> File \"log.txt\" checking"

        if os.path.exists(PATH + "log.txt") is False:

            datamanager.write_text(sender, PATH, "log", "")

    def data_checking(self, sender, PATH):
        sender += " -> File \"data.json\" checking"

        datafile_was_created = False

        if os.path.exists(PATH + "data.json") is False:

            mess_for_log = "WARNING! File \"data.json\" not found!"
            logger.message_output(sender, mess_for_log)

            data_json = {
                "bot_token": "",
                "admin_token": "",
                "subjects": [
                    {
                        "name": "",
                        "file_name": "",
                        "path": "",
                        "check_subject": 0,
                        "interval": 60
                    }
                ]
            }

            datamanager.write_json(sender, PATH, "data", data_json)

            subject_json = {
                "name": "",
                "wiki_database_id": "-0_0",
                "total_last_date": "0",
                "topics": [],
                "owner_id": 0,
                "post_checker_settings": {
                    "check_posts": 1,
                    "posts_count": 1,
                    "last_date": "0",
                    "filter": "all",
                    "send_to": 0
                },
                "topic_checker_settings": {
                    "post_count": 1,
                    "check_topics": 0,
                    "send_to": 0
                },
                "photo_checker_settings": {
                    "last_date": "0",
                    "photo_count": 1,
                    "check_photo": 0,
                    "send_to": 0
                },
                "photo_comments_checker_settings": {
                    "comment_count": 1,
                    "check_comments": 0,
                    "last_date": "0",
                    "send_to": 0
                },
                "post_comments_checker_settings": {
                    "comment_count": 1,
                    "filter": "all",
                    "posts_count": 1,
                    "check_comments": 0,
                    "last_date": "0",
                    "send_to": 0,
                    "check_by_attachments": 0,
                    "check_by_keywords": 0,
                    "keywords": []
                }
            }

            datamanager.write_json(sender, PATH, "template", subject_json)

            datafile_was_created = True

        return datafile_was_created

    def tokens_checking(self, sender, PATH):
        sender += " -> Tokens checking"

        #  Получение данных из файла JSON
        data_json = datamanager.read_json(sender, PATH, "data")

        admin_token_validity = True
        bot_token_validity = True

        if len(str(data_json["admin_token"])) <= 1:
            admin_token_validity = False
        if len(str(data_json["bot_token"])) <= 1:
            bot_token_validity = False

        token_validity = {
            "admin_token": admin_token_validity,
            "bot_token": bot_token_validity
        }

        return data_json, token_validity


def update_token(sender, PATH, data_json, token_validity, tokens):

    admin_token_validity = token_validity["admin_token"]
    bot_token_validity = token_validity["bot_token"]

    admin_token = tokens["admin_token"]
    bot_token = tokens["bot_token"]

    if not admin_token_validity:
        data_json["admin_token"] = admin_token
    if not bot_token_validity:
        data_json["bot_token"] = bot_token

    datamanager.write_json(sender, PATH, "data", data_json)

    return data_json


def autorization(sender, data_access, auth_type):
    sender += " -> Authorization"

    try:

        if auth_type == "token":
            #  Авторизация по токену
            access_token = data_access["token"]
            vk_session = vk_api.VkApi(token=access_token)
            vk_session._auth_token()

        elif auth_type == "login":
            #  Авторизация по имени пользователя и паролю
            vk_login = data_access["login"]
            vk_passwd = data_access["password"]
            vk_session = vk_api.VkApi(login=vk_login, password=vk_passwd)
            vk_session.auth()

        else:
            mess_for_log = "Error of authorization. Exit from program..."
            logger.message_output(sender, mess_for_log)

            exit(0)

        return vk_session

    except Exception as var_except:
        logger.exception_handler(sender, var_except)
        return autorization(data_access, auth_type)
