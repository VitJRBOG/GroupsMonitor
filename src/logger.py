# coding: utf-8


import time


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
        exit(0)


def message_output(sender, message):
    try:
        print("COMPUTER [" + sender + "]: " + message)

    except Exception as var_except:
        sender += " -> Message output"
        print(
            "COMPUTER [" + sender + "]: Error, " +
            str(var_except) +
            ". Exit from program...")
        exit(0)
