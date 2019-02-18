# coding: utf8
u"""Модуль проверки комментариев под видео."""


import tools
import data_manager
import request_handler
import output_data


def run_monitoring_video_comments(sender, res_filename, subject_data, monitor_data):
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
                    owner_signature = tools.make_group_signature(
                        sender, subject_data, owner_signature, owner_id)
                else:
                    owner_id = subject_data["owner_id"]
                    owner_signature = tools.make_user_signature(
                        sender, subject_data, owner_signature, owner_id)
                return owner_signature

            def select_author_signature(sender, subject_data, video_comment):
                u"""Выбирает из словаря данные и формирует гиперссылку на автора."""
                author_signature = ""
                author_id = video_comment["from_id"]
                if str(author_id)[0] == "-":
                    author_signature = tools.make_group_signature(
                        sender, subject_data, author_signature, author_id)
                else:
                    author_id = video_comment["from_id"]
                    author_signature += tools.make_user_signature(
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
                str_date = tools.ts_date_to_str(ts_date, "%d.%m.%Y %H:%M:%S")
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
            str_date = tools.ts_date_to_str(
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
                    video_comments_data[i].update(
                        {"owner_id": video["owner_id"]})
                videos_comments_data.extend(video_comments_data)

        last_date = int(monitor_data["last_date"])

        new_videos_comments = []

        for item in reversed(videos_comments_data):
            if item["date"] > last_date:
                new_videos_comments.append(item)

        if len(new_videos_comments) > 0:
            videos_comments = tools.sort_items(new_videos_comments)

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
