# coding: utf8


import model.datamanager as datamanager
import time


def status_changer(status_list, new_status):
    status_list["offline"].hide()
    status_list["waiting"].hide()
    status_list["processing"].hide()

    status_list[new_status].show()


def logger_changer(objLabelLogger):
    last_len = 0
    while True:
        PATH = datamanager.read_path()
        log_text = datamanager.read_text(PATH + "bot_notificator/", "log")
        if len(log_text) > last_len:
            log_text = log_text.split('\n')
            text = ""
            i = 0
            for line in log_text:
                text += line
                if i < 99:
                    text += "\n"
                i += 1
                if i == 100:
                    break
            objLabelLogger.set_text(text)
            last_len = len(text)
        time.sleep(1)
