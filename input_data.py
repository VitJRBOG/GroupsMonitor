# coding: utf8
u"""Модуль ввода данных пользователем."""


def get_vk_user_token(token_owner):
    u"""Ввод пользователем токена доступа в консоль."""
    access_token = raw_input("USER [" + token_owner + "'s new access token]: ")
    return access_token
