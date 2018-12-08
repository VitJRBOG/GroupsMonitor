# coding: utf8
u"""Модуль алгоритмов проверки."""


import time
import data_manager
import output_data


def before_start_operations(sender):
    u"""Операции перед началом проверки."""
    message = "Started..."
    output_data.output_text_row(sender, message)


def before_end_operations(sender):
    u"""Операции перед окончанием проверки."""
    message = "Stopped..."
    output_data.output_text_row(sender, message)


def run_post_checker(subject_name, subject_data, end_flag):
    checker_name = "post_checker_settings"
    operation_name = "new post checking"
    sender = subject_name + "'s " + operation_name
    PATH = data_manager.read_path()
    subject_path = PATH + subject_data["path"]
    checker_data = data_manager.read_json(subject_path, checker_name)
    if checker_data["check_flag"] != 1:
        return
    while True:
        before_start_operations(sender)
        if end_flag.isSet():
            before_end_operations(sender)
            return
        #### СТАРЫЙ ФРАГМЕНТ, НУЖНО МЕНЯТЬ
        checker.check_for_posts(sender, path_to_json, subject, subject_data,
                                subject_section_data, sessions_list)
        #### СТАРЫЙ ФРАГМЕНТ, НУЖНО МЕНЯТЬ
        interval = checker_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                before_end_operations(sender)
                return


def run_album_photo_checker(subject_name, subject_data, end_flag):
    checker_name = "photo_checker_settings"
    operation_name = "new album photo checking"
    sender = subject_name + "'s " + operation_name
    PATH = data_manager.read_path()
    subject_path = PATH + subject_data["path"]
    checker_data = data_manager.read_json(subject_path, checker_name)
    if checker_data["check_flag"] != 1:
        return
    while True:
        before_start_operations(sender)
        if end_flag.isSet():
            before_end_operations(sender)
            return
        #### СТАРЫЙ ФРАГМЕНТ, НУЖНО МЕНЯТЬ
        checker.check_for_albums(sender, path_to_json, subject, subject_data,
                                 subject_section_data, sessions_list)
        #### СТАРЫЙ ФРАГМЕНТ, НУЖНО МЕНЯТЬ
        interval = checker_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                before_end_operations(sender)
                return


def run_video_checker(subject_name, subject_data, end_flag):
    checker_name = "video_checker_settings"
    operation_name = "new video checking"
    sender = subject_name + "'s " + operation_name
    PATH = data_manager.read_path()
    subject_path = PATH + subject_data["path"]
    checker_data = data_manager.read_json(subject_path, checker_name)
    if checker_data["check_flag"] != 1:
        return
    while True:
        before_start_operations(sender)
        if end_flag.isSet():
            before_end_operations(sender)
            return
        #### СТАРЫЙ ФРАГМЕНТ, НУЖНО МЕНЯТЬ
        checker.check_for_videos(sender, path_to_json, subject, subject_data,
                                 subject_section_data, sessions_list)
        #### СТАРЫЙ ФРАГМЕНТ, НУЖНО МЕНЯТЬ
        interval = checker_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                before_end_operations(sender)
                return


def run_photo_comments_checker(subject_name, subject_data, end_flag):
    checker_name = "photo_comments_checker_settings"
    operation_name = "new photo comment checking"
    sender = subject_name + "'s " + operation_name
    PATH = data_manager.read_path()
    subject_path = PATH + subject_data["path"]
    checker_data = data_manager.read_json(subject_path, checker_name)
    if checker_data["check_flag"] != 1:
        return
    while True:
        before_start_operations(sender)
        if end_flag.isSet():
            before_end_operations(sender)
            return
        #### СТАРЫЙ ФРАГМЕНТ, НУЖНО МЕНЯТЬ
        checker.check_for_comments_photo(sender, path_to_json, subject,
                                         subject_data,
                                         subject_section_data, sessions_list)
        #### СТАРЫЙ ФРАГМЕНТ, НУЖНО МЕНЯТЬ
        interval = checker_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                before_end_operations(sender)
                return


def run_video_comments_checker(subject_name, subject_data, end_flag):
    checker_name = "video_comments_checker_settings"
    operation_name = "new video comment checking"
    sender = subject_name + "'s " + operation_name
    PATH = data_manager.read_path()
    subject_path = PATH + subject_data["path"]
    checker_data = data_manager.read_json(subject_path, checker_name)
    if checker_data["check_flag"] != 1:
        return
    while True:
        before_start_operations(sender)
        if end_flag.isSet():
            before_end_operations(sender)
            return
        #### СТАРЫЙ ФРАГМЕНТ, НУЖНО МЕНЯТЬ
        checker.check_for_comments_video(sender, path_to_json, subject,
                                         subject_data,
                                         subject_section_data, sessions_list)
        #### СТАРЫЙ ФРАГМЕНТ, НУЖНО МЕНЯТЬ
        interval = checker_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                before_end_operations(sender)
                return


def run_topic_comments_checker(subject_name, subject_data, end_flag):
    checker_name = "topic_checker_settings"
    operation_name = "new topic comment checking"
    sender = subject_name + "'s " + operation_name
    PATH = data_manager.read_path()
    subject_path = PATH + subject_data["path"]
    checker_data = data_manager.read_json(subject_path, checker_name)
    if checker_data["check_flag"] != 1:
        return
    while True:
        before_start_operations(sender)
        if end_flag.isSet():
            before_end_operations(sender)
            return
        #### СТАРЫЙ ФРАГМЕНТ, НУЖНО МЕНЯТЬ
        checker.check_for_topics(sender, path_to_json, subject, subject_data,
                                 subject_section_data, sessions_list)
        #### СТАРЫЙ ФРАГМЕНТ, НУЖНО МЕНЯТЬ
        interval = checker_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                before_end_operations(sender)
                return


def run_post_comments_checker(subject_name, subject_data, end_flag):
    checker_name = "post_comments_checker_settings"
    operation_name = "new post comment checking"
    sender = subject_name + "'s " + operation_name
    PATH = data_manager.read_path()
    subject_path = PATH + subject_data["path"]
    checker_data = data_manager.read_json(subject_path, checker_name)
    if checker_data["check_flag"] != 1:
        return
    while True:
        before_start_operations(sender)
        if end_flag.isSet():
            before_end_operations(sender)
            return
        #### СТАРЫЙ ФРАГМЕНТ, НУЖНО МЕНЯТЬ
        checker.check_for_comments_post(sender, path_to_json, subject,
                                        subject_data,
                                        subject_section_data, sessions_list)
        #### СТАРЫЙ ФРАГМЕНТ, НУЖНО МЕНЯТЬ
        interval = checker_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                before_end_operations(sender)
                return
