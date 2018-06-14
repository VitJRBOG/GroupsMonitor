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

        threads_list = []

        out_flag = {
            "time_to_end": False
        }

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

                objCommunitiChecker = CommunitiChecker(out_flag, total_sender,
                                                       PATH, interval,
                                                       subjects[i],
                                                       sessions_list)
                objCommunitiChecker.start()

                thread_name = "Checker of " + subjects[i]["name"]

                values = {
                    "thread_name": thread_name,
                    "thread": objCommunitiChecker
                }

                threads_list.append(values)

            i += 1

        objThreadListener = ThreadListener(out_flag, total_sender,
                                           threads_list)
        objThreadListener.start()

        return out_flag

    except Exception as var_except:
        logger.exception_handler(sender, var_except)
        return main(vk_admin_session, vk_bot_session)


def algorithm_checker(total_sender, PATH, subject,
                      sessions_list, time_for_backup):

    path_to_subject_json = subject["path"]

    if len(path_to_subject_json) > 0 and path_to_subject_json[0] != "/":
        path_to_subject_json = PATH + "/" + path_to_subject_json
    else:
        path_to_subject_json = PATH + path_to_subject_json

    subject_data = datamanager.read_json(total_sender,
                                         path_to_subject_json,
                                         subject["file_name"])

    sender = total_sender + " -> " + subject_data["name"]

    if time_for_backup:
        datamanager.save_backup(sender, PATH, sessions_list["admin"], subject)

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
    def __init__(self, out_flag, total_sender, PATH, interval, subject, sessions_list):
        Thread.__init__(self)
        self.out_flag = out_flag
        self.total_sender = total_sender
        self.PATH = PATH
        self.interval = interval
        self.subject = subject
        self.sessions_list = sessions_list

    def run(self):
        time_for_backup = True
        next_time = 60 * 10 + int(datetime.datetime.now().strftime('%s'))

        def flag_checker(out_flag):
            if out_flag["time_to_end"]:
                sender = "Thread " + self.subject["name"]
                message = "Completion implementation a thread of" +\
                    " listening for " + self.subject["name"] + "."
                logger.message_output(sender, message)
                return True
        while True:
            if flag_checker(self.out_flag):
                return
            algorithm_checker(self.total_sender, self.PATH, self.subject,
                              self.sessions_list, time_for_backup)
            now_time = int(datetime.datetime.now().strftime('%s'))
            if now_time < next_time:
                time_for_backup = False
            else:
                time_for_backup = True
                next_time = 60 * 10 + now_time

            i = 0
            while i < self.interval:
                time.sleep(1)
                if flag_checker(self.out_flag):
                    return
                i += 1


class ThreadListener(Thread):
    def __init__(self, out_flag, sender, threads_list):
        Thread.__init__(self)
        self.out_flag = out_flag
        self.sender = sender
        self.threads_list = threads_list

    def run(self):
        def thread_listener(sender, threads_list):
            count_threads = len(threads_list)
            i = 0
            while i < count_threads:
                if not threads_list[i]["thread"].is_alive():
                    message = "Thread " + threads_list[i]["thread_name"] +\
                        " is not responding."
                    logger.message_output(sender, message)
                    threads_list[i].pop()
                    count_threads = len(threads_list)
                i += 1

        while True:
            if self.out_flag["time_to_end"]:
                sender = "Thread listener"
                message = "Completion implementation a thread of" +\
                    " listening for a thread listener."
                logger.message_output(sender, message)
                return
            thread_listener(self.sender, self.threads_list)
            time.sleep(5)
