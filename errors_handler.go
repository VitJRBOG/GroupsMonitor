package main

import "strings"

// RequestError проверяет ошибки при отправке запроса
func RequestError(errorMessage string) (string, string) {
	timeoutErrors := []string{
		"connection reset by peer",
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
		"captcha needed", "failed to establish a new connection",
		"connection aborted", "internal server error", "response code 504",
		"response code 502", "many requests per second",
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
