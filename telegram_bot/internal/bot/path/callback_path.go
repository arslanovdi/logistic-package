// Package path - парсинг команд бота
package path

import (
	"errors"
	"fmt"
	"strings"
)

const callbackCount = 4 // "Domain__Subdomain__CallbackName__CallbackData"

// CallbackPath содержит параметры кнопок
type CallbackPath struct {
	Domain       string
	Subdomain    string
	CallbackName string
	CallbackData string
}

// ErrUnknownCallback некорректная команда кнопки
var ErrUnknownCallback = errors.New("unknown callback")

// ParseCallback парсинг строки вида: "Domain__Subdomain__CallbackName__CallbackData" в структуру CallbackPath
func ParseCallback(callbackData string) (CallbackPath, error) {
	callbackParts := strings.SplitN(callbackData, "__", callbackCount)
	if len(callbackParts) != callbackCount {
		return CallbackPath{}, ErrUnknownCallback
	}

	return CallbackPath{
		Domain:       callbackParts[0],
		Subdomain:    callbackParts[1],
		CallbackName: callbackParts[2],
		CallbackData: callbackParts[3],
	}, nil
}

// String строковое представление структуры CallbackPath в виде "Domain__Subdomain__CallbackName__CallbackData"
func (p CallbackPath) String() string {
	return fmt.Sprintf("%s__%s__%s__%s", p.Domain, p.Subdomain, p.CallbackName, p.CallbackData)
}
