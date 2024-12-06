package path

import (
	"errors"
	"fmt"
	"strings"
)

// Файл с обработкой команд

// CommandPath содержит параметры команд
type CommandPath struct {
	CommandName string
	Domain      string
	Subdomain   string
}

// ErrUnknownCommand некорректная команда
var ErrUnknownCommand = errors.New("unknown command")

// ParseCommand парсинг строки вида: "CommandName__Domain__Subdomain" в структуру CommandPath
func ParseCommand(commandText string) (CommandPath, error) {
	commandParts := strings.SplitN(commandText, "__", 3)
	if len(commandParts) != 3 {
		return CommandPath{}, ErrUnknownCommand
	}

	return CommandPath{
		CommandName: commandParts[0],
		Domain:      commandParts[1],
		Subdomain:   commandParts[2],
	}, nil
}

// WithCommandName NO usages
func (c CommandPath) WithCommandName(commandName string) CommandPath {
	c.CommandName = commandName
	// TODO no usages
	return c
}

// String строковое представление структуры CommandPath в виде "CommandName__Domain__Subdomain"
func (c CommandPath) String() string {
	return fmt.Sprintf("/%s__%s__%s", c.CommandName, c.Domain, c.Subdomain)
}
