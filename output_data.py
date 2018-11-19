# coding: utf8
u"""Модуль вывода данных."""


import datetime


def output_text_row(sender, message):
    u"""Запускает вывод текстовой строки."""
    def to_console(sender, message):
        u"""Выводит текстовую строку в консоль."""
        date = datetime.datetime.now().strftime("%d.%m.%Y %H:%M:%S")
        output = "[" + str(date) + "] " + "[" + str(sender) + \
            "]: " + str(message.encode("utf8"))
        print(output)

    to_console(sender, message)
