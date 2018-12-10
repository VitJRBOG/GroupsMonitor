# coding: utf8
u"""Модуль работы с потоками."""


import threading
import time
import data_manager
import output_data
import monitoring_executor


def run_thread_starter(dict_sessions):
    u"""Запуск функций старта потоков."""
    def run_threads(data_threads):
        u"""Запускает потоки."""
        for data_thread in data_threads:
            data_thread["thread"].start()
    def thread_checker(data_threads):
        u"""Запускает поток проверки статуса потоков."""
        def thread_checker_algorithm(data_threads):
            u"""Алгоритм проверки статуса потоков."""
            while True:
                for data_thread in data_threads:
                    if not data_thread["thread"].isAlive() and\
                       not data_thread["end_flag"].isSet():
                        sender = data_thread["sender"]
                        message = "WARNING! Monitoring is stopped..."
                        output_data.output_text_row(sender, message)
                time.sleep(30)
        objThread = threading.Thread(target=thread_checker_algorithm, args=(data_threads,))
        objThread.daemon = True
        objThread.start()
    data_threads = thread_creator(dict_sessions)
    run_threads(data_threads)
    thread_checker(data_threads)
    return data_threads


def thread_creator(dict_sessions):
    u"""Создает алгоритмы потоков."""
    def select_subjects_data():
        u"""Получение данных о проверяемых субъектах."""
        PATH = data_manager.read_path()
        dict_data = data_manager.read_json(PATH, "data")
        subjects_names = []
        subjects_path = {}
        for subject in dict_data["subjects"]:
            if subject["monitor_subject"] == 1:
                subjects_names.append(subject["name"])
                subjects_path.update({subject["name"]: subject["path"]})
        return subjects_names, subjects_path
    
    def subjects_data_dict_creator(subjects_names, dict_sessions, subjects_path):
        u"""Собирает словарь с данными о проверяемых субъектах."""
        subjects_data = {}
        for subject_name in subjects_names:
            PATH = data_manager.read_path()
            subject_path = PATH + subjects_path[subject_name]
            external_subject_data = data_manager.read_json(
                subject_path, "subject_data")
            values = {
                "sender_session": dict_sessions[subject_name],
                "admin_session": dict_sessions["Admin"],
                "path": subjects_path[subject_name],
                "owner_id": external_subject_data["owner_id"]
            }
            subjects_data.update({subject_name: values})

        return subjects_data
    
    def preparation_thread(subject_name, subject_data):
        u"""Подготовка потоков."""
        def thread_post_monitoring():
            u"""Создание потока проверки постов."""
            end_flag = threading.Event()
            res_filename = "post_monitor_settings"
            operation_name = "post monitoring"
            sender = subject_name + "'s " + operation_name
            thread_data = {
                "res_filename": res_filename,
                "sender": sender,
                "end_flag": end_flag
            }
            objThread =\
                threading.Thread(target=monitoring_executor.run_wall_posts_monitor,
                                 args=(subject_name, subject_data, thread_data,))
            objThread.daemon = True
            thread_data.update({"thread": objThread})
            return thread_data

        def thread_album_photo_monitoring():
            u"""Создание потока проверки фотографий."""
            end_flag = threading.Event()
            res_filename = "photo_monitor_settings"
            operation_name = "album photo monitoring"
            sender = subject_name + "'s " + operation_name
            thread_data = {
                "res_filename": res_filename,
                "sender": sender,
                "end_flag": end_flag
            }
            objThread =\
                threading.Thread(target=monitoring_executor.run_album_photos_monitor,
                                 args=(subject_name, subject_data, thread_data,))
            objThread.daemon = True
            thread_data.update({"thread": objThread})
            return thread_data

        def thread_video_monitoring():
            u"""Создание потока проверки видеозаписей."""
            end_flag = threading.Event()
            res_filename = "video_monitor_settings"
            operation_name = "video monitoring"
            sender = subject_name + "'s " + operation_name
            thread_data = {
                "res_filename": res_filename,
                "sender": sender,
                "end_flag": end_flag
            }
            objThread =\
                threading.Thread(target=monitoring_executor.run_videos_monitor,
                                 args=(subject_name, subject_data, thread_data,))
            objThread.daemon = True
            thread_data.update({"thread": objThread})
            return thread_data

        def thread_photo_comments_monitoring():
            u"""Создание потока проверки комментов под фотками."""
            end_flag = threading.Event()
            res_filename = "photo_comments_monitor_settings"
            operation_name = "photo comment monitoring"
            sender = subject_name + "'s " + operation_name
            thread_data = {
                "res_filename": res_filename,
                "sender": sender,
                "end_flag": end_flag
            }
            objThread =\
                threading.Thread(target=monitoring_executor.run_photo_comments_monitor,
                                 args=(subject_name, subject_data, thread_data,))
            objThread.daemon = True
            thread_data.update({"thread": objThread})
            return thread_data

        def thread_video_comments_monitoring():
            u"""Создание потока проверки комментов под видео."""
            end_flag = threading.Event()
            res_filename = "video_comments_monitor_settings"
            operation_name = "video comment monitoring"
            sender = subject_name + "'s " + operation_name
            thread_data = {
                "res_filename": res_filename,
                "sender": sender,
                "end_flag": end_flag
            }
            objThread =\
                threading.Thread(target=monitoring_executor.run_video_comments_monitor,
                                 args=(subject_name, subject_data, thread_data,))
            objThread.daemon = True
            thread_data.update({"thread": objThread})
            return thread_data

        def thread_topic_comments_monitoring():
            u"""Создание потока проверки комментов в обсуждениях."""
            end_flag = threading.Event()
            res_filename = "topic_monitor_settings"
            operation_name = "topic comment monitoring"
            sender = subject_name + "'s " + operation_name
            thread_data = {
                "res_filename": res_filename,
                "sender": sender,
                "end_flag": end_flag
            }
            objThread =\
                threading.Thread(target=monitoring_executor.run_topic_comments_monitor,
                                 args=(subject_name, subject_data, thread_data,))
            objThread.daemon = True
            thread_data.update({"thread": objThread})
            return thread_data

        def thread_post_comments_monitoring():
            u"""Создание потока проверки комментов под постами."""
            end_flag = threading.Event()
            res_filename = "post_comments_monitor_settings"
            operation_name = "post comment monitoring"
            sender = subject_name + "'s " + operation_name
            thread_data = {
                "res_filename": res_filename,
                "sender": sender,
                "end_flag": end_flag
            }
            objThread =\
                threading.Thread(target=monitoring_executor.run_wall_post_comments_monitor,
                                 args=(subject_name, subject_data, thread_data,))
            objThread.daemon = True
            thread_data.update({"thread": objThread})
            return thread_data

        threads_list = []
        threads_list.append(thread_post_monitoring)
        threads_list.append(thread_album_photo_monitoring)
        threads_list.append(thread_video_monitoring)
        threads_list.append(thread_photo_comments_monitoring)
        threads_list.append(thread_video_comments_monitoring)
        threads_list.append(thread_topic_comments_monitoring)
        threads_list.append(thread_post_comments_monitoring)

        return threads_list
    
    subjects_names, subjects_path = select_subjects_data()
    subjects_data = subjects_data_dict_creator(
        subjects_names, dict_sessions, subjects_path)
    data_threads = []
    for subject_name in subjects_names:
        data_threads.extend(preparation_thread(
            subject_name, subjects_data[subject_name]))
    return data_threads
