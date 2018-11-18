# coding: utf8


import model.checker as checker
import model.vk_authorization as vk_authorization
import controller.tools as tools
import controller.checker_run as checker_run
import threading
import time


def new_post_checking_thread(sender, path_to_json, subject,
                             subject_data, subject_section_data,
                             sessions_list):
    while True:
        checker.check_for_posts(sender, path_to_json, subject,
                                subject_data, subject_section_data,
                                sessions_list)
        time.sleep(subject_section_data["interval"])


def new_album_photo_checking_thread(sender, path_to_json, subject,
                                    subject_data, subject_section_data,
                                    sessions_list):
    while True:
        checker.check_for_albums(sender, path_to_json, subject,
                                 subject_data, subject_section_data,
                                 sessions_list)
        time.sleep(subject_section_data["interval"])


def new_video_checking_thread(sender, path_to_json, subject, subject_data,
                              subject_section_data,
                              sessions_list):
    while True:
        checker.check_for_videos(sender, path_to_json, subject,
                                 subject_data, subject_section_data,
                                 sessions_list)
        time.sleep(subject_section_data["interval"])


def new_photo_comment_checking_thread(sender, path_to_json, subject,
                                      subject_data, subject_section_data,
                                      sessions_list):
    while True:
        checker.check_for_comments_photo(sender, path_to_json, subject,
                                         subject_data, subject_section_data,
                                         sessions_list)
        time.sleep(subject_section_data["interval"])


def new_video_comment_checking_thread(sender, path_to_json, subject,
                                      subject_data, subject_section_data,
                                      sessions_list):
    while True:
        checker.check_for_comments_video(sender, path_to_json, subject,
                                         subject_data, subject_section_data,
                                         sessions_list)
        time.sleep(subject_section_data["interval"])


def new_topic_comment_checking_thread(sender, path_to_json, subject,
                                      subject_data, subject_section_data,
                                      sessions_list):
    while True:
        checker.check_for_topics(sender, path_to_json, subject,
                                 subject_data, subject_section_data,
                                 sessions_list)
        time.sleep(subject_section_data["interval"])


def new_post_comment_checking_thread(sender, path_to_json, subject,
                                     subject_data, subject_section_data,
                                     sessions_list):
    while True:
        checker.check_for_comments_post(sender, path_to_json, subject,
                                        subject_data, subject_section_data,
                                        sessions_list)
        time.sleep(subject_section_data["interval"])


def start_bot():
    def authorization_bot():
        sessions_list = tools.create_sessions_list()
        vk_authorization.update_sessions_list(sessions_list)

        return sessions_list

    def run_publics_checker(sessions_list):
        subjects = tools.load_public_list()

        operations_names = [
            "new post checking", "new album photo checking",
            "new video checking", "new photo comment checking",
            "new video comment checking", "new topic comment checking",
            "new post comment checking"
        ]

        checkers_names = [
            "post_checker_settings", "photo_checker_settings",
            "video_checker_settings", "photo_comments_checker_settings",
            "video_comments_checker_settings", "topic_checker_settings",
            "post_comments_checker_settings"
        ]

        threads = []

        def ref_thread(operation_name, arg):
            if operation_name == "new post checking":
                objThread = threading.Thread(target=new_post_checking_thread,
                                             args=(arg["sender"],
                                                   arg["path_to_json"],
                                                   arg["subject"],
                                                   arg["subject_data"],
                                                   arg["subject_section_data"],
                                                   arg["sessions_list"],))
                objThread.daemon = True
                return objThread

            elif operation_name == "new album photo checking":
                objThread = threading.Thread(target=new_album_photo_checking_thread,
                                             args=(arg["sender"],
                                                   arg["path_to_json"],
                                                   arg["subject"],
                                                   arg["subject_data"],
                                                   arg["subject_section_data"],
                                                   arg["sessions_list"],))
                objThread.daemon = True
                return objThread

            elif operation_name == "new video checking":
                objThread = threading.Thread(target=new_video_checking_thread,
                                             args=(arg["sender"],
                                                   arg["path_to_json"],
                                                   arg["subject"],
                                                   arg["subject_data"],
                                                   arg["subject_section_data"],
                                                   arg["sessions_list"],))
                objThread.daemon = True
                return objThread

            elif operation_name == "new photo comment checking":
                objThread = threading.Thread(target=new_photo_comment_checking_thread,
                                             args=(arg["sender"],
                                                   arg["path_to_json"],
                                                   arg["subject"],
                                                   arg["subject_data"],
                                                   arg["subject_section_data"],
                                                   arg["sessions_list"],))
                objThread.daemon = True
                return objThread

            elif operation_name == "new video comment checking":
                objThread = threading.Thread(target=new_video_comment_checking_thread,
                                             args=(arg["sender"],
                                                   arg["path_to_json"],
                                                   arg["subject"],
                                                   arg["subject_data"],
                                                   arg["subject_section_data"],
                                                   arg["sessions_list"],))
                objThread.daemon = True
                return objThread

            elif operation_name == "new topic comment checking":
                objThread = threading.Thread(target=new_topic_comment_checking_thread,
                                             args=(arg["sender"],
                                                   arg["path_to_json"],
                                                   arg["subject"],
                                                   arg["subject_data"],
                                                   arg["subject_section_data"],
                                                   arg["sessions_list"],))
                objThread.daemon = True
                return objThread

            elif operation_name == "new post comment checking":
                objThread = threading.Thread(target=new_post_comment_checking_thread,
                                             args=(arg["sender"],
                                                   arg["path_to_json"],
                                                   arg["subject"],
                                                   arg["subject_data"],
                                                   arg["subject_section_data"],
                                                   arg["sessions_list"],))
                objThread.daemon = True
                return objThread

        for subject in subjects:
            for i in range(len(operations_names)):
                sender = subject["name"] + " " + operations_names[i]

                PATH, path_to_json, subject_data, subject_section_data = \
                    checker_run.load_data_for_checker(checkers_names[i],
                                                      subject,
                                                      operations_names[i])
                if subject_section_data["check_flag"] == 1:
                    arg = {
                        "sender": sender,
                        "path_to_json": path_to_json,
                        "subject": subject,
                        "subject_data": subject_data,
                        "subject_section_data": subject_section_data,
                        "sessions_list": sessions_list
                    }
                    objThread = ref_thread(operations_names[i], arg)
                    threads.append(objThread)

        for objThread in threads:
            objThread.start()

    sessions_list = authorization_bot()
    run_publics_checker(sessions_list)

    while True:
        user_asnwer = raw_input()
        if user_asnwer == "quit":
            exit(0)
