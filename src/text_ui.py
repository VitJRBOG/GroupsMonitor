# coding: utf-8


import starter
import logger
import datamanager


def main_menu():
    SENDER = "Main menu"

    print("\nCOMPUTER [" + SENDER + "]: You are in Main menu.")
    print("COMPUTER [" + SENDER + "]: 1 == Start bot.")
    print("COMPUTER [" + SENDER + "]: 2 == Settings.")
    print("COMPUTER [" + SENDER + "]: 3 == Log output.")
    print("COMPUTER [" + SENDER + "]: 4 == Save backups.")
    print("COMPUTER [" + SENDER + "]: 5 == Stop bot.")
    print("COMPUTER [" + SENDER + "]: 0 == Close program.")

    print("COMPUTER [" + SENDER + "]: Enter digit for next action.")
    user_answer = raw_input("USER [" + SENDER + "]: (1-5/0) ")

    if user_answer == "0":
        print("COMPUTER [" + SENDER + "]: Exit from program...")
        exit(0)
    elif user_answer == "1":
        start_bot()
    elif user_answer == "2":
        settings()
    elif user_answer == "3":
        log_output()
    elif user_answer == "4":
        save_backups()
    elif user_answer == "5":
        stop_bot()
    else:
        print("COMPUTER [" + SENDER + "]: Unknown command. Retry query...")
        main_menu()


def start_bot():
    SENDER = "Main menu -> Start bot"

    objStart = starter.Start()

    objStart.path_checking(SENDER)

    PATH = datamanager.read_path(SENDER)

    objStart.log_file_checking(SENDER, PATH)

    datafile_was_created = objStart.data_checking(SENDER, PATH)

    if datafile_was_created:
        mess_for_log = "WARNING! Data file has been created. " +\
            "List of subjects is empty. Check menu Settings."
        logger.message_output(SENDER, mess_for_log)
        main_menu()
    else:
        data_json, token_validity = objStart.tokens_checking(SENDER, PATH)
        if not token_validity["admin_token"] or not token_validity["bot_token"]:
            admin_token_validity = token_validity["admin_token"]
            bot_token_validity = token_validity["bot_token"]

            admin_token = ""
            bot_token = ""

            if not admin_token_validity:
                admin_token = raw_input("USER [" + SENDER + "]: ")
            if not bot_token_validity:
                bot_token = raw_input("USER [" + SENDER + "]: ")

            tokens = {
                "admin_token": admin_token,
                "bot_token": bot_token
            }

            data_json = starter.update_token(SENDER, PATH, data_json, token_validity, tokens)

        objStart.starting(SENDER, data_json)

    main_menu()


def settings():

    SENDER = "Main menu -> Settings"

    def data_settings(SENDER, PATH):
        sender = SENDER + " -> Data settings"

        def show_settings(sender, data_json):

            print("\n")
            print("Token of admin: " + str(data_json["admin_token"]))
            print("Token of bot: " + str(data_json["bot_token"]))
            print("\n")

        data_json = datamanager.read_json(sender, PATH, "data")

        show_settings(sender, data_json)

        print("COMPUTER [" + sender + "]: 1 == Update token of admin.")
        print("COMPUTER [" + sender + "]: 2 == Update token of bot.")
        print("COMPUTER [" + sender + "]: 0 == Step back.")

        print("COMPUTER [" + sender + "]: Enter digit for next action.")
        user_answer = raw_input("USER [" + sender + "]: (1-2/0) ")

        if user_answer == "0":
            settings()
        elif user_answer == "1":
            sender += " -> New admin token"

            token_validity = {
                "admin_token": False,
                "bot_token": True
            }

            admin_token_validity = token_validity["admin_token"]
            bot_token_validity = token_validity["bot_token"]

            admin_token = ""
            bot_token = ""

            if not admin_token_validity:
                admin_token = raw_input("USER [" + sender + "]: ")
            if not bot_token_validity:
                bot_token = raw_input("USER [" + sender + "]: ")

            tokens = {
                "admin_token": admin_token,
                "bot_token": bot_token
            }

            starter.update_token(sender, PATH, data_json, token_validity, tokens)
        elif user_answer == "2":
            sender += " -> New admin token"

            token_validity = {
                "admin_token": True,
                "bot_token": False
            }

            admin_token_validity = token_validity["admin_token"]
            bot_token_validity = token_validity["bot_token"]

            admin_token = ""
            bot_token = ""

            if not admin_token_validity:
                admin_token = raw_input("USER [" + sender + "]: ")
            if not bot_token_validity:
                bot_token = raw_input("USER [" + sender + "]: ")

            tokens = {
                "admin_token": admin_token,
                "bot_token": bot_token
            }

            starter.update_token(sender, PATH, data_json, token_validity, tokens)
        else:
            print("COMPUTER [" + sender + "]: Unknown command. Retry query...")
            settings()

        data_settings(SENDER, PATH)

    PATH = datamanager.read_path(SENDER)

    print("\nCOMPUTER [" + SENDER + "]: You are in menu Settings.")
    print("COMPUTER [" + SENDER + "]: 1 == Open data settings.")
    print("COMPUTER [" + SENDER + "]: 0 == Step back.")

    print("COMPUTER [" + SENDER + "]: Enter digit for next action.")
    user_answer = raw_input("USER [" + SENDER + "]: (1/0) ")

    if user_answer == "0":
        main_menu()
    elif user_answer == "1":
        data_settings(SENDER, PATH)
    else:
        print("COMPUTER [" + SENDER + "]: Unknown command. Retry query...")
        settings()

    main_menu()


def log_output():
    SENDER = "Main menu -> Log output"
    print("COMPUTER [" + SENDER + "]: Here is empty....")
    main_menu()


def save_backups():
    try:
        SENDER = "Main menu -> Save backups"

        print("\nCOMPUTER [" + SENDER + "]: You are in menu Save backups.")

        PATH = datamanager.read_path(SENDER)

        data_json = datamanager.read_json(SENDER, PATH, "data")
        subjects = data_json["subjects"]

        i = 0
        while i < len(subjects):
            print("COMPUTER [" + SENDER + "]: " + str(i + 1) + " == " + subjects[i]["name"])
            i += 1

        if i > 1:
            print("COMPUTER [" + SENDER + "]: " + str(len(subjects) + 1) + " == Backup all.")
        print("COMPUTER [" + SENDER + "]: 0 == Step back.")
        print("COMPUTER [" + SENDER + "]: Enter digit for next action.")
        if i > 1:
            user_answer = raw_input("USER [" + SENDER + "]: (1-" + str(i + 1) + "/0) ")
        else:
            user_answer = raw_input("USER [" + SENDER + "]: (1/0) ")

        data_access_admin = {
            "token": data_json["admin_token"]
        }

        vk_admin_session = starter.autorization(SENDER, data_access_admin, "token")

        if user_answer == "0":
            main_menu()
        elif user_answer == str(len(subjects) + 1):
            j = 0
            while j < len(subjects):
                datamanager.save_backup(SENDER + " -> " + subjects[j]["name"], PATH, vk_admin_session, subjects[j])
                j += 1
        else:
            i = int(user_answer) - 1
            datamanager.save_backup(SENDER + " -> " + subjects[i]["name"], PATH, vk_admin_session, subjects[i])

        save_backups()

    except Exception as var_except:
        logger.exception_handler(SENDER, var_except)
        return main_menu()

    main_menu()


def stop_bot():
    SENDER = "Main menu -> Stop bot"
    print("COMPUTER [" + SENDER + "]: Here is empty....")
    main_menu()


main_menu()
