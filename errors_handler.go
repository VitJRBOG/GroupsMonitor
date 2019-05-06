package main

import "strings"

// DBIOError проверяет ошибки при вводе/выводе данных из БД
func DBIOError(errorMessage string) (string, string) {
	timeoutErrors := []string{
		"database is locked",
	}

	// проверяем текст ошибки на наличие похожих в списке
	for _, item := range timeoutErrors {
		if strings.Contains(strings.ToLower(errorMessage), item) {
			return "timeout error", item
		}
	}

	return "unknown error", ""
}

// RequestError проверяет ошибки при отправке запроса
func RequestError(errorMessage string) (string, string) {
	timeoutErrors := []string{
		"connection reset by peer", "read: connection reset by peer",
		"operation timed out",
		"server sent goaway and closed the connection",
		"failed to establish a new connection",
		"connection aborted",
	}

	// проверяем текст ошибки на наличие похожих в списке
	for _, item := range timeoutErrors {
		if strings.Contains(strings.ToLower(errorMessage), item) {
			return "timeout error", item
		}
	}

	return "unknown error", ""
}

// VkAPIError проверяет ошибки от VK API
func VkAPIError(errorMessage string) (string, string) {
	// список ошибок для вызова таймаута
	timeoutErrors := []string{
		"captcha needed", "many requests per second",
		"internal server error",
	}

	// список ошибок для запроса нового токена доступа
	accessTokenErrors := []string{
		"invalid access_token",
		"access_token was given to another ip address",
		"access_token has expired",
		"no access_token passed",
	}

	// список пропускаемых ошибок
	skipErrors := []string{
		"post was deleted", "post was not found",
		"comment was not found",
	}

	// проверяем текст ошибки на наличие похожих в трех списках
	for _, item := range timeoutErrors {
		if strings.Contains(strings.ToLower(errorMessage), item) {
			return "timeout error", item
		}
	}
	for _, item := range accessTokenErrors {
		if strings.Contains(strings.ToLower(errorMessage), item) {
			return "access token error", item
		}
	}
	for _, item := range skipErrors {
		if strings.Contains(strings.ToLower(errorMessage), item) {
			return "skip error", item
		}
	}

	// если в списках этой ошибки нет, то так и пишем
	return "unknown error", ""
}
