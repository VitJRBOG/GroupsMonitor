# coding: utf8
u"""Модуль работы с потоками."""


import threading
import time
import data_manager
import output_data
import checker_algorithms


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
                        message = "WARNING! Operation is stopped..."
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
            if subject["check_subject"]:
                subjects_names.append(subject["name"])
                subjects_path.update({subject["name"]: subject["path"]})
        return subjects_names, subjects_path
    
    def subjects_data_dict_creator(subjects_names, dict_sessions, subjects_path):
        u"""Собирает словарь с данными о проверяемых субъектах."""
        subjects_data = {}
        for subject_name in subjects_names:
            values = {
                "sender_session": dict_sessions[subject_name],
                "admin_session": dict_sessions["Admin"],
                "path": subjects_path[subject_name]
            }
            subjects_data.update({subject_name: values})

        return subjects_data
    
    def preparation_thread(subject_name, subject_data):
        u"""Подготовка потоков."""
        def thread_post_checker():
            u"""Создание потока проверки новых постов."""
            end_flag = threading.Event()
            objThread =\
                threading.Thread(target=checker_algorithms.run_post_checker,
                                 args=(subject_name, subject_data, end_flag,))
            objThread.daemon = True
            operation_name = "new post checking"
            sender = subject_name + "'s " + operation_name
            thread_data = {
                "thread": objThread,
                "sender": sender,
                "end_flag": end_flag
            }
            return thread_data

        def thread_album_photo_checker():
            u"""Создание потока проверки новых фотографий."""
            end_flag = threading.Event()
            objThread =\
                threading.Thread(target=checker_algorithms.run_album_photo_checker,
                                 args=(subject_name, subject_data, end_flag,))
            objThread.daemon = True
            operation_name = "new album photo checking"
            sender = subject_name + "'s " + operation_name
            thread_data = {
                "thread": objThread,
                "sender": sender,
                "end_flag": end_flag
            }
            return thread_data

        def thread_video_checker():
            u"""Создание потока проверки новых видеозаписей."""
            end_flag = threading.Event()
            objThread =\
                threading.Thread(target=checker_algorithms.run_video_checker,
                                 args=(subject_name, subject_data, end_flag,))
            objThread.daemon = True
            operation_name = "new video checking"
            sender = subject_name + "'s " + operation_name
            thread_data = {
                "thread": objThread,
                "sender": sender,
                "end_flag": end_flag
            }
            return thread_data

        def thread_photo_comments_checker():
            u"""Создание потока проверки новых комментов под фотками."""
            end_flag = threading.Event()
            objThread =\
                threading.Thread(target=checker_algorithms.run_photo_comments_checker,
                                 args=(subject_name, subject_data, end_flag,))
            objThread.daemon = True
            operation_name = "new photo comment checking"
            sender = subject_name + "'s " + operation_name
            thread_data = {
                "thread": objThread,
                "sender": sender,
                "end_flag": end_flag
            }
            return thread_data

        def thread_video_comments_checker():
            u"""Создание потока проверки новых комментов под видео."""
            end_flag = threading.Event()
            objThread =\
                threading.Thread(target=checker_algorithms.run_video_comments_checker,
                                 args=(subject_name, subject_data, end_flag,))
            objThread.daemon = True
            operation_name = "new video comment checking"
            sender = subject_name + "'s " + operation_name
            thread_data = {
                "thread": objThread,
                "sender": sender,
                "end_flag": end_flag
            }
            return thread_data

        def thread_topic_comments_checker():
            u"""Создание потока проверки новых комментов в обсуждениях."""
            end_flag = threading.Event()
            objThread =\
                threading.Thread(target=checker_algorithms.run_topic_comments_checker,
                                 args=(subject_name, subject_data, end_flag,))
            objThread.daemon = True
            operation_name = "new topic comment checking"
            sender = subject_name + "'s " + operation_name
            thread_data = {
                "thread": objThread,
                "sender": sender,
                "end_flag": end_flag
            }
            return thread_data

        def thread_post_comments_checker():
            u"""Создание потока проверки новых комментов под постами."""
            end_flag = threading.Event()
            objThread =\
                threading.Thread(target=checker_algorithms.run_post_comments_checker,
                                 args=(subject_name, subject_data, end_flag,))
            objThread.daemon = True
            operation_name = "new post comment checking"
            sender = subject_name + "'s " + operation_name
            thread_data = {
                "thread": objThread,
                "sender": sender,
                "end_flag": end_flag
            }
            return thread_data

        threads_list = []
        threads_list.append(thread_post_checker)
        threads_list.append(thread_album_photo_checker)
        threads_list.append(thread_video_checker)
        threads_list.append(thread_photo_comments_checker)
        threads_list.append(thread_video_comments_checker)
        threads_list.append(thread_topic_comments_checker)
        threads_list.append(thread_post_comments_checker)

        return threads_list
    
    subjects_names, subjects_path = select_subjects_data()
    subjects_data = subjects_data_dict_creator(
        subjects_names, dict_sessions, subjects_path)
    data_threads = []
    for subject_name in subjects_names:
        data_threads.extend(preparation_thread(
            subject_name, subjects_data[subject_name]))
    return data_threads
