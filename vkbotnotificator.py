# coding: utf-8


import vk_api
import time
import os
import json
import copy
import urllib

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

        #  Файлы с токеном, временная мера

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


def main(vk_admin_session, vk_bot_session):

    PATH = read_path()

    # Вечный цикл

    while True:

        data_json = read_json("Main", PATH, "data")

        publics = copy.deepcopy(data_json["publics"])

        i = 0

        while i < len(publics):

            response = get_post(vk_admin_session, publics[i]["id"])

            j = len(response["items"]) - 1

            while j >= 0:

                item = response["items"][j]

                last_date = int(publics[i]["last_date"])

                if item["date"] > last_date:

                    last_date = send_message(vk_bot_session, item,
                                             publics[i]["id"], 
                                             last_date)

                    data_json["publics"][i]["last_date"] = str(last_date)

                    write_json("Main", PATH, "data", data_json)

                    print(str("New last date " + publics[i]["name"] +
                              ": " + str(last_date)))

                    time.sleep(1)

                j -= 1

            i += 1

        # Задержка на 60 секунд

        time.sleep(60)


def get_post(vk_admin_session, owner_id):

    values = {
        'owner_id': int(owner_id),
        'count': 10,
        'filter': "suggests"
    }

    response = vk_admin_session.method("wall.get", values)

    return response


def send_message(vk_bot_session, item, send_to, last_date):

    text = ""

    if len(item["text"]) > 1000:
        text = item["text"][0:1000] + "... \n [long text]"
    else:
        text = item["text"]

    post_info = {
            "text": text,
            "from_id": item["from_id"],
            "id": item["id"],
            "owner_id": item["owner_id"]
        }

    response_author = vk_bot_session.method("users.get", {
                                            "user_ids": item["from_id"]})

    first_name = response_author[0]["first_name"]
    last_name = response_author[0]["last_name"]

    author_full_name = first_name + " " + last_name

    author_url = "*id" + str(item["from_id"]) + " (" + author_full_name + ")"

    id_post = str(post_info["owner_id"]) + "_" + str(post_info["id"])

    text_post = author_url + "\n" +\
                "\n" + post_info["text"] + "\n" +\
                "\n" + "https://vk.com/wall" + id_post

    values = {
            "user_id": send_to,
            "message": text_post
        }

    vk_bot_session.method("messages.send", values)

    return item["date"]


starter()
