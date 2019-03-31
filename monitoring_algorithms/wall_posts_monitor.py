# coding: utf8
u"""Модуль проверки постов на стене."""


import tools
import output_data
import request_handler
import data_manager


def run_monitoring_wall_posts(sender, res_filename, subject_data, monitor_data):
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
                    owner_signature = tools.make_group_signature(
                        sender, subject_data, owner_signature, owner_id)
                else:
                    owner_id = post["owner_id"]
                    owner_signature = tools.make_user_signature(
                        sender, subject_data, owner_signature, owner_id)
                return owner_signature

            def select_author_signature(sender, subject_data, post):
                u"""Выбирает из словаря данные и формирует гиперссылку на автора."""
                author_signature = ""
                owner_id = post["owner_id"]
                if str(owner_id)[0] == "-":
                    if "signer_id" in post:
                        author_id = post["signer_id"]
                        author_signature = tools.make_user_signature(
                            sender, subject_data, author_signature, author_id)
                    else:
                        if post["post_type"] == "suggest":
                            author_id = post["from_id"]
                            author_signature = tools.make_user_signature(
                                sender, subject_data, author_signature, author_id)
                        else:
                            author_signature += "[no data]"
                else:
                    author_id = post["owner_id"]
                    author_signature += tools.make_user_signature(
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
                    str(copy_history["owner_id"]) + \
                    "_" + str(copy_history["id"])
                return media_items

            def select_post_url(post):
                u"""Выбирает из словаря данные и формирует URL на пост."""
                post_url = "https://vk.com/wall"
                post_url += str(post["owner_id"]) + "_" + str(post["id"])
                return post_url

            def select_date(post):
                u"""Выбирает из словаря данные и определяет дату поста."""
                ts_date = post["date"]
                str_date = tools.ts_date_to_str(ts_date, "%d.%m.%Y %H:%M:%S")
                return str_date

            data_for_message = {}

            owner_signature = select_owner_signature(
                sender, subject_data, post)
            author_signature = select_author_signature(
                sender, subject_data, post)
            post_url = select_post_url(post)
            publication_date = select_date(post)

            text = "New " + post["post_type"].encode("utf8") + "\n"
            text += "Location: " + owner_signature.encode("utf8") + "\n"
            text += "Author: " + author_signature.encode("utf8") + "\n"
            text += "Created: " + str(publication_date).encode("utf8") + "\n\n"
            if len(post["text"].encode("utf8")) > 3500:
                text += post["text"].encode("utf8")[0:3500] + "\n"
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

            request_handler.send_message(
                sender, data_for_message, access_token)

        def show_message_about_new_post(sender, post):
            u"""Алгоритмы отображения сообщения о новом посте."""
            str_date = tools.ts_date_to_str(post["date"], "%d.%m.%Y %H:%M:%S")
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
            monitor_data = data_manager.read_json(
                path_to_res_file, res_filename)

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
            monitor_data = data_manager.read_json(
                path_to_res_file, res_filename)
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
            wall_posts = tools.sort_items(new_posts)

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
