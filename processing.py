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
        sender = "Check " + token_owner + "'s access token"
        try:
            # КОСТЫЛЬ: проверка идет по id Павла Дурова
            values = {
                "user_ids": "1"
            }
            vk_session.method("users.get", values)
            message = token_owner + "'s access token is valid."
            output_data.output_text_row(sender, message)
            return vk_session
        except Exception as message_error:
            if str(message_error).lower().find("invalid access_token") > -1:
                message = token_owner + "'s access token is invalid. Need another..."
                output_data.output_text_row(sender, message)
                access_token = request_new_access_token(token_owner)
                update_access_token(token_owner, access_token)
                return check_session(token_owner, access_token)

    def request_new_access_token(token_owner):
        u"""Запрос нового токена доступа."""
        access_token = input_data.get_vk_user_token(token_owner)
        return access_token

    def update_access_token(token_owner, access_token):
        u"""Обновляет токен доступа в файле с данными."""
        PATH = data_manager.read_path()
        dict_data = data_manager.read_json(PATH, "data")
        if token_owner == "Admin":
            dict_data["admin_access_token"] = access_token
        else:
            subjects = dict_data["subjects"]
            for i, subject in enumerate(subjects):
                if subject["name"] == token_owner:
                    dict_data["subjects"][i]["sender_access_token"] = access_token
        data_manager.write_json(PATH, "data", dict_data)

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
