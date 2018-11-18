# coding: utf8


import model.datamanager as datamanager
import model.vk_authorization as vk_authorization
import checker_run
import view_update
import threading


def load_public_list():
    PATH = datamanager.read_path()
    subjects = (datamanager.read_json(PATH + "bot_notificator/", "data"))["subjects"]
    return subjects


def load_subject_data(subject, section_name):
    PATH = datamanager.read_path()
    path_to_json = subject["path"]
    if len(path_to_json) > 0 and path_to_json[len(path_to_json) - 1] != "/":
        path_to_json = PATH + "bot_notificator/" + path_to_json + "/"
    else:
        path_to_json = PATH + "bot_notificator/" + path_to_json
    subject_section_data = datamanager.read_json(path_to_json, section_name)
    return subject_section_data


def load_tokens():
    PATH = datamanager.read_path()
    loads_json = datamanager.read_json(PATH + "bot_notificator/", "data")
    list_tokens = {
        "admin_token": loads_json["admin_token"],
        "bot_token": loads_json["bot_token"]
    }
    return list_tokens


def create_sessions_list():
    sessions_list = {
        "admin_session": None,
        "bot_session": None
    }
    return sessions_list


def thread_logger_update(objLabelLogger):
    objThread = threading.Thread(target=view_update.logger_changer,
                                 args=(objLabelLogger,))
    objThread.daemon = True
    objThread.start()


def thread_authorization(sessions_list):
    objThread = threading.Thread(target=vk_authorization.update_sessions_list,
                                 args=(sessions_list,))
    objThread.daemon = True
    objThread.start()


def out_flag_creator():
    out_flag = threading.Event()

    return out_flag


def thread_creator(thread_name, status_list,
                   out_flag, sessions_list, subject):
    def tbutton_posts():
        objThread =\
            threading.Thread(target=checker_run.run_post_checker,
                             args=(out_flag,
                                   status_list,
                                   sessions_list,
                                   subject,))
        objThread.daemon = True
        return objThread

    def tbutton_album_photos():
        objThread =\
            threading.Thread(target=checker_run.run_album_photo_checker,
                             args=(out_flag,
                                   status_list,
                                   sessions_list,
                                   subject,))
        objThread.daemon = True
        return objThread

    def tbutton_videos():
        objThread =\
            threading.Thread(target=checker_run.run_video_checker,
                             args=(out_flag,
                                   status_list,
                                   sessions_list,
                                   subject,))
        objThread.daemon = True
        return objThread

    def tbutton_photo_comments():
        objThread =\
            threading.Thread(target=checker_run.run_photo_comments_checker,
                             args=(out_flag,
                                   status_list,
                                   sessions_list,
                                   subject,))
        objThread.daemon = True
        return objThread

    def tbutton_video_comments():
        objThread =\
            threading.Thread(target=checker_run.run_video_comments_checker,
                             args=(out_flag,
                                   status_list,
                                   sessions_list,
                                   subject,))
        objThread.daemon = True
        return objThread

    def tbutton_topic_comments():
        objThread =\
            threading.Thread(target=checker_run.run_topic_comments_checker,
                             args=(out_flag,
                                   status_list,
                                   sessions_list,
                                   subject,))
        objThread.daemon = True
        return objThread

    def tbutton_post_comments():
        objThread =\
            threading.Thread(target=checker_run.run_post_comments_checker,
                             args=(out_flag,
                                   status_list,
                                   sessions_list,
                                   subject,))
        objThread.daemon = True
        return objThread

    if thread_name == "Posts":
        objThread = tbutton_posts()

    if thread_name == "Album photos":
        objThread = tbutton_album_photos()

    if thread_name == "Videos":
        objThread = tbutton_videos()

    if thread_name == "Photo comments":
        objThread = tbutton_photo_comments()

    if thread_name == "Video comments":
        objThread = tbutton_video_comments()

    if thread_name == "Topic comments":
        objThread = tbutton_topic_comments()

    if thread_name == "Post comments":
        objThread = tbutton_post_comments()

    return objThread
