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
                    if type_attachment == "photo" or\
                       type_attachment == "video" or\
                       type_attachment == "audio" or\
                       type_attachment == "doc" or\
                       type_attachment == "poll":
                        values_attachments = {
                            "owner_id": attachment[type_attachment]["owner_id"],
                            "id": attachment[type_attachment]["id"],
                            "type": type_attachment
                        }
                        if "access_key" in attachment[type_attachment]:
                            values_attachments.update(
                                {"access_key": attachment[type_attachment]["access_key"]})
                        attachments.append(values_attachments)
                if len(attachments) > 0:
                    values.update({"attachments": attachments})
            wall_posts_data.append(values)
        return wall_posts_data

    sender += " -> Get wall posts"

    data_for_request = make_data_for_request(subject_data, monitor_data)
    response = send_request(sender, data_for_request)
    wall_posts_data = select_data_from_response(response)
    return wall_posts_data


def request_album_photos(sender, subject_data, monitor_data):
    u"""Запрос фотографий в альбомах."""
    def make_data_for_request(subject_data, monitor_data):
        u"""Подготовка данных для отправки запроса."""
        values = {
            "access_token": subject_data["access_tokens"]["admin"],
            "method": "photos.getAll",
            "values": {
                "owner_id": subject_data["owner_id"],
                "count": monitor_data["photo_count"],
                "no_service_albums": 1,
                "v": 5.92
            }
        }
        return values
    
    def select_data_from_response(response):
        u"""Извлекает данные из словаря с результатами запроса."""
        album_photos_data = []
        items = response["items"]
        for item in items:
            values = {
                "id": item["id"],
                "album_id": item["album_id"],
                "owner_id": item["owner_id"],
                "date": item["date"],
                "text": item["text"]
            }
            if "user_id" in item:
                values.update({"user_id": item["user_id"]})
            album_photos_data.append(values)
        return album_photos_data

    sender += " -> Get album photos"

    data_for_request = make_data_for_request(subject_data, monitor_data)
    response = send_request(sender, data_for_request)
    album_photos_data = select_data_from_response(response)
    return album_photos_data


def request_videos(sender, subject_data, monitor_data):
    u"""Запрос видеороликов."""
    def make_data_for_request(subject_data, monitor_data):
        u"""Подготовка данных для отправки запроса."""
        values = {
            "access_token": subject_data["access_tokens"]["admin"],
            "method": "video.get",
            "values": {
                "owner_id": subject_data["owner_id"],
                "count": monitor_data["video_count"],
                "v": 5.92
            }
        }
        return values

    def select_data_from_response(response):
        u"""Извлекает данные из словаря с результатами запроса."""
        videos_data = []
        items = response["items"]
        for item in items:
            values = {
                "id": item["id"],
                "owner_id": item["owner_id"],
                "date": item["date"],
                "description": item["description"]
            }
            videos_data.append(values)
        return videos_data
    
    sender += " -> Get videos"

    data_for_request = make_data_for_request(subject_data, monitor_data)
    response = send_request(sender, data_for_request)
    videos_data = select_data_from_response(response)

    return videos_data


def request_photo_comments(sender, subject_data, monitor_data):
    u"""Запрос комментариев под фотографиями."""
    def make_data_for_request(subject_data, monitor_data):
        u"""Подготовка данных для отправки запроса."""
        values = {
            "access_token": subject_data["access_tokens"]["admin"],
            "method": "photos.getAllComments",
            "values": {
                "owner_id": subject_data["owner_id"],
                "count": monitor_data["comment_count"],
                "v": 5.92
            }
        }
        return values

    def select_data_from_response(response):
        u"""Извлекает данные из словаря с результатами запроса."""
        photo_comments_data = []
        items = response["items"]
        for item in items:
            values = {
                "id": item["id"],
                "from_id": item["from_id"],
                "pid": item["pid"],
                "date": item["date"],
                "text": item["text"]
            }
            if "attachments" in item:
                attachments = []
                for attachment in item["attachments"]:
                    type_attachment = attachment["type"]
                    if type_attachment == "photo" or\
                       type_attachment == "video" or\
                       type_attachment == "audio" or\
                       type_attachment == "doc" or\
                       type_attachment == "poll":
                        values_attachments = {
                            "owner_id": attachment[type_attachment]["owner_id"],
                            "id": attachment[type_attachment]["id"],
                            "type": type_attachment
                        }
                        if "access_key" in attachment[type_attachment]:
                            values_attachments.update(
                                {"access_key": attachment[type_attachment]["access_key"]})
                        attachments.append(values_attachments)
                if len(attachments) > 0:
                    values.update({"attachments": attachments})
            photo_comments_data.append(values)
        return photo_comments_data

    sender += " -> Get photo comments"

    data_for_request = make_data_for_request(subject_data, monitor_data)
    response = send_request(sender, data_for_request)
    photo_comments_data = select_data_from_response(response)

    return photo_comments_data


def request_video_comments(sender, subject_data, monitor_data, video):
    u"""Запрос комментариев под видеороликами."""
    def make_data_for_request(subject_data, monitor_data, video):
        u"""Подготовка данных для отправки запроса."""
        values = {
            "access_token": subject_data["access_tokens"]["admin"],
            "method": "video.getComments",
            "values": {
                "owner_id": video["owner_id"],
                "video_id": video["id"],
                "count": monitor_data["comment_count"],
                "v": 5.92
            }
        }
        return values

    def select_data_from_response(response):
        u"""Извлекает данные из словаря с результатами запроса."""
        video_comments_data = []
        items = response["items"]
        for item in items:
            values = {
                "id": item["id"],
                "from_id": item["from_id"],
                "date": item["date"],
                "text": item["text"]
            }
            if "attachments" in item:
                attachments = []
                for attachment in item["attachments"]:
                    type_attachment = attachment["type"]
                    if type_attachment == "photo" or\
                       type_attachment == "video" or\
                       type_attachment == "audio" or\
                       type_attachment == "doc" or\
                       type_attachment == "poll":
                        values_attachments = {
                            "owner_id": attachment[type_attachment]["owner_id"],
                            "id": attachment[type_attachment]["id"],
                            "type": type_attachment
                        }
                        if "access_key" in attachment[type_attachment]:
                            values_attachments.update(
                                {"access_key": attachment[type_attachment]["access_key"]})
                        attachments.append(values_attachments)
                if len(attachments) > 0:
                    values.update({"attachments": attachments})
            video_comments_data.append(values)
        return video_comments_data

    sender += " -> Get video comments"

    data_for_request = make_data_for_request(subject_data, monitor_data, video)
    response = send_request(sender, data_for_request)
    if response == "access error":
        return []
    video_comments_data = select_data_from_response(response)

    return video_comments_data


def request_topic_comments(sender, subject_data, monitor_data, topic):
    u"""Запрос комментариев в обсуждениях."""
    def make_data_for_request(subject_data, monitor_data, topic):
        u"""Подготовка данных для отправки запроса."""
        values = {
            "access_token": subject_data["access_tokens"]["admin"],
            "method": "board.getComments",
            "values": {
                "group_id": topic["owner_id"],
                "topic_id": topic["id"],
                "sort": "desc",
                "count": monitor_data["post_count"],
                "v": 5.92
            }
        }
        return values
    
    def select_data_from_response(response):
        u"""Извлекает данные из словаря с результатами запроса."""
        topic_comments_data = []
        items = response["items"]
        for item in items:
            values = {
                "id": item["id"],
                "from_id": item["from_id"],
                "date": item["date"],
                "text": item["text"]
            }
            if "attachments" in item:
                attachments = []
                for attachment in item["attachments"]:
                    type_attachment = attachment["type"]
                    if type_attachment == "photo" or\
                       type_attachment == "video" or\
                       type_attachment == "audio" or\
                       type_attachment == "doc" or\
                       type_attachment == "poll":
                        values_attachments = {
                            "owner_id": attachment[type_attachment]["owner_id"],
                            "id": attachment[type_attachment]["id"],
                            "type": type_attachment
                        }
                        if "access_key" in attachment[type_attachment]:
                            values_attachments.update(
                                {"access_key": attachment[type_attachment]["access_key"]})
                        attachments.append(values_attachments)
                if len(attachments) > 0:
                    values.update({"attachments": attachments})
            topic_comments_data.append(values)
        return topic_comments_data
    
    sender += " -> Get topic comments"

    data_for_request = make_data_for_request(subject_data, monitor_data, topic)
    response = send_request(sender, data_for_request)
    topic_comments_data = select_data_from_response(response)

    return topic_comments_data


def request_wall_post_comments(sender, subject_data, monitor_data, wall_post):
    u"""Запрос комментариев под постами на стене."""
    def make_data_for_request(subject_data, monitor_data, wall_post):
        u"""Подготовка данных для отправки запроса."""
        values = {
            "access_token": subject_data["access_tokens"]["admin"],
            "method": "wall.getComments",
            "values": {
                "owner_id": wall_post["owner_id"],
                "post_id": wall_post["id"],
                "sort": "desc",
                "count": monitor_data["comment_count"],
                "thread_items_count": monitor_data["thread_items_count"],
                "v": 5.92
            }
        }
        return values
    
    def select_data_from_response(response):
        u"""Извлекает данные из словаря с результатами запроса."""
        wall_post_comments_data = []
        items = response["items"]
        for item in items:
            values = {
                "id": item["id"],
                "owner_id": item["owner_id"],
                "date": item["date"],
                "from_id": item["from_id"],
                "post_id": item["post_id"],
                "text": item["text"]
            }
            if "attachments" in item:
                attachments_data = select_attachments(item["attachments"])
                if len(attachments_data) > 0:
                    values.update({"attachments": attachments_data})
            wall_post_comments_data.append(values)

            if len(item["thread"]["items"]) > 0:
                thread_items = []
                for thread_item in item["thread"]["items"]:
                    values = {
                        "id": thread_item["id"],
                        "owner_id": thread_item["owner_id"],
                        "parents_stack": thread_item["parents_stack"],
                        "date": thread_item["date"],
                        "from_id": thread_item["from_id"],
                        "post_id": thread_item["post_id"],
                        "text": thread_item["text"]
                    }
                    if "attachments" in thread_item:
                        thread_attachments_data = select_attachments(
                            thread_item["attachments"])
                        if len(thread_attachments_data) > 0:
                            values.update({"attachments": thread_attachments_data})
                    thread_items.append(values)
                wall_post_comments_data.extend(thread_items)
        return wall_post_comments_data

    sender += " -> Get wall post comments"

    data_for_request = make_data_for_request(subject_data, monitor_data, wall_post)
    response = send_request(sender, data_for_request)
    wall_post_comments_data = select_data_from_response(response)

    return wall_post_comments_data


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


def request_photo_album_info(sender, subject_data, data_for_request):
    u"""Запрос информации об альбоме с фотографиями."""
    def make_data_for_request(subject_data, data_for_request):
        u"""Подготовка данных для отправки запроса."""
        values = {
            "access_token": subject_data["access_tokens"]["admin"],
            "method": "photos.getAlbums",
            "values": {
                "owner_id": data_for_request["owner_id"],
                "album_ids": data_for_request["album_ids"],
                "v": 5.92
            }
        }
        return values
    
    def select_data_from_response(response):
        u"""Извлекает данные из словаря с результатами запроса."""
        albums_info = []
        items = response["items"]
        for item in items:
            values = {
                "id": item["id"],
                "title": item["title"]
            }
            albums_info.append(values)
        return albums_info

    data_for_request = make_data_for_request(subject_data, data_for_request)
    response = send_request(sender, data_for_request)
    albums_info = select_data_from_response(response)
    return albums_info


def request_topics_info(sender, subject_data, monitor_data):
    u"""Запрос информации о топиках обсуждений."""
    def make_data_for_request(subject_data, monitor_data):
        u"""Подготовка данных для отправки запроса."""
        values = {
            "access_token": subject_data["access_tokens"]["admin"],
            "method": "board.getTopics",
            "values": {
                "group_id": int(str(subject_data["owner_id"])[1:]),
                "count": monitor_data["topics_count"],
                "v": 5.92
            }
        }
        return values

    def select_data_from_response(response, subject_data):
        u"""Извлекает данные из словаря с результатами запроса."""
        topics_data = []
        items = response["items"]
        for item in items:
            values = {
                "id": item["id"],
                "owner_id": int(str(subject_data["owner_id"])[1:]),
                "title": item["title"]
            }
            topics_data.append(values)
        return topics_data

    sender += " -> Get topics info"

    data_for_request = make_data_for_request(subject_data, monitor_data)
    response = send_request(sender, data_for_request)
    topics_data = select_data_from_response(response, subject_data)

    return topics_data


def select_attachments(attachments):
    u"""Извлекает медиаконтент из элемента."""
    attachments_data = []
    for attachment in attachments:
        type_attachment = attachment["type"]
        if type_attachment == "photo" or\
            type_attachment == "video" or\
            type_attachment == "audio" or\
            type_attachment == "doc" or\
            type_attachment == "poll":
            values = {
                "owner_id": attachment[type_attachment]["owner_id"],
                "id": attachment[type_attachment]["id"],
                "type": type_attachment
            }
            if "access_key" in attachment[type_attachment]:
                values.update(
                    {"access_key": attachment[type_attachment]["access_key"]})
            attachments_data.append(values)
    return attachments_data



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
            access_errors = [
                "comments for this video are closed"
            ]
            message_error = result["error"]["error_msg"]
            for access_error in access_errors:
                if message_error.lower().find(access_error) != -1:
                    return "access error"
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
