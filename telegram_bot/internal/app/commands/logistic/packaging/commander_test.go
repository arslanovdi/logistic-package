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
		name    string
		message tgbotapi.Message
	}{
		{
			name: "create with title",
			message: tgbotapi.Message{
				Text: command + " title",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // create with title
		{
			name: "create with title and weight",
			message: tgbotapi.Message{
				Text: command + " title 50",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // create with title and weight
	}

	failTests := []struct {
		name    string
		message tgbotapi.Message
	}{
		{
			name: "create without arguments",
			message: tgbotapi.Message{
				Text: command + "",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // create without arguments
		{
			name: "create with title and weight where weight is not a number",
			message: tgbotapi.Message{
				Text: command + " title low",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // create with title and weight where weight is not a number
		{
			name: "create with title and weight where weight very big",
			message: tgbotapi.Message{
				Text: command + " title 999999999999999999999999",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // create with title and weight where weight very big
		{
			name: "create with title and weight where weight < 0",
			message: tgbotapi.Message{
				Text: command + " title -99",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // create with title and weight where weight < 0
		{
			name: "create with 3 arguments",
			message: tgbotapi.Message{
				Text: command + " title 50 20",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
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

			handler.New(&tt.message)

		})
	}

	for _, tt := range failTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			grpcMock := mocks.NewClient(t)

			grpcMock.AssertNotCalled(t, "Create", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*model.Package"))

			handler := NewCommander(&tgbotapi.BotAPI{}, service.NewPackageService(grpcMock))

			handler.New(&tt.message)

		})
	}
}

func TestCommander_Get(t *testing.T) {
	t.Parallel()

	command := "/get__logistic__package"

	sucessTests := []struct {
		name    string
		message tgbotapi.Message
	}{
		{
			name: "get with package id",
			message: tgbotapi.Message{
				Text: command + " 5",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // get with package id
	}

	failTests := []struct {
		name    string
		message tgbotapi.Message
	}{
		{
			name: "get without arguments",
			message: tgbotapi.Message{
				Text: command + "",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // get without arguments
		{
			name: "get with package id is not a number",
			message: tgbotapi.Message{
				Text: command + " id",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // get with package id is not a number
		{
			name: "get where package id very big",
			message: tgbotapi.Message{
				Text: command + " 999999999999999999999999",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // get where package id very big
		{
			name: "get where package id < 0",
			message: tgbotapi.Message{
				Text: command + " -99",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // get where package id < 0
		{
			name: "get with 2 arguments",
			message: tgbotapi.Message{
				Text: command + " 50 20",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
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

			handler.Get(&tt.message)

		})
	}

	for _, tt := range failTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			grpcMock := mocks.NewClient(t)

			grpcMock.AssertNotCalled(t, "Get", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("uint64"))

			handler := NewCommander(&tgbotapi.BotAPI{}, service.NewPackageService(grpcMock))

			handler.Get(&tt.message)

		})
	}
}

func TestCommander_Delete(t *testing.T) {
	t.Parallel()

	command := "/delete__logistic__package"

	sucessTests := []struct {
		name    string
		message tgbotapi.Message
	}{
		{
			name: "delete with package id",
			message: tgbotapi.Message{
				Text: command + " 5",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // delete with package id
	}

	failTests := []struct {
		name    string
		message tgbotapi.Message
	}{
		{
			name: "delete without arguments",
			message: tgbotapi.Message{
				Text: command + "",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // delete without arguments
		{
			name: "delete with package id is not a number",
			message: tgbotapi.Message{
				Text: command + " id",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // delete with package id is not a number
		{
			name: "delete where package id very big",
			message: tgbotapi.Message{
				Text: command + " 999999999999999999999999",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // delete where package id very big
		{
			name: "delete where package id < 0",
			message: tgbotapi.Message{
				Text: command + " -99",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // delete where package id < 0
		{
			name: "delete with 2 arguments",
			message: tgbotapi.Message{
				Text: command + " 50 20",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
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

			handler.Delete(&tt.message)

		})
	}

	for _, tt := range failTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			grpcMock := mocks.NewClient(t)

			grpcMock.AssertNotCalled(t, "Get", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("uint64"))

			handler := NewCommander(&tgbotapi.BotAPI{}, service.NewPackageService(grpcMock))

			handler.Delete(&tt.message)

		})
	}
}

func TestCommander_Edit(t *testing.T) {
	t.Parallel()

	command := "/edit__logistic__package"

	sucessTests := []struct {
		name    string
		message tgbotapi.Message
	}{
		{
			name: "edit with package id and title",
			message: tgbotapi.Message{
				Text: command + " 5 title",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // edit with package id and title
		{
			name: "edit with package id, title and weight",
			message: tgbotapi.Message{
				Text: command + " 5 title 100",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // edit with package id, title and weight
	}

	failTests := []struct {
		name    string
		message tgbotapi.Message
	}{
		{
			name: "edit without arguments",
			message: tgbotapi.Message{
				Text: command + "",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // edit without arguments
		{
			name: "edit where package id is not a number",
			message: tgbotapi.Message{
				Text: command + " id title",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // edit where package id is not a number
		{
			name: "edit where package id very big",
			message: tgbotapi.Message{
				Text: command + " 999999999999999999999999 title",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // edit where package id very big
		{
			name: "edit where package id < 0",
			message: tgbotapi.Message{
				Text: command + " -99 title",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // edit where package id < 0
		{
			name: "edit with 4 arguments",
			message: tgbotapi.Message{
				Text: command + " 50 title 20 20",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // edit with 4 arguments
		{
			name: "edit where weight is not a number",
			message: tgbotapi.Message{
				Text: command + " 50 title low",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // edit where weight is not a number
		{
			name: "edit where weight very big",
			message: tgbotapi.Message{
				Text: command + " 50 title 999999999999999999999999",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
		}, // edit where weight very big
		{
			name: "edit where weight < 0",
			message: tgbotapi.Message{
				Text: command + " 50 title -99",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
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

			handler.Edit(&tt.message)

		})
	}

	for _, tt := range failTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			grpcMock := mocks.NewClient(t)

			grpcMock.AssertNotCalled(t, "Update", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*model.Package"))

			handler := NewCommander(&tgbotapi.BotAPI{}, service.NewPackageService(grpcMock))

			handler.Edit(&tt.message)

		})
	}
}

func TestCommander_List(t *testing.T) {
	t.Parallel()

	command := "/list__logistic__package"

	sucessTests := []struct {
		name    string
		message tgbotapi.Message
	}{
		{
			name: "list",
			message: tgbotapi.Message{
				Text: command + "",
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(command),
					},
				},
			},
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

			handler.List(&tt.message)

		})
	}
}
