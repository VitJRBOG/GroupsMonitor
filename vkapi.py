# coding: utf8
u"""Модуль формирования запросов к VK API."""


import json
import requests
import vk_api


def method(method_name, values, access_token):
    u"""Отправка запроса к VK API."""
    request = "https://api.vk.com/method/"

    request += method_name
    request += "?access_token=" + access_token

    server_answer = requests.post(request, values)

    str_result = server_answer.text

    result = json.loads(str_result)

    return result
