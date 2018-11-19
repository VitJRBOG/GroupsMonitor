# coding: utf8
u"""Модуль обработки данных."""


import vk_api
import data_manager
import input_data
import output_data
import thread_starter


def run_processing():
    u"""Запуск функций обработки."""
    dict_sessions = check_access_tokens()
    dict_threads = thread_starter.run_thread_starter(dict_sessions)
    user_answer_checker(dict_threads)


def check_access_tokens():
    u"""Проверяет валидность токенов доступа."""
    def authorization(access_token):
        u"""Авторизация в ВК."""
        vk_session = vk_api.VkApi(token=access_token)
        vk_session._auth_token()

        return vk_session

    def check_session(token_owner, access_token):
        u"""Проверяет валидность сессии."""
        vk_session = authorization(access_token)
        try:
            # КОСТЫЛЬ: проверка идет по id Павла Дурова
            values = {
                "user_ids": "1"
            }
            vk_session.method("users.get", values)
            sender = "Check " + token_owner + "'s access token"
            message = token_owner + "'s access token is valid."
            output_data.output_text_row(sender, message)
            return vk_session
        except Exception as message_error:
            if message_error == "invalid access_token":
                message = token_owner + "'s access token is invalid."
                access_token = request_new_access_token(token_owner)
                return check_session(token_owner, access_token)

    def request_new_access_token(token_owner):
        u"""Запрос нового токена доступа."""
        access_token = input_data.get_vk_user_token(token_owner)
        return access_token

    PATH = data_manager.read_path()
    dict_data = data_manager.read_json(PATH, "data")
    subjects = dict_data["subjects"]
    admin_session = check_session("Admin", dict_data["admin_access_token"])
    dict_sessions = {}
    for subject in subjects:
        vk_session = check_session(
            subject["name"], subject["sender_access_token"])
        dict_sessions.update({subject["name"]: vk_session})
    dict_sessions.update({"Admin": admin_session})

    return dict_sessions


def user_answer_checker(dict_threads):
    u"""Проверка команд пользователя."""
    while True:
        user_asnwer = raw_input()
        if user_asnwer == "quit":
            exit(0)
        if user_asnwer == "stop":
            for thread_data in dict_threads:
                thread_data["flag"] = False
                # пока временный пример. Изменю, когда опишу работу с потоками.
