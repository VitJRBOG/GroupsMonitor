# coding: utf-8


import json
import time
import bughandler


def read_path(sender):
    if sender != "":
        sender += " -> Read \"path.txt\""
    else:
        sender += "Read \"path.txt\""

    try:
        path = str(open("path.txt", "r").read())

        if len(path) > 0 and path[len(path) - 1] != "/":
            path += "/"

        return path

    except Exception as var_except:
        bughandler.exception_handler(sender, var_except)
        return read_path(sender)


def read_json(sender, path, file_name):
    sender += " -> Read JSON"

    try:
        loads_json = json.loads(open(str(path) + str(file_name) +
                                     ".json", 'r').read())  # dict

        return loads_json
    except Exception as var_except:
        bughandler.exception_handler(sender, var_except)
        return read_json(sender, path, file_name)


def write_json(sender, path, file_name, loads_json):
    sender += " -> Write JSON"

    try:
        file_json = open(str(path) + str(file_name) + ".json", "w")
        file_json.write(json.dumps(loads_json, indent=4, ensure_ascii=True))
        file_json.close()

    except Exception as var_except:
        bughandler.exception_handler(sender, var_except)
        return write_json(sender, path, file_name, loads_json)


def read_wiki(sender, vk_admin_session, wiki_full_id):
    sender += " -> Read Wiki"

    try:
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

    except Exception as var_except:
        bughandler.exception_handler(sender, var_except)
        return save_wiki(sender, vk_admin_session, wiki_id, text)


def save_wiki(sender, vk_admin_session, wiki_full_id, data_json):
    sender += " -> Save Wiki"

    try:
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

    except Exception as var_except:
        bughandler.exception_handler(sender, var_except)
        return save_wiki(sender, vk_admin_session, wiki_full_id, text)
