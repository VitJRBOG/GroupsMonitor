# coding: utf8
u"""Модуль проверки комментариев под постами на стене."""


import re
import tools
import output_data
import request_handler
import data_manager


def run_monitoring_wall_post_comments(sender, res_filename, subject_data, monitor_data):
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
                    owner_signature = tools.make_group_signature(
                        sender, subject_data, owner_signature, owner_id)
                else:
                    owner_id = subject_data["owner_id"]
                    owner_signature = tools.make_user_signature(
                        sender, subject_data, owner_signature, owner_id)
                return owner_signature

            def select_author_signature(sender, subject_data, wall_post_comment):
                u"""Выбирает из словаря данные и формирует гиперссылку на автора."""
                author_signature = ""
                author_id = wall_post_comment["from_id"]
                if str(author_id)[0] == "-":
                    author_signature = tools.make_group_signature(
                        sender, subject_data, author_signature, author_id)
                else:
                    author_id = wall_post_comment["from_id"]
                    author_signature += tools.make_user_signature(
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
                    post_url += "&thread=" + \
                        str(wall_post_comment["parents_stack"][0])
                return post_url

            def select_date(wall_post_comment):
                u"""Выбирает из словаря данные и определяет дату комментария."""
                ts_date = wall_post_comment["date"]
                str_date = tools.ts_date_to_str(ts_date, "%d.%m.%Y %H:%M:%S")
                return str_date

            data_for_message = {}

            owner_signature = select_owner_signature(
                sender, subject_data)
            author_signature = select_author_signature(
                sender, subject_data, wall_post_comment)
            wall_post_comment_url = select_wall_post_comment_url(
                wall_post_comment, subject_data)
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
            str_date = tools.ts_date_to_str(
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
            monitor_data = data_manager.read_json(
                path_to_res_file, res_filename)
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
                        user_name = users_info[0]["first_name"] + \
                            " " + users_info[0]["last_name"]
                        if user_name == name:
                            suspicious = True
                            return suspicious
            return suspicious

        def check_by_empty_comment(suspicious, item):
            u"""Проверка по пустому комментарию."""
            if len(item["text"]) < 1:
                if "attachments" in item:
                    return suspicious
                else:
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
                if len(re.findall(r"[0-9]", sym)) > 0:
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
                        if i < len(symbs_from_text) - 1:
                            if len(re.findall(r"[0-9]", symbs_from_text[i + 1])) == 0:
                                suspicious = True
                                return suspicious
                        else:
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
            monitor_data = data_manager.read_json(
                path_to_res_file, res_filename)

            small_messages = monitor_data["check_by_keywords"]["small_messages"]
            suspicious = check_small_messages_algorithm(
                suspicious, text, small_messages)
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
            wall_posts_comments = tools.sort_items(new_wall_post_comments)

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
                suspicious = suspicion_check(
                    sender, values, subject_data, item)
                if suspicious:
                    found_new_wall_post_comment(
                        sender, values, subject_data, item)
                update_last_date(values, item)  
