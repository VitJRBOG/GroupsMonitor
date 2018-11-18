# coding: utf8


import controller.tools as tools
import model.datamanager as datamanager
import model.logger as logger
import model.backup as backup


def changer_workplace(objTButton, objVBoxWorkplace, objButtonBox):
    if objTButton.get_active():
        objVBoxWorkplace.show()
        for item in objButtonBox:
            if item is not objTButton:
                item.set_active(False)
    else:
        objVBoxWorkplace.hide()


def changer_public_workplace(objTButton, objHBoxPublicWorkplace,
                             objHButtonBoxPublics):
    if objTButton.get_active():
        objHBoxPublicWorkplace.show()
        for item in objHButtonBoxPublics:
            if item is not objTButton:
                item.set_active(False)
    else:
        objHBoxPublicWorkplace.hide()


def settings_update(objButton, objTButtonOperation,
                    widgets_blocs, subject_name, settings_name):
    PATH = datamanager.read_path()
    path_to_json = PATH + "bot_notificator/" + subject_name + "/"
    subject_section_data =\
        datamanager.read_json(path_to_json, settings_name)
    if "check_flag" in widgets_blocs:
        if widgets_blocs["check_flag"].get_active():
            subject_section_data["check_flag"] = 1
            objTButtonOperation.set_sensitive(True)
        else:
            subject_section_data["check_flag"] = 0
            if objTButtonOperation.get_active():
                objTButtonOperation.set_active(False)
            objTButtonOperation.set_sensitive(False)
    if "interval" in widgets_blocs:
        subject_section_data["interval"] =\
            int(widgets_blocs["interval"].get_text())
    if "send_to" in widgets_blocs:
        subject_section_data["send_to"] =\
            int(widgets_blocs["send_to"].get_text())
    if "posts_count" in widgets_blocs:
        subject_section_data["posts_count"] =\
            int(widgets_blocs["posts_count"].get_text())
    if "videos_count" in widgets_blocs:
        subject_section_data["videos_count"] =\
            int(widgets_blocs["videos_count"].get_text())
    if "comment_count" in widgets_blocs:
        subject_section_data["comment_count"] =\
            int(widgets_blocs["comment_count"].get_text())
    # !!!! ВКЛЮЧАТЬ ТОЛЬКО ДЛЯ ТЕСТОВ !!!!
    # if "last_date" in widgets_blocs:
    #     subject_section_data["last_date"] =\
    #         widgets_blocs["last_date"].get_text()
    # !!!! ВКЛЮЧАТЬ ТОЛЬКО ДЛЯ ТЕСТОВ !!!!
    if "keywords" in widgets_blocs:
        begin = widgets_blocs["keywords"].get_start_iter()
        end = widgets_blocs["keywords"].get_end_iter()
        keywords = widgets_blocs["keywords"].get_text(begin, end,
                                                      include_hidden_chars=True)
        if len(keywords) > 3:
            keywords = keywords.split(", ")
        subject_section_data["keywords"] = keywords
    #
    # обновление топиков не требуется, т.к. программа обновляет их сама
    # данными из ВК
    #
    # if "topics" in widgets_blocs:
    #     topics = widgets_blocs["topics"].get_text()
    #     topics = topics.split("\n")
    #     for i in range(len(topics)):
    #         subject_section_data["topics"][i]["title"] = topics[i]
    if "post_count" in widgets_blocs:
        subject_section_data["post_count"] =\
            int(widgets_blocs["post_count"].get_text())
    if "photo_count" in widgets_blocs:
        subject_section_data["photo_count"] =\
            int(widgets_blocs["photo_count"].get_text())
    if "video_count" in widgets_blocs:
        subject_section_data["video_count"] =\
            int(widgets_blocs["video_count"].get_text())
    if "filter" in widgets_blocs:
        subject_section_data["filter"] =\
            widgets_blocs["filter"].get_text()
    if "check_by_attachments" in widgets_blocs:
        if widgets_blocs["check_by_attachments"].get_active():
            subject_section_data["check_by_attachments"] = 1
        else:
            subject_section_data["check_by_attachments"] = 0
    if "check_by_communities" in widgets_blocs:
        if widgets_blocs["check_by_communities"].get_active():
            subject_section_data["check_by_communities"] = 1
        else:
            subject_section_data["check_by_communities"] = 0
    if "check_by_keywords" in widgets_blocs:
        if widgets_blocs["check_by_keywords"].get_active():
            subject_section_data["check_by_keywords"] = 1
        else:
            subject_section_data["check_by_keywords"] = 0
    if "check_by_charchange" in widgets_blocs:
        if widgets_blocs["check_by_charchange"].get_active():
            subject_section_data["check_by_charchange"] = 1
        else:
            subject_section_data["check_by_charchange"] = 0
    datamanager.write_json(path_to_json, settings_name, subject_section_data)


def setting_textfields_update(objButton, objTButtonOperation,
                              widgets_blocs, subject_name, settings_name):
    PATH = datamanager.read_path()
    path_to_json = PATH + "bot_notificator/" + subject_name + "/"
    subject_section_data =\
        datamanager.read_json(path_to_json, settings_name)
    if "check_flag" in widgets_blocs:
        if subject_section_data["check_flag"] == 1:
            widgets_blocs["check_flag"].set_active(True)
            objTButtonOperation.set_sensitive(True)
        elif subject_section_data["check_flag"] == 0:
            if objTButtonOperation.get_active():
                objTButtonOperation.set_active(False)
            objTButtonOperation.set_sensitive(False)
            widgets_blocs["check_flag"].set_active(False)
    if "interval" in widgets_blocs:
        widgets_blocs["interval"].set_text(str(subject_section_data["interval"]))
    if "send_to" in widgets_blocs:
        widgets_blocs["send_to"].set_text(str(subject_section_data["send_to"]))
    if "posts_count" in widgets_blocs:
        widgets_blocs["posts_count"].set_text(str(subject_section_data["posts_count"]))
    if "videos_count" in widgets_blocs:
        widgets_blocs["videos_count"].set_text(str(subject_section_data["videos_count"]))
    if "comment_count" in widgets_blocs:
        widgets_blocs["comment_count"].set_text(str(subject_section_data["comment_count"]))
    if "last_date" in widgets_blocs:
        widgets_blocs["last_date"].set_text(subject_section_data["last_date"])
    if "keywords" in widgets_blocs:
        keywords = ""
        if len(subject_section_data["keywords"]) > 0:
            for i in range(len(subject_section_data["keywords"])):
                keyword = subject_section_data["keywords"][i]
                keywords += keyword
                if i < len(subject_section_data["keywords"]) - 1:
                    keywords += ", "
        widgets_blocs["keywords"].set_text(keywords)
    if "topics" in widgets_blocs:
        topics = ""
        for i in range(len(subject_section_data["topics"])):
            topic = subject_section_data["topics"][i]["title"]
            topics += topic
            if i < len(subject_section_data["topics"]) - 1:
                topics += "\n"
        widgets_blocs["topics"].set_text(topics)
    if "post_count" in widgets_blocs:
        widgets_blocs["post_count"].set_text(str(subject_section_data["post_count"]))
    if "photo_count" in widgets_blocs:
        widgets_blocs["photo_count"].set_text(str(subject_section_data["photo_count"]))
    if "video_count" in widgets_blocs:
        widgets_blocs["video_count"].set_text(str(subject_section_data["video_count"]))
    if "filter" in widgets_blocs:
        widgets_blocs["filter"].set_text(subject_section_data["filter"])
    if "check_by_attachments" in widgets_blocs:
        if subject_section_data["check_by_attachments"] == 1:
            widgets_blocs["check_by_attachments"].set_active(True)
        elif subject_section_data["check_by_attachments"] == 0:
            widgets_blocs["check_by_attachments"].set_active(False)
    if "check_by_communities" in widgets_blocs:
        if subject_section_data["check_by_communities"] == 1:
            widgets_blocs["check_by_communities"].set_active(True)
        elif subject_section_data["check_by_communities"] == 0:
            widgets_blocs["check_by_communities"].set_active(False)
    if "check_by_keywords" in widgets_blocs:
        if subject_section_data["check_by_keywords"] == 1:
            widgets_blocs["check_by_keywords"].set_active(True)
        elif subject_section_data["check_by_keywords"] == 0:
            widgets_blocs["check_by_keywords"].set_active(False)
    if "check_by_charchange" in widgets_blocs:
        if subject_section_data["check_by_charchange"] == 1:
            widgets_blocs["check_by_charchange"].set_active(True)
        elif subject_section_data["check_by_charchange"] == 0:
            widgets_blocs["check_by_charchange"].set_active(False)


def switch_checker(objTButton, threads_list, thread_name, out_flag,
                   status_list, sessions_list, subject):
    objThread = threads_list[thread_name]
    if objTButton.get_active():
        if not sessions_list["admin_session"] is None and\
          not sessions_list["bot_session"] is None:
            if not objThread.is_alive():
                out_flag.clear()
                objThread.start()
            else:
                print("Thread already exist")
        else:
            sender = "Switch checker"
            message = "Sessions is not created."
            logger.message_output(sender, message)
            objTButton.set_active(False)
    else:
        out_flag.set()
        threads_list[thread_name] = tools.thread_creator(thread_name,
                                                         status_list,
                                                         out_flag,
                                                         sessions_list,
                                                         subject)


def update_tokens(objButton, list_entry):
    PATH = datamanager.read_path()
    loads_json = datamanager.read_json(PATH + "bot_notificator/", "data")
    loads_json["admin_token"] = list_entry["admin_token"].get_text()
    loads_json["bot_token"] = list_entry["bot_token"].get_text()
    datamanager.write_json(PATH + "bot_notificator/", "data", loads_json)


def update_entry_tokens(objButton, list_entry):
    list_tokens = tools.load_tokens()

    list_entry["admin_token"].set_text(list_tokens["admin_token"])
    list_entry["bot_token"].set_text(list_tokens["bot_token"])


def authorization(objButton, sessions_list):
    tools.thread_authorization(sessions_list)


def save_data_to_cloud(objButton, sessions_list, subject):
    admin_session = sessions_list["admin_session"]
    backup.save_backup(admin_session, subject)


def load_data_from_cloud(objButton, sessions_list, subject):
    admin_session = sessions_list["admin_session"]
    backup.load_backup(admin_session, subject)
