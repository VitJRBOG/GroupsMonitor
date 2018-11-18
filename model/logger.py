# coding: utf-8


import time
import datetime
import model.datamanager as datamanager


def exception_handler(sender, var_except):
    try:
        if str(var_except).lower().find("captcha needed") != -1:
            message = "Error, " +\
                str(var_except) + ". " +\
                "Timeout: 60 sec."
            message_output(sender, message)
            time.sleep(60)

            return
        elif str(var_except).lower().find("failed to establish " +
                                          "a new connection") != -1:
            message = "Error, " +\
                str(var_except) + ". " +\
                "Timeout: 60 sec."
            message_output(sender, message)
            time.sleep(60)

            return
        elif str(var_except).lower().find("connection aborted") != -1:
            message = "Error, " +\
                str(var_except) + ". " +\
                "Timeout: 60 sec."
            message_output(sender, message)
            time.sleep(60)

            return
        elif str(var_except).lower().find("internal server error") != -1:
            message = "Error, " +\
                str(var_except) + ". " +\
                "Timeout: 60 sec."
            message_output(sender, message)
            time.sleep(60)

            return
        elif str(var_except).lower().find("access_token was " +
                                          "given to another " +
                                          "ip address") != -1:
            message = "Error, " +\
                str(var_except) + "."
            message_output(sender, message)

            return
        elif str(var_except).lower().find("invalid access_token") != -1:
            message = "Error, " +\
                str(var_except) + "."
            message_output(sender, message)

            return
        elif str(var_except).lower().find("response code 504") != -1:
            message = "Error, " +\
                str(var_except) + ". " +\
                "Timeout: 60 sec."
            message_output(sender, message)
            time.sleep(60)

            return
        elif str(var_except).lower().find("response code 502") != -1:
            message = "Error, " +\
                str(var_except) + ". " +\
                "Timeout: 60 sec."
            message_output(sender, message)
            time.sleep(60)

            return
        else:
            message = "Error, " +\
                str(var_except) +\
                ". Exit from program..."
            message_output(sender, message)
            exit(0)

    except Exception as var_except:
        sender += " -> Exception handler"
        message = "Error, " +\
            str(var_except) +\
            ". Exit from program..."
        message_output(sender, message)
        exit(0)


def message_output(sender, message):
    def to_console(sender, message):
        date = datetime.datetime.now().strftime("%d.%m.%Y %H:%M:%S")
        message = "[" + str(date) + "] " + "[" + str(sender) + "]: " + message.encode("utf8")
        print(message)

    def to_textfile(sender, message):
        try:
            PATH = datamanager.read_path()
            text = datamanager.read_text(PATH + "bot_notificator/", "log")

            date = datetime.datetime.now().strftime("%d.%m.%Y %H:%M:%S")

            message = "[" + str(date) + "] " + "[" + str(sender) + "]: " + str(message.encode("utf8"))

            if len(text) > 2:
                text = message + "\n" + text
            else:
                text += message

            datamanager.write_text(PATH + "bot_notificator/", "log", text)

        except Exception as var_except:
            sender += " -> Message output -> To text file"
            print(
                "[" + sender + "]: Error, " +
                str(var_except) +
                ". Exit from program...")
            exit(0)

    to_console(sender, message)
    # to_textfile(sender, message)
