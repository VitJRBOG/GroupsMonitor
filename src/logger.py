# coding: utf-8


import time
import datetime
import datamanager


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
    def to_console(sender, message):
        try:
            print("COMPUTER [" + sender + "]: " + message)

        except Exception as var_except:
            sender += " -> Message output -> To console"
            print(
                "COMPUTER [" + sender + "]: Error, " +
                str(var_except) +
                ". Exit from program...")
            exit(0)

    def to_textfile(sender, message):
        def get_unix_date(time_stamp, junk, type_operation):
            try:
                if type_operation == "from_line":
                    idx = junk.find("] ")
                    if idx != -1:
                        str_date = junk[1:idx]
                    else:
                        return "date not found"
                elif type_operation == "from_str_date":
                    str_date = junk

                date = datetime.datetime.fromtimestamp(time.mktime(time.strptime(str_date, time_stamp)))
                unix_date = int(time.mktime(date.timetuple()))

                return unix_date
            except Exception as var_except:
                sender = "HANYA"
                print(
                    "COMPUTER [" + sender + "]: Error, " +
                    str(var_except) +
                    ". Exit from program...")
                exit(0)

        try:
            PATH = datamanager.read_path(sender)
            # text = datamanager.read_text(sender, PATH, "log")
            # в готовой функции у метода чтения присутствует аргумент, который всё портит
            file = open(PATH + "log.txt")

            time_stamp = "%d.%m.%Y %H:%M:%S"

            str_date = datetime.datetime.now().strftime(time_stamp)

            unix_date_now = get_unix_date(time_stamp, str_date, "from_str_date")

            message = "[" + str(str_date) + "] " + "[" + str(sender) + "]: " + str(message.encode("utf8"))

            new_text = ""

            for line in file:
                unix_date_line = get_unix_date(time_stamp, line, "from_line")

                if unix_date_line > (unix_date_now - 86400):
                    new_text = line + new_text

            if len(new_text) > 0:
                new_text = message + "\n" + new_text
            else:
                new_text = message

            text_array = new_text.split('\n')
            new_text = ""

            i = 0
            while i < len(text_array):

                step = False

                j = 0
                while j < len(text_array) - i - 1:

                    current_time = get_unix_date(time_stamp, text_array[j], "from_line")
                    next_time = get_unix_date(time_stamp, text_array[j + 1], "from_line")

                    if current_time < next_time:
                        next_line = text_array[j + 1]
                        current_line = text_array[j]
                        text_array[j + 1] = current_line
                        text_array[j] = next_line
                        step = True
                    j += 1

                if not step:
                    break

                i += 1

            if len(text_array[0]) == 0:
                text_array.pop(0)
            new_text = '\n'.join(text_array)

            datamanager.write_text(sender, PATH, "log", new_text)

        except Exception as var_except:
            sender += " -> Message output -> To text file"
            print(
                "COMPUTER [" + sender + "]: Error, " +
                str(var_except) +
                ". Exit from program...")
            exit(0)

    # to_console(sender, message)
    to_textfile(sender, message)
