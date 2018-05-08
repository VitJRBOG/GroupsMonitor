# coding: utf8


import time
import copy
import datetime
import logger
import notificator
import datamanager


def main(vk_admin_session, vk_bot_session):

    # TODO: Отрефакторить эту функцию, здесь слишком много хлама

    sender = "Core"

    try:
        PATH = datamanager.read_path(sender)

        data_file = datamanager.read_json(sender, PATH, "data")
        wiki_full_id = data_file["wiki_database_id"]
        data_wiki = datamanager.read_wiki(sender, vk_admin_session, wiki_full_id)

        if int(data_wiki["total_last_date"]) >\
           int(data_file["total_last_date"]):
            datamanager.write_json(sender, PATH, "data", data_wiki)

            date = datetime.datetime.now().strftime("%d.%m.%Y %H:%M:%S")

            mess_for_log = "Backup has been saved in file at " +\
                str(date) + "."
            logger.message_output(sender, mess_for_log)

        elif int(data_file["total_last_date"]) > int(data_wiki["total_last_date"]):
            datamanager.save_wiki(sender, vk_admin_session, wiki_full_id, data_file)

            date = datetime.datetime.now().strftime("%d.%m.%Y %H:%M:%S")

            mess_for_log = "Backup has been saved in wiki-page at " +\
                str(date) + "."
            logger.message_output(sender, mess_for_log)
        else:

            mess_for_log = "Data in wiki-page and data in file are identical."
            logger.message_output(sender, mess_for_log)

        data_file = None
        data_wiki = None

        delay = 0

        while True:
            data_json = datamanager.read_json(sender, PATH, "data")

            if delay >= 10:
                wiki_full_id = data_json["wiki_database_id"]
                datamanager.save_wiki(sender, vk_admin_session, wiki_full_id, data_json)

                date = datetime.datetime.now().strftime("%d.%m.%Y %H:%M:%S")

                mess_for_log = "Backup has been saved in wiki-page at " +\
                    str(date) + "."
                logger.message_output(sender, mess_for_log)

                delay = 0

            subjects = copy.deepcopy(data_json["subjects"])

            i = 0

            while i < len(subjects):
                sessions_list = {
                    "admin": vk_admin_session,
                    "bot": vk_bot_session
                }

                subject_data = copy.deepcopy(subjects[i])

                objNewPost = notificator.NewPost()

                response = objNewPost.new_post(sender, sessions_list, subject_data)

                last_date = int(subject_data["last_date"])

                n = len(response["items"]) - 1

                while n >= 0:
                    item = response["items"][n]

                    if item["date"] > last_date:

                        message, post_attachments =\
                            objNewPost.make_message(sender, vk_admin_session, item)

                        message_object = {
                            "message": message,
                            "post_attachments": post_attachments
                        }

                        objNewPost.send_message(sender, vk_bot_session, subject_data, message_object)

                        last_date = item["date"]

                        data_json["subjects"][i]["last_date"] = str(last_date)
                        if int(last_date) > int(data_json["total_last_date"]):
                            data_json["total_last_date"] = str(last_date)

                        datamanager.write_json(sender, PATH, "data", data_json)

                        date = datetime.datetime.fromtimestamp(
                                    int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

                        mess_for_log = subject_data["name"] +\
                            "'s new " +\
                            subject_data["filter"] +\
                            ": " + str(date)
                        logger.message_output(sender, mess_for_log)

                    n -= 1

                if subject_data["topic_notificator_settings"]["check_topics"] == 1:

                    subject_data = copy.deepcopy(data_json["subjects"][i])

                    objNewTopicMessage = notificator.NewTopicMessage()

                    response, subject_data, list_response = objNewTopicMessage.new_topic_message(sender, sessions_list, subject_data)

                    n = 0

                    while n < len(list_response):

                        comments_values = list_response[n]

                        j = len(comments_values["comments"]) - 1

                        while j >= 0:

                            item = comments_values["comments"][j]
                            last_date = comments_values["last_date"]

                            if item["date"] > int(last_date):

                                message, post_attachments =\
                                    objNewTopicMessage.make_message(sender, vk_admin_session, subject_data, comments_values, item)

                                message_object = {
                                    "message": message,
                                    "post_attachments": post_attachments
                                }

                                objNewTopicMessage.send_message(sender, vk_bot_session, subject_data, message_object)

                                last_date = item["date"]

                                k = 0

                                while k < len(subject_data["topics"]):

                                    if comments_values["topic_id"] ==\
                                      subject_data["topics"][k]["id"]:
                                        subject_data["topics"][k]["last_date"] = last_date

                                    k += 1

                                data_json["subjects"][i] = copy.deepcopy(subject_data)

                                x = 0

                                while x < len(subject_data["topics"]):
                                    topic = subject_data["topics"][x]

                                    if int(topic["last_date"]) > int(data_json["total_last_date"]):
                                        data_json["total_last_date"] = str(topic["last_date"])

                                    x += 1

                                datamanager.write_json(sender, PATH, "data", data_json)

                                date = datetime.datetime.fromtimestamp(
                                            int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

                                mess_for_log = comments_values["topic_title"] +\
                                    "'s new comment" + ": " + str(date)
                                logger.message_output(sender, mess_for_log)

                            j -= 1

                        n += 1

                if subject_data["photo_notificator_settings"]["check_photo"] == 1:

                    subject_data = copy.deepcopy(data_json["subjects"][i])

                    objNewAlbumPhoto = notificator.NewAlbumPhoto()

                    response = objNewAlbumPhoto.new_album_photo(sender, sessions_list, subject_data)

                    last_date = int(subject_data["photo_notificator_settings"]["last_date"])

                    n = len(response["items"]) - 1

                    while n >= 0:
                        item = response["items"][n]

                        if item["date"] > last_date:

                            album_response = objNewAlbumPhoto.get_album(sender, vk_admin_session, item)

                            album = {
                                "album_title": album_response["items"][0]["title"],
                                "album_id": album_response["items"][0]["id"]
                            }

                            item.update(album)

                            message, post_attachments =\
                                objNewAlbumPhoto.make_message(sender, vk_admin_session, item)

                            message_object = {
                                "message": message,
                                "post_attachments": post_attachments
                            }

                            objNewAlbumPhoto.send_message(sender, vk_bot_session, subject_data, message_object)

                            last_date = item["date"]

                            data_json["subjects"][i]["photo_notificator_settings"]["last_date"] = str(last_date)

                            if int(last_date) > int(data_json["total_last_date"]):
                                data_json["total_last_date"] = str(last_date)

                            datamanager.write_json(sender, PATH, "data", data_json)

                            date = datetime.datetime.fromtimestamp(
                                        int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

                            mess_for_log = album["album_title"] +\
                                "'s new photo" + ": " + str(date)
                            logger.message_output(sender, mess_for_log)

                        n -= 1

                i += 1

            delay += 1

            time.sleep(60)

    except Exception as var_except:
        logger.exception_handler(sender, var_except)
        return main(vk_admin_session, vk_bot_session)
