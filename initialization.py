# coding: utf8


import os
import output_data
import data_manager


def check_res_files():
    u"""Проверка существования пути."""
    def create_file_path(sender):
        u"""Создание файла для хранение пути."""
        file_text = open("path.txt", "w")
        file_text.write("db/")
        file_text.close()
        message = "Was created file \"path.txt\"."
        output_data.output_text_row(sender, message)

    def create_folder(sender, path, new_folder):
        u"""Создание каталога."""
        os.mkdir(path + new_folder)
        message = "Was created directory \"" + new_folder + "\"."
        output_data.output_text_row(sender, message)

    def check_db_folder(sender, path):
        u"""Проверка содержимого каталога с данными."""
        for tree in os.walk(path):
            tree_path, folders, files = tree
            if len(folders) == 0:
                create_example_subject(sender, path)
                return True
        return False

    sender = "Initialization"

    need_presetting = False

    if os.path.exists("path.txt") is not True:
        create_file_path(sender)

    PATH = data_manager.read_path()

    if os.path.exists(PATH) is not True:
        create_folder(sender, PATH, "")
        need_presetting = True

    if os.path.exists(PATH + "data.json") is not True:
        data = {"subjects": []}
        data_manager.write_json(PATH, "data", data)
        message = "Was created file \"data.json\"."
        output_data.output_text_row(sender, message)
        need_presetting = True

    if need_presetting:
        check_db_folder(sender, PATH)

    return need_presetting


def create_example_subject(sender, path):
    u"""Создание всех файлов субъекта-примера."""
    def make_dicts_settings():
        u"""Описание словарей с настройками и данными."""
        data_dicts = {
            "subject_data": {
                "name": "Example subject",
                "wiki_database_id": "-123_456",
                "total_last_date": "0",
                "owner_id": -123
            },
            "wall_posts_monitor": {
                "check_by_attachments": {
                    "check": 0,
                    "types": [
                        "photo",
                        "video",
                        "audio",
                        "doc",
                        "poll"
                    ]
                },
                "check_by_keywords": {
                    "keywords": [],
                    "check": 0
                },
                "check_by_exclude": {
                    "check": 0,
                    "authors": []
                },
                "check_all": {
                    "check": 1,
                    "check_exclude": 0
                },
                "need_monitoring": 1,
                "interval": 5,
                "send_to": [
                    123
                ],
                "filter": "all",
                "last_date": "0",
                "posts_count": 5
            },
            "album_photos_monitor": {
                "last_date": "0",
                "photo_count": 5,
                "send_to": [
                    123
                ],
                "interval": 5,
                "need_monitoring": 1
            },
            "videos_monitor": {
                "send_to": [
                    123
                ],
                "last_date": "0",
                "interval": 5,
                "video_count": 1,
                "need_monitoring": 1
            },
            "photo_comments_monitor": {
                "last_date": "0",
                "interval": 5,
                "need_monitoring": 1,
                "send_to": [
                    123
                ],
                "comment_count": 5
            },
            "video_comments_monitor": {
                "need_monitoring": 1,
                "video_count": 5,
                "interval": 5,
                "comment_count": 1,
                "send_to": [
                    123
                ],
                "last_date": "0"
            },
            "topic_comments_monitor": {
                "topics_count": 5,
                "need_monitoring": 1,
                "post_count": 5,
                "interval": 5,
                "send_to": [
                    123
                ],
                "last_date": "0"
            },
            "wall_post_comments_monitor": {
                "check_by_card_number": {
                    "check": 0,
                    "digits_count": []
                },
                "check_by_phone_number": {
                    "check": 0,
                    "digits_count": []
                },
                "check_by_profile_data": {
                    "check": 0,
                    "names": [],
                    "ids": []
                },
                "check_by_attachments": {
                    "check": 0,
                    "types": [
                        "photo",
                        "video",
                        "audio",
                        "doc",
                        "poll"
                    ]
                },
                "check_all": {
                    "check": 0,
                    "check_exclude": 0
                },
                "need_monitoring": 1,
                "check_by_exclude": {
                    "check": 0,
                    "authors": []
                },
                "interval": 5,
                "send_to": [
                    123
                ],
                "thread_items_count": 5,
                "filter": "all",
                "comment_count": 5,
                "check_by_communities": {
                    "check": 1
                },
                "last_date": "0",
                "posts_count": 5,
                "check_by_keywords": {
                    "check_charchange": 0,
                    "chars_cyr_to_lat": {},
                    "chars_lat_to_cyr": {},
                    "small_messages": [],
                    "keywords": [],
                    "check": 0
                }
            }
        }
        return data_dicts

    def create_subject_folder(path):
        u"""Создание каталога субъекта-примера."""
        folder_example_subject = "example_subject"
        os.mkdir(path + folder_example_subject)
        message = "Was created directory \"" + folder_example_subject + "\"."
        output_data.output_text_row(sender, message)

    def create_data_files(sender, path, data_dicts):
        u"""Создание файлов с данными субъекта-примера."""
        data_files = data_dicts.keys()
        for data_file in data_files:
            data_manager.write_json(path, data_file, data_dicts[data_file])
            message = "Was created file \"" + data_file + ".json\"."
            output_data.output_text_row(sender, message)

    def add_data_subject_to_data_dict(sender, path):
        u"""Добавление данных о субъекте-примере в главный файл данных."""
        subject = {
            "path": "example_subject",
            "access_tokens": {
                "admin": "token",
                "album_photos_monitor": "token",
                "wall_posts_monitor": "token",
                "topic_comments_monitor": "token",
                "photo_comments_monitor": "token",
                "wall_post_comments_monitor": "token",
                "videos_monitor": "token",
                "video_comments_monitor": "token"
            },
            "interval": 5,
            "monitor_subject": 1,
            "name": "Example_subject"
        }
        data = data_manager.read_json(path, "data")
        data["subjects"].append(subject)
        data_manager.write_json(path, "data", data)
        message = "Data by example_subject was added to \"data.json\"."
        output_data.output_text_row(sender, message)

    path_to_example_folder = path + "example_subject/"
    data_dicts = make_dicts_settings()
    create_subject_folder(path)
    create_data_files(sender, path_to_example_folder, data_dicts)
    add_data_subject_to_data_dict(sender, path)
