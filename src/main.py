# coding: utf8


import time
import copy
import datetime
import bughandler
import notificator
import datamanager


def core(vk_admin_session, vk_bot_session):

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
            print("COMPUTER [" + sender + "]: Backup has been saved " +
                  "in file at " + str(date) + ".")
        elif int(data_file["total_last_date"]) >\
           int(data_wiki["total_last_date"]):
            datamanager.save_wiki(sender, vk_admin_session, wiki_full_id, data_file)

            date = datetime.datetime.now().strftime("%d.%m.%Y %H:%M:%S")
            print("COMPUTER [" + sender + "]: Backup has been saved in " +
                  "wiki-page at " + str(date) + ".")
        else:
            print("COMPUTER [" + sender + "]: Data in wiki-page and " +
                  "data in file are identical.")

        data_file = None
        data_wiki = None

        delay = 0

        while True:
            data_json = datamanager.read_json(sender, PATH, "data")

            if delay >= 10:
                wiki_full_id = data_json["wiki_database_id"]
                datamanager.save_wiki(sender, vk_admin_session, wiki_full_id, data_json)

                date = datetime.datetime.now().strftime("%d.%m.%Y %H:%M:%S")
                print("COMPUTER [" + sender + "]: Backup has been saved in " +
                      "wiki-page at " + str(date) + ".")

                delay = 0

            subjects = copy.deepcopy(data_json["subjects"])

            i = 0

            while i < len(subjects):
                sessions_list = {
                    "admin": vk_admin_session,
                    "bot": vk_bot_session
                }

                subject_data = copy.deepcopy(subjects[i])

                last_date = notificator.new_post(sender, sessions_list, subject_data)

                data_json["subjects"][i]["last_date"] = str(last_date)
                if int(last_date) > int(data_json["total_last_date"]):
                    data_json["total_last_date"] = str(last_date)

                datamanager.write_json(sender, PATH, "data", data_json)

                if subject_data["topic_notificator_settings"]["check_topics"] == 1:

                    subject_data = copy.deepcopy(data_json["subjects"][i])

                    subject_data = notificator.new_topic_message(sender, sessions_list, subject_data)

                    data_json["subjects"][i] = copy.deepcopy(subject_data)

                    j = 0

                    while j < len(subject_data["topics"]):
                        topic = subject_data["topics"][j]

                        if int(topic["last_date"]) > int(data_json["total_last_date"]):
                            data_json["total_last_date"] = str(topic["last_date"])

                        j += 1

                    datamanager.write_json(sender, PATH, "data", data_json)

                if subject_data["photo_notificator_settings"]["check_photo"] == 1:

                    subject_data = copy.deepcopy(data_json["subjects"][i])

                    last_date = notificator.new_album_photo(sender, sessions_list, subject_data)

                    data_json["subjects"][i]["photo_notificator_settings"]["last_date"] = str(last_date)

                    if int(last_date) > int(data_json["total_last_date"]):
                        data_json["total_last_date"] = str(last_date)

                    datamanager.write_json(sender, PATH, "data", data_json)

                i += 1

            delay += 1

            time.sleep(60)

    except Exception as var_except:
        bughandler.exception_handler(sender, var_except)
        return core(vk_admin_session, vk_bot_session)
