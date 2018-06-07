# coding: utf8


import time
import copy
import datetime
import logger
import dataloader
import datamanager
from threading import Thread


def main(vk_admin_session, vk_bot_session):

    # TODO: Отрефакторить эту функцию, здесь слишком много хлама

    total_sender = "Core"

    try:
        PATH = datamanager.read_path(total_sender)

        data_json = datamanager.read_json(total_sender, PATH, "data")

        subjects = copy.deepcopy(data_json["subjects"])

        i = 0

        while i < len(subjects):
            sender = total_sender + " -> " + subjects[i]["name"]

            if subjects[i]["check_subject"] == 1:
                if i > 0:
                    time.sleep(5)

                sessions_list = {
                    "admin": vk_admin_session,
                    "bot": vk_bot_session
                }

                interval = subjects[i]["interval"]

                objCommunitiChecker = CommunitiChecker(total_sender, PATH, interval, subjects[i], sessions_list)
                objCommunitiChecker.start()

            i += 1

    except Exception as var_except:
        logger.exception_handler(sender, var_except)
        return main(vk_admin_session, vk_bot_session)


def algorithm_checker(total_sender, PATH, subject, sessions_list, delay):

    def check_for_posts(total_sender, subject_data):
        sender = total_sender + " -> " + subject_data["name"] + " -> Post checking"

        objNewPost = dataloader.NewPost()

        response = objNewPost.new_post(sender, sessions_list, subject_data)

        last_date = int(subject_data["post_checker_settings"]["last_date"])

        def sort_items(response):
            items = response["items"]
            for j in range(len(items) - 1):
                f = 0
                for i in range(len(items) - 1 - j):
                    if items[i]["date"] < items[i + 1]["date"]:
                        x = items[i]
                        y = items[i + 1]
                        items[i + 1] = x
                        items[i] = y
                        f = 1
                if f == 0:
                    break
            return items
        items = sort_items(response)

        n = len(items) - 1

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

                subject_data["post_checker_settings"]["last_date"] = str(last_date)
                if int(last_date) > int(subject_data["total_last_date"]):
                    subject_data["total_last_date"] = str(last_date)

                datamanager.write_json(sender, PATH, subject["file_name"], subject_data)

                date = datetime.datetime.fromtimestamp(
                            int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

                mess_for_log = subject_data["name"] +\
                    "'s new " +\
                    subject_data["post_checker_settings"]["filter"] +\
                    ": " + str(date)
                logger.message_output(sender, mess_for_log)

            n -= 1

    def check_for_topics(total_sender, subject_data):
        sender = total_sender + " -> " + subject_data["name"] + " -> Topic checking"

        subject_data = datamanager.read_json(sender,
                                             path_to_subject_json,
                                             subject["file_name"])

        objNewTopicMessage = dataloader.NewTopicMessage()

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

                    x = 0

                    while x < len(subject_data["topics"]):
                        topic = subject_data["topics"][x]

                        if int(topic["last_date"]) > int(subject_data["total_last_date"]):
                            subject_data["total_last_date"] = str(topic["last_date"])

                        x += 1

                    datamanager.write_json(sender, PATH, subject["file_name"], subject_data)

                    date = datetime.datetime.fromtimestamp(
                                int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

                    mess_for_log = comments_values["topic_title"] +\
                        "'s new comment" + ": " + str(date)
                    logger.message_output(sender, mess_for_log)

                j -= 1

            n += 1

    def check_for_albums(total_sender, subject_data):
        sender = total_sender + " -> " + subject_data["name"] + " -> Photo checking"

        subject_data = datamanager.read_json(sender,
                                             path_to_subject_json,
                                             subject["file_name"])

        objNewAlbumPhoto = dataloader.NewAlbumPhoto()

        response = objNewAlbumPhoto.new_album_photo(sender, sessions_list, subject_data)

        last_date = int(subject_data["photo_checker_settings"]["last_date"])

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

                subject_data["photo_checker_settings"]["last_date"] = str(last_date)

                if int(last_date) > int(subject_data["total_last_date"]):
                    subject_data["total_last_date"] = str(last_date)

                datamanager.write_json(sender, PATH, subject["file_name"], subject_data)

                date = datetime.datetime.fromtimestamp(
                            int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

                mess_for_log = album["album_title"] +\
                    "'s new photo" + ": " + str(date)
                logger.message_output(sender, mess_for_log)

            n -= 1

    def check_for_comments_photo(total_sender, subject_data):
        sender = total_sender + " -> " + subject_data["name"] + " -> Photo comments checking"

        subject_data = datamanager.read_json(sender,
                                             path_to_subject_json,
                                             subject["file_name"])

        objNewPhotoComment = dataloader.NewPhotoComment()

        response = objNewPhotoComment.new_photo_comment(sender, sessions_list, subject_data)

        last_date = int(subject_data["photo_comments_checker_settings"]["last_date"])

        n = len(response["items"]) - 1
        while n >= 0:
            item = response["items"][n]

            if item["date"] > last_date:

                message, comment_attachments =\
                    objNewPhotoComment.make_message(sender, vk_admin_session, item, subject_data)

                message_object = {
                    "message": message,
                    "comment_attachments": comment_attachments
                }

                objNewPhotoComment.send_message(sender, vk_bot_session, subject_data, message_object)

                last_date = item["date"]

                subject_data["photo_comments_checker_settings"]["last_date"] = str(last_date)

                if int(last_date) > int(subject_data["total_last_date"]):
                    subject_data["total_last_date"] = str(last_date)

                datamanager.write_json(sender, PATH, subject["file_name"], subject_data)

                date = datetime.datetime.fromtimestamp(
                            int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

                mess_for_log = subject_data["name"] +\
                    "'s new comment under photo: " + str(date)
                logger.message_output(sender, mess_for_log)

            n -= 1

    def check_for_comments_post(total_sender, subject_data):
        sender = total_sender + " -> " + subject_data["name"] + " -> Post comments checking"

        objNewPostComment = dataloader.NewPostComment()

        response = objNewPostComment.get_posts(sender,
                                               sessions_list["admin"],
                                               subject_data)

        def sort_posts(posts):
            for j in range(len(posts) - 1):
                f = 0
                for i in range(len(posts) - 1 - j):
                    if posts[i] < posts[i + 1]:
                        x = posts[i]
                        y = posts[i + 1]
                        posts[i + 1] = x
                        posts[i] = y
                        f = 1
                if f == 0:
                    break
            return posts

        posts = response["items"]
        posts = sort_posts(posts)

        comments = []

        n = len(posts) - 1

        while n >= 0:
            post = posts[n]

            post_id = {
                "post_id": post["id"],
                "post_owner_id": post["owner_id"]
            }

            response = objNewPostComment.new_post_comment(sender,
                                                          sessions_list,
                                                          post, subject_data)
            items = response["items"]

            i = 0

            while i < len(items):
                items[i].update(post_id)

                i += 1

            comments.extend(items)

            n -= 1

        def sort_comments(comments):
            array = comments

            left = []
            equals = []
            right = []

            s = int((array[0]["date"] + array[int(len(array) / 2)]["date"] + array[len(array) - 1]["date"]) / 3)

            for item in array:
                if item["date"] > s:
                    left.append(item)
                elif item["date"] < s:
                    right.append(item)
                else:
                    equals.append(item)

            if len(left) > 1:
                left = sort_comments(left)
            if len(right) > 1:
                right = sort_comments(right)

            array = []
            array.extend(left)
            array.extend(equals)
            array.extend(right)

            return array

        if len(comments) > 1:
            comments = sort_comments(comments)

        last_date = int(subject_data["post_comments_checker_settings"]["last_date"])

        n = len(comments) - 1

        while n >= 0:
            item = comments[n]

            if item["date"] > last_date:

                check = False

                if subject_data["post_comments_checker_settings"]["check_by_attachments"] == 1:

                    if "attachments" in item:
                        attachments = item["attachments"]

                        i = 0

                        while i < len(attachments):

                            media_item = attachments[i]
                            if media_item["type"] == "photo" or\
                               media_item["type"] == "video" or\
                               media_item["type"] == "doc" or\
                               media_item["type"] == "link":
                                check = True

                            i += 1

                if subject_data["post_comments_checker_settings"]["check_by_keywords"] == 1:

                    if len(item["text"]) > 0 and\
                       len(subject_data["post_comments_checker_settings"]["keywords"]) == 0:
                        check = True

                    if len(item["text"]) > 0 and\
                       len(subject_data["post_comments_checker_settings"]["keywords"]) > 0:
                        text_array = item["text"].split(' ')
                        keywords = subject_data["post_comments_checker_settings"]["keywords"]

                        def search(line, underline):  # вместо find, у которого траблы с кодировками
                            last_i = -1
                            j = 0
                            i = 0
                            while i < len(line):

                                if line[i] == underline[j]:
                                    if last_i == -1:
                                        last_i = i
                                    j += 1
                                else:
                                    if last_i != -1:
                                        last_i = -1
                                    if j > 0:
                                        j = 0

                                if j == len(underline):
                                    return last_i

                                i += 1

                            if j < len(underline) - 1:
                                last_i = -1

                            return last_i

                        for word in text_array:
                            for keyword in keywords:
                                search_result = search(word, keyword)
                                if search_result != -1:
                                    check = True
                                    break

                if check:
                    message, comment_attachments =\
                        objNewPostComment.make_message(sender, vk_admin_session, item, subject_data)

                    message_object = {
                        "message": message,
                        "comment_attachments": comment_attachments
                    }

                    objNewPostComment.send_message(sender, vk_bot_session, subject_data, message_object)

                    last_date = item["date"]

                    subject_data["post_comments_checker_settings"]["last_date"] = str(last_date)

                    if int(last_date) > int(subject_data["total_last_date"]):
                        subject_data["total_last_date"] = str(last_date)

                    datamanager.write_json(sender, PATH, subject["file_name"], subject_data)

                    date = datetime.datetime.fromtimestamp(
                                int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

                    mess_for_log = subject_data["name"] +\
                        "'s new comment under post: " + str(date)
                    logger.message_output(sender, mess_for_log)

            n -= 1

    vk_admin_session = sessions_list["admin"]
    vk_bot_session = sessions_list["bot"]

    path_to_subject_json = subject["path"]

    if len(path_to_subject_json) > 0 and path_to_subject_json[0] != "/":
        path_to_subject_json = PATH + "/" + path_to_subject_json
    else:
        path_to_subject_json = PATH + path_to_subject_json

    subject_data = datamanager.read_json(total_sender,
                                         path_to_subject_json,
                                         subject["file_name"])

    sender = total_sender + " -> " + subject_data["name"]

    if delay == 0:
        datamanager.save_backup(sender, PATH, vk_admin_session, subject)

    if delay >= 100:
        datamanager.save_backup(sender, PATH, vk_admin_session, subject)

    if subject_data["post_checker_settings"]["check_posts"] == 1:
        check_for_posts(total_sender, subject_data)

    if subject_data["topic_checker_settings"]["check_topics"] == 1:
        check_for_topics(total_sender, subject_data)

    if subject_data["photo_checker_settings"]["check_photo"] == 1:
        check_for_albums(total_sender, subject_data)

    if subject_data["photo_comments_checker_settings"]["check_comments"] == 1:
        check_for_comments_photo(total_sender, subject_data)

    if subject_data["post_comments_checker_settings"]["check_comments"] == 1:
        check_for_comments_post(total_sender, subject_data)


class CommunitiChecker(Thread):
    def __init__(self, total_sender, PATH, interval, subject, sessions_list):
        Thread.__init__(self)
        self.total_sender = total_sender
        self.PATH = PATH
        self.interval = interval
        self.subject = subject
        self.sessions_list = sessions_list

    def run(self):
        delay = 0
        while True:
            algorithm_checker(self.total_sender, self.PATH, self.subject, self.sessions_list, delay)
            if delay >= 100:
                delay = 1
            else:
                if self.interval < 60:
                    delay += 1
                else:
                    delay += 10

            time.sleep(self.interval)
