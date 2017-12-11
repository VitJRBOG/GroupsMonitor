# coding: utf-8


import json


def read_path():
    try:
        path = str(open("path.txt", "r").read())

        if len(path) > 0 and path[len(path) - 1] != "/":
            path += "/"

        return path

    except Exception as var_except:
        print(
            "COMPUTER [.. -> Read \"path.txt\"]: Error, " + str(var_except) +
            ". Exit from program...")
        exit(0)


def read_json(sender, path, file_name):
    try:
        loads_json = json.loads(open(str(path) + str(file_name) +
                                     ".json", 'r').read())  # dict

        return loads_json
    except Exception as var_except:
        print(
            "COMPUTER [.. -> " + str(sender) +
            " -> Read JSON]: Error, " + str(var_except) +
            ". Exit from program...")
        exit(0)


def write_json(sender, path, file_name, loads_json):
    try:
        file_json = open(str(path) + str(file_name) + ".json", "w")
        file_json.write(json.dumps(loads_json, indent=4, ensure_ascii=True))
        file_json.close()

    except Exception as var_except:
        print(
            "COMPUTER [.. -> " + str(sender) +
            " -> Write JSON]: Error, " + str(var_except) +
            ". Exit from program...")
        exit(0)


def main():
    try:
        PATH = read_path()

        loads_json = read_json("Main", PATH, "data")

        user_answer = raw_input("USER [Main -> New token]: ")

        loads_json["bot_token"] = user_answer
        loads_json["admin_token"] = user_answer

        write_json("Main", PATH, "data", loads_json)

    except Exception as var_except:
        print(
            "COMPUTER [Main]: Error, " + str(var_except) +
            ". Exit from program...")
        exit(0)


main()
