# coding: utf-8


import starter
import logger
import datamanager


def main_menu():
    sender = "Main menu"

    print("\nCOMPUTER [" + sender + "]: You are in Main menu.")
    print("COMPUTER [" + sender + "]: 1 == Start bot.")
    print("COMPUTER [" + sender + "]: 2 == Settings.")
    print("COMPUTER [" + sender + "]: 3 == Log output.")
    print("COMPUTER [" + sender + "]: 4 == Save backups.")
    print("COMPUTER [" + sender + "]: 5 == Stop bot.")
    print("COMPUTER [" + sender + "]: 0 == Close program.")

    print("COMPUTER [" + sender + "]: Enter digit for next action.")
    user_answer = raw_input("USER [" + sender + "]: (1-5/0) ")

    if user_answer == "0":
        close_program(sender)
    elif user_answer == "1":
        start_bot(sender)
    elif user_answer == "2":
        settings(sender)
    elif user_answer == "3":
        log_output(sender)
    elif user_answer == "4":
        save_backups(sender)
    elif user_answer == "5":
        stop_bot(sender)
    else:
        print("COMPUTER [" + sender + "]: Unknown command. Retry query...")
        main_menu()


def start_bot(sender):
    sender += " -> Start bot"

    objStart = starter.Start()

    objStart.path_checking(sender)

    PATH = datamanager.read_path(sender)

    datafile_was_created = objStart.data_checking(sender, PATH)

    if datafile_was_created:
        mess_for_log = "WARNING! Data file has been created. " +\
            "List of subjects is empty. Check menu Settings."
        logger.message_output(sender, mess_for_log)
        main_menu()
    else:
        data_json, token_validity = objStart.tokens_checking(sender, PATH)
        if not token_validity["admin_token"] or not token_validity["bot_token"]:
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

            data_json = objStart.update_token(sender, PATH, data_json, token_validity, tokens)

        objStart.starting(sender, data_json)

    main_menu()


def settings(sender):
    sender += " -> Settings"
    print("COMPUTER [" + sender + "]: Here is empty....")
    main_menu()


def log_output(sender):
    sender += " -> Log output"
    print("COMPUTER [" + sender + "]: Here is empty....")
    main_menu()


def save_backups(sender):
    try:
        sender += " -> Save backups"

        print("\nCOMPUTER [" + sender + "]: You are in menu Save backups.")

        PATH = datamanager.read_path(sender)

        data_json = datamanager.read_json(sender, PATH, "data")
        subjects = data_json["subjects"]

        i = 0
        while i < len(subjects):
            print("COMPUTER [" + sender + "]: " + str(i + 1) + " == " + subjects[i]["name"])
            i += 1

        if i > 1:
            print("COMPUTER [" + sender + "]: " + str(len(subjects) + 1) + " == Backup all.")
        print("COMPUTER [" + sender + "]: 0 == Step back.")
        print("COMPUTER [" + sender + "]: Enter digit for next action.")
        if i > 1:
            user_answer = raw_input("USER [" + sender + "]: (1-" + str(i + 1) + "/0) ")
        else:
            user_answer = raw_input("USER [" + sender + "]: (1/0) ")

        data_access_admin = {
            "token": data_json["admin_token"]
        }

        vk_admin_session = starter.autorization(sender, data_access_admin, "token")

        if user_answer == "0":
            main_menu()
        elif user_answer == str(len(subjects) + 1):
            j = 0
            while j < len(subjects):
                datamanager.save_backup(sender + " -> " + subjects[j]["name"], PATH, vk_admin_session, subjects[j])
                j += 1
        else:
            i = int(user_answer) - 1
            datamanager.save_backup(sender + " -> " + subjects[i]["name"], PATH, vk_admin_session, subjects[i])

    except Exception as var_except:
        logger.exception_handler(sender, var_except)
        return main_menu()

    main_menu()


def stop_bot(sender):
    sender += " -> Stop bot"
    print("COMPUTER [" + sender + "]: Here is empty....")
    main_menu()


def close_program(sender):
    print("COMPUTER [" + sender + "]: Exit from program...")
    exit(0)


main_menu()
