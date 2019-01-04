# coding: utf8
u"""Модуль обработки запросов к VK API."""

import urllib2
import vkapi
import exception_handler


def request_wall_posts(sender, subject_data, monitor_data):
    u"""Запрос постов на стене."""
    def make_data_for_request(subject_data, monitor_data):
        u"""Подготовка данных для отправки запроса."""
        values = {
            "access_token": subject_data["access_tokens"]["admin"],
            "method": "wall.get",
            "values": {
                "owner_id": subject_data["owner_id"],
                "count": monitor_data["posts_count"],
                "filter": monitor_data["filter"],
                "v": 5.92
            }
        }
        return values
    
    def select_data_from_response(response):
        u"""Извлекает данные из словаря с результатами запроса."""
        wall_posts_data = []
        items = response["items"]
        for item in items:
            values = {
                "id": item["id"],
                "owner_id": item["owner_id"],
                "from_id": item["from_id"],
                "date": item["date"],
                "post_type": item["post_type"],
                "text": item["text"]
            }
            if "signer_id" in item:
                values.update({"signer_id": item["signer_id"]})
            if "attachments" in item:
                attachments = []
                for attachment in item["attachments"]:
                    type_attachment = attachment["type"]
                    values_attachments = {
                        "owner_id": attachment[type_attachment]["owner_id"],
                        "id": attachment[type_attachment]["id"],
                        "type": type_attachment
                    }
                    if "access_key" in attachment[type_attachment]:
                        values_attachments.update(
                            {"access_key": attachment[type_attachment]["access_key"]})
                    attachments.append(values_attachments)
                values.update({"attachments": attachments})
            wall_posts_data.append(values)
        return wall_posts_data

    sender += " -> Get wall posts"

    data_for_request = make_data_for_request(subject_data, monitor_data)
    response = send_request(sender, data_for_request)
    wall_posts_data = select_data_from_response(response)
    return wall_posts_data


# def request_album_photos():
#     u"""Запрос фотографий в альбомах."""
# def request_videos():
#     u"""Запрос видеороликов."""
# def request_photo_comments():
#     u"""Запрос комментариев под фотографиями."""
# def request_video_comments():
#     u"""Запрос комментариев под видеороликами."""
# def request_topic_comments():
#     u"""Запрос комментариев в обсуждениях."""
# def request_wall_post_comments():
#     u"""Запрос комментариев под постами на стене."""


def request_user_info(sender, subject_data, data_for_request):
    u"""Запрос информации о пользователе."""
    def make_data_for_request(subject_data, data_for_request):
        u"""Подготовка данных для отправки запроса."""
        values = {
            "access_token": subject_data["access_tokens"]["admin"],
            "method": "users.get",
            "values": {
                "user_ids": data_for_request["user_ids"],
                "v": 5.92
            }
        }
        return values

    def select_data_from_response(response):
        u"""Извлекает данные из словаря с результатами запроса."""
        users_info = []
        items = response
        for item in items:
            values = {
                "id": item["id"],
                "first_name": item["first_name"],
                "last_name": item["last_name"]
            }
            users_info.append(values)
        return users_info

    data_for_request = make_data_for_request(subject_data, data_for_request)
    response = send_request(sender, data_for_request)
    users_info = select_data_from_response(response)
    return users_info


def request_group_info(sender, subject_data, data_for_request):
    u"""Запрос информации о сообществе."""
    def make_data_for_request(subject_data, data_for_request):
        u"""Подготовка данных для отправки запроса."""
        values = {
            "access_token": subject_data["access_tokens"]["admin"],
            "method": "groups.getById",
            "values": {
                "group_ids": data_for_request["group_ids"],
                "v": 5.92
            }
        }
        return values

    def select_data_from_response(response):
        u"""Извлекает данные из словаря с результатами запроса."""
        groups_info = []
        items = response
        for item in items:
            values = {
                "id": item["id"],
                "name": item["name"],
                "screen_name": item["screen_name"]
            }
            groups_info.append(values)
        return groups_info
    
    data_for_request = make_data_for_request(subject_data, data_for_request)
    response = send_request(sender, data_for_request)
    groups_info = select_data_from_response(response)
    return groups_info


def send_request(sender, data_for_request):
    u"""Алгоритмы отправки запроса к VK API."""
    def req(sender, data_for_request, error_repeats):
        u"""Отправка запроса к VK API."""
        access_token = data_for_request["access_token"]
        method = data_for_request["method"]
        values = data_for_request["values"]
        result = vkapi.method(method, values, access_token)
        if "response" in result:
            return result["response"]
        else:
            message_error = result["error"]["error_msg"]
            if error_repeats < 5:
                error_repeats += 1
            timeout = error_repeats * 2
            exception_handler.handling(sender, message_error, timeout)

            return req(sender, data_for_request, error_repeats)

    error_repeats = 0
    response = req(sender, data_for_request, error_repeats)
    return response


def send_message(sender, data_for_message, access_token):
    u"""Алгоритмы отправки сообщения в ВК."""
    def make_message(data_for_message):
        u"""Сборка сообщения перед отправкой."""
        text = urllib2.quote(data_for_message["text"])
        values = {
            "message": text,
            "v": 5.68
        }
        if "attachment" in data_for_message:
            values.update({"attachment": data_for_message["attachment"]})
        return values

    def send(sender, values, access_token):
        u"""Отправка сообщения."""
        error_repeats = 0
        result = vkapi.method("messages.send", values, access_token)
        if "response" in result:
            return result["response"]
        else:
            message_error = result["error"]["error_msg"]
            if error_repeats < 5:
                error_repeats += 1
            timeout = error_repeats * 2
            exception_handler.handling(sender, message_error, timeout)

            return send(sender, values, access_token)

    sender += " -> Send message"
    values = make_message(data_for_message)
    send_to = data_for_message["send_to"]
    for addressee in send_to:
        values.update({"peer_id": addressee})
        send(sender, values, access_token)
