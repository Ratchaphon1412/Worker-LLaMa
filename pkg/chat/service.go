package chat

import (
	"github.com/Ratchaphon1412/worker-llama/pkg/entities"
)

type Service interface {
	CreateChat(chat *entities.Chat) (*entities.Chat, error)
	UpdateChat(chat *entities.Chat) (*entities.Chat, error)
	FindChatByID(id uint) (*entities.Chat, error)
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}
func (s *service) CreateChat(chat *entities.Chat) (*entities.Chat, error) {
	chat, err := s.repository.Create(chat)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (s *service) UpdateChat(chatData *entities.Chat) (*entities.Chat, error) {
	chat, err := s.repository.Update(chatData)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (s *service) FindChatByID(id uint) (*entities.Chat, error) {
	chat, err := s.repository.GetChatByID(id)
	if err != nil {
		return nil, err
	}
	return chat, nil
}
