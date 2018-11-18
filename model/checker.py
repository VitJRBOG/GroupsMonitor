# coding: utf8


import datetime
import model.datamanager as datamanager
import dataloader
import logger


def check_for_posts(sender, path_to_json, subject, subject_data,
                    subject_section_data, sessions_list):
    objNewPost = dataloader.NewPost()

    response = objNewPost.new_post(sender, sessions_list, subject_data,
                                   subject_section_data)

    last_date = int(subject_section_data["last_date"])

    def sort_posts(posts):
        for j in range(len(posts) - 1):
            f = 0
            for i in range(len(posts) - 1 - j):
                if posts[i]["date"] < posts[i + 1]["date"]:
                    x = posts[i]
                    y = posts[i + 1]
                    posts[i + 1] = x
                    posts[i] = y
                    f = 1
            if f == 0:
                break
        return posts

    items = response["items"]

    if len(items) > 1:
        items = sort_posts(items)

    n = len(items) - 1

    while n >= 0:
        item = response["items"][n]

        if item["date"] > last_date:

            message, post_attachments =\
                objNewPost.make_message(sender,
                                        sessions_list["admin_session"],
                                        item)

            message_object = {
                "message": message,
                "post_attachments": post_attachments
            }

            objNewPost.send_message(sender, sessions_list["bot_session"],
                                    subject_data,
                                    subject_section_data, message_object)

            last_date = item["date"]

            subject_section_data["last_date"] = str(last_date)
            # пока не знаю, что делать с total_last_date
            # if int(last_date) > int(subject_data["total_last_date"]):
            #     subject_data["total_last_date"] = str(last_date)

            datamanager.write_json(path_to_json,
                                   "post_checker_settings",
                                   subject_section_data)

            date = datetime.datetime.fromtimestamp(
                        int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

            post_type = subject_section_data["filter"]

            if post_type == "all":
                post_type = "post"

            mess_for_log = subject["name"] +\
                "'s new " +\
                post_type +\
                ": " + str(date)
            logger.message_output(sender, mess_for_log)

        n -= 1


def check_for_albums(sender, path_to_json, subject, subject_data,
                     subject_section_data, sessions_list):
    objNewAlbumPhoto = dataloader.NewAlbumPhoto()

    response = objNewAlbumPhoto.new_album_photo(sender, sessions_list,
                                                subject_data,
                                                subject_section_data)

    last_date = int(subject_section_data["last_date"])

    n = len(response["items"]) - 1

    while n >= 0:
        item = response["items"][n]

        if item["date"] > last_date:

            album_response =\
                objNewAlbumPhoto.get_album(sender,
                                           sessions_list["admin_session"],
                                           item)

            album = {
                "album_title": album_response["items"][0]["title"],
                "album_id": album_response["items"][0]["id"]
            }

            item.update(album)

            subject_id = {
                "subject_id": subject_data["owner_id"]
            }

            item.update(subject_id)

            message, post_attachments =\
                objNewAlbumPhoto.make_message(sender,
                                              sessions_list["admin_session"],
                                              item)

            message_object = {
                "message": message,
                "post_attachments": post_attachments
            }

            objNewAlbumPhoto.send_message(sender,
                                          sessions_list["bot_session"],
                                          subject_data, subject_section_data,
                                          message_object)

            last_date = item["date"]

            subject_section_data["last_date"] = str(last_date)
            # пока не знаю, что делать с total_last_date
            # if int(last_date) > int(subject_data["total_last_date"]):
            #     subject_data["total_last_date"] = str(last_date)

            datamanager.write_json(path_to_json,
                                   "photo_checker_settings",
                                   subject_section_data)

            date = datetime.datetime.fromtimestamp(
                        int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

            mess_for_log = album["album_title"] +\
                "'s new photo" + ": " + str(date)
            logger.message_output(sender, mess_for_log)

        n -= 1

    return subject_data


def check_for_videos(sender, path_to_json, subject, subject_data,
                     subject_section_data, sessions_list):
    objNewVideo = dataloader.NewVideo()

    response = objNewVideo.new_video(sender, sessions_list,
                                     subject_data, subject_section_data)

    last_date = int(subject_section_data["last_date"])

    n = len(response["items"]) - 1

    while n >= 0:
        item = response["items"][n]

        if item["date"] > last_date:

            subject_id = {
                "subject_id": item["owner_id"]
            }

            item.update(subject_id)

            message, post_attachments =\
                objNewVideo.make_message(sender,
                                         sessions_list["admin_session"],
                                         item)

            message_object = {
                "message": message,
                "post_attachments": post_attachments
            }

            objNewVideo.send_message(sender,
                                     sessions_list["bot_session"],
                                     subject_data, subject_section_data,
                                     message_object)

            last_date = item["adding_date"]

            subject_section_data["last_date"] = str(last_date)
            # пока не знаю, что делать с total_last_date
            # if int(last_date) > int(subject_data["total_last_date"]):
            #     subject_data["total_last_date"] = str(last_date)

            datamanager.write_json(path_to_json,
                                   "video_checker_settings",
                                   subject_section_data)

            date = datetime.datetime.fromtimestamp(
                        int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

            mess_for_log = subject_data["name"] +\
                "'s new video" + ": " + str(date)
            logger.message_output(sender, mess_for_log)

        n -= 1

    return subject_data


def check_for_comments_photo(sender, path_to_json, subject, subject_data,
                             subject_section_data, sessions_list):
    objNewPhotoComment = dataloader.NewPhotoComment()

    response = objNewPhotoComment.new_photo_comment(sender,
                                                    sessions_list,
                                                    subject_data,
                                                    subject_section_data)

    last_date = int(subject_section_data["last_date"])

    n = len(response["items"]) - 1
    while n >= 0:
        item = response["items"][n]

        if item["date"] > last_date:

            message, comment_attachments =\
                objNewPhotoComment.make_message(sender,
                                                sessions_list["admin_session"],
                                                item, subject_data)

            message_object = {
                "message": message,
                "comment_attachments": comment_attachments
            }

            objNewPhotoComment.send_message(sender,
                                            sessions_list["bot_session"],
                                            subject_data,
                                            subject_section_data,
                                            message_object)

            last_date = item["date"]

            subject_section_data["last_date"] = str(last_date)
            # пока не знаю, что делать с total_last_date
            # if int(last_date) > int(subject_data["total_last_date"]):
            #     subject_data["total_last_date"] = str(last_date)

            datamanager.write_json(path_to_json,
                                   "photo_comments_checker_settings",
                                   subject_section_data)

            date = datetime.datetime.fromtimestamp(
                        int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

            mess_for_log = subject_data["name"] +\
                "'s new comment under photo: " + str(date)
            logger.message_output(sender, mess_for_log)

        n -= 1

    return subject_data


def check_for_comments_video(sender, path_to_json, subject, subject_data,
                             subject_section_data, sessions_list):
    objNewVideoComment = dataloader.NewVideoComment()

    response = objNewVideoComment.get_videos(sender,
                                             sessions_list["admin_session"],
                                             subject_data, subject_section_data)

    last_date = int(subject_section_data["last_date"])

    videos = response["items"]

    comments = []

    n = len(videos) - 1

    while n >= 0:
        video = videos[n]

        response = []
        items = []

        if video["owner_id"] == subject_data["owner_id"]:
            response = objNewVideoComment.new_video_comment(sender,
                                                            sessions_list,
                                                            video, subject_data,
                                                            subject_section_data)
            items = response["items"]

        if len(items) > 0 and last_date < items[0]["date"]:

            video_id = {
                "video_id": video["id"],
                "video_owner_id": video["owner_id"]
            }

            i = 0

            while i < len(items):
                items[i].update(video_id)

                i += 1

            comments.extend(items)

        n -= 1

    def sort_comments(comments):
        array = comments

        left = []
        equals = []
        right = []

        s = int((array[0]["date"] + array[int(len(array) / 2)]["date"] +
                 array[len(array) - 1]["date"]) / 3)

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

    n = len(comments) - 1

    while n >= 0:
        item = comments[n]

        if item["date"] > last_date:

            message, comment_attachments =\
                objNewVideoComment.make_message(sender,
                                                sessions_list["admin_session"],
                                                item, subject_data)

            message_object = {
                "message": message,
                "comment_attachments": comment_attachments
            }

            objNewVideoComment.send_message(sender,
                                            sessions_list["bot_session"],
                                            subject_data,
                                            subject_section_data,
                                            message_object)

            last_date = item["date"]

            subject_section_data["last_date"] = str(last_date)
            # пока не знаю, что делать с total_last_date
            # if int(last_date) > int(subject_data["total_last_date"]):
            #     subject_data["total_last_date"] = str(last_date)

            datamanager.write_json(path_to_json,
                                   "video_comments_checker_settings",
                                   subject_section_data)

            date = datetime.datetime.fromtimestamp(
                        int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

            mess_for_log = subject_data["name"] +\
                "'s new comment under video: " + str(date)
            logger.message_output(sender, mess_for_log)

        n -= 1

    return subject_data


def check_for_topics(sender, path_to_json, subject, subject_data,
                     subject_section_data, sessions_list):
    objNewTopicMessage = dataloader.NewTopicMessage()

    response, subject_data, list_response =\
        objNewTopicMessage.new_topic_message(sender,
                                             sessions_list,
                                             subject_data,
                                             subject_section_data)

    n = 0

    while n < len(list_response):

        comments_values = list_response[n]

        j = len(comments_values["comments"]) - 1

        while j >= 0:

            item = comments_values["comments"][j]
            last_date = comments_values["last_date"]

            if item["date"] > int(last_date):

                message, post_attachments =\
                    objNewTopicMessage.make_message(sender,
                                                    sessions_list["admin_session"],
                                                    subject_data,
                                                    comments_values, item)

                message_object = {
                    "message": message,
                    "post_attachments": post_attachments
                }

                objNewTopicMessage.send_message(sender,
                                                sessions_list["bot_session"],
                                                subject_data,
                                                subject_section_data,
                                                message_object)

                last_date = item["date"]

                k = 0

                while k < len(subject_section_data["topics"]):

                    if comments_values["topic_id"] ==\
                      subject_section_data["topics"][k]["id"]:
                        subject_section_data["topics"][k]["last_date"] =\
                            last_date

                    k += 1

                # пока не знаю, что делать с total_last_date
                # x = 0
                # while x < len(subject_data["topics"]):
                #     topic = subject_data["topics"][x]
                #     if int(topic["last_date"]) >\
                #        int(subject_data["total_last_date"]):
                #         subject_data["total_last_date"] =\
                #             str(topic["last_date"])
                #     x += 1

                datamanager.write_json(path_to_json,
                                       "topic_checker_settings",
                                       subject_section_data)

                date = datetime.datetime.fromtimestamp(
                            int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

                mess_for_log = comments_values["topic_title"] +\
                    "'s new comment" + ": " + str(date)
                logger.message_output(sender, mess_for_log)

            j -= 1

        n += 1

    return subject_data


def check_for_comments_post(sender, path_to_json, subject, subject_data,
                            subject_section_data, sessions_list):
    objNewPostComment = dataloader.NewPostComment()

    response = objNewPostComment.get_posts(sender,
                                           sessions_list["admin_session"],
                                           subject_data, subject_section_data)

    last_date = int(subject_section_data["last_date"])

    posts = response["items"]

    comments = []

    n = len(posts) - 1

    while n >= 0:
        post = posts[n]

        response = objNewPostComment.new_post_comment(sender,
                                                      sessions_list,
                                                      post, subject_data,
                                                      subject_section_data)
        items = response["items"]

        if len(items) > 0 and last_date < items[0]["date"]:

            post_id = {
                "post_id": post["id"],
                "post_owner_id": post["owner_id"]
            }

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

        s = int((array[0]["date"] + array[int(len(array) / 2)]["date"] +
                 array[len(array) - 1]["date"]) / 3)

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

    n = len(comments) - 1

    while n >= 0:
        item = comments[n]

        if item["date"] > last_date:

            check = False

            if not check:

                if subject_section_data["check_by_communities"] == 1:

                    if str(item["from_id"])[0] == "-":
                        check = True

            if not check:

                if subject_section_data["check_by_attachments"] == 1:

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

            if not check:

                if subject_section_data["check_by_keywords"] == 1:

                    def char_changer(chars_for_changer, text):
                        chars = list(text)
                        for i in range(len(chars)):
                            if chars[i] in chars_for_changer:
                                chars[i] = chars_for_changer[chars[i]]
                        text = ''.join(chars)

                        return text

                    def check_algorithm(subject_section_data, text, check):

                        def check_by_len(condition, check,
                                         text, subject_section_data):
                            if len(text) == condition and\
                               len(subject_section_data["keywords"]) > 0:
                                text = text.lower()
                                keywords = subject_section_data["keywords"]

                                for keyword in keywords:
                                    if len(keyword) == condition:
                                        search_result =\
                                            text.find(keyword.lower())
                                        if search_result != -1:
                                            check = True
                                            break
                            return check

                        if len(text) > 0 and\
                           len(subject_section_data["keywords"]) == 0:
                            check = True

                        len_small_keywords = 2

                        for i in range(len_small_keywords):
                            check = check_by_len(i + 1, check,
                                                 text, subject_section_data)
                            if check:
                                break

                        if len(text) > len_small_keywords and\
                           len(subject_section_data["keywords"]) > 0:
                            text = text.lower()
                            keywords = subject_section_data["keywords"]

                            for keyword in keywords:
                                if len(keyword) > len_small_keywords:
                                    search_result = text.find(keyword.lower())
                                    if search_result != -1:
                                        check = True
                                        break
                        return check

                    if subject_section_data["check_by_charchange"] == 1:
                        text = item["text"]
                        chars_cyr_to_lat =\
                            subject_section_data["chars_cyr_to_lat"]
                        text = char_changer(chars_cyr_to_lat, text)
                        check = check_algorithm(subject_section_data, text, check)
                        if not check:
                            chars_lat_to_cyr =\
                                subject_section_data["chars_lat_to_cyr"]
                            text = char_changer(chars_lat_to_cyr, text)
                            check = check_algorithm(subject_section_data,
                                                    text, check)
                    else:
                        text = item["text"]
                        check = check_algorithm(subject_section_data,
                                                text, check)

            if check:
                message, comment_attachments =\
                    objNewPostComment.make_message(sender,
                                                   sessions_list["admin_session"],
                                                   item, subject_data)

                message_object = {
                    "message": message,
                    "comment_attachments": comment_attachments
                }

                objNewPostComment.send_message(sender,
                                               sessions_list["bot_session"],
                                               subject_data,
                                               subject_section_data,
                                               message_object)

                last_date = item["date"]

                subject_section_data["last_date"] = str(last_date)
                # пока не знаю, что делать с total_last_date
                # if int(last_date) > int(subject_data["total_last_date"]):
                #     subject_data["total_last_date"] = str(last_date)

                datamanager.write_json(path_to_json,
                                       "post_comments_checker_settings",
                                       subject_section_data)

                date = datetime.datetime.fromtimestamp(
                            int(last_date)).strftime("%d.%m.%Y %H:%M:%S")

                mess_for_log = subject_data["name"] +\
                    "'s new comment under post: " + str(date)
                logger.message_output(sender, mess_for_log)

        if n == 0 and last_date < item["date"]:
            last_date = item["date"]
            subject_section_data["last_date"] = str(last_date)

            datamanager.write_json(path_to_json,
                                   "post_comments_checker_settings",
                                   subject_section_data)

        n -= 1

    return subject_data
