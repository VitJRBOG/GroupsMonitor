# coding: utf8
u"""Модуль проверки новых постов в топиках обсуждений."""


import tools
import data_manager
import output_data
import request_handler


def run_monitoring_topic_comments(sender, res_filename, subject_data, monitor_data):
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
                owner_signature = tools.make_group_signature(
                    sender, subject_data, owner_signature, owner_id)
                return owner_signature

            def select_author_signature(sender, subject_data, topic_comment):
                u"""Выбирает из словаря данные и формирует гиперссылку на автора."""
                author_signature = ""
                author_id = topic_comment["from_id"]
                if str(author_id)[0] == "-":
                    author_signature = tools.make_group_signature(
                        sender, subject_data, author_signature, author_id)
                else:
                    author_id = topic_comment["from_id"]
                    author_signature += tools.make_user_signature(
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
                str_date = tools.ts_date_to_str(ts_date, "%d.%m.%Y %H:%M:%S")
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
            str_date = tools.ts_date_to_str(
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
                    topic_comments_data[i].update(
                        {"topic_name": topic["title"]})
                topics_comments_data.extend(topic_comments_data)

        last_date = int(monitor_data["last_date"])

        new_topic_comments = []

        for item in reversed(topics_comments_data):
            if item["date"] > last_date:
                new_topic_comments.append(item)

        if len(new_topic_comments) > 0:
            topics_comments = tools.sort_items(new_topic_comments)

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
