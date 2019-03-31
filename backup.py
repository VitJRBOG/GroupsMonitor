# coding: utf8
u"""Модуль для сохранения резервной копии данных из файлов с настройками."""


import json
import os
import vkapi
import data_manager
import output_data
import exception_handler


def read_wiki(access_token, wiki_full_id):

    wiki_owner_id = int(wiki_full_id[0:wiki_full_id.rfind('_')])
    wiki_id = int(wiki_full_id[wiki_full_id.rfind('_') + 1:])

    values = {
        "owner_id": wiki_owner_id,
        "page_id": wiki_id,
        "need_html": 1,
        "v": 5.92
    }

    result = send_request("Read wiki", "pages.get", values, access_token)

    text = result["response"]["html"][8:]
    data_json = json.loads(text)

    return data_json


def save_wiki(access_token, wiki_full_id, data_json):
    wiki_owner_id = int(wiki_full_id[1:wiki_full_id.rfind('_')])
    wiki_id = int(wiki_full_id[wiki_full_id.rfind('_') + 1:])

    text = json.dumps(data_json, indent=4, ensure_ascii=False)

    values = {
        "group_id": wiki_owner_id,
        "page_id": wiki_id,
        "text": text,
        "v": 5.92
    }

    result = send_request("Save wiki", "pages.save", values, access_token)


def send_request(sender, method_name, values, access_token):
    error_repeats = 0
    result = vkapi.method(method_name, values, access_token)
    if "response" in result:
        return result["response"]
    else:
        message_error = result["error"]["error_msg"]
        if error_repeats < 5:
            error_repeats += 1
        timeout = error_repeats * 2
        exception_handler.handling(sender, message_error, timeout)
        return send_request(sender, method_name, values, access_token)

    return result


def make_json_text(path_to_subject_json):
    subject_data = data_manager.read_json(path_to_subject_json, "subject_data")

    files = os.listdir(path_to_subject_json)
    jsonfiles = filter(lambda x: x.endswith(".json"), files)

    for file_name in jsonfiles:
        file_name = file_name.replace(".json", "")
        if file_name != "subject_data":
            subject_section_data = data_manager.read_json(path_to_subject_json,
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


def save_backup(access_token, subject):
    PATH = data_manager.read_path()
    subject_name = subject["name"]
    path_to_subject_json = subject["path"]

    if len(path_to_subject_json) > 0 and path_to_subject_json[0] != "/":
        path_to_subject_json = PATH + "/" + path_to_subject_json + "/"
    else:
        path_to_subject_json = PATH + path_to_subject_json + "/"

    subject_data = make_json_text(path_to_subject_json)
    wiki_full_id = subject_data["wiki_database_id"]

    save_wiki(access_token, wiki_full_id, subject_data)

    mess_for_log = "Data has been saved in wiki-page."
    output_data.output_text_row(
        "Save backup for " + subject_name, mess_for_log)


def load_backup(access_token, subject):
    PATH = data_manager.read_path()
    subject_name = subject["name"]
    path_to_subject_json = subject["path"]

    if len(path_to_subject_json) > 0 and path_to_subject_json[0] != "/":
        path_to_subject_json = PATH + "/" + path_to_subject_json + "/"
    else:
        path_to_subject_json = PATH + path_to_subject_json + "/"

    subject_data = data_manager.read_json(path_to_subject_json, "subject_data")
    wiki_full_id = subject_data["wiki_database_id"]

    json_text_from_wiki = read_wiki(access_token, wiki_full_id)
    json_files_list = split_json_text(path_to_subject_json,
                                      json_text_from_wiki)

    for values in json_files_list:
        json_text = values["json_text"]
        section_name = values["section_name"]
        data_manager.write_json(path_to_subject_json, section_name,
                               json_text)

    mess_for_log = "Data has been saved in files."
    output_data.output_text_row(
        "Load backup for " + subject_name, mess_for_log)
