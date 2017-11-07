# coding: utf8


import vk_api
import time


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


def starter():

    #  Токен аккаунта 404386296. Временная мера
    vk_user_token = "b4e7aa2e35dd290404eef768d72132e1454b9d" +\
                    "2150b7ad960c12f4bbf52d56dfcf88685cd7cafa25fe7fc"

    data_access = {
        "token": vk_user_token
    }

    vk_session = autorization(data_access, "token")

    main(vk_session)


def main(vk_session):

    vk_adminroom_token = "fec0fef8222c5a09d76b44674ca3fef0eecc602647" +\
                         "defb2978a5ee03ddf7da886bc35629f35ebd0f77781"


    last_date = int(open("last_date.txt", "r").read())

    # Вечный цикл

    while True:

        response = get_post(vk_session, -61061413)

        i = len(response["items"]) - 1

        while i >= 0:

            item = response["items"][i]

            if item["date"] > last_date:

                # ПЕРЕД ПРОДАКШЕНОМ ПОМЕНЯТЬ vk_adminroom_token

                last_date = send_message(vk_adminroom_token, item, 404386296, last_date)

                file_text = open("last_date.txt", "w")
                file_text.write(str(last_date))
                file_text.close()

                print(str("New last date: " + str(last_date)))

                time.sleep(1)

            i -= 1

        # Задержка на 60 секунд

        time.sleep(60)


def get_post(vk_session, owner_id):

    values = {
        'owner_id': int(owner_id),
        'count': 10,
        'filter': "suggests"
    }

    response = vk_session.method("wall.get", values)

    return response


def send_message(access_token, item, send_to, last_date):

    post_info = {
            "text": item["text"],
            "from_id": item["from_id"],
            "id": item["id"],
            "owner_id": item["owner_id"]
        }

    author = str(post_info["from_id"])
    id_post = str(post_info["owner_id"]) + "_" + str(post_info["id"])

    text_post = post_info["text"] +\
                "\n|----------" +\
                "\n|" + "https://vk.com/id" + author +\
                "\n|" + "https://vk.com/wall" + id_post +\
                "\n|----------"

    values = {
            "user_id": send_to,
            "message": text_post
        }


    data_access = {
        "token": access_token
    }

    #  Авторизация
    vk_message_session = autorization(data_access, "token")

    vk_message_session.method("messages.send", values)

    return item["date"]


starter()
