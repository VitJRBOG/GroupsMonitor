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


def through_vk_api(method_name, values, access_token):
    u"""Отравка запроса к VK API через стороннюю библиотеку."""
    def get_session(access_token):
        vk_session = vk_api.VkApi(token=access_token)
        vk_session._auth_token()

        return vk_session

    vk_session = get_session(access_token)
    result = vk_session.method(method_name, values)

    return result
