# coding: utf8
u"""Модуль формирования запросов к VK API."""


import requests


def method(method_name, values, access_token):
    u"""Отправка запроса к VK API."""
    request = "https://api.vk.com/method/"

    request += method_name
    request += "?access_token=" + access_token

    keys_values = values.keys()
    for key in keys_values:
        request += "&" + key + "=" + str(values[key])

    response = requests.post(request)

    return response.text
