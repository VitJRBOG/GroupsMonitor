# coding: utf8
u"""Модуль вспомогательных функций модулей проверки."""


import datetime
import request_handler


def sort_items(items):
    u"""Сортировка постов методом пузырька."""
    for j in range(len(items) - 1):
        f = 0
        for i in range(len(items) - 1 - j):
            if items[i]["date"] < items[i + 1]["date"]:
                x = items[i]
                y = items[i + 1]
                items[i + 1] = x
                items[i] = y
                f = 1
        if f == 0:
            break
    return items


def make_user_signature(sender, subject_data, user_signature, user_id):
    u"""Собирает подпись пользователя."""
    data_for_request = {
        "user_ids": user_id
    }
    users_info = request_handler.request_user_info(
        sender, subject_data, data_for_request)
    user_signature += "*id" + str(users_info[0]["id"])
    user_signature += " (" + users_info[0]["first_name"]
    user_signature += " " + users_info[0]["last_name"] + ")"

    return user_signature


def make_group_signature(sender, subject_data, group_signature, group_id):
    u"""Собирает подпись сообщества."""
    data_for_request = {
        "group_ids": int(str(group_id)[1:])
    }
    groups_info = request_handler.request_group_info(
        sender, subject_data, data_for_request)
    group_signature += "*" + groups_info[0]["screen_name"]
    group_signature += " (" + groups_info[0]["name"] + ")"

    return group_signature


def ts_date_to_str(ts_date, date_format):
    u"""Получение даты в читабельном формате."""
    str_date = datetime.datetime.fromtimestamp(
        ts_date).strftime(date_format)
    return str_date
