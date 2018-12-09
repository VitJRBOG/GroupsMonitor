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


def read_res_files(subject_data, res_filename):
    u"""Читает ресурсные файлы проверяльщика и возвращает словарь с ними,"""
    PATH = data_manager.read_path()
    subject_path = PATH + subject_data["path"]
    monitor_data = data_manager.read_json(subject_path, res_filename)
    return monitor_data


def run_post_monitor(subject_name, subject_data, thread_data):
    res_filename = thread_data["res_filename"]
    sender = thread_data["sender"]
    end_flag = thread_data["end_flag"]
    monitor_data = read_res_files(subject_data, res_filename)
    if monitor_data["need_monitoring"] != 1:
        return
    while True:
        before_start_operations(sender)
        if end_flag.isSet():
            before_end_operations(sender)
            return
        #### НАБРОСОК
        monitoring_algorithms.check_for_posts(
            sender, res_filename, subject_data, monitor_data)
        #### НАБРОСОК
        interval = monitor_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                before_end_operations(sender)
                return


def run_album_photo_monitor(subject_name, subject_data, thread_data):
    res_filename = thread_data["res_filename"]
    sender = thread_data["sender"]
    end_flag = thread_data["end_flag"]
    monitor_data = read_res_files(subject_data, res_filename)
    if monitor_data["need_monitoring"] != 1:
        return
    while True:
        before_start_operations(sender)
        if end_flag.isSet():
            before_end_operations(sender)
            return
        #### НАБРОСОК
        monitoring_algorithms.check_for_albums(
            sender, res_filename, subject_data, monitor_data)
        #### НАБРОСОК
        interval = monitor_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                before_end_operations(sender)
                return


def run_video_monitor(subject_name, subject_data, thread_data):
    res_filename = thread_data["res_filename"]
    sender = thread_data["sender"]
    end_flag = thread_data["end_flag"]
    monitor_data = read_res_files(subject_data, res_filename)
    if monitor_data["need_monitoring"] != 1:
        return
    while True:
        before_start_operations(sender)
        if end_flag.isSet():
            before_end_operations(sender)
            return
        #### НАБРОСОК
        monitoring_algorithms.check_for_videos(
            sender, res_filename, subject_data, monitor_data)
        #### НАБРОСОК
        interval = monitor_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                before_end_operations(sender)
                return


def run_photo_comments_monitor(subject_name, subject_data, thread_data):
    res_filename = thread_data["res_filename"]
    sender = thread_data["sender"]
    end_flag = thread_data["end_flag"]
    monitor_data = read_res_files(subject_data, res_filename)
    if monitor_data["need_monitoring"] != 1:
        return
    while True:
        before_start_operations(sender)
        if end_flag.isSet():
            before_end_operations(sender)
            return
        #### НАБРОСОК
        monitoring_algorithms.check_for_comments_photo(
            sender, res_filename, subject_data, monitor_data)
        #### НАБРОСОК
        interval = monitor_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                before_end_operations(sender)
                return


def run_video_comments_monitor(subject_name, subject_data, thread_data):
    res_filename = thread_data["res_filename"]
    sender = thread_data["sender"]
    end_flag = thread_data["end_flag"]
    monitor_data = read_res_files(subject_data, res_filename)
    if monitor_data["need_monitoring"] != 1:
        return
    while True:
        before_start_operations(sender)
        if end_flag.isSet():
            before_end_operations(sender)
            return
        #### НАБРОСОК
        monitoring_algorithms.check_for_comments_video(
            sender, res_filename, subject_data, monitor_data)
        #### НАБРОСОК
        interval = monitor_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                before_end_operations(sender)
                return


def run_topic_comments_monitor(subject_name, subject_data, thread_data):
    res_filename = thread_data["res_filename"]
    sender = thread_data["sender"]
    end_flag = thread_data["end_flag"]
    monitor_data = read_res_files(subject_data, res_filename)
    if monitor_data["need_monitoring"] != 1:
        return
    while True:
        before_start_operations(sender)
        if end_flag.isSet():
            before_end_operations(sender)
            return
        #### НАБРОСОК
        monitoring_algorithms.check_for_topics(
            sender, res_filename, subject_data, monitor_data)
        #### НАБРОСОК
        interval = monitor_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                before_end_operations(sender)
                return


def run_post_comments_monitor(subject_name, subject_data, thread_data):
    res_filename = thread_data["res_filename"]
    sender = thread_data["sender"]
    end_flag = thread_data["end_flag"]
    monitor_data = read_res_files(subject_data, res_filename)
    if monitor_data["need_monitoring"] != 1:
        return
    while True:
        before_start_operations(sender)
        if end_flag.isSet():
            before_end_operations(sender)
            return
        #### НАБРОСОК
        monitoring_algorithms.check_for_comments_post(
            sender, res_filename, subject_data, monitor_data)
        #### НАБРОСОК
        interval = monitor_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                before_end_operations(sender)
                return
