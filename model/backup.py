# coding: utf8


import logger
import model.datamanager as datamanager
import time
import os
import json


def read_wiki(vk_admin_session, wiki_full_id):

    wiki_owner_id = int(wiki_full_id[0:wiki_full_id.rfind('_')])
    wiki_id = int(wiki_full_id[wiki_full_id.rfind('_') + 1:])

    values = {
        "owner_id": wiki_owner_id,
        "page_id": wiki_id,
        "need_html": 1
    }

    time.sleep(1)
    response = vk_admin_session.method("pages.get", values)

    text = response["html"][8:]
    data_json = json.loads(text)

    return data_json


def save_wiki(vk_admin_session, wiki_full_id, data_json):
    wiki_owner_id = int(wiki_full_id[1:wiki_full_id.rfind('_')])
    wiki_id = int(wiki_full_id[wiki_full_id.rfind('_') + 1:])

    text = json.dumps(data_json)

    values = {
        "group_id": wiki_owner_id,
        "page_id": wiki_id,
        "text": text
    }

    time.sleep(1)
    vk_admin_session.method("pages.save", values)


def make_json_text(path_to_subject_json):
    subject_data = datamanager.read_json(path_to_subject_json, "subject_data")

    files = os.listdir(path_to_subject_json)
    jsonfiles = filter(lambda x: x.endswith(".json"), files)

    for file_name in jsonfiles:
        file_name = file_name.replace(".json", "")
        if file_name != "subject_data":
            subject_section_data = datamanager.read_json(path_to_subject_json,
                                                         file_name)
            subject_data.update({file_name: subject_section_data})

    return subject_data


def split_json_text(path_to_subject_json, subject_data):
    files = os.listdir(path_to_subject_json)
    jsonfiles = filter(lambda x: x.endswith(".json"), files)

    json_files_list = []

    for file_name in jsonfiles:
        file_name = file_name.replace(".json", "")
        if file_name in subject_data:
            subject_section_data = subject_data[file_name]
            values = {
                "section_name": file_name,
                "json_text": subject_section_data
            }
            json_files_list.append(values)
            subject_data.pop(file_name)

    values = {
        "section_name": "subject_data",
        "json_text": subject_data
    }

    json_files_list.append(values)

    return json_files_list


def save_backup(vk_admin_session, subject):
    PATH = datamanager.read_path()
    path_to_subject_json = subject["path"]

    if len(path_to_subject_json) > 0 and path_to_subject_json[0] != "/":
        path_to_subject_json = PATH + "/bot_notificator/" + path_to_subject_json + "/"
    else:
        path_to_subject_json = PATH + "bot_notificator/" + path_to_subject_json + "/"

    subject_data = make_json_text(path_to_subject_json)
    wiki_full_id = subject_data["wiki_database_id"]

    save_wiki(vk_admin_session, wiki_full_id, subject_data)

    mess_for_log = "Data has been saved in wiki-page."
    logger.message_output("Save backup", mess_for_log)


def load_backup(vk_admin_session, subject):
    PATH = datamanager.read_path()
    path_to_subject_json = subject["path"]

    if len(path_to_subject_json) > 0 and path_to_subject_json[0] != "/":
        path_to_subject_json = PATH + "/bot_notificator/" + path_to_subject_json + "/"
    else:
        path_to_subject_json = PATH + "bot_notificator/" + path_to_subject_json + "/"

    subject_data = datamanager.read_json(path_to_subject_json, "subject_data")
    wiki_full_id = subject_data["wiki_database_id"]

    json_text_from_wiki = read_wiki(vk_admin_session, wiki_full_id)
    json_files_list = split_json_text(path_to_subject_json,
                                      json_text_from_wiki)

    for values in json_files_list:
        json_text = values["json_text"]
        section_name = values["section_name"]
        datamanager.write_json(path_to_subject_json, section_name,
                               json_text)

    mess_for_log = "Data has been saved in files."
    logger.message_output("Load backup", mess_for_log)
