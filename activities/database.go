package activities

import (
	"context"
	"log"

	"github.com/Ratchaphon1412/worker-llama/configs"
	"github.com/Ratchaphon1412/worker-llama/pkg/chat"
	"github.com/Ratchaphon1412/worker-llama/pkg/entities"
	"github.com/Ratchaphon1412/worker-llama/worker/drivers/database"
)

type Answer struct {
	ChatID   uint     `json:"chat_id"`
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
	Media    string   `json:"media"`
	Search   []Search `json:"search"`
}

type Search struct {
	Kind        string      `json:"kind"`
	Title       string      `json:"title"`
	DisplayLink string      `json:"display_link"`
	Link        string      `json:"link"`
	Image       string      `json:"image"`
	ChatID      uint        `json:"chat_id"`
	Thumbnails  []Thumbnail `json:"thumbnails"`
}
type Thumbnail struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Src    string `json:"src"`
}

func UpdateAnswer(ctx context.Context, conf configs.Config, answer Answer) error {

	// create a chat repository
	chatRepo := chat.NewRepository(database.DB.Db)
	// create a service
	chatService := chat.NewService(chatRepo)
	// convert []Search to []entities.Research
	var researches []entities.Research
	for _, s := range answer.Search {
		research := entities.Research{
			Kind:        s.Kind,
			Title:       s.Title,
			DisplayLink: s.DisplayLink,
			Link:        s.Link,
			Image:       s.Image,
			ChatID:      s.ChatID,
			Thumbnails: func(ths []Thumbnail) []entities.Thumbnail {
				var result []entities.Thumbnail
				for _, t := range ths {
					result = append(result, entities.Thumbnail{
						Width:  t.Width,
						Height: t.Height,
						Src:    t.Src,
					})
				}
				return result
			}(s.Thumbnails),
		}
		researches = append(researches, research)
	}
	// find the chat by ID
	chatOld, err := chatService.FindChatByID(answer.ChatID)
	log.Println("chatOld", chatOld)
	if err != nil {
		log.Println("Error finding chat by ID:", err)
		return err
	}

	chatOld.Answer = answer.Answer
	chatOld.AnswerMedia = answer.Media
	chatOld.Status = entities.STATUS_COMPLETE
	chatOld.Prompt = answer.Question
	chatOld.Researches = researches

	// update the answer in the database
	if _, err := chatService.UpdateChat(chatOld); err != nil {
		return err
	}

	return nil
}
