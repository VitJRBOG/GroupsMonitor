# coding: utf8
u"""Модуль алгоритмов проверки."""


import datetime
import data_manager
import request_handler
import output_data


def wall_posts_monitor(sender, res_filename, subject_data, monitor_data):
    u"""Проверяет посты на стене."""
    def found_new_post(sender, values, subject_data, post):
        u"""Алгоритмы обработки целевого поста."""
        def send_post(sender, values, subject_data, post):
            u"""Отправка данных из целевого поста."""
            def select_owner_signature(sender, subject_data, post):
                u"""Выбирает из словаря данные и формирует гиперссылку на владельца поста."""
                owner_signature = ""
                owner_id = post["owner_id"]
                if str(owner_id)[0] == "-":
                    owner_id = post["owner_id"]
                    owner_signature = make_group_signature(
                        sender, subject_data, owner_signature, owner_id)
                else:
                    owner_id = post["owner_id"]
                    owner_signature = make_user_signature(
                        sender, subject_data, owner_signature, owner_id)
                return owner_signature

            def select_author_signature(sender, subject_data, post):
                u"""Выбирает из словаря данные и формирует гиперссылку на автора."""
                author_signature = ""
                owner_id = post["owner_id"]
                if str(owner_id)[0] == "-":
                    if "signer_id" in post:
                        author_id = post["signer_id"]
                        author_signature = make_user_signature(
                            sender, subject_data, author_signature, author_id)
                    else:
                        author_signature += "[no data]"
                else:
                    author_id = post["owner_id"]
                    author_signature += make_user_signature(
                        sender, subject_data, author_signature, author_id)
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
                ts_date = post["date"]
                str_date = ts_date_to_str(ts_date, "%d.%m.%Y %H:%M:%S")
                return str_date

            data_for_message = {}

            owner_signature = select_owner_signature(sender, subject_data, post)
            author_signature = select_author_signature(
                sender, subject_data, post)
            post_url = select_post_url(post)
            publication_date = select_date(post)

            text = "New " + post["post_type"].encode("utf8") + "\n"
            text += "Location: " + owner_signature.encode("utf8") + "\n"
            text += "Author: " + author_signature.encode("utf8") + "\n"
            text += "Created: " + str(publication_date).encode("utf8") + "\n\n"
            text += post["text"].encode("utf8") + "\n\n"
            text += post_url

            # БАГ: максимальная длина сообщения - 4096 знаков

            send_to = values["send_to"]
            access_token = values["access_token"]

            data_for_message = {
                "send_to": send_to,
                "text": text
            }

            if "attachments" in post:
                media_items = select_attachments(post)
                if len(media_items) > 0:
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
            str_date = ts_date_to_str(post["date"], "%d.%m.%Y %H:%M:%S")
            message = "New " + post["post_type"] + " at " + str_date + "."
            output_data.output_text_row(sender, message)

        send_post(sender, values, subject_data, post)
        update_last_date(values, post)
        show_message_about_new_post(sender, post)

    sender += " -> Wall posts monitor"

    wall_posts_data = request_handler.request_wall_posts(
        sender, subject_data, monitor_data)
    wall_posts = sort_items(wall_posts_data)

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
    

def album_photos_monitor(sender, res_filename, subject_data, monitor_data):
    u"""Проверяет фотографии в альбомах."""
    def found_new_photo(sender, values, subject_data, photo):
        u"""Алгоритмы обработки целевой фотографии."""
        def send_photo(sender, values, subject_data, photo):
            u"""Отправка данных из целевой фотографии."""
            def select_album_name(sender, subject_data, photo):
                u"""Выбирает из словаря данные и подписывает название альбома."""
                data_for_request = {
                    "owner_id": photo["owner_id"],
                    "album_ids": photo["album_id"]
                }
                albums_info = request_handler.request_photo_album_info(
                    sender, subject_data, data_for_request)
                album_name = albums_info[0]["title"]
                return album_name

            def select_owner_signature(sender, subject_data, photo):
                u"""Выбирает из словаря данные и формирует гиперссылку на владельца фото."""
                owner_signature = ""
                owner_id = photo["owner_id"]
                if str(owner_id)[0] == "-":
                    owner_signature = make_group_signature(
                        sender, subject_data, owner_signature, owner_id)
                else:
                    owner_signature = make_user_signature(
                        sender, subject_data, owner_signature, owner_id)
                return owner_signature

            def select_author_signature(sender, subject_data, photo):
                u"""Выбирает из словаря данные и формирует гиперссылку на автора."""
                author_signature = ""
                owner_id = photo["owner_id"]
                if str(owner_id)[0] == "-":
                    if "user_id" in photo:
                        if photo["user_id"] == 100:
                            author_id = photo["owner_id"]
                            author_signature += make_group_signature(
                                sender, subject_data, author_signature, author_id)
                        else:
                            author_id = photo["user_id"]
                            author_signature = make_user_signature(
                                sender, subject_data, author_signature, author_id)
                    else:
                        author_signature += "[no data]"
                else:
                    author_id = photo["owner_id"]
                    author_signature += make_user_signature(
                        sender, subject_data, author_signature, author_id)
                return author_signature
            
            def select_attachments(photo):
                u"""Выбирает из словаря данные и формирует список прикреплений."""
                media_items = "photo"
                media_items += str(photo["owner_id"]) + "_" + str(photo["id"])
                return media_items
            
            def select_photo_url(photo):
                u"""Выбирает из словаря данные и формирует URL на фото."""
                photo_url = "https://vk.com/photo"
                photo_url += str(photo["owner_id"]) + "_" + str(photo["id"])
                return photo_url

            def select_date(photo):
                u"""Выбирает из словаря данные и определяет дату фото."""
                ts_date = photo["date"]
                str_date = ts_date_to_str(ts_date, "%d.%m.%Y %H:%M:%S")
                return str_date
            
            data_for_message = {}

            album_name = select_album_name(sender, subject_data, photo)
            owner_signature = select_owner_signature(sender, subject_data, photo)
            author_signature = select_author_signature(
                sender, subject_data, photo)
            post_url = select_photo_url(photo)
            publication_date = select_date(photo)
            media_items = select_attachments(photo)

            text = "New photo" + "\n"
            text += "Album: " + album_name.encode("utf8") + "\n"
            text += "Location: " + owner_signature.encode("utf8") + "\n"
            text += "Author: " + author_signature.encode("utf8") + "\n"
            text += "Created: " + str(publication_date).encode("utf8") + "\n\n"
            if len(photo["text"]) > 0:
                text += photo["text"].encode("utf8") + "\n\n"
            text += post_url.encode("utf8")

            send_to = values["send_to"]
            access_token = values["access_token"]

            data_for_message = {
                "send_to": send_to,
                "text": text,
                "attachment": media_items
            }

            request_handler.send_message(sender, data_for_message, access_token)

            return album_name
        
        def update_last_date(values, photo):
            u"""Обновление даты последнего проверенного фото."""
            res_filename = values["res_filename"]
            path_to_res_file = values["path_to_res_file"]
            dict_json = data_manager.read_json(path_to_res_file, res_filename)
            dict_json["last_date"] = str(photo["date"])
            data_manager.write_json(path_to_res_file, res_filename, dict_json)

        def show_message_about_new_photo(sender, photo, album_name):
            u"""Алгоритмы отображения сообщения о новом фото."""
            str_date = ts_date_to_str(photo["date"], "%d.%m.%Y %H:%M:%S")
            message = "New photo in \"" + album_name + "\" at " + str_date + "."
            output_data.output_text_row(sender, message)

        album_name = send_photo(sender, values, subject_data, photo)
        update_last_date(values, photo)
        show_message_about_new_photo(sender, photo, album_name)

    sender += " -> Album photos monitor"

    album_photos_data = request_handler.request_album_photos(
        sender, subject_data, monitor_data)
    album_photos = sort_items(album_photos_data)

    last_date = int(monitor_data["last_date"])

    new_photos = []

    for item in reversed(album_photos):
        if item["date"] > last_date:
            new_photos.append(item)

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

    for item in new_photos:
        found_new_photo(sender, values, subject_data, item)


def videos_monitor(sender, res_filename, subject_data, monitor_data):
    u"""Проверяет видео."""
    def found_new_video(sender, values, subject_data, video):
        u"""Алгоритмы обработки целевого видео."""
        def send_video(sender, values, subject_data, video):
            u"""Отправка данных из целевого видео."""
            def select_owner_signature(sender, subject_data):
                u"""Выбирает из словаря данные и формирует гиперссылку на владельца видео."""
                owner_signature = ""
                owner_id = subject_data["owner_id"]
                if str(owner_id)[0] == "-":
                    owner_signature = make_group_signature(
                        sender, subject_data, owner_signature, owner_id)
                else:
                    owner_signature = make_user_signature(
                        sender, subject_data, owner_signature, owner_id)
                return owner_signature

            def select_author_signature(sender, subject_data, video):
                u"""Выбирает из словаря данные и формирует гиперссылку на автора."""
                author_signature = ""
                owner_id = video["owner_id"]
                if str(owner_id)[0] == "-":
                    if "user_id" in video:
                        author_id = video["user_id"]
                        author_signature = make_user_signature(
                            sender, subject_data, author_signature, author_id)
                    else:
                        author_id = video["owner_id"]
                        author_signature = make_group_signature(
                            sender, subject_data, author_signature, author_id)
                else:
                    author_id = video["owner_id"]
                    author_signature += make_user_signature(
                        sender, subject_data, author_signature, author_id)
                return author_signature

            def select_attachments(video):
                u"""Выбирает из словаря данные и формирует список прикреплений."""
                media_items = "video"
                media_items += str(video["owner_id"]) + "_" + str(video["id"])
                return media_items

            def select_video_url(video):
                u"""Выбирает из словаря данные и формирует URL на видео."""
                video_url = "https://vk.com/video"
                video_url += str(video["owner_id"]) + "_" + str(video["id"])
                return video_url

            def select_date(video):
                u"""Выбирает из словаря данные и определяет дату видео."""
                ts_date = video["date"]
                str_date = ts_date_to_str(ts_date, "%d.%m.%Y %H:%M:%S")
                return str_date

            data_for_message = {}

            owner_signature = select_owner_signature(sender, subject_data)
            author_signature = select_author_signature(
                sender, subject_data, video)
            post_url = select_video_url(video)
            publication_date = select_date(video)
            media_items = select_attachments(video)

            text = "New video" + "\n"
            text += "Location: " + owner_signature.encode("utf8") + "\n"
            text += "Author: " + author_signature.encode("utf8") + "\n"
            text += "Created: " + str(publication_date).encode("utf8") + "\n\n"
            if len(video["description"]) > 0:
                text += video["description"].encode("utf8") + "\n\n"
            text += post_url.encode("utf8")

            send_to = values["send_to"]
            access_token = values["access_token"]

            data_for_message = {
                "send_to": send_to,
                "text": text,
                "attachment": media_items
            }

            request_handler.send_message(
                sender, data_for_message, access_token)
            
        def update_last_date(values, video):
            u"""Обновление даты последнего проверенного видео."""
            res_filename = values["res_filename"]
            path_to_res_file = values["path_to_res_file"]
            dict_json = data_manager.read_json(path_to_res_file, res_filename)
            dict_json["last_date"] = str(video["date"])
            data_manager.write_json(path_to_res_file, res_filename, dict_json)

        def show_message_about_new_video(sender, video):
            u"""Алгоритмы отображения сообщения о новом видео."""
            str_date = ts_date_to_str(video["date"], "%d.%m.%Y %H:%M:%S")
            message = "New video at " + str_date + "."
            output_data.output_text_row(sender, message)

        send_video(sender, values, subject_data, video)
        update_last_date(values, video)
        show_message_about_new_video(sender, video)

    sender += " -> Videos monitor"

    videos_data = request_handler.request_videos(
        sender, subject_data, monitor_data)
    videos = sort_items(videos_data)

    last_date = int(monitor_data["last_date"])

    new_videos = []

    for item in reversed(videos):
        if item["date"] > last_date:
            new_videos.append(item)

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

    for item in new_videos:
        found_new_video(sender, values, subject_data, item)


# def photo_comments_monitor():
#     u"""Проверяет комментарии под фотографиями."""
# def video_comments_monitor():
#     u"""Проверяет комментарии под видео."""
# def topic_comments_monitor():
#     u"""Проверяет комментарии в обсуждениях."""
# def wall_post_comments_monitor():
#     u"""Проверяет комментарии под постами на стене."""


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
