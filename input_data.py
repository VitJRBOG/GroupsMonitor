# coding: utf8
u"""Модуль ввода данных пользователем."""


def get_vk_user_token(token_owner, token_purpose):
    u"""Ввод пользователем токена доступа в консоль."""
    access_token = raw_input("USER [" + token_owner +\
        "'s new access token for " + token_purpose + "]: ")
    return access_token
