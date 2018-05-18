# coding: utf-8


import logger
import time
import datetime
import copy


class NewPost:

    def get_posts(self, sender, vk_admin_session, subject_data):
        sender += " -> Get post"

        try:
            owner_id = int(subject_data["owner_id"])
            count = int(subject_data["post_checker_settings"]["posts_count"])
            post_filter = str(subject_data["post_checker_settings"]["filter"])

            values = {
                'owner_id': owner_id,
                'count': count,
                'filter': post_filter
            }

            time.sleep(1)

            response = vk_admin_session.method("wall.get", values)

            return response

        except Exception as var_except:
            logger.exception_handler(sender, var_except)
            return self.get_posts(sender, vk_admin_session, subject_data)

    def make_message(self, sender, vk_admin_session, item):
        sender += " -> Make message"

        message = ""

        try:

            # Функция возвращает автора поста
            # и время размещения в предложке
            def get_signature(sender, vk_admin_session, item):
                sender += " -> Get signature"

                try:
                    post_signature = ""

                    if str(item["from_id"])[0] == "-":
                        author_values = {
                            "group_id": int(str(item["from_id"])[1:])
                        }

                        time.sleep(1)

                        response_author =\
                            vk_admin_session.method("groups.getById",
                                                    author_values)

                        author_name = response_author[0]["name"]

                        author_url = "*" +\
                            response_author[0]["screen_name"] + " " +\
                            "(" + author_name + ")"

                        date = datetime.datetime.fromtimestamp(
                            int(item["date"])).strftime("%d.%m.%Y %H:%M:%S")

                        post_signature = author_url + "\n" + str(date)

                    else:
                        author_values = {
                            "user_ids": item["from_id"]
                        }

                        time.sleep(1)

                        response_author =\
                            vk_admin_session.method("users.get",
                                                    author_values)

                        first_name = response_author[0]["first_name"]
                        last_name = response_author[0]["last_name"]

                        author_full_name = first_name + " " + last_name

                        author_url = "*id" + str(item["from_id"]) +\
                            " (" + author_full_name + ")"

                        date = datetime.datetime.fromtimestamp(
                            int(item["date"])).strftime("%d.%m.%Y %H:%M:%S")

                        post_signature = author_url + "\n" + str(date)

                    return post_signature

                except Exception as var_except:
                    logger.exception_handler(sender, var_except)
                    return get_signature(sender, vk_admin_session,
                                         item)

            # Функция возвращает текст из поста
            def get_text(sender, item):
                sender += " -> Get text"

                try:
                    post_text = ""

                    post_text = item["text"]

                    return post_text

                except Exception as var_except:
                    logger.exception_handler(sender, var_except)
                    return get_text(sender, item)

            # Функция возвращает URL поста
            def get_url(sender, item):
                sender += " -> Get URL"

                try:
                    post_url = ""

                    id_post = str(item["owner_id"]) + "_" + str(item["id"])

                    post_url = "https://vk.com/wall" + id_post

                    return post_url

                except Exception as var_except:
                    logger.exception_handler(sender, var_except)
                    return get_url(sender, item)

            # Функция возвращает прикрепления к посту
            def get_attachments(sender, item):
                sender += " -> Get attachments"

                try:
                    list_media = []

                    if "attachments" in item:
                        attachments = item["attachments"]

                        i = 0
                        while i < len(attachments):
                            media_item = attachments[i]

                            if media_item["type"] == "photo" or\
                               media_item["type"] == "video" or\
                               media_item["type"] == "audio" or\
                               media_item["type"] == "doc":

                                media = media_item[media_item["type"]]

                                id_media = media_item["type"] +\
                                    str(media["owner_id"]) +\
                                    "_" + str(media["id"])

                                if "access_key" in media:
                                    id_media += "_" + media["access_key"]

                                list_media.append(id_media)

                            i += 1

                    if "copy_history" in item:
                        repost = item["copy_history"][0]

                        post_url = "wall" +\
                            str(repost["owner_id"]) + "_" +\
                            str(repost["id"])

                        if "access_key" in repost:
                            post_url += "_" + repost["access_key"]

                        list_media.append(post_url)

                    if len(list_media) > 0:
                        return ",".join(list_media)

                    else:
                        return ""

                except Exception as var_except:
                    logger.exception_handler(sender, var_except)
                    return get_attachments(sender, item)

            post_signature = get_signature(sender,
                                           vk_admin_session,
                                           item)
            post_text = get_text(sender, item)
            post_url = get_url(sender, item)
            post_attachments = get_attachments(sender, item)

            mes_long_text = "...\n[long text]"

            post_length = len(post_signature + "\n\n" +
                              post_text +
                              mes_long_text + "\n\n" +
                              post_url)

            limit_symbols = 3900

            if post_length > limit_symbols:
                count_symbols = post_length -\
                    (post_length - limit_symbols) - 1
                post_text = post_text[0:count_symbols]

                message = post_signature + "\n\n" +\
                    post_text +\
                    mes_long_text + "\n\n" +\
                    post_url
            else:
                message = post_signature + "\n\n" +\
                    post_text + "\n\n" +\
                    post_url

            return message, post_attachments

        except Exception as var_except:
            logger.exception_handler(sender, var_except)
            return self.make_message(sender, vk_admin_session, item)

    def send_message(self, sender, vk_bot_session,
                     subject_data, message_object):
        sender += " -> Send message"

        try:
            peer_id = subject_data["post_checker_settings"]["send_to"]
            message = message_object["message"]
            post_attachments = message_object["post_attachments"]

            if post_attachments != "":
                values = {
                    "peer_id": peer_id,
                    "message": message,
                    "attachment": post_attachments
                }
            else:
                values = {
                    "peer_id": peer_id,
                    "message": message
                }

            time.sleep(1)

            vk_bot_session.method("messages.send", values)

        except Exception as var_except:
            logger.exception_handler(sender, var_except)
            return self.send_message(sender, vk_bot_session, subject_data, message_object)

    def new_post(self, sender, sessions_list, subject_data):

        try:

            sender += " -> Notificator -> New post"

            vk_admin_session = sessions_list["admin"]
            # vk_bot_session = sessions_list["bot"]

            response = self.get_posts(sender, vk_admin_session, subject_data)

            return response

        except Exception as var_except:
            logger.exception_handler(sender, var_except)
            return self.new_post(sender, sessions_list, subject_data)


class NewTopicMessage:

    def get_topics(self, sender, vk_admin_session, subject_data):
        sender += " -> Get topics"

        try:
            owner_id = int(subject_data["owner_id"])

            if str(owner_id)[0] == "-":
                owner_id = int(str(owner_id)[1:])

            values = {
                'group_id': owner_id
            }

            time.sleep(1)

            response = vk_admin_session.method("board.getTopics",
                                               values)

            # items, потому что response постов,
            # и response топиков отличается

            return response["items"]

        except Exception as var_except:
            logger.exception_handler(sender, var_except)
            return self.get_topics(sender, vk_admin_session, subject_data)

    def checking_existence(self, sender, subject_data, response):
        sender += " -> Checking existence"

        try:

            # Проверка существования топика в базе

            topics_subject = copy.deepcopy(subject_data["topics"])

            if len(topics_subject) > 0:
                i = 0

                while i < len(response):

                    response_item = response[i]
                    not_exist = False

                    j = 0

                    while j < len(topics_subject):

                        topics_subject_item = topics_subject[j]

                        if response_item["id"] ==\
                           topics_subject_item["id"]:
                            not_exist = False
                            if subject_data["topics"][j]["title"] !=\
                               response_item["title"]:
                                subject_data["topics"][j]["title"] =\
                                    response_item["title"]
                            break
                        else:
                            not_exist = True

                        j += 1

                    if not_exist:
                        topic_values = {
                            "last_date": "0",
                            "id": response_item["id"],
                            "title": response_item["title"]
                        }
                        subject_data["topics"].append(copy.deepcopy(topic_values))

                    i += 1

                # Проверка существования топика в группе

                topics_subject = copy.deepcopy(subject_data["topics"])

                i = 0

                while i < len(topics_subject):

                    topics_subject_item = topics_subject[i]
                    not_exist = False

                    j = 0

                    while j < len(response):

                        response_item = response[j]

                        if topics_subject_item["id"] ==\
                           response_item["id"]:
                            not_exist = False
                            break
                        else:
                            not_exist = True

                        j += 1

                    if not_exist:
                        subject_data["topics"].pop(i)

                    i += 1

            else:
                i = 0

                while i < len(response):

                    response_item = response[i]

                    topic_values = {
                        "last_date": "0",
                        "id": response_item["id"],
                        "title": response_item["title"]
                    }
                    subject_data["topics"].append(copy.deepcopy(topic_values))

                    i += 1

            return subject_data

        except Exception as var_except:
            logger.exception_handler(sender, var_except)
            return self.checking_existence(sender, subject_data, response)

    def get_comments(self, sender, vk_admin_session, subject_data):
        sender += " -> Get comments"

        try:

            list_response = []

            owner_id = int(subject_data["owner_id"])

            if str(owner_id)[0] == "-":
                owner_id = int(str(owner_id)[1:])

            topic_checker_settings = \
                subject_data["topic_checker_settings"]

            i = 0

            while i < len(subject_data["topics"]):

                topic_id = int(subject_data["topics"][i]["id"])

                values = {
                    "count": topic_checker_settings["post_count"],
                    "group_id": owner_id,
                    "topic_id": topic_id,
                    "sort": "desc"
                }

                response = vk_admin_session.method("board.getComments",
                                                   values)

                comments_values = {
                    "owner_id": subject_data["owner_id"],
                    "topic_id": topic_id,
                    "topic_title": subject_data["topics"][i]["title"],
                    "last_date": subject_data["topics"][i]["last_date"],
                    "comments": copy.deepcopy(response["items"])
                }

                list_response.append(copy.deepcopy(comments_values))

                i += 1

            return list_response

        except Exception as var_except:
            logger.exception_handler(sender, var_except)
            return self.get_comments(sender, vk_admin_session, subject_data)

    def make_message(self, sender, vk_admin_session,
                     subject_data, comments_values, item):
        sender += " -> Make message"

        try:
            def get_signature(sender, vk_admin_session,
                              comments_values, item):
                sender += " -> Get signature"

                try:

                    post_signature = "Topic: "

                    post_signature += comments_values["topic_title"] +\
                        "\n"

                    if str(item["from_id"])[0] == "-":
                        author_values = {
                            "group_id": int(str(item["from_id"])[1:])
                        }

                        time.sleep(1)

                        response_author =\
                            vk_admin_session.method("groups.getById",
                                                    author_values)

                        author_name = response_author[0]["name"]

                        author_url = "*" +\
                            response_author[0]["screen_name"] + " " +\
                            "(" + author_name + ")"

                        date = datetime.datetime.fromtimestamp(
                            int(item["date"])).strftime("%d.%m.%Y %H:%M:%S")

                        post_signature = author_url + "\n" + str(date)

                    else:
                        author_values = {
                            "user_ids": item["from_id"]
                        }

                        time.sleep(1)

                        response_author =\
                            vk_admin_session.method("users.get",
                                                    author_values)

                        first_name = response_author[0]["first_name"]
                        last_name = response_author[0]["last_name"]

                        author_full_name = first_name + " " + last_name

                        author_url = "*id" + str(item["from_id"]) +\
                            " (" + author_full_name + ")"

                        date = datetime.datetime.fromtimestamp(
                            int(item["date"])).strftime("%d.%m.%Y %H:%M:%S")

                        post_signature += author_url + "\n" + str(date)

                    return post_signature

                except Exception as var_except:
                    logger.exception_handler(sender, var_except)
                    return get_signature(sender, vk_admin_session,
                                         comments_values, item)

            def get_text(sender, item):
                sender += " -> Get text"

                try:

                    post_text = ""

                    post_text = item["text"]

                    return post_text

                except Exception as var_except:
                    logger.exception_handler(sender, var_except)
                    return get_text(sender, item)

            def get_url(sender, comments_values, item):
                sender += " -> Get URL"

                try:

                    comment_url = ""

                    topic_id = str(comments_values["owner_id"]) + "_" +\
                        str(comments_values["topic_id"])

                    comment_url = "https://vk.com/topic" + topic_id +\
                        "?post=" + str(item["id"])

                    return comment_url

                except Exception as var_except:
                    logger.exception_handler(sender, var_except)
                    return get_url(sender, comments_values, item)

            def get_attachments(sender, item):
                sender += " -> Get attachments"

                try:

                    list_media = []

                    if "attachments" in item:
                        attachments = item["attachments"]

                        i = 0
                        while i < len(attachments):
                            media_item = attachments[i]

                            if media_item["type"] == "photo" or\
                               media_item["type"] == "video" or\
                               media_item["type"] == "audio" or\
                               media_item["type"] == "doc":

                                media = media_item[media_item["type"]]

                                id_media = media_item["type"] +\
                                    str(media["owner_id"]) +\
                                    "_" + str(media["id"])

                                if "access_key" in media:
                                    id_media += "_" + media["access_key"]

                                list_media.append(id_media)

                            i += 1

                    if len(list_media) > 0:
                        return ",".join(list_media)

                    else:
                        return ""

                except Exception as var_except:
                    logger.exception_handler(sender, var_except)
                    return get_attachments(sender, item)

            post_signature = get_signature(sender, vk_admin_session,
                                           comments_values, item)
            post_text = get_text(sender, item)
            post_url = get_url(sender, comments_values, item)
            post_attachments = get_attachments(sender, item)

            mes_long_text = "...\n[long text]"

            post_length = len(post_signature + "\n\n" +
                              post_text +
                              mes_long_text + "\n\n" +
                              post_url)

            limit_symbols = 3900

            if post_length > limit_symbols:
                count_symbols = post_length -\
                    (post_length - limit_symbols) - 1
                post_text = post_text[0:count_symbols]

                message = post_signature + "\n\n" +\
                    post_text +\
                    mes_long_text + "\n\n" +\
                    post_url
            else:
                message = post_signature + "\n\n" +\
                    post_text + "\n\n" +\
                    post_url

            return message, post_attachments

        except Exception as var_except:
            logger.exception_handler(sender, var_except)
            return self.make_message(sender, vk_admin_session, subject_data, comments_values, item)

    def send_message(self, sender, vk_bot_session,
                     subject_data, message_object):
        sender += " -> Send message"

        try:
            peer_id = subject_data["topic_checker_settings"]["send_to"]
            message = message_object["message"]
            post_attachments = message_object["post_attachments"]

            if post_attachments != "":
                values = {
                    "peer_id": peer_id,
                    "message": message,
                    "attachment": post_attachments
                }
            else:
                values = {
                    "peer_id": peer_id,
                    "message": message
                }

            time.sleep(1)

            vk_bot_session.method("messages.send", values)

        except Exception as var_except:
            logger.exception_handler(sender, var_except)
            return self.send_message(sender, vk_bot_session, subject_data, message_object)

    def new_topic_message(self, sender, sessions_list, subject_data):

        try:

            vk_admin_session = sessions_list["admin"]
            # vk_bot_session = sessions_list["bot"]

            sender += " -> Notificator -> New topic message"

            response = self.get_topics(sender, vk_admin_session, subject_data)
            subject_data = self.checking_existence(sender, subject_data, response)
            list_response = self.get_comments(sender, vk_admin_session, subject_data)

            return response, subject_data, list_response

        except Exception as var_except:
            logger.exception_handler(sender, var_except)
            return self.new_topic_message(sender, sessions_list, subject_data)


class NewAlbumPhoto:

    def get_photo(self, sender, vk_admin_session, subject_data):
        sender += " -> Get photo"

        try:
            settings = subject_data["photo_checker_settings"]

            owner_id = int(subject_data["owner_id"])
            count = int(settings["photo_count"])

            values = {
                "owner_id": owner_id,
                "count": count,
                "no_service_albums": 1
            }

            time.sleep(1)

            response = vk_admin_session.method("photos.getAll", values)

            return response

        except Exception as var_except:
            logger.exception_handler(sender, var_except)
            return self.get_photo(sender, vk_admin_session, subject_data)

    def get_album(self, sender, vk_admin_session, item):
        sender += " -> Get album"

        try:
            owner_id = int(item["owner_id"])
            album_id = int(item["album_id"])

            values = {
                "owner_id": owner_id,
                "album_ids": album_id
            }

            time.sleep(1)

            response = vk_admin_session.method("photos.getAlbums", values)

            return response

        except Exception as var_except:
            logger.exception_handler(sender, var_except)
            return self.get_album(sender, vk_admin_session, item)

    def make_message(self, sender, vk_admin_session, item):
        sender += " -> Make message"

        message = ""

        try:

            # Функция возвращает автора поста
            # и время размещения в предложке
            def get_signature(sender, vk_admin_session, item):
                sender += " -> Get signature"

                try:
                    post_signature = "Album: "

                    post_signature += item["album_title"] +\
                        "\n"

                    author_values = {
                            "user_ids": item["user_id"]
                        }

                    time.sleep(1)

                    response_author =\
                        vk_admin_session.method("users.get",
                                                author_values)

                    first_name = response_author[0]["first_name"]
                    last_name = response_author[0]["last_name"]

                    author_full_name = first_name + " " + last_name

                    author_url = "*id" + str(item["user_id"]) +\
                        " (" + author_full_name + ")"

                    date = datetime.datetime.fromtimestamp(
                        int(item["date"])).strftime("%d.%m.%Y %H:%M:%S")

                    post_signature += author_url + "\n" + str(date)

                    return post_signature

                except Exception as var_except:
                    logger.exception_handler(sender, var_except)
                    return get_signature(sender, vk_admin_session,
                                         item)

            # Функция возвращает текст из поста
            def get_text(sender, item):
                sender += " -> Get text"

                try:
                    post_text = ""

                    post_text = item["text"]

                    return post_text

                except Exception as var_except:
                    logger.exception_handler(sender, var_except)
                    return get_text(sender, item)

            # Функция возвращает URL поста
            def get_url(sender, item):
                sender += " -> Get URL"

                try:
                    post_url = ""

                    id_post = str(item["owner_id"]) + "_" + str(item["id"])

                    post_url = "https://vk.com/photo" + id_post

                    return post_url

                except Exception as var_except:
                    logger.exception_handler(sender, var_except)
                    return get_url(sender, item)

            # Функция возвращает прикрепления к посту
            def get_attachments(sender, item):
                sender += " -> Get attachments"

                try:
                    media = "photo" + str(item["owner_id"]) +\
                        "_" + str(item["id"])

                    return media

                except Exception as var_except:
                    logger.exception_handler(sender, var_except)
                    return get_attachments(sender, item)

            post_signature = get_signature(sender,
                                           vk_admin_session,
                                           item)
            post_text = get_text(sender, item)
            post_url = get_url(sender, item)
            post_attachments = get_attachments(sender, item)

            mes_long_text = "...\n[long text]"

            post_length = len(post_signature + "\n\n" +
                              post_text +
                              mes_long_text + "\n\n" +
                              post_url)

            limit_symbols = 3900

            if post_length > limit_symbols:
                count_symbols = post_length -\
                    (post_length - limit_symbols) - 1
                post_text = post_text[0:count_symbols]

                message = post_signature + "\n\n" +\
                    post_text +\
                    mes_long_text + "\n\n" +\
                    post_url
            else:
                message = post_signature + "\n\n" +\
                    post_text + "\n\n" +\
                    post_url

            return message, post_attachments

        except Exception as var_except:
            logger.exception_handler(sender, var_except)
            return self.make_message(sender, vk_admin_session, item)

    def send_message(self, sender, vk_bot_session,
                     subject_data, message_object):
        sender += " -> Send message"

        try:
            peer_id = subject_data["photo_checker_settings"]["send_to"]
            message = message_object["message"]
            post_attachments = message_object["post_attachments"]

            if post_attachments != "":
                values = {
                    "peer_id": peer_id,
                    "message": message,
                    "attachment": post_attachments
                }
            else:
                values = {
                    "peer_id": peer_id,
                    "message": message
                }

            time.sleep(1)

            vk_bot_session.method("messages.send", values)

        except Exception as var_except:
            logger.exception_handler(sender, var_except)
            return self.send_message(sender, vk_bot_session, subject_data, message_object)

    def new_album_photo(self, sender, sessions_list, subject_data):

        try:

            vk_admin_session = sessions_list["admin"]
            # vk_bot_session = sessions_list["bot"]

            sender += " -> Notificator -> New album photo"

            response = self.get_photo(sender, vk_admin_session, subject_data)

            return response

        except Exception as var_except:
            logger.exception_handler(sender, var_except)
            return self.new_album_photo(sender, sessions_list, subject_data)
