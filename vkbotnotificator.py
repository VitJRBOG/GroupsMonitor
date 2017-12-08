# coding: utf-8


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
                "publics": [],
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
    try:
        loads_json = json.loads(open(str(path) + str(file_name) +
                                     ".json", 'r').read())  # dict

        return loads_json
    except Exception as var_except:
        print(
            "COMPUTER [.. -> " + str(sender) +
            " -> Read JSON]: Error, " + str(var_except) +
            ". Exit from program...")
        exit(0)


def write_json(sender, path, file_name, loads_json):
    try:
        file_json = open(str(path) + str(file_name) + ".json", "w")
        file_json.write(json.dumps(loads_json, indent=4, ensure_ascii=False))
        file_json.close()

    except Exception as var_except:
        print(
            "COMPUTER [.. -> " + str(sender) +
            " -> Write JSON]: Error, " + str(var_except) +
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
            "COMPUTER [.. -> Authorization" + "]: Error, " +
            str(var_except) +
            ". Exit from program...")
        exit(0)


def notificator(sender, sessions_list, subject_data):
    sender += " -> Notificator"

    try:
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
                                int(item["date"])).strftime("%d.%m.%Y \
                                                            %H:%M:%S")

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
                                int(item["date"])).strftime("%d.%m.%Y \
                                                            %H:%M:%S")

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

                post_signature = get_signature(sender, vk_admin_session, item)
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

        def send_message(sender, vk_bot_session, subject_data, message_object):
            sender += " -> Send message"

            try:
                user_id = subject_data["send_to"]
                message = message_object["message"]
                post_attachments = message_object["post_attachments"]

                if post_attachments != "":
                    values = {
                        "user_id": user_id,
                        "message": message,
                        "attachment": post_attachments
                    }
                else:
                    values = {
                        "user_id": user_id,
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

                elif str(var_except).lower().find("failed to establish " +
                                                  "a new connection") != -1:
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

                message, post_attachments = make_message(sender,
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


def main(vk_admin_session, vk_bot_session):
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

                last_date = notificator(sender, sessions_list, subject_data)

                data_json["subjects"][i]["last_date"] = str(last_date)

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
