# coding: utf8
u"""Модуль обработки данных."""


import os
import time
import initialization
import vkapi
import exception_handler
import data_manager
import input_data
import output_data
import thread_starter


def run_processing():
    u"""Запуск функций обработки."""
    sender = "Starting"
    need_presetting = initialization.check_res_files()
    if need_presetting:
        message = "Data base is not full. Need presetting. Quit..."
        output_data.output_text_row(sender, message)
    else:
        dict_tokens = check_access_tokens()
        data_threads = thread_starter.run_thread_starter(dict_tokens)
        user_answer_checker(data_threads)


def check_access_tokens():
    u"""Алгоритм проверки валидности токенов доступа."""

    def check_token(token_owner, token_purpose, access_token):
        u"""Проверяет валидность токена."""
        sender = "Check " + token_owner + "'s access token for " + token_purpose

        # КОСТЫЛЬ. 22822305 - id паблика команды ВК
        values = {
            "group_id": 22822305,
            "v": 5.92
        }
        result = vkapi.method("groups.getOnlineStatus", values, access_token)
        # КОСТЫЛЬ. Проверка валидности токена с помощью определения статуса другого паблика.

        if "response" in result:
            return access_token
        else:
            message_error = result["error"]["error_msg"]
            invalid_token_errors = [
                "invalid access_token",
                "access_token was given to another ip address",
                "access_token has expired"
            ]
            for i, text_error in enumerate(invalid_token_errors):
                if str(message_error).lower().find(text_error) > -1:
                    message = "Access token is invalid. Need another..."
                    output_data.output_text_row(sender, message)
                    access_token = request_new_access_token(
                        token_owner, token_purpose)
                    update_access_token(token_owner, token_purpose, access_token)
                    return check_token(token_owner, token_purpose, access_token)
                else:
                    if i == len(invalid_token_errors) - 1:
                        exception_handler.handling(sender, message_error, 0)

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

    sender = "Checking access tokens"
    message = "Please stand by..."
    output_data.output_text_row(sender, message)

    PATH = data_manager.read_path()
    dict_data = data_manager.read_json(PATH, "data")
    subjects = dict_data["subjects"]
    dict_tokens = {}
    for subject in subjects:
        if subject["monitor_subject"] == 1:
            token_purposes = subject["access_tokens"].keys()
            values = {}
            for token_purpose in token_purposes:
                need_check_token = False
                if token_purpose != "admin":
                    path_to_subject = PATH + subject["path"] + "/"
                    monitor_settings = data_manager.read_json(
                        path_to_subject, token_purpose)
                    if monitor_settings["need_monitoring"] == 1:
                        need_check_token = True
                else:
                    need_check_token = True
                if need_check_token is True:
                    access_token = check_token(
                        subject["name"], token_purpose,
                        subject["access_tokens"][token_purpose])
                    values.update({token_purpose: access_token})
            dict_tokens.update({subject["name"]: values})

    return dict_tokens


def user_answer_checker(data_threads):
    u"""Проверка команд пользователя."""
    while True:
        sender = "Main menu"
        user_asnwer = raw_input()
        if user_asnwer == "quit":
            message = "Force quit..."
            output_data.output_text_row(sender, message)
            exit(0)
        if user_asnwer == "stop":
            message = "Stopping threads..."
            output_data.output_text_row(sender, message)
            waiting_time = 10
            for data_thread in data_threads:
                thread_sender = data_thread["sender"]
                data_thread["end_flag"].set()
                if data_thread["end_flag"].isSet() and\
                   data_thread["was_turned_on"]:
                    for i in range(waiting_time):
                        time.sleep(1)
                        if not data_thread["thread"].isAlive():
                            message = "OK! Monitoring is stopped..."
                            output_data.output_text_row(thread_sender, message)
                            break
                        else:
                            if i == waiting_time:
                                message = "WARNING! Monitoring cannot be stopped..."
                                output_data.output_text_row(thread_sender, message)
            message = "Quit..."
            output_data.output_text_row(sender, message)
            exit(0)
