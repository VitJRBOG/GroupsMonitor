# coding: utf8


import view_update
import time
import model.datamanager as datamanager
import model.checker as checker
import model.logger as logger


def load_data_for_checker(checker_name, subject, operation_name):
    PATH = datamanager.read_path()
    path_to_json = subject["path"]
    if len(path_to_json) > 0 and path_to_json[len(path_to_json) - 1] != "/":
        path_to_json = PATH + "bot_notificator/" + path_to_json + "/"
    else:
        path_to_json = PATH + "bot_notificator/" + path_to_json
    subject_section_data = datamanager.read_json(path_to_json, checker_name)
    subject_data = datamanager.read_json(path_to_json, "subject_data")
    sender = subject["name"]
    message = "Starting an operation " + operation_name + "."
    logger.message_output(sender, message)

    return PATH, path_to_json, subject_data, subject_section_data


# проверка значения out_flag
def time_for_end(PATH, subject_name, operation_name, out_flag, status_list):
    if out_flag.isSet():
        sender = subject_name
        message = "Completion of an operation " + operation_name + "."
        logger.message_output(sender, message)
        view_update.status_changer(status_list, "offline")
        return True
    else:
        return False


def run_post_checker(out_flag, status_list, sessions_list, subject):
    checker_name = "post_checker_settings"
    operation_name = "new post checking"
    PATH, path_to_json, subject_data, subject_section_data =\
        load_data_for_checker(checker_name, subject, operation_name)
    while True:
        view_update.status_changer(status_list, "processing")
        if time_for_end(PATH, subject["name"], operation_name,
                        out_flag, status_list):
            return
        sender = subject["name"] + " " + operation_name
        checker.check_for_posts(sender, path_to_json, subject, subject_data,
                                subject_section_data, sessions_list)
        view_update.status_changer(status_list, "waiting")
        interval = subject_section_data["interval"]
        i = 0
        while i < interval:
            time.sleep(1)
            if time_for_end(PATH, subject["name"], operation_name,
                            out_flag, status_list):
                return
            i += 1


def run_album_photo_checker(out_flag, status_list, sessions_list, subject):
    checker_name = "photo_checker_settings"
    operation_name = "new album photo checking"
    PATH, path_to_json, subject_data, subject_section_data =\
        load_data_for_checker(checker_name, subject, operation_name)
    while True:
        view_update.status_changer(status_list, "processing")
        if time_for_end(PATH, subject["name"], operation_name,
                        out_flag, status_list):
            return
        sender = subject["name"] + " " + operation_name
        checker.check_for_albums(sender, path_to_json, subject, subject_data,
                                 subject_section_data, sessions_list)
        view_update.status_changer(status_list, "waiting")
        interval = subject_section_data["interval"]
        i = 0
        while i < interval:
            time.sleep(1)
            if time_for_end(PATH, subject["name"], operation_name,
                            out_flag, status_list):
                return
            i += 1


def run_video_checker(out_flag, status_list, sessions_list, subject):
    checker_name = "video_checker_settings"
    operation_name = "new video checking"
    PATH, path_to_json, subject_data, subject_section_data =\
        load_data_for_checker(checker_name, subject, operation_name)
    while True:
        view_update.status_changer(status_list, "processing")
        if time_for_end(PATH, subject["name"], operation_name,
                        out_flag, status_list):
            return
        sender = subject["name"] + " " + operation_name
        checker.check_for_videos(sender, path_to_json, subject, subject_data,
                                 subject_section_data, sessions_list)
        view_update.status_changer(status_list, "waiting")
        interval = subject_section_data["interval"]
        i = 0
        while i < interval:
            time.sleep(1)
            if time_for_end(PATH, subject["name"], operation_name,
                            out_flag, status_list):
                return
            i += 1


def run_photo_comments_checker(out_flag, status_list, sessions_list, subject):
    checker_name = "photo_comments_checker_settings"
    operation_name = "new photo comment checking"
    PATH, path_to_json, subject_data, subject_section_data =\
        load_data_for_checker(checker_name, subject, operation_name)
    while True:
        view_update.status_changer(status_list, "processing")
        if time_for_end(PATH, subject["name"], operation_name,
                        out_flag, status_list):
            return
        sender = subject["name"] + " " + operation_name
        checker.check_for_comments_photo(sender, path_to_json, subject,
                                         subject_data,
                                         subject_section_data, sessions_list)
        view_update.status_changer(status_list, "waiting")
        interval = subject_section_data["interval"]
        i = 0
        while i < interval:
            time.sleep(1)
            if time_for_end(PATH, subject["name"], operation_name,
                            out_flag, status_list):
                return
            i += 1


def run_video_comments_checker(out_flag, status_list, sessions_list, subject):
    checker_name = "video_comments_checker_settings"
    operation_name = "new video comment checking"
    PATH, path_to_json, subject_data, subject_section_data =\
        load_data_for_checker(checker_name, subject, operation_name)
    while True:
        view_update.status_changer(status_list, "processing")
        if time_for_end(PATH, subject["name"], operation_name,
                        out_flag, status_list):
            return
        sender = subject["name"] + " " + operation_name
        checker.check_for_comments_video(sender, path_to_json, subject,
                                         subject_data,
                                         subject_section_data, sessions_list)
        view_update.status_changer(status_list, "waiting")
        interval = subject_section_data["interval"]
        i = 0
        while i < interval:
            time.sleep(1)
            if time_for_end(PATH, subject["name"], operation_name,
                            out_flag, status_list):
                return
            i += 1


def run_topic_comments_checker(out_flag, status_list, sessions_list, subject):
    checker_name = "topic_checker_settings"
    operation_name = "new topic comment checking"
    PATH, path_to_json, subject_data, subject_section_data =\
        load_data_for_checker(checker_name, subject, operation_name)
    while True:
        view_update.status_changer(status_list, "processing")
        if time_for_end(PATH, subject["name"], operation_name,
                        out_flag, status_list):
            return
        sender = subject["name"] + " " + operation_name
        checker.check_for_topics(sender, path_to_json, subject, subject_data,
                                 subject_section_data, sessions_list)
        view_update.status_changer(status_list, "waiting")
        interval = subject_section_data["interval"]
        i = 0
        while i < interval:
            time.sleep(1)
            if time_for_end(PATH, subject["name"], operation_name,
                            out_flag, status_list):
                return
            i += 1


def run_post_comments_checker(out_flag, status_list, sessions_list, subject):
    checker_name = "post_comments_checker_settings"
    operation_name = "new post comment checking"
    PATH, path_to_json, subject_data, subject_section_data =\
        load_data_for_checker(checker_name, subject, operation_name)
    while True:
        view_update.status_changer(status_list, "processing")
        if time_for_end(PATH, subject["name"], operation_name,
                        out_flag, status_list):
            return
        sender = subject["name"] + " " + operation_name
        checker.check_for_comments_post(sender, path_to_json, subject,
                                        subject_data,
                                        subject_section_data, sessions_list)
        view_update.status_changer(status_list, "waiting")
        interval = subject_section_data["interval"]
        i = 0
        while i < interval:
            time.sleep(1)
            if time_for_end(PATH, subject["name"], operation_name,
                            out_flag, status_list):
                return
            i += 1
