package chat

import (
	"github.com/Ratchaphon1412/worker-llama/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	Create(chat *entities.Chat) (*entities.Chat, error)
	Update(chat *entities.Chat) (*entities.Chat, error)
	Delete(id uint) error
	GetChatByID(id uint) (*entities.Chat, error)
	GetChatByAccountID(accountID uint) ([]entities.Chat, error)
	GetAllChats() ([]entities.Chat, error)
	GetChatByWorkflowID(workflowID uint) ([]entities.Chat, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(chat *entities.Chat) (*entities.Chat, error) {
	if err := r.db.Create(chat).Error; err != nil {
		return nil, err
	}
	return chat, nil
}
func (r *repository) Update(chat *entities.Chat) (*entities.Chat, error) {
	if err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(chat).Error; err != nil {
		return nil, err
	}
	return chat, nil
}
func (r *repository) Delete(id uint) error {
	if err := r.db.Delete(&entities.Chat{}, id).Error; err != nil {
		return err
	}
	return nil
}
func (r *repository) GetChatByID(id uint) (*entities.Chat, error) {
	var chat entities.Chat
	if err := r.db.First(&chat, id).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}
func (r *repository) GetChatByAccountID(accountID uint) ([]entities.Chat, error) {
	var chats []entities.Chat
	if err := r.db.Where("account_id = ?", accountID).Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}
func (r *repository) GetAllChats() ([]entities.Chat, error) {
	var chats []entities.Chat
	if err := r.db.Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}
func (r *repository) GetChatByWorkflowID(workflowID uint) ([]entities.Chat, error) {
	var chats []entities.Chat
	if err := r.db.Where("workflow_id = ?", workflowID).Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}
