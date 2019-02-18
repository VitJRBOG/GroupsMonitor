# coding: utf8
u"""Модуль запуска алгоритмов проверки."""


import time
import data_manager
import output_data
import monitoring_algorithms.wall_posts_monitor
import monitoring_algorithms.album_photos_monitor
import monitoring_algorithms.videos_monitor
import monitoring_algorithms.photo_comments_monitor
import monitoring_algorithms.video_comments_monitor
import monitoring_algorithms.topic_comments_monitor
import monitoring_algorithms.wall_post_comments_monitor


def before_start_operations(sender):
    u"""Операции перед началом проверки."""
    message = "Started..."
    output_data.output_text_row(sender, message)


def read_res_files(subject_data, res_filename):
    u"""Читает ресурсные файлы проверяльщика и возвращает словарь с ними,"""
    PATH = data_manager.read_path()
    subject_path = PATH + subject_data["path"] + "/"
    monitor_data = data_manager.read_json(subject_path, res_filename)
    return monitor_data


def run_wall_posts_monitor(subject_name, subject_data, thread_data):
    u"""Запускает алгоритмы проверки постов на стене."""
    res_filename = thread_data["res_filename"]
    sender = thread_data["sender"]
    end_flag = thread_data["end_flag"]
    monitor_data = read_res_files(subject_data, res_filename)
    if monitor_data["need_monitoring"] != 1:
        return
    thread_data["was_turned_on"] = True
    before_start_operations(sender)
    while True:
        if end_flag.isSet():
            return
        monitoring_algorithms.wall_posts_monitor.run_monitoring_wall_posts(
            sender, res_filename, subject_data, monitor_data)
        monitor_data = read_res_files(subject_data, res_filename)
        interval = monitor_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                return


def run_album_photos_monitor(subject_name, subject_data, thread_data):
    u"""Запускает алгоритмы проверки фотографий в альбомах."""
    res_filename = thread_data["res_filename"]
    sender = thread_data["sender"]
    end_flag = thread_data["end_flag"]
    monitor_data = read_res_files(subject_data, res_filename)
    if monitor_data["need_monitoring"] != 1:
        return
    thread_data["was_turned_on"] = True
    before_start_operations(sender)
    while True:
        if end_flag.isSet():
            return
        monitoring_algorithms.album_photos_monitor.run_monitoring_album_photos(
            sender, res_filename, subject_data, monitor_data)
        monitor_data = read_res_files(subject_data, res_filename)
        interval = monitor_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                return


def run_videos_monitor(subject_name, subject_data, thread_data):
    u"""Запускает алгоритмы проверки видеороликов."""
    res_filename = thread_data["res_filename"]
    sender = thread_data["sender"]
    end_flag = thread_data["end_flag"]
    monitor_data = read_res_files(subject_data, res_filename)
    if monitor_data["need_monitoring"] != 1:
        return
    thread_data["was_turned_on"] = True
    before_start_operations(sender)
    while True:
        if end_flag.isSet():
            return
        monitoring_algorithms.videos_monitor.run_monitoring_videos(
            sender, res_filename, subject_data, monitor_data)
        monitor_data = read_res_files(subject_data, res_filename)
        interval = monitor_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                return


def run_photo_comments_monitor(subject_name, subject_data, thread_data):
    u"""Запускает алгоритмы проверки комментариев под фотографиями."""
    res_filename = thread_data["res_filename"]
    sender = thread_data["sender"]
    end_flag = thread_data["end_flag"]
    monitor_data = read_res_files(subject_data, res_filename)
    if monitor_data["need_monitoring"] != 1:
        return
    thread_data["was_turned_on"] = True
    before_start_operations(sender)
    while True:
        if end_flag.isSet():
            return
        monitoring_algorithms.photo_comments_monitor.run_monitoring_photo_comments(
            sender, res_filename, subject_data, monitor_data)
        monitor_data = read_res_files(subject_data, res_filename)
        interval = monitor_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                return


def run_video_comments_monitor(subject_name, subject_data, thread_data):
    u"""Запускает алгоритмы проверки комментариев под видеороликами."""
    res_filename = thread_data["res_filename"]
    sender = thread_data["sender"]
    end_flag = thread_data["end_flag"]
    monitor_data = read_res_files(subject_data, res_filename)
    if monitor_data["need_monitoring"] != 1:
        return
    thread_data["was_turned_on"] = True
    before_start_operations(sender)
    while True:
        if end_flag.isSet():
            return
        monitoring_algorithms.video_comments_monitor.run_monitoring_video_comments(
            sender, res_filename, subject_data, monitor_data)
        monitor_data = read_res_files(subject_data, res_filename)
        interval = monitor_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                return


def run_topic_comments_monitor(subject_name, subject_data, thread_data):
    u"""Запускает алгоритмы проверки комментариев в обсуждениях."""
    res_filename = thread_data["res_filename"]
    sender = thread_data["sender"]
    end_flag = thread_data["end_flag"]
    monitor_data = read_res_files(subject_data, res_filename)
    if monitor_data["need_monitoring"] != 1:
        return
    if str(subject_data["owner_id"])[0] != "-":
        return
    thread_data["was_turned_on"] = True
    before_start_operations(sender)
    while True:
        if end_flag.isSet():
            return
        monitoring_algorithms.topic_comments_monitor.run_monitoring_topic_comments(
            sender, res_filename, subject_data, monitor_data)
        monitor_data = read_res_files(subject_data, res_filename)
        interval = monitor_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                return


def run_wall_post_comments_monitor(subject_name, subject_data, thread_data):
    u"""Запускает алгоритмы проверки комментариев под постами на стене."""
    res_filename = thread_data["res_filename"]
    sender = thread_data["sender"]
    end_flag = thread_data["end_flag"]
    monitor_data = read_res_files(subject_data, res_filename)
    if monitor_data["need_monitoring"] != 1:
        return
    thread_data["was_turned_on"] = True
    before_start_operations(sender)
    while True:
        if end_flag.isSet():
            return
        monitoring_algorithms.wall_post_comments_monitor.run_monitoring_wall_post_comments(
            sender, res_filename, subject_data, monitor_data)
        monitor_data = read_res_files(subject_data, res_filename)
        interval = monitor_data["interval"]
        for i in range(interval):
            time.sleep(1)
            if end_flag.isSet():
                return
