# coding: utf-8


import json


def read_path():
    path = str(open("path.txt", "r").read())

    if len(path) > 0 and path[len(path) - 1] != "/":
        path += "/"

    return path


def read_json(path, file_name):
    loads_json = json.loads(open(str(path) + str(file_name) +
                                 ".json", 'r').read())

    return loads_json


def write_json(path, file_name, loads_json):
    file_json = open(str(path) + str(file_name) + ".json", "w")
    file_json.write(json.dumps(loads_json, indent=4, ensure_ascii=True))
    file_json.close()


def read_text(path, file_name):
    file_text = open(str(path) + str(file_name) + ".txt", "r")
    text = file_text.read()
    file_text.close()

    return text


def write_text(path, file_name, text_output):
    file_text = open(str(path) + str(file_name) + ".txt", "w")
    file_text.write(text_output)
    file_text.close()
