package packaging

import (
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/arslanovdi/logistic-package/telegram_bot/internal/service"
	"github.com/arslanovdi/logistic-package/telegram_bot/mocks"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCommander_New(t *testing.T) {
	t.Parallel()

	command := "/new__logistic__package"

	sucessTests := []struct {
		name string
		args string
	}{
		{
			name: "create with title",
			args: " title",
		}, // create with title
		{
			name: "create with title and weight",
			args: " title 50",
		}, // create with title and weight
	}

	failTests := []struct {
		name string
		args string
	}{
		{
			name: "create without arguments",
			args: "",
		}, // create without arguments
		{
			name: "create with title and weight where weight is not a number",
			args: " title low",
		}, // create with title and weight where weight is not a number
		{
			name: "create with title and weight where weight very big",
			args: " title 999999999999999999999999",
		}, // create with title and weight where weight very big
		{
			name: "create with title and weight where weight < 0",
			args: " title -99",
		}, // create with title and weight where weight < 0
		{
			name: "create with 3 arguments",
			args: " title 50 20",
		}, // create with 3 arguments
	}

	for _, tt := range sucessTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			grpcMock := mocks.NewClient(t)

			var (
				err error
				id  uint64
			)
			grpcMock.EXPECT().Create(mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*model.Package")).Return(&id, err)

			handler := NewCommander(&tgbotapi.BotAPI{}, service.NewPackageService(grpcMock))

			message := tgbotapi.Message{
				Text: command + tt.args,
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			}

			handler.New(&message)

		})
	}

	for _, tt := range failTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			grpcMock := mocks.NewClient(t)

			grpcMock.AssertNotCalled(t, "Create", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*model.Package"))

			handler := NewCommander(&tgbotapi.BotAPI{}, service.NewPackageService(grpcMock))

			message := tgbotapi.Message{
				Text: command + tt.args,
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			}
			handler.New(&message)

		})
	}
}

func TestCommander_Get(t *testing.T) {
	t.Parallel()

	command := "/get__logistic__package"

	sucessTests := []struct {
		name string
		args string
	}{
		{
			name: "get with package id",
			args: " 5",
		}, // get with package id
	}

	failTests := []struct {
		name string
		args string
	}{
		{
			name: "get without arguments",
			args: "",
		}, // get without arguments
		{
			name: "get with package id is not a number",
			args: " id",
		}, // get with package id is not a number
		{
			name: "get where package id very big",
			args: " 999999999999999999999999",
		}, // get where package id very big
		{
			name: "get where package id < 0",
			args: " -99",
		}, // get where package id < 0
		{
			name: "get with 2 arguments",
			args: " 50 20",
		}, // get with 2 arguments
	}

	for _, tt := range sucessTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			grpcMock := mocks.NewClient(t)

			var (
				err error
				pkg model.Package
			)
			grpcMock.EXPECT().Get(mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("uint64")).Return(&pkg, err)

			handler := NewCommander(&tgbotapi.BotAPI{}, service.NewPackageService(grpcMock))

			message := tgbotapi.Message{
				Text: command + tt.args,
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			}
			handler.Get(&message)

		})
	}

	for _, tt := range failTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			grpcMock := mocks.NewClient(t)

			grpcMock.AssertNotCalled(t, "Get", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("uint64"))

			handler := NewCommander(&tgbotapi.BotAPI{}, service.NewPackageService(grpcMock))

			message := tgbotapi.Message{
				Text: command + tt.args,
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			}

			handler.Get(&message)

		})
	}
}

func TestCommander_Delete(t *testing.T) {
	t.Parallel()

	command := "/delete__logistic__package"

	sucessTests := []struct {
		name string
		args string
	}{
		{
			name: "delete with package id",
			args: " 5",
		}, // delete with package id
	}

	failTests := []struct {
		name string
		args string
	}{
		{
			name: "delete without arguments",
			args: "",
		}, // delete without arguments
		{
			name: "delete with package id is not a number",
			args: " id",
		}, // delete with package id is not a number
		{
			name: "delete where package id very big",
			args: " 999999999999999999999999",
		}, // delete where package id very big
		{
			name: "delete where package id < 0",
			args: " -99",
		}, // delete where package id < 0
		{
			name: "delete with 2 arguments",
			args: " 50 20",
		}, // delete with 2 arguments
	}

	for _, tt := range sucessTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			grpcMock := mocks.NewClient(t)

			var (
				err error
			)
			grpcMock.EXPECT().Delete(mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("uint64")).Return(err)

			handler := NewCommander(&tgbotapi.BotAPI{}, service.NewPackageService(grpcMock))

			message := tgbotapi.Message{
				Text: command + tt.args,
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			}

			handler.Delete(&message)

		})
	}

	for _, tt := range failTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			grpcMock := mocks.NewClient(t)

			grpcMock.AssertNotCalled(t, "Get", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("uint64"))

			handler := NewCommander(&tgbotapi.BotAPI{}, service.NewPackageService(grpcMock))

			message := tgbotapi.Message{
				Text: command + tt.args,
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			}

			handler.Delete(&message)

		})
	}
}

func TestCommander_Edit(t *testing.T) {
	t.Parallel()

	command := "/edit__logistic__package"

	sucessTests := []struct {
		name string
		args string
	}{
		{
			name: "edit with package id and title",
			args: " 5 title",
		}, // edit with package id and title
		{
			name: "edit with package id, title and weight",
			args: " 5 title 100",
		}, // edit with package id, title and weight
	}

	failTests := []struct {
		name string
		args string
	}{
		{
			name: "edit without arguments",
			args: "",
		}, // edit without arguments
		{
			name: "edit where package id is not a number",
			args: " id title",
		}, // edit where package id is not a number
		{
			name: "edit where package id very big",
			args: " 999999999999999999999999 title",
		}, // edit where package id very big
		{
			name: "edit where package id < 0",
			args: " -99 title",
		}, // edit where package id < 0
		{
			name: "edit with 4 arguments",
			args: " 50 title 20 20",
		}, // edit with 4 arguments
		{
			name: "edit where weight is not a number",
			args: " 50 title low",
		}, // edit where weight is not a number
		{
			name: "edit where weight very big",
			args: " 50 title 999999999999999999999999",
		}, // edit where weight very big
		{
			name: "edit where weight < 0",
			args: " 50 title -99",
		}, // edit where weight < 0
	}

	for _, tt := range sucessTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			grpcMock := mocks.NewClient(t)

			var (
				err error
			)
			grpcMock.EXPECT().Update(mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*model.Package")).Return(err)

			handler := NewCommander(&tgbotapi.BotAPI{}, service.NewPackageService(grpcMock))

			message := tgbotapi.Message{
				Text: command + tt.args,
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			}

			handler.Edit(&message)

		})
	}

	for _, tt := range failTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			grpcMock := mocks.NewClient(t)

			grpcMock.AssertNotCalled(t, "Update", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*model.Package"))

			handler := NewCommander(&tgbotapi.BotAPI{}, service.NewPackageService(grpcMock))

			message := tgbotapi.Message{
				Text: command + tt.args,
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			}

			handler.Edit(&message)

		})
	}
}

func TestCommander_List(t *testing.T) {
	t.Parallel()

	command := "/list__logistic__package"

	sucessTests := []struct {
		name string
		args string
	}{
		{
			name: "list",
			args: "",
		}, // list
	}

	for _, tt := range sucessTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			grpcMock := mocks.NewClient(t)

			var (
				err error
				pkg []model.Package
			)
			grpcMock.EXPECT().List(mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("uint64"), mock.AnythingOfType("uint64")).Return(pkg, err)

			handler := NewCommander(&tgbotapi.BotAPI{}, service.NewPackageService(grpcMock))

			message := tgbotapi.Message{
				Text: command + tt.args,
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			}

			handler.List(&message)

		})
	}
}
