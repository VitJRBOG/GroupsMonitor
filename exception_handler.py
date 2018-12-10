# coding: utf-8
u"""Модуль обработчика исключений."""


import time
import datetime
import output_data


def handling(sender, var_except, timeout):
    u"""Функция обработки исключений."""
    def timeout_errors_handler(sender, var_except, timeout):
        u"""Обработчик ошибок с таймаутом."""
        timeout_errors = [
            "captcha needed", "failed to establish a new connection",
            "connection aborted", "internal server error", "response code 504",
            "response code 502"
        ]
        for text_error in timeout_errors:
            if str(var_except).lower().find(text_error) != -1:
                message = "Error, " +\
                    str(var_except) + ". " +\
                    "Timeout: " + str(timeout) + " sec."
                output_data.output_text_row(sender, message)
                time.sleep(timeout)

                return True

    def fatal_errors_handler(sender, var_except):
        u"""Обработчик фатальных ошибок."""
        fatal_errors = [
            "invalid access_token",
            "access_token was given to another ip address"
        ]
        for text_error in fatal_errors:
            if str(var_except).lower().find(text_error) != -1:
                message = "Error, " +\
                    str(var_except) + ". " +\
                    "It's a fatal error. Stop operation..."
                output_data.output_text_row(sender, message)

                return True

    if timeout_errors_handler(sender, var_except, timeout) is True:
        return
    elif fatal_errors_handler(sender, var_except) is True:
        exit(0)
    else:
        message = "Error, " +\
            str(var_except) + ". " +\
            "It's a unknown error. Stop operation..."
        output_data.output_text_row(sender, message)
        exit(0)
