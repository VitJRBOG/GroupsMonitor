# coding: utf8
u"""Модуль проверки новых фотографий в альбомах."""


import tools
import request_handler
import data_manager
import output_data


def run_monitoring_album_photos(sender, res_filename, subject_data, monitor_data):
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
                    owner_signature = tools.make_group_signature(
                        sender, subject_data, owner_signature, owner_id)
                else:
                    owner_signature = tools.make_user_signature(
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
                            author_signature += tools.make_group_signature(
                                sender, subject_data, author_signature, author_id)
                        else:
                            author_id = photo["user_id"]
                            author_signature = tools.make_user_signature(
                                sender, subject_data, author_signature, author_id)
                    else:
                        author_signature += "[no data]"
                else:
                    author_id = photo["owner_id"]
                    author_signature += tools.make_user_signature(
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
                str_date = tools.ts_date_to_str(ts_date, "%d.%m.%Y %H:%M:%S")
                return str_date

            data_for_message = {}

            album_name = select_album_name(sender, subject_data, photo)
            owner_signature = select_owner_signature(
                sender, subject_data, photo)
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

            request_handler.send_message(
                sender, data_for_message, access_token)

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
            str_date = tools.ts_date_to_str(photo["date"], "%d.%m.%Y %H:%M:%S")
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
            album_photos = tools.sort_items(new_photos)

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
