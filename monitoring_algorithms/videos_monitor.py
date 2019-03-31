# coding: utf8
u"""Модуль проверки новых видео."""


import tools
import request_handler
import output_data
import data_manager


def run_monitoring_videos(sender, res_filename, subject_data, monitor_data):
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
                    owner_signature = tools.make_group_signature(
                        sender, subject_data, owner_signature, owner_id)
                else:
                    owner_signature = tools.make_user_signature(
                        sender, subject_data, owner_signature, owner_id)
                return owner_signature

            def select_author_signature(sender, subject_data, video):
                u"""Выбирает из словаря данные и формирует гиперссылку на автора."""
                author_signature = ""
                owner_id = video["owner_id"]
                if str(owner_id)[0] == "-":
                    if "user_id" in video:
                        author_id = video["user_id"]
                        author_signature = tools.make_user_signature(
                            sender, subject_data, author_signature, author_id)
                    else:
                        author_id = video["owner_id"]
                        author_signature = tools.make_group_signature(
                            sender, subject_data, author_signature, author_id)
                else:
                    author_id = video["owner_id"]
                    author_signature += tools.make_user_signature(
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
                str_date = tools.ts_date_to_str(ts_date, "%d.%m.%Y %H:%M:%S")
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
                if len(video["description"].encode("utf8")) > 3500:
                    text += video["description"].encode("utf8")[0:3500] + "\n"
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
            str_date = tools.ts_date_to_str(video["date"], "%d.%m.%Y %H:%M:%S")
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
            videos = tools.sort_items(new_videos)

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
