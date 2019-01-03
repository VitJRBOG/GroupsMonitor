# coding: utf8
u"""Модуль обработки данных."""


import vkapi
import data_manager
import input_data
import output_data
import thread_starter


def run_processing():
    u"""Запуск функций обработки."""
    dict_tokens = check_access_tokens()
    data_threads = thread_starter.run_thread_starter(dict_tokens)
    user_answer_checker(data_threads)


def check_access_tokens():
    u"""Алгоритм проверки валидности токенов доступа."""

    def check_token(token_owner, token_purpose, access_token):
        u"""Проверяет валидность токена."""
        sender = "Check " + token_owner + "'s access token for " + token_purpose
        try:
            # КОСТЫЛЬ: проверка идет по id Павла Дурова
            values = {
                "user_ids": "1"
            }
            vkapi.method("users.get", values, access_token)
            return access_token
        except Exception as message_error:
            if str(message_error).lower().find("invalid access_token") > -1:
                message = token_owner + "'s access token for " + \
                    token_purpose.replace("_", " ") + " is invalid. Need another..."
                output_data.output_text_row(sender, message)
                access_token = request_new_access_token(
                    token_owner, token_purpose)
                update_access_token(token_owner, token_purpose, access_token)
                return check_token(token_owner, token_purpose, access_token)

    def request_new_access_token(token_owner, token_purpose):
        u"""Запрос нового токена доступа."""
        access_token = input_data.get_vk_user_token(
            token_owner, token_purpose.replace("_", " "))
        return access_token

    def update_access_token(token_owner, token_purpose, access_token):
        u"""Обновляет токен доступа в файле с данными."""
        PATH = data_manager.read_path()
        dict_data = data_manager.read_json(PATH, "data")
        subjects = dict_data["subjects"]
        for i, subject in enumerate(subjects):
            if subject["name"] == token_owner:
                dict_data["subjects"][i]["access_tokens"][token_purpose] = \
                    access_token
        data_manager.write_json(PATH, "data", dict_data)

    PATH = data_manager.read_path()
    dict_data = data_manager.read_json(PATH, "data")
    subjects = dict_data["subjects"]
    dict_tokens = {}
    for subject in subjects:
        if subject["monitor_subject"] == 1:
            token_purposes = subject["access_tokens"].keys()
            values = {}
            for token_purpose in token_purposes:
                access_token = check_token(
                    subject["name"], token_purpose, subject["access_tokens"][token_purpose])
                values.update({token_purpose: access_token})
            dict_tokens.update({subject["name"]: values})

    return dict_tokens


def user_answer_checker(data_threads):
    u"""Проверка команд пользователя."""
    while True:
        sender = "[Main menu]"
        user_asnwer = raw_input()
        if user_asnwer == "quit":
            message = "Force quit..."
            output_data.output_text_row(sender, message)
            exit(0)
        if user_asnwer == "stop":
            message = "Stopping threads..."
            output_data.output_text_row(sender, message)
            for data_thread in data_threads:
                thread_sender = data_thread["sender"]
                data_thread["end_flag"].set()
                if not data_thread["thread"].isAlive() and\
                   data_thread["end_flag"].isSet():
                    message = "OK! Monitoring is stopped..."
                    output_data.output_text_row(thread_sender, message)
                if data_thread["thread"].isAlive() and\
                   data_thread["end_flag"].isSet():
                    message = "WARNING! Monitoring cannot be stopped..."
            message = "Quit..."
            output_data.output_text_row(sender, message)
            exit(0)
