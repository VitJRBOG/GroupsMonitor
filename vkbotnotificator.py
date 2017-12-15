# coding: utf8


import vk_api
import time
import os
import json
import copy
import datetime


def starter():

    try:

        if os.path.exists("path.txt") is False:
            file_text = open("path.txt", "w")
            file_text.write("")
            file_text.close()

            print("COMPUTER: Was created file \"path.txt\".")

        PATH = read_path()

        if os.path.exists(PATH + "data.json") is False:
            print("\nCOMPUTER: WARNING! File \"data.json\" not found!")

            data_json = {
                "subjects": [
                    {
                        "name": "",
                        "topics": [],
                        "topic_notificator_settings": {
                            "post_count": 0,
                            "send_to": 0,
                            "check_topics": 0
                        },
                        "photo_notificator_settings": {
                            "photo_count": 0,
                            "check_photo": 0,
                            "send_to": 0,
                            "last_date": "0"
                        },
                        "send_to": 0,
                        "filter": "",
                        "last_date": "0",
                        "posts_count": 0,
                        "owner_id": 0
                    }
                ],
                "admin_token": "",
                "bot_token": ""
            }

            write_json("Starter", PATH, "data", data_json)

            print("COMPUTER: Please, enter the necessary data to " +
                  "file \"data.json\". Exit from program...")
            exit(0)

        #  Получение данных из файла JSON

        data_json = read_json("Starter", PATH, "data")

        vk_admin_token = data_json["admin_token"]
        vk_bot_token = data_json["bot_token"]

        data_access_admin = {
            "token": vk_admin_token
        }
        data_access_bot = {
            "token": vk_bot_token
        }

        vk_admin_session = autorization(data_access_admin, "token")
        vk_bot_session = autorization(data_access_bot, "token")

        print("COMPUTER [Starter]: Program was started.")
        main(vk_admin_session, vk_bot_session)

    except Exception as var_except:
        print(
            "COMPUTER: Error, " + str(var_except) + ". Exit from program...")
        exit(0)


def read_path():
    try:
        path = str(open("path.txt", "r").read())

        if len(path) > 0 and path[len(path) - 1] != "/":
            path += "/"

        return path

    except Exception as var_except:
        print(
            "COMPUTER [.. -> Read \"path.txt\"]: Error, " + str(var_except) +
            ". Exit from program...")
        exit(0)


def read_json(sender, path, file_name):
    sender += " -> Read JSON"

    try:
        loads_json = json.loads(open(str(path) + str(file_name) +
                                     ".json", 'r').read())  # dict

        return loads_json
    except Exception as var_except:
        print(
            "COMPUTER [" + str(sender) +
            "]: Error, " + str(var_except) +
            ". Exit from program...")
        exit(0)


def write_json(sender, path, file_name, loads_json):
    sender += " -> Write JSON"

    try:
        file_json = open(str(path) + str(file_name) + ".json", "w")
        file_json.write(json.dumps(loads_json, indent=4, ensure_ascii=True))
        file_json.close()

    except Exception as var_except:
        print(
            "COMPUTER [" + str(sender) +
            "]: Error, " + str(var_except) +
            ". Exit from program...")
        exit(0)


def autorization(data_access, auth_type):

    try:

        if auth_type == "token":

            #  Авторизация по токену
            access_token = data_access["token"]
            vk_session = vk_api.VkApi(token=access_token)
            vk_session._auth_token()

        if auth_type == "login":

            #  Авторизация по имени пользователя и паролю
            vk_login = data_access["login"]
            vk_passwd = data_access["password"]
            vk_session = vk_api.VkApi(login=vk_login, password=vk_passwd)
            vk_session.auth()

        if auth_type != "token" and auth_type != "login":

            print("COMPUTER: Error of authorization. Exit from program...")
            exit(0)

        return vk_session

    except Exception as var_except:
        print(
            "COMPUTER [.. -> Authorization]: Error, " +
            str(var_except) +
            ". Exit from program...")
        exit(0)


class Notificator():

    def new_post(self, sender, sessions_list, subject_data):

        try:
            self.sender = sender
            self.sessions_list = sessions_list
            self.subject_data = subject_data

            sender += " -> Notificator -> New post"

            vk_admin_session = sessions_list["admin"]
            vk_bot_session = sessions_list["bot"]

            def get_posts(sender, vk_admin_session, subject_data):
                sender += " -> Get post"

                try:
                    owner_id = int(subject_data["owner_id"])
                    count = int(subject_data["posts_count"])
                    post_filter = str(subject_data["filter"])

                    values = {
                        'owner_id': owner_id,
                        'count': count,
                        'filter': post_filter
                    }

                    time.sleep(1)

                    response = vk_admin_session.method("wall.get", values)

                    return response

                # Общий обработчик исключений
                # TODO: Добавить частные обработчики

                except Exception as var_except:
                    if str(var_except).lower().find("captcha needed") !=\
                       -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return get_posts(sender, vk_admin_session, subject_data)

                    elif str(var_except).lower().find("failed to establish " +
                                                      "a new connection") != -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return get_posts(sender, vk_admin_session, subject_data)

                    elif str(var_except).lower().find("connection aborted") != -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return get_posts(sender, vk_admin_session, subject_data)

                    else:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) +
                            ". Exit from program...")
                        exit(0)

            def make_message(sender, vk_admin_session, item):
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
                            print(
                                "COMPUTER [" + sender + "]: Error, " +
                                str(var_except) +
                                ". Exit from program...")
                            exit(0)

                    # Функция возвращает текст из поста
                    def get_text(sender, item):
                        sender += " -> Get text"

                        try:
                            post_text = ""

                            post_text = item["text"]

                            return post_text

                        except Exception as var_except:
                            print(
                                "COMPUTER [" + sender + "]: Error, " +
                                str(var_except) +
                                ". Exit from program...")
                            exit(0)

                    # Функция возвращает URL поста
                    def get_url(sender, item):
                        sender += " -> Get URL"

                        try:
                            post_url = ""

                            id_post = str(item["owner_id"]) + "_" + str(item["id"])

                            post_url = "https://vk.com/wall" + id_post

                            return post_url

                        except Exception as var_except:
                            print(
                                "COMPUTER [" + sender + "]: Error, " +
                                str(var_except) +
                                ". Exit from program...")
                            exit(0)

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
                            print(
                                "COMPUTER [" + sender + "]: Error, " +
                                str(var_except) +
                                ". Exit from program...")
                            exit(0)

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
                    print(
                        "COMPUTER [" + sender + "]: Error, " +
                        str(var_except) +
                        ". Exit from program...")
                    exit(0)

            def send_message(sender, vk_bot_session,
                             subject_data, message_object):
                sender += " -> Send message"

                try:
                    peer_id = subject_data["send_to"]
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
                    if str(var_except).lower().find("captcha needed") !=\
                       -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return send_message(sender, vk_bot_session,
                                            subject_data, message_object)

                    elif str(var_except).lower().find("connection aborted") != -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return send_message(sender, vk_bot_session,
                                            subject_data, message_object)

                    elif str(var_except).lower().find("failed to establish a new connection") != -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return send_message(sender, vk_bot_session,
                                            subject_data, message_object)

                    else:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) +
                            ". Exit from program...")
                        exit(0)

            response = get_posts(sender, vk_admin_session, subject_data)

            last_date = int(subject_data["last_date"])

            i = len(response["items"]) - 1

            while i >= 0:
                item = response["items"][i]

                if item["date"] > last_date:

                    message, post_attachments =\
                        make_message(sender,
                                     vk_admin_session,
                                     item)

                    message_object = {
                        "message": message,
                        "post_attachments": post_attachments
                    }

                    send_message(sender, vk_bot_session,
                                 subject_data, message_object)

                    last_date = item["date"]

                    date = datetime.datetime.fromtimestamp(
                                int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

                    print(subject_data["name"] + "'s new " +
                          subject_data["filter"] + ": " + str(date))

                i -= 1

            return last_date

        except Exception as var_except:
            print(
                "COMPUTER [" + str(sender) + "]: Error, " +
                str(var_except) +
                ". Exit from program...")
            exit(0)

    def new_topic_message(self, sender, sessions_list, subject_data):

        try:
            self.sender = sender
            self.sessions_list = sessions_list
            self.subject_data = subject_data

            vk_admin_session = sessions_list["admin"]
            vk_bot_session = sessions_list["bot"]

            sender += " -> Notificator -> New topic message"

            def get_topics(sender, vk_admin_session, subject_data):
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
                    if str(var_except).lower().find("captcha needed") !=\
                       -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return get_posts(sender, vk_admin_session, subject_data)

                    elif str(var_except).lower().find("failed to establish " +
                                                      "a new connection") != -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return get_posts(sender, vk_admin_session, subject_data)

                    elif str(var_except).lower().find("connection aborted") != -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return get_posts(sender, vk_admin_session, subject_data)

                    else:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) +
                            ". Exit from program...")
                        exit(0)

            def checking_existence(sender, subject_data, response):
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
                    print(
                        "COMPUTER [" + str(sender) + "]: Error, " +
                        str(var_except) +
                        ". Exit from program...")
                    exit(0)

            def get_comments(sender, vk_admin_session, subject_data):
                sender += " -> Get comments"

                try:

                    list_response = []

                    owner_id = int(subject_data["owner_id"])

                    if str(owner_id)[0] == "-":
                        owner_id = int(str(owner_id)[1:])

                    topic_notificator_settings = \
                        subject_data["topic_notificator_settings"]

                    i = 0

                    while i < len(subject_data["topics"]):

                        topic_id = int(subject_data["topics"][i]["id"])

                        values = {
                            "count": topic_notificator_settings["post_count"],
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
                    print(
                        "COMPUTER [" + str(sender) + "]: Error, " +
                        str(var_except) +
                        ". Exit from program...")
                    exit(0)

            def make_message(sender, vk_admin_session,
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
                            print(
                                "COMPUTER [" + str(sender) + "]: Error, " +
                                str(var_except) +
                                ". Exit from program...")
                            exit(0)

                    def get_text(sender, item):
                        sender += " -> Get text"

                        try:

                            post_text = ""

                            post_text = item["text"]

                            return post_text

                        except Exception as var_except:
                            print(
                                "COMPUTER [" + str(sender) + "]: Error, " +
                                str(var_except) +
                                ". Exit from program...")
                            exit(0)

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
                            print(
                                "COMPUTER [" + str(sender) + "]: Error, " +
                                str(var_except) +
                                ". Exit from program...")
                            exit(0)

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
                            print(
                                "COMPUTER [" + str(sender) + "]: Error, " +
                                str(var_except) +
                                ". Exit from program...")
                            exit(0)

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
                    print(
                        "COMPUTER [" + str(sender) + "]: Error, " +
                        str(var_except) +
                        ". Exit from program...")
                    exit(0)

            def send_message(sender, vk_bot_session,
                             subject_data, message_object):
                sender += " -> Send message"

                try:
                    peer_id = subject_data["topic_notificator_settings"]["send_to"]
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
                    if str(var_except).lower().find("captcha needed") !=\
                       -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return send_message(sender, vk_bot_session,
                                            subject_data, message_object)

                    elif str(var_except).lower().find("connection aborted") != -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return send_message(sender, vk_bot_session,
                                            subject_data, message_object)

                    elif str(var_except).lower().find("failed to establish a new connection") != -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return send_message(sender, vk_bot_session,
                                            subject_data, message_object)

                    else:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) +
                            ". Exit from program...")
                        exit(0)

            response = get_topics(sender, vk_admin_session, subject_data)
            subject_data = checking_existence(sender, subject_data, response)
            list_response = get_comments(sender,
                                         vk_admin_session, subject_data)

            i = 0

            while i < len(list_response):

                comments_values = list_response[i]

                j = len(comments_values["comments"]) - 1

                while j >= 0:

                    item = comments_values["comments"][j]
                    last_date = comments_values["last_date"]

                    if item["date"] > int(last_date):

                        message, post_attachments =\
                            make_message(sender, vk_admin_session,
                                         subject_data, comments_values, item)

                        message_object = {
                            "message": message,
                            "post_attachments": post_attachments
                        }

                        send_message(sender, vk_bot_session,
                                     subject_data, message_object)

                        last_date = item["date"]

                        n = 0

                        while n < len(subject_data["topics"]):

                            if comments_values["topic_id"] ==\
                               subject_data["topics"][n]["id"]:
                               subject_data["topics"][n]["last_date"] = last_date

                            n += 1

                        date = datetime.datetime.fromtimestamp(
                                    int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

                        print(comments_values["topic_title"] + "'s new " +
                              "comment" + ": " + str(date))

                    j -= 1

                i += 1

            return subject_data

        except Exception as var_except:
            print(
                "COMPUTER [" + str(sender) + "]: Error, " +
                str(var_except) +
                ". Exit from program...")
            exit(0)

    def new_album_photo(self, sender, sessions_list, subject_data):

        try:
            self.sender = sender
            self.sessions_list = sessions_list
            self.subject_data = subject_data

            vk_admin_session = sessions_list["admin"]
            vk_bot_session = sessions_list["bot"]

            sender += " -> Notificator -> New album photo"

            def get_photo(sender, vk_admin_session, subject_data):
                sender += " -> Get photo"

                try:
                    settings = subject_data["photo_notificator_settings"]

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
                    if str(var_except).lower().find("captcha needed") !=\
                       -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return get_posts(sender, vk_admin_session, subject_data)

                    elif str(var_except).lower().find("failed to establish " +
                                                      "a new connection") != -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return get_posts(sender, vk_admin_session, subject_data)

                    elif str(var_except).lower().find("connection aborted") != -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return get_posts(sender, vk_admin_session, subject_data)

                    else:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) +
                            ". Exit from program...")
                        exit(0)

            def get_album(sender, vk_admin_session, item):
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
                    if str(var_except).lower().find("captcha needed") !=\
                       -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return get_posts(sender, vk_admin_session, subject_data)

                    elif str(var_except).lower().find("failed to establish " +
                                                      "a new connection") != -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return get_posts(sender, vk_admin_session, subject_data)

                    elif str(var_except).lower().find("connection aborted") != -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return get_posts(sender, vk_admin_session, subject_data)

                    else:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) +
                            ". Exit from program...")
                        exit(0)


            def make_message(sender, vk_admin_session, item):
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
                            print(
                                "COMPUTER [" + sender + "]: Error, " +
                                str(var_except) +
                                ". Exit from program...")
                            exit(0)

                    # Функция возвращает текст из поста
                    def get_text(sender, item):
                        sender += " -> Get text"

                        try:
                            post_text = ""

                            post_text = item["text"]

                            return post_text

                        except Exception as var_except:
                            print(
                                "COMPUTER [" + sender + "]: Error, " +
                                str(var_except) +
                                ". Exit from program...")
                            exit(0)

                    # Функция возвращает URL поста
                    def get_url(sender, item):
                        sender += " -> Get URL"

                        try:
                            post_url = ""

                            id_post = str(item["owner_id"]) + "_" + str(item["id"])

                            post_url = "https://vk.com/photo" + id_post

                            return post_url

                        except Exception as var_except:
                            print(
                                "COMPUTER [" + sender + "]: Error, " +
                                str(var_except) +
                                ". Exit from program...")
                            exit(0)

                    # Функция возвращает прикрепления к посту
                    def get_attachments(sender, item):
                        sender += " -> Get attachments"

                        try:
                            media = "photo" + str(item["owner_id"]) +\
                                "_" + str(item["id"])

                            return media

                        except Exception as var_except:
                            print(
                                "COMPUTER [" + sender + "]: Error, " +
                                str(var_except) +
                                ". Exit from program...")
                            exit(0)

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
                    print(
                        "COMPUTER [" + sender + "]: Error, " +
                        str(var_except) +
                        ". Exit from program...")
                    exit(0)

            def send_message(sender, vk_bot_session,
                             subject_data, message_object):
                sender += " -> Send message"

                try:
                    peer_id = subject_data["photo_notificator_settings"]["send_to"]
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
                    if str(var_except).lower().find("captcha needed") !=\
                       -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return send_message(sender, vk_bot_session,
                                            subject_data, message_object)

                    elif str(var_except).lower().find("connection aborted") != -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return send_message(sender, vk_bot_session,
                                            subject_data, message_object)

                    elif str(var_except).lower().find("failed to establish a new connection") != -1:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) + ". " +
                            "Timeout: 60 sec.")
                        time.sleep(60)

                        return send_message(sender, vk_bot_session,
                                            subject_data, message_object)

                    else:
                        print(
                            "COMPUTER [" + sender + "]: Error, " +
                            str(var_except) +
                            ". Exit from program...")
                        exit(0)

            response = get_photo(sender, vk_admin_session, subject_data)

            last_date = int(subject_data["photo_notificator_settings"]["last_date"])

            i = len(response["items"]) - 1

            while i >= 0:
                item = response["items"][i]

                if item["date"] > last_date:

                    album_response = get_album(sender, vk_admin_session, item)

                    album = {
                        "album_title": album_response["items"][0]["title"],
                        "album_id": album_response["items"][0]["id"]
                    }

                    item.update(album)

                    message, post_attachments =\
                        make_message(sender,
                                     vk_admin_session,
                                     item)

                    message_object = {
                        "message": message,
                        "post_attachments": post_attachments
                    }

                    send_message(sender, vk_bot_session,
                                 subject_data, message_object)

                    last_date = item["date"]

                    date = datetime.datetime.fromtimestamp(
                                int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

                    print(album["album_title"] + "'s new " +
                          "photo" + ": " + str(date))

                i -= 1

            return last_date

        except Exception as var_except:
            print(
                "COMPUTER [" + str(sender) + "]: Error, " +
                str(var_except) +
                ". Exit from program...")
            exit(0)


def main(vk_admin_session, vk_bot_session):

    # TODO: Отрефакторить эту функцию, здесь слишком много хлама

    sender = "Main"

    try:
        PATH = read_path()

        while True:
            data_json = read_json(sender, PATH, "data")

            subjects = copy.deepcopy(data_json["subjects"])

            i = 0

            while i < len(subjects):
                sessions_list = {
                    "admin": vk_admin_session,
                    "bot": vk_bot_session
                }

                subject_data = copy.deepcopy(subjects[i])

                objNotificator = Notificator()

                last_date = objNotificator.new_post(sender, sessions_list,
                                                    subject_data)

                data_json["subjects"][i]["last_date"] = str(last_date)

                write_json(sender, PATH, "data", data_json)

                if subject_data["topic_notificator_settings"]["check_topics"] == 1:

                    subject_data = copy.deepcopy(data_json["subjects"][i])

                    subject_data = objNotificator.new_topic_message(sender,
                                                                    sessions_list,
                                                                    subject_data)

                    data_json["subjects"][i] = copy.deepcopy(subject_data)

                    write_json(sender, PATH, "data", data_json)

                if subject_data["photo_notificator_settings"]["check_photo"] == 1:

                    subject_data = copy.deepcopy(data_json["subjects"][i])

                    last_date = objNotificator.new_album_photo(sender,
                                                               sessions_list,
                                                               subject_data)

                    data_json["subjects"][i]["photo_notificator_settings"]["last_date"] = str(last_date)

                    write_json(sender, PATH, "data", data_json)

                i += 1

            time.sleep(60)

    except Exception as var_except:
        print(
            "COMPUTER [" + str(sender) + "]: Error, " +
            str(var_except) +
            ". Exit from program...")
        exit(0)


starter()
