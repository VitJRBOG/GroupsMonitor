# coding: utf8
u"""Модуль алгоритмов проверки."""


import datetime
import data_manager
import request_handler
import output_data


def wall_posts_monitor(sender, res_filename, subject_data, monitor_data):
    u"""Проверяет посты на стене."""
    def sort_posts(posts):
        u"""Сортировка постов методом пузырька."""
        for j in range(len(posts) - 1):
            f = 0
            for i in range(len(posts) - 1 - j):
                if posts[i]["date"] < posts[i + 1]["date"]:
                    x = posts[i]
                    y = posts[i + 1]
                    posts[i + 1] = x
                    posts[i] = y
                    f = 1
            if f == 0:
                break
        return posts

    def found_new_post(sender, values, subject_data, post):
        u"""Алгоритмы обработки целевого поста."""
        def send_post(sender, values, post):
            u"""Отправка данных из целевого поста."""
            def make_user_signature(user_signature, user_id):
                u"""Собирает подпись пользователя."""
                data_for_request = {
                    "user_ids": user_id
                }
                users_info = request_handler.request_user_info(
                    sender, subject_data, data_for_request)
                user_signature += "*id" + str(users_info[0]["id"])
                user_signature += " (" + users_info[0]["first_name"]
                user_signature += " " + \
                    users_info[0]["last_name"] + ")"
                return user_signature

            def make_group_signature(group_signature, group_id):
                u"""Собирает подпись сообщества."""
                data_for_request = {
                    "group_ids": int(str(group_id)[1:])
                }
                groups_info = request_handler.request_group_info(
                    sender, subject_data, data_for_request)
                group_signature += "*" + groups_info[0]["screen_name"]
                group_signature += " (" + groups_info[0]["name"] + ")"
                return group_signature

            def select_owner_signature(post):
                u"""Выбирает из словаря данные и формирует гиперссылку на владельца поста."""
                owner_signature = ""
                owner_id = post["owner_id"]
                if str(owner_id)[0] == "-":
                    owner_id = post["owner_id"]
                    owner_signature = make_group_signature(owner_signature, owner_id)
                else:
                    owner_id = post["owner_id"]
                    owner_signature = make_user_signature(
                        owner_signature, owner_id)
                return owner_signature

            def select_author_signature(post):
                u"""Выбирает из словаря данные и формирует гиперссылку на автора."""
                author_signature = ""
                owner_id = post["owner_id"]
                if str(owner_id)[0] == "-":
                    if "signer_id" in post:
                        author_id = post["signer_id"]
                        author_signature = make_user_signature(
                            author_signature, author_id)
                    else:
                        author_signature += "[no data]"
                else:
                    author_id = post["owner_id"]
                    author_signature += make_user_signature(
                        author_signature, author_id)
                return author_signature
            def select_attachments(post):
                u"""Выбирает из словаря данные и формирует список прикреплений."""
                attachments = post["attachments"]
                media_items = ""
                for i, attachment in enumerate(attachments):
                    if attachment["type"] == "photo" or\
                       attachment["type"] == "video" or\
                       attachment["type"] == "audio" or\
                       attachment["type"] == "doc" or\
                       attachment["type"] == "poll":
                        media_items += attachment["type"]
                        media_items += str(attachment["owner_id"])
                        media_items += "_" + str(attachment["id"])
                        if "access_key" in attachment:
                            media_items += "?access_key=" + attachment["access_key"]
                        if len(attachments) > 1 and i < len(attachments) - 1:
                            media_items += ","
                return media_items
            def select_post_url(post):
                u"""Выбирает из словаря данные и формирует URL на пост."""
                post_url = "https://vk.com/wall"
                post_url += str(post["owner_id"]) + "_" + str(post["id"])
                return post_url
            def select_date(post):
                u"""Выбирает из словаря данные и определяет дату поста."""
                def ts_date_to_str(ts_date, date_format):
                    u"""Получение даты в читабельном формате."""
                    str_date = datetime.datetime.fromtimestamp(ts_date).strftime(date_format)
                    return str_date
                ts_date = post["date"]
                str_date = ts_date_to_str(ts_date, "%d.%m.%Y %H:%M:%S")
                return str_date

            data_for_message = {}

            owner_signature = select_owner_signature(post)
            author_signature = select_author_signature(post)
            post_url = select_post_url(post)
            publication_date = select_date(post)

            text = ""
            text += "Location: " + owner_signature.encode("utf8") + "\n"
            text += "Author: " + author_signature.encode("utf8") + "\n"
            text += "Created: " + str(publication_date).encode("utf8") + "\n\n"
            text += post["text"].encode("utf8") + "\n\n"
            text += post_url

            send_to = values["send_to"]
            access_token = values["access_token"]

            data_for_message = {
                "send_to": send_to,
                "text": text
            }

            if "attachments" in post:
                media_items = select_attachments(post)
                data_for_message.update({"attachment": media_items})

            request_handler.send_message(sender, data_for_message, access_token)

        def update_last_date(values, post):
            u"""Обновление даты последнего проверенного поста."""
            res_filename = values["res_filename"]
            path_to_res_file = values["path_to_res_file"]
            dict_json = data_manager.read_json(path_to_res_file, res_filename)
            dict_json["last_date"] = str(item["date"])
            data_manager.write_json(path_to_res_file, res_filename, dict_json)

        def show_message_about_new_post(sender, post):
            u"""Алгоритмы отображения сообщения о новом посте."""
            def ts_date_to_str(ts_date, date_format):
                u"""Получение даты в читабельном формате."""
                str_date = datetime.datetime.fromtimestamp(
                    ts_date).strftime(date_format)
                return str_date
            str_date = ts_date_to_str(post["date"], "%d.%m.%Y %H:%M:%S")
            message = "New " + post["post_type"] + " for " + str_date + "."
            output_data.output_text_row(sender, message)

        send_post(sender, values, post)
        update_last_date(values, post)
        show_message_about_new_post(sender, post)

    sender += " -> Wall posts monitor"

    wall_posts_data = request_handler.request_wall_posts(
        sender, subject_data, monitor_data)
    wall_posts = sort_posts(wall_posts_data)

    last_date = int(monitor_data["last_date"])

    new_posts = []

    for item in reversed(wall_posts):
        if item["date"] > last_date:
            new_posts.append(item)

    PATH = data_manager.read_path()
    path_to_res_file = PATH + subject_data["path"] + "/"
    access_token = subject_data["access_tokens"][res_filename]
    send_to = monitor_data["send_to"]

    values = {
        "res_filename": res_filename,
        "path_to_res_file": path_to_res_file,
        "send_to": send_to,
        "access_token": access_token
    }
    
    for item in new_posts:
        found_new_post(sender, values, subject_data, item)
    

# def album_photos_monitor():
#     u"""Проверяет фотографии в альбомах."""
# def videos_monitor():
#     u"""Проверяет видео."""
# def photo_comments_monitor():
#     u"""Проверяет комментарии под фотографиями."""
# def video_comments_monitor():
#     u"""Проверяет комментарии под видео."""
# def topic_comments_monitor():
#     u"""Проверяет комментарии в обсуждениях."""
# def wall_post_comments_monitor():
#     u"""Проверяет комментарии под постами на стене."""
