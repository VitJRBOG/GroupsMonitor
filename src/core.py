# coding: utf8


import time
import copy
import datetime
import logger
import datamanager
import checker
from threading import Thread


def main(vk_admin_session, vk_bot_session):

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


                }








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
        subject_data = checker.check_for_posts(total_sender,
                                               PATH, path_to_subject_json,
                                               subject, subject_data,
                                               sessions_list)

    if subject_data["topic_checker_settings"]["check_topics"] == 1:
        subject_data = checker.check_for_topics(total_sender,
                                                PATH, path_to_subject_json,
                                                subject, subject_data,
                                                sessions_list)

    if subject_data["photo_checker_settings"]["check_photo"] == 1:
        subject_data = checker.check_for_albums(total_sender,
                                                PATH, path_to_subject_json,
                                                subject, subject_data,
                                                sessions_list)

    if subject_data["photo_comments_checker_settings"]["check_comments"] == 1:
        subject_data = checker.check_for_comments_photo(total_sender,
                                                        PATH,
                                                        path_to_subject_json,
                                                        subject, subject_data,
                                                        sessions_list)

    if subject_data["post_comments_checker_settings"]["check_comments"] == 1:
        subject_data = checker.check_for_comments_post(total_sender,
                                                       PATH,
                                                       path_to_subject_json,
                                                       subject, subject_data,
                                                       sessions_list)


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
