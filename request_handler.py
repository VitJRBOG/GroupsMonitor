# coding: utf8
u"""Модуль обработки запросов к VK API."""


import exception_handler


# def request_wall_posts():
#     u"""Запрос постов на стене."""
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

def send_request(sender, data_for_request):
    u"""Отправка запроса к VK API."""
    error_repeats = 0
    try:
        vk_session = data_for_request["vk_session"]
        request = data_for_request["request"]
        values = data_for_request["values"]
        response = vk_session.method(request, values)
        return response
    except Exception as var_except:
        if error_repeats < 5:
            error_repeats += 1
        timeout = error_repeats * 2
        exception_handler.handling(sender, var_except, timeout)

        return send_request(sender, data_for_request)
