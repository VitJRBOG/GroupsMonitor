# coding: utf8
u"""Модуль алгоритмов проверки."""


import datetime
import re
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
                        if post["post_type"] == "suggest":
                            author_id = post["from_id"]
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
                            media_items += "_" + attachment["access_key"]
                        if len(attachments) > 1 and i < len(attachments) - 1:
                            media_items += ","
                return media_items
            def select_copy_history(post):
                u"""Выбирает из словаря данные о репосте."""
                copy_history = post["copy_history"]
                media_items = "wall" + \
                    str(copy_history["owner_id"]) + "_" + str(copy_history["id"])
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
            if len(post["text"].encode("utf8")) > 1000:
                text += post["text"].encode("utf8")[0:1000] + "\n"
                text += "<..>\n[long text]\n\n"
            else:
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
                if len(media_items) > 0:
                    data_for_message.update({"attachment": media_items})

            if "copy_history" in post:
                copy_history = select_copy_history(post)
                if len(copy_history) > 0:
                    if "attachment" in data_for_message:
                        data_for_message["attachment"] += "," + copy_history
                    else:
                        data_for_message.update({"attachment": copy_history})

            request_handler.send_message(sender, data_for_message, access_token)
        
        def show_message_about_new_post(sender, post):
            u"""Алгоритмы отображения сообщения о новом посте."""
            str_date = ts_date_to_str(post["date"], "%d.%m.%Y %H:%M:%S")
            message = "New " + post["post_type"] + " at " + str_date + "."
            output_data.output_text_row(sender, message)

        send_post(sender, values, subject_data, post)
        show_message_about_new_post(sender, post)

    def update_last_date(values, post):
        u"""Обновление даты последнего проверенного поста."""
        res_filename = values["res_filename"]
        path_to_res_file = values["path_to_res_file"]
        dict_json = data_manager.read_json(path_to_res_file, res_filename)
        dict_json["last_date"] = str(item["date"])
        data_manager.write_json(path_to_res_file, res_filename, dict_json)

    def target_check(sender, values, subject_data, item):
        u"""Проверка целевого поста."""
        def check_by_attachments(target, values, item):
            u"""Проверка по наличию медиаконтента."""
            if "attachments" in item:
                res_filename = values["res_filename"]
                path_to_res_file = values["path_to_res_file"]
                monitor_data = data_manager.read_json(
                    path_to_res_file, res_filename)
                attachments_types = monitor_data["check_by_attachments"]["types"]
                if len(attachments_types) > 0:
                    for attachment in item["attachments"]:
                        for attachment_type in attachments_types:
                            if attachment_type == attachment["type"]:
                                target = True
                                return target
            return target

        def check_by_keywords(target, item):
            u"""Проверка по наличию ключевых фраз."""
            def check_keywords_algorithm(target, text, keywords):
                u"""Алгоритм проверки слов."""
                if len(keywords) > 0:
                    for keyword in keywords:
                        if text.lower().find(keyword) > -1:
                            target = True
                            return target
                return target

            text = item["text"]

            res_filename = values["res_filename"]
            path_to_res_file = values["path_to_res_file"]
            monitor_data = data_manager.read_json(path_to_res_file, res_filename)

            keywords = monitor_data["check_by_keywords"]["keywords"]
            target = check_keywords_algorithm(target, text, keywords)
            if target:
                return target

            return target

        def check_by_exclude(item):
            u"""Проверка автора по списку игнорируемых."""
            ignore = False
            res_filename = values["res_filename"]
            path_to_res_file = values["path_to_res_file"]
            monitor_data = data_manager.read_json(path_to_res_file, res_filename)
            ignore_users = monitor_data["check_by_exclude"]["authors"]
            if len(ignore_users) > 0:
                author_id = ""
                if "from_id" in item:
                    if str(item["owner_id"])[0] == "-" and \
                    str(item["from_id"])[0] != "-":
                        author_id = item["from_id"]
                elif "signer_id" in item:
                    author_id = item["signer_id"]
                else:
                    ignore = False
                    return ignore

                for ignore_user in ignore_users:
                    if ignore_user == author_id:
                        ignore = True
                        return ignore
            return ignore

        target = False

        if monitor_data["check_all"]["check"] == 1:
            if monitor_data["check_all"]["check_exclude"] == 0:
                is_exclude = check_by_exclude(item)
                if is_exclude:
                    target = False
                else:
                    target = True
            else:
                target = True
            return target

        if monitor_data["check_by_exclude"]["check"] == 1:
            is_exclude = check_by_exclude(item)
            if is_exclude:
                target = False
                return target
        
        if monitor_data["check_by_attachments"]["check"] == 1:
            target = check_by_attachments(target, values, item)
            if target:
                return target

        if len(item["text"]) > 0:
            if monitor_data["check_by_keywords"]["check"] == 1:
                target = check_by_keywords(target, item)
                if target:
                    return target
        return target

    sender += " -> Wall posts monitor"

    wall_posts_data = request_handler.request_wall_posts(
        sender, subject_data, monitor_data)

    if len(wall_posts_data) > 0:
        last_date = int(monitor_data["last_date"])

        new_posts = []

        for item in reversed(wall_posts_data):
            if item["date"] > last_date:
                new_posts.append(item)

        if len(new_posts) > 0:
            wall_posts = sort_items(new_posts)

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
            
            for item in reversed(wall_posts):
                target = target_check(sender, values, subject_data, item)
                if target:
                    found_new_post(sender, values, subject_data, item)
                update_last_date(values, item)



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
                if len(photo["text"].encode("utf8")) > 1000:
                    text += photo["text"].encode("utf8")[0:1000] + "\n"
                    text += "<..>\n[long text]\n\n"
                else:
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

    if len(album_photos_data) > 0:
        last_date = int(monitor_data["last_date"])

        new_photos = []

        for item in reversed(album_photos_data):
            if item["date"] > last_date:
                new_photos.append(item)

        if len(new_photos) > 0:
            album_photos = sort_items(new_photos)

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

            for item in reversed(album_photos):
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
                if len(video["description"].encode("utf8")) > 1000:
                    text += video["description"].encode("utf8")[0:1000] + "\n"
                    text += "<..>\n[long text]\n\n"
                else:
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

    if len(videos_data) > 0:
        last_date = int(monitor_data["last_date"])

        new_videos = []

        for item in reversed(videos_data):
            if item["date"] > last_date:
                new_videos.append(item)

        if len(new_videos) > 0:
            videos = sort_items(new_videos)

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

            for item in reversed(videos):
                found_new_video(sender, values, subject_data, item)


def photo_comments_monitor(sender, res_filename, subject_data, monitor_data):
    u"""Проверяет комментарии под фотографиями."""
    def found_new_photo_comment(sender, values, subject_data, photo_comment):
        u"""Алгоритмы обработки целевого комментария."""
        def send_photo_comment(sender, values, subject_data, photo_comment):
            u"""Отправка данных из целевого комментария."""
            def select_owner_signature(sender, subject_data):
                u"""Выбирает из словаря данные и формирует гиперссылку на владельца комментария."""
                owner_signature = ""
                owner_id = subject_data["owner_id"]
                if str(owner_id)[0] == "-":
                    owner_id = subject_data["owner_id"]
                    owner_signature = make_group_signature(
                        sender, subject_data, owner_signature, owner_id)
                else:
                    owner_id = subject_data["owner_id"]
                    owner_signature = make_user_signature(
                        sender, subject_data, owner_signature, owner_id)
                return owner_signature

            def select_author_signature(sender, subject_data, photo_comment):
                u"""Выбирает из словаря данные и формирует гиперссылку на автора."""
                author_signature = ""
                author_id = photo_comment["from_id"]
                if str(author_id)[0] == "-":
                    author_signature = make_group_signature(
                        sender, subject_data, author_signature, author_id)
                else:
                    author_id = photo_comment["from_id"]
                    author_signature += make_user_signature(
                        sender, subject_data, author_signature, author_id)
                return author_signature

            def select_attachments(photo_comment):
                u"""Выбирает из словаря данные и формирует список прикреплений."""
                attachments = photo_comment["attachments"]
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
                            media_items += "_" + \
                                attachment["access_key"]
                        if len(attachments) > 1 and i < len(attachments) - 1:
                            media_items += ","
                return media_items

            def select_photo_url(photo_comment, subject_data):
                u"""Выбирает из словаря данные и формирует URL на фото."""
                post_url = "https://vk.com/photo"
                post_url += str(subject_data["owner_id"]) + \
                    "_" + str(photo_comment["pid"])
                return post_url

            def select_date(photo_comment):
                u"""Выбирает из словаря данные и определяет дату комментария."""
                ts_date = photo_comment["date"]
                str_date = ts_date_to_str(ts_date, "%d.%m.%Y %H:%M:%S")
                return str_date

            data_for_message = {}

            owner_signature = select_owner_signature(
                sender, subject_data)
            author_signature = select_author_signature(
                sender, subject_data, photo_comment)
            photo_url = select_photo_url(photo_comment, subject_data)
            publication_date = select_date(photo_comment)

            text = "New comment under photo\n"
            text += "Location: " + owner_signature.encode("utf8") + "\n"
            text += "Author: " + author_signature.encode("utf8") + "\n"
            text += "Created: " + str(publication_date).encode("utf8") + "\n\n"
            if len(photo_comment["text"].encode("utf8")) > 1000:
                text += photo_comment["text"].encode("utf8")[0:1000] + "\n"
                text += "<..>\n[long text]\n\n"
            else:
                text += photo_comment["text"].encode("utf8") + "\n\n"
            text += photo_url

            send_to = values["send_to"]
            access_token = values["access_token"]

            data_for_message = {
                "send_to": send_to,
                "text": text
            }

            if "attachments" in photo_comment:
                media_items = select_attachments(photo_comment)
                if len(media_items) > 0:
                    data_for_message.update({"attachment": media_items})

            request_handler.send_message(sender, data_for_message, access_token)

        def update_last_date(values, photo_comment):
            u"""Обновление даты последнего проверенного комментария."""
            res_filename = values["res_filename"]
            path_to_res_file = values["path_to_res_file"]
            dict_json = data_manager.read_json(path_to_res_file, res_filename)
            dict_json["last_date"] = str(item["date"])
            data_manager.write_json(path_to_res_file, res_filename, dict_json)

        def show_message_about_new_photo_comment(sender, photo_comment):
            u"""Алгоритмы отображения сообщения о новом комментарии."""
            str_date = ts_date_to_str(
                photo_comment["date"], "%d.%m.%Y %H:%M:%S")
            message = "New comment under photo at " + str_date + "."
            output_data.output_text_row(sender, message)

        send_photo_comment(sender, values, subject_data, photo_comment)
        update_last_date(values, photo_comment)
        show_message_about_new_photo_comment(sender, photo_comment)

    sender += " -> Photo comments monitor"

    photo_comments_data = request_handler.request_photo_comments(
        sender, subject_data, monitor_data)
    
    if len(photo_comments_data) > 0:
        last_date = int(monitor_data["last_date"])

        new_photo_comments = []

        for item in reversed(photo_comments_data):
            if item["date"] > last_date:
                new_photo_comments.append(item)

        if len(new_photo_comments) > 0:
            photo_comments = sort_items(new_photo_comments)

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

            for item in reversed(photo_comments):
                found_new_photo_comment(sender, values, subject_data, item)


def video_comments_monitor(sender, res_filename, subject_data, monitor_data):
    u"""Проверяет комментарии под видео."""
    def found_new_video_comment(sender, values, subject_data, video_comment):
        u"""Алгоритмы обработки целевого комментария."""
        def send_video_comment(sender, values, subject_data, video_comment):
            u"""Отправка данных из целевого комментария."""
            def select_owner_signature(sender, subject_data):
                u"""Выбирает из словаря данные и формирует гиперссылку на владельца комментария."""
                owner_signature = ""
                owner_id = subject_data["owner_id"]
                if str(owner_id)[0] == "-":
                    owner_id = subject_data["owner_id"]
                    owner_signature = make_group_signature(
                        sender, subject_data, owner_signature, owner_id)
                else:
                    owner_id = subject_data["owner_id"]
                    owner_signature = make_user_signature(
                        sender, subject_data, owner_signature, owner_id)
                return owner_signature

            def select_author_signature(sender, subject_data, video_comment):
                u"""Выбирает из словаря данные и формирует гиперссылку на автора."""
                author_signature = ""
                author_id = video_comment["from_id"]
                if str(author_id)[0] == "-":
                    author_signature = make_group_signature(
                        sender, subject_data, author_signature, author_id)
                else:
                    author_id = video_comment["from_id"]
                    author_signature += make_user_signature(
                        sender, subject_data, author_signature, author_id)
                return author_signature

            def select_attachments(video_comment):
                u"""Выбирает из словаря данные и формирует список прикреплений."""
                attachments = video_comment["attachments"]
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
                            media_items += "_" + \
                                attachment["access_key"]
                        if len(attachments) > 1 and i < len(attachments) - 1:
                            media_items += ","
                return media_items

            def select_video_url(video_comment, subject_data):
                u"""Выбирает из словаря данные и формирует URL на видео."""
                post_url = "https://vk.com/video"
                post_url += str(video_comment["owner_id"]) + \
                    "_" + str(video_comment["vid"])
                return post_url

            def select_date(video_comment):
                u"""Выбирает из словаря данные и определяет дату комментария."""
                ts_date = video_comment["date"]
                str_date = ts_date_to_str(ts_date, "%d.%m.%Y %H:%M:%S")
                return str_date

            data_for_message = {}

            owner_signature = select_owner_signature(
                sender, subject_data)
            author_signature = select_author_signature(
                sender, subject_data, video_comment)
            video_url = select_video_url(video_comment, subject_data)
            publication_date = select_date(video_comment)

            text = "New comment under video\n"
            text += "Location: " + owner_signature.encode("utf8") + "\n"
            text += "Author: " + author_signature.encode("utf8") + "\n"
            text += "Created: " + str(publication_date).encode("utf8") + "\n\n"
            if len(video_comment["text"].encode("utf8")) > 1000:
                text += video_comment["text"].encode("utf8")[0:1000] + "\n"
                text += "<..>\n[long text]\n\n"
            else:
                text += video_comment["text"].encode("utf8") + "\n\n"
            text += video_url

            send_to = values["send_to"]
            access_token = values["access_token"]

            data_for_message = {
                "send_to": send_to,
                "text": text
            }

            if "attachments" in video_comment:
                media_items = select_attachments(video_comment)
                if len(media_items) > 0:
                    data_for_message.update({"attachment": media_items})

            request_handler.send_message(
                sender, data_for_message, access_token)

        def update_last_date(values, video_comment):
            u"""Обновление даты последнего проверенного комментария."""
            res_filename = values["res_filename"]
            path_to_res_file = values["path_to_res_file"]
            dict_json = data_manager.read_json(path_to_res_file, res_filename)
            dict_json["last_date"] = str(item["date"])
            data_manager.write_json(path_to_res_file, res_filename, dict_json)

        def show_message_about_new_video_comment(sender, video_comment):
            u"""Алгоритмы отображения сообщения о новом комментарии."""
            str_date = ts_date_to_str(
                video_comment["date"], "%d.%m.%Y %H:%M:%S")
            message = "New comment under video at " + str_date + "."
            output_data.output_text_row(sender, message)

        send_video_comment(sender, values, subject_data, video_comment)
        update_last_date(values, video_comment)
        show_message_about_new_video_comment(sender, video_comment)

    sender += " -> Video comments monitor"

    videos_data = request_handler.request_videos(
        sender, subject_data, monitor_data)

    if len(videos_data) > 0:
        videos_comments_data = []

        for video in videos_data:
            video_comments_data = request_handler.request_video_comments(
                sender, subject_data, monitor_data, video)
            if len(video_comments_data) > 0:
                for i, comment in enumerate(video_comments_data):
                    video_comments_data[i].update({"vid": video["id"]})
                    video_comments_data[i].update({"owner_id": video["owner_id"]})
                videos_comments_data.extend(video_comments_data)

        last_date = int(monitor_data["last_date"])

        new_videos_comments = []

        for item in reversed(videos_comments_data):
            if item["date"] > last_date:
                new_videos_comments.append(item)

        if len(new_videos_comments) > 0:
            videos_comments = sort_items(new_videos_comments)

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

            for item in reversed(videos_comments):
                found_new_video_comment(sender, values, subject_data, item)


def topic_comments_monitor(sender, res_filename, subject_data, monitor_data):
    u"""Проверяет комментарии в обсуждениях."""
    def found_new_topic_comment(sender, values, subject_data, topic_comment):
        u"""Алгоритмы обработки целевого комментария."""
        def send_topic_comment(sender, values, subject_data, topic_comment):
            u"""Отправка данных из целевого комментария."""
            def select_topic_name(sender, topic_comment):
                u"""Выбирает из словаря данные и подписывает название альбома."""
                topic_name = topic_comment["topic_name"]
                return topic_name

            def select_owner_signature(sender, subject_data):
                u"""Выбирает из словаря данные и формирует гиперссылку на владельца комментария."""
                owner_signature = ""
                owner_id = subject_data["owner_id"]
                owner_signature = make_group_signature(
                    sender, subject_data, owner_signature, owner_id)
                return owner_signature

            def select_author_signature(sender, subject_data, topic_comment):
                u"""Выбирает из словаря данные и формирует гиперссылку на автора."""
                author_signature = ""
                author_id = topic_comment["from_id"]
                if str(author_id)[0] == "-":
                    author_signature = make_group_signature(
                        sender, subject_data, author_signature, author_id)
                else:
                    author_id = topic_comment["from_id"]
                    author_signature += make_user_signature(
                        sender, subject_data, author_signature, author_id)
                return author_signature

            def select_attachments(topic_comment):
                u"""Выбирает из словаря данные и формирует список прикреплений."""
                attachments = topic_comment["attachments"]
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
                            media_items += "_" + \
                                attachment["access_key"]
                        if len(attachments) > 1 and i < len(attachments) - 1:
                            media_items += ","
                return media_items

            def select_topic_url(topic_comment, subject_data):
                u"""Выбирает из словаря данные и формирует URL на комментарий."""
                post_url = "https://vk.com/topic"
                post_url += str(topic_comment["owner_id"]) + \
                    "_" + str(topic_comment["tid"]) + "?post=" + \
                    str(topic_comment["id"])
                return post_url

            def select_date(topic_comment):
                u"""Выбирает из словаря данные и определяет дату комментария."""
                ts_date = topic_comment["date"]
                str_date = ts_date_to_str(ts_date, "%d.%m.%Y %H:%M:%S")
                return str_date

            data_for_message = {}

            topic_name = select_topic_name(sender, topic_comment)
            owner_signature = select_owner_signature(
                sender, subject_data)
            author_signature = select_author_signature(
                sender, subject_data, topic_comment)
            topic_url = select_topic_url(topic_comment, subject_data)
            publication_date = select_date(topic_comment)

            text = "New comment under topic\n"
            text += "Topic: " + topic_name.encode("utf8") + "\n"
            text += "Location: " + owner_signature.encode("utf8") + "\n"
            text += "Author: " + author_signature.encode("utf8") + "\n"
            text += "Created: " + str(publication_date).encode("utf8") + "\n\n"
            # BUG: длинный запрос. Поставлен костыль на ограничение количества символов.
            # SOLUTION: проверка длины всего URI и обрезание текстовой части сообщения.
            if len(topic_comment["text"].encode("utf8")) > 800:
                text += topic_comment["text"].encode("utf8")[0:800] + "\n"
                text += "<..>\n[long text]\n\n"
            else:
                text += topic_comment["text"].encode("utf8") + "\n\n"
            text += topic_url

            send_to = values["send_to"]
            access_token = values["access_token"]

            data_for_message = {
                "send_to": send_to,
                "text": text
            }

            if "attachments" in topic_comment:
                media_items = select_attachments(topic_comment)
                if len(media_items) > 0:
                    data_for_message.update({"attachment": media_items})

            request_handler.send_message(
                sender, data_for_message, access_token)

        def update_last_date(values):
            u"""Обновление даты последнего проверенного комментария."""
            res_filename = values["res_filename"]
            path_to_res_file = values["path_to_res_file"]
            dict_json = data_manager.read_json(path_to_res_file, res_filename)
            dict_json["last_date"] = str(item["date"])
            data_manager.write_json(path_to_res_file, res_filename, dict_json)

        def show_message_about_new_topic_comment(sender, topic_comment):
            u"""Алгоритмы отображения сообщения о новом комментарии."""
            str_date = ts_date_to_str(
                topic_comment["date"], "%d.%m.%Y %H:%M:%S")
            message = "New comment under topic at " + str_date + "."
            output_data.output_text_row(sender, message)

        send_topic_comment(sender, values, subject_data, topic_comment)
        update_last_date(values)
        show_message_about_new_topic_comment(sender, topic_comment)

    sender += " -> Topic comments monitor"

    topics_data = request_handler.request_topics_info(
        sender, subject_data, monitor_data)

    if len(topics_data) > 0:
        topics_comments_data = []

        for topic in topics_data:
            topic_comments_data = request_handler.request_topic_comments(
                sender, subject_data, monitor_data, topic)
            if len(topic_comments_data) > 0:
                for i, comment in enumerate(topic_comments_data):
                    topic_comments_data[i].update({"tid": topic["id"]})
                    topic_comments_data[i].update(
                        {"owner_id": int("-" + str(topic["owner_id"]))})
                    topic_comments_data[i].update({"topic_name": topic["title"]})
                topics_comments_data.extend(topic_comments_data)

        last_date = int(monitor_data["last_date"])

        new_topic_comments = []

        for item in reversed(topics_comments_data):
            if item["date"] > last_date:
                new_topic_comments.append(item)

        if len(new_topic_comments) > 0:
            topics_comments = sort_items(new_topic_comments)

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

            for item in reversed(topics_comments):
                found_new_topic_comment(sender, values, subject_data, item)


def wall_post_comments_monitor(sender, res_filename, subject_data, monitor_data):
    u"""Проверяет комментарии под постами на стене."""
    def found_new_wall_post_comment(sender, values, subject_data, wall_post_comment):
        u"""Алгоритмы обработки целевого комментария."""
        def send_wall_post_comment(sender, values, subject_data, wall_post_comment):
            u"""Отправка данных из целевого комментария."""
            def select_owner_signature(sender, subject_data):
                u"""Выбирает из словаря данные и формирует гиперссылку на владельца комментария."""
                owner_signature = ""
                owner_id = subject_data["owner_id"]
                if str(owner_id)[0] == "-":
                    owner_id = subject_data["owner_id"]
                    owner_signature = make_group_signature(
                        sender, subject_data, owner_signature, owner_id)
                else:
                    owner_id = subject_data["owner_id"]
                    owner_signature = make_user_signature(
                        sender, subject_data, owner_signature, owner_id)
                return owner_signature

            def select_author_signature(sender, subject_data, wall_post_comment):
                u"""Выбирает из словаря данные и формирует гиперссылку на автора."""
                author_signature = ""
                author_id = wall_post_comment["from_id"]
                if str(author_id)[0] == "-":
                    author_signature = make_group_signature(
                        sender, subject_data, author_signature, author_id)
                else:
                    author_id = wall_post_comment["from_id"]
                    author_signature += make_user_signature(
                        sender, subject_data, author_signature, author_id)
                return author_signature

            def select_attachments(wall_post_comment):
                u"""Выбирает из словаря данные и формирует список прикреплений."""
                attachments = wall_post_comment["attachments"]
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
                            media_items += "_" + \
                                attachment["access_key"]
                        if len(attachments) > 1 and i < len(attachments) - 1:
                            media_items += ","
                return media_items

            def select_wall_post_comment_url(wall_post_comment, subject_data):
                u"""Выбирает из словаря данные и формирует URL на комментарий."""
                post_url = "https://vk.com/wall"
                post_url += str(wall_post_comment["owner_id"]) + \
                    "_" + str(wall_post_comment["post_id"]) + "?reply=" + \
                    str(wall_post_comment["id"])
                if "parents_stack" in wall_post_comment:
                    post_url += "&thread=" + str(wall_post_comment["parents_stack"][0])
                return post_url

            def select_date(wall_post_comment):
                u"""Выбирает из словаря данные и определяет дату комментария."""
                ts_date = wall_post_comment["date"]
                str_date = ts_date_to_str(ts_date, "%d.%m.%Y %H:%M:%S")
                return str_date

            data_for_message = {}

            owner_signature = select_owner_signature(
                sender, subject_data)
            author_signature = select_author_signature(
                sender, subject_data, wall_post_comment)
            wall_post_comment_url = select_wall_post_comment_url(wall_post_comment, subject_data)
            publication_date = select_date(wall_post_comment)

            text = "New comment under wall post\n"
            text += "Location: " + owner_signature.encode("utf8") + "\n"
            text += "Author: " + author_signature.encode("utf8") + "\n"
            text += "Created: " + str(publication_date).encode("utf8") + "\n\n"
            if len(wall_post_comment["text"].encode("utf8")) > 1000:
                text += wall_post_comment["text"].encode("utf8")[0:1000] + "\n"
                text += "<..>\n[long text]\n\n"
            else:
                text += wall_post_comment["text"].encode("utf8") + "\n\n"
            text += wall_post_comment_url

            send_to = values["send_to"]
            access_token = values["access_token"]

            data_for_message = {
                "send_to": send_to,
                "text": text
            }

            if "attachments" in wall_post_comment:
                media_items = select_attachments(wall_post_comment)
                if len(media_items) > 0:
                    data_for_message.update({"attachment": media_items})

            request_handler.send_message(
                sender, data_for_message, access_token)

        def show_message_about_new_wall_post_comment(sender, wall_post_comment):
            u"""Алгоритмы отображения сообщения о новом комментарии."""
            str_date = ts_date_to_str(
                wall_post_comment["date"], "%d.%m.%Y %H:%M:%S")
            message = "New comment under wall post at " + str_date + "."
            output_data.output_text_row(sender, message)

        send_wall_post_comment(sender, values, subject_data, wall_post_comment)
        show_message_about_new_wall_post_comment(sender, wall_post_comment)

    def update_last_date(values, item):
        u"""Обновление даты последнего проверенного комментария."""
        res_filename = values["res_filename"]
        path_to_res_file = values["path_to_res_file"]
        dict_json = data_manager.read_json(path_to_res_file, res_filename)
        dict_json["last_date"] = str(item["date"])
        data_manager.write_json(path_to_res_file, res_filename, dict_json)

    def suspicion_check(sender, values, subject_data, item):
        u"""Проверка подозрительности комментария."""
        def check_by_exclude(item):
            u"""Проверка по комментариям от авторов исключенных из проверки."""
            is_exclude = False
            res_filename = values["res_filename"]
            path_to_res_file = values["path_to_res_file"]
            monitor_data = data_manager.read_json(path_to_res_file, res_filename)
            if len(monitor_data["check_by_exclude"]["authors"]) > 0:
                exclude_authors = monitor_data["check_by_exclude"]["authors"]
                for exclude in exclude_authors:
                    if item["from_id"] != exclude:
                        is_exclude = False
                    else:
                        is_exclude = True
                        break
            return is_exclude

        def check_by_communities(suspicious, item):
            u"""Проверка по комментариям от имени сообщества."""
            if str(item["from_id"])[0] == "-":
                suspicious = True
            return suspicious

        def check_by_profile_data(suspicious, item, subject_data):
            u"""Проверка по данным пользовательского профиля."""
            res_filename = values["res_filename"]
            path_to_res_file = values["path_to_res_file"]
            monitor_data = data_manager.read_json(
                path_to_res_file, res_filename)

            if len(monitor_data["check_by_profile_data"]["ids"]) > 0:
                ids = monitor_data["check_by_profile_data"]["ids"]
                for from_id in ids:
                    if item["from_id"] == from_id:
                        suspicious = True
                        return suspicious

            if str(item["from_id"])[0] != "-":
                if len(monitor_data["check_by_profile_data"]["names"]) > 0:
                    names = monitor_data["check_by_profile_data"]["names"]
                    for name in names:
                        data_for_request = {
                            "user_ids": item["from_id"]
                        }
                        users_info = request_handler.request_user_info(
                            sender, subject_data, data_for_request)
                        user_name = users_info[0]["first_name"] + " " + users_info[0]["last_name"]
                        if user_name == name:
                            suspicious = True
                            return suspicious
            return suspicious

        def check_by_empty_comment(suspicious, item):
            u"""Проверка по пустому комментарию."""
            if len(item["text"]) < 1:
                if "attachments" not in item:
                    suspicious = True
                    return suspicious
            return suspicious

        def check_by_attachments(suspicious, values, item):
            u"""Проверка по наличию медиаконтента."""
            if "attachments" in item:
                res_filename = values["res_filename"]
                path_to_res_file = values["path_to_res_file"]
                monitor_data = data_manager.read_json(
                    path_to_res_file, res_filename)
                attachments_types = monitor_data["check_by_attachments"]["types"]
                if len(attachments_types) > 0:
                    for attachment in item["attachments"]:
                        for attachment_type in attachments_types:
                            if attachment_type == attachment["type"]:
                                suspicious = True
                                return suspicious
            return suspicious

        def check_by_phone_number(suspicious, item):
            u"""Проверка по наличию номера телефона в тексте."""
            res_filename = values["res_filename"]
            path_to_res_file = values["path_to_res_file"]
            monitor_data = data_manager.read_json(
                path_to_res_file, res_filename)
            need_repeats = monitor_data["check_by_phone_number"]["digits_count"]
            symbs_from_text = item["text"]
            repeats = 0
            interrupts = 0
            for i, sym in enumerate(symbs_from_text):
                if len(re.findall(r"[0123456789]", sym)) > 0:
                    repeats += 1
                    interrupts = 0
                elif len(re.findall(r"[()\- ]", sym)) > 0:
                    if interrupts == 0:
                        interrupts += 1
                    else:
                        interrupts = 0
                        repeats = 0
                else:
                    repeats = 0
                    interrupts = 0
            for need_repeat in need_repeats:
                if repeats == need_repeat:
                    suspicious = True
                    return suspicious

        def check_by_card_number(suspicious, item):
            u"""Проверка по наличию номера карты в тексте."""
            res_filename = values["res_filename"]
            path_to_res_file = values["path_to_res_file"]
            monitor_data = data_manager.read_json(
                path_to_res_file, res_filename)
            need_repeats = monitor_data["check_by_card_number"]["digits_count"]
            symbs_from_text = item["text"]
            repeats = 0
            interrupts = 0
            for sym in symbs_from_text:
                if len(re.findall(r"[0123456789]", sym)) > 0:
                    repeats += 1
                    interrupts = 0
                elif len(re.findall(r"[ ]", sym)) > 0:
                    if interrupts == 0:
                        interrupts += 1
                    else:
                        interrupts = 0
                        repeats = 0
                else:
                    repeats = 0
                    interrupts = 0
            for need_repeat in need_repeats:
                if repeats == need_repeat:
                    suspicious = True
                    return suspicious

        def check_by_keywords(suspicious, item):
            u"""Проверка по наличию ключевых фраз."""
            def charchange(text, chan_symbs_dict):
                u"""Алгоритм замены символов."""
                if len(chan_symbs_dict) > 0:
                    symbs_from_text = list(text)
                    for i, sym in enumerate(symbs_from_text):
                        if sym in chan_symbs_dict:
                            symbs_from_text[i] = chan_symbs_dict[sym]
                    text = ''.join(symbs_from_text)
                return text

            def check_small_messages_algorithm(suspicious, text, small_messages):
                u"""Алгоритм проверки небольших сообщений."""
                if len(small_messages) > 0:
                    for small_message in small_messages:
                        if len(text) == len(small_message):
                            if text.lower().find(small_message) > -1:
                                suspicious = True
                                return suspicious
                return suspicious

            def check_keywords_algorithm(suspicious, text, keywords):
                u"""Алгоритм проверки слов."""
                if len(keywords) > 0:
                    for keyword in keywords:
                        if text.lower().find(keyword) > -1:
                            suspicious = True
                            return suspicious
                return suspicious

            text = item["text"]

            res_filename = values["res_filename"]
            path_to_res_file = values["path_to_res_file"]
            monitor_data = data_manager.read_json(path_to_res_file, res_filename)

            small_messages = monitor_data["check_by_keywords"]["small_messages"]
            suspicious = check_small_messages_algorithm(suspicious, text, small_messages)
            if suspicious:
                return suspicious

            keywords = monitor_data["check_by_keywords"]["keywords"]
            suspicious = check_keywords_algorithm(suspicious, text, keywords)
            if suspicious:
                return suspicious

            if monitor_data["check_by_keywords"]["check_charchange"] == 1:
                chan_symbs_dict = monitor_data["check_by_keywords"]["chars_lat_to_cyr"]
                text = charchange(text, chan_symbs_dict)

                suspicious = check_small_messages_algorithm(
                    suspicious, text, small_messages)
                if suspicious:
                    return suspicious

                suspicious = check_keywords_algorithm(
                    suspicious, text, keywords)
                if suspicious:
                    return suspicious

            return suspicious

        suspicious = False

        if monitor_data["check_all"]["check"] == 1:
            if monitor_data["check_all"]["check_exclude"] == 0:
                is_exclude = check_by_exclude(item)
                if is_exclude:
                    suspicious = False
                else:
                    suspicious = True
            else:
                suspicious = True
            return suspicious

        if monitor_data["check_by_exclude"]["check"] == 1:
            is_exclude = check_by_exclude(item)
            if is_exclude:
                suspicious = False
                return suspicious

        if monitor_data["check_by_communities"]["check"] == 1:
            suspicious = check_by_communities(suspicious, item)
            if suspicious:
                return suspicious

        if monitor_data["check_by_profile_data"]["check"] == 1:
            suspicious = check_by_profile_data(suspicious, item, subject_data)
            if suspicious:
                return suspicious

        if monitor_data["check_by_empty_comment"]["check"] == 1:
            suspicious = check_by_empty_comment(suspicious, item)
            if suspicious:
                return suspicious

        if monitor_data["check_by_attachments"]["check"] == 1:
            suspicious = check_by_attachments(suspicious, values, item)
            if suspicious:
                return suspicious

        if monitor_data["check_by_phone_number"]["check"] == 1:
            suspicious = check_by_phone_number(suspicious, item)
            if suspicious:
                return suspicious

        if monitor_data["check_by_card_number"]["check"] == 1:
            suspicious = check_by_card_number(suspicious, item)
            if suspicious:
                return suspicious

        if len(item["text"]) > 0:
            if monitor_data["check_by_keywords"]["check"] == 1:
                suspicious = check_by_keywords(suspicious, item)
                if suspicious:
                    return suspicious
        return suspicious

    sender += " -> Wall post comments monitor"

    wall_posts_data = request_handler.request_wall_posts(
        sender, subject_data, monitor_data)

    if len(wall_posts_data) > 0:
        wall_posts_comments_data = []

        for wall_post in wall_posts_data:
            wall_post_comments_data = request_handler.request_wall_post_comments(
                sender, subject_data, monitor_data, wall_post)
            if len(wall_post_comments_data) > 0:
                wall_posts_comments_data.extend(wall_post_comments_data)

        last_date = int(monitor_data["last_date"])

        new_wall_post_comments = []

        for item in reversed(wall_posts_comments_data):
            if item["date"] > last_date:
                new_wall_post_comments.append(item)

        if len(new_wall_post_comments) > 0:
            wall_posts_comments = sort_items(new_wall_post_comments)

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

            for item in reversed(wall_posts_comments):
                suspicious = suspicion_check(sender, values, subject_data, item)
                if suspicious:
                    found_new_wall_post_comment(sender, values, subject_data, item)
                update_last_date(values, item)


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
