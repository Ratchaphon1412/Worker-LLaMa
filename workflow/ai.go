package workflow

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Ratchaphon1412/worker-llama/activities"
	"github.com/Ratchaphon1412/worker-llama/configs"
	"github.com/caarlos0/env/v11"
	"go.temporal.io/sdk/workflow"
)

type LLMResponse struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Prompt  []interface{} `json:"prompt"`
	Choices []struct {
		FinishReason string      `json:"finish_reason"`
		Seed         uint64      `json:"seed"`
		Logprobs     interface{} `json:"logprobs"`
		Index        int         `json:"index"`
		Message      struct {
			Role      string        `json:"role"`
			Content   string        `json:"content"`
			ToolCalls []interface{} `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Workflow is a Hello World workflow definition.
func AIWorkflow(ctx workflow.Context, chatID uint, redisChanel string, prompt string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	workflowID := workflow.GetInfo(ctx).WorkflowExecution.ID

	logger := workflow.GetLogger(ctx)
	logger.Info(" workflow AI Workflow started")
	logger.Info("Prompt: ", prompt)

	var conf configs.Config
	// Load environment variables
	if err := env.Parse(&conf); err != nil {
		logger.Error("Failed to parse env vars: %v", err)
	}

	// Research the topic
	var researchResult activities.ResearchResult
	if err := workflow.ExecuteActivity(ctx, activities.Research, conf, prompt).Get(ctx, &researchResult); err != nil {
		logger.Error("Activity Research Failed Error: ", err)
		return "", err
	}

	logger.Info("Research completed.", "result : ", researchResult)

	// Scrap the result
	var scrapResult []string
	if err := workflow.ExecuteActivity(ctx, activities.WebScrap, conf, researchResult).Get(ctx, &scrapResult); err != nil {
		logger.Error("Activity WebScrap failed.", "Error", err)
		return "", err
	}
	logger.Info("Scrap completed.", "result : ", scrapResult)

	// LLM Answer
	llmparam := activities.LLMParam{
		SystemPrompt: fmt.Sprintf(`
		เมื่อได้รับคำถามจากผู้ใช้และบริบทบางส่วน โปรดเขียนคำตอบที่ชัดเจน กระชับ
		และถูกต้องตามบริบทสำหรับคำถามนั้น คุณจะได้รับ
		ชุดบริบทที่เกี่ยวข้องกับคำถามที่ได้จาก internet โปรดใช้
		บริบทเมื่อร่างคำตอบของคุณ และวิเคราะห์ว่าควรใช้คำตอบจาก บริบท หรือ ชุดข้อมูลระบบ

		นี่คือชุดบริบทจาก internet:
		<|context|>
			%s
		<|endcontext|>

		<|system|>
			คุณคือผู้ช่วย AI ที่มีความรู้สูง ที่มีชื่อว่า "หนึ่ง" คุณจะต้องตอบคำถามของผู้ใช้ 

			ซึ่งตอบเป็นภาษาไทย คำใหนเป็นภาษาอังกฤษให้แปลเป็นภาษาไทย แบบสะกดคำ
			ให้ถูกต้องตามหลักภาษาไทย
		<|system|>
	`, scrapResult),
		Prompt: prompt,
	}

	var answer []byte
	err := workflow.ExecuteActivity(ctx, activities.LLM, conf, llmparam).Get(ctx, &answer)
	if err != nil {
		logger.Error("Activity LLM failed.", "Error : ", err)
		return "", err
	}
	logger.Info("LLM completed.", "result : ", answer)

	// Convert the result to json
	var answerObj LLMResponse
	if err := json.Unmarshal(answer, &answerObj); err != nil {
		logger.Error("Failed to unmarshal answer to JSON.", "Error", err)
		return "", err
	}

	// TTS Answer
	if err := workflow.ExecuteActivity(ctx, activities.TTS, conf, answerObj.Choices[0].Message.Content, workflowID).Get(ctx, &answer); err != nil {
		logger.Error("Activity TTS failed.", "Error : ", err)
		return "", err
	}

	// Upload to TTS MinIO
	pathFile := fmt.Sprintf(conf.TTSSaveToLocal+"/%s.mp3", workflowID)
	fileName := fmt.Sprintf("%s.mp3", workflowID)
	var mediaPublicURL string
	if err := workflow.ExecuteActivity(ctx, activities.Storage, conf, fileName, pathFile).Get(ctx, &mediaPublicURL); err != nil {
		logger.Error("Activity Storage failed.", "Error : ", err)
		return "", err
	}
	logger.Info("Storage completed.", "result: ", pathFile)

	// Convert the searchResult to Object for update

	var researchResultCompleted []activities.Search

	for _, item := range researchResult.Items {
		var image string
		var thumbnails []activities.Thumbnail

		// ดึงภาพหลักจาก pagemap.cse_image
		if cseImages, ok := item.Pagemap["cse_image"].([]interface{}); ok && len(cseImages) > 0 {
			if imgMap, ok := cseImages[0].(map[string]interface{}); ok {
				if src, ok := imgMap["src"].(string); ok {
					image = src
				}
			}
		}

		// ดึง thumbnail จาก pagemap.cse_thumbnail
		if cseThumbs, ok := item.Pagemap["cse_thumbnail"].([]interface{}); ok {
			for _, thumb := range cseThumbs {
				if thumbMap, ok := thumb.(map[string]interface{}); ok {
					width, _ := strconv.Atoi(fmt.Sprintf("%v", thumbMap["width"]))
					height, _ := strconv.Atoi(fmt.Sprintf("%v", thumbMap["height"]))
					src, _ := thumbMap["src"].(string)
					thumbnails = append(thumbnails, activities.Thumbnail{
						Width:  width,
						Height: height,
						Src:    src,
					})
				}
			}
		}

		researchResultCompleted = append(researchResultCompleted, activities.Search{
			Kind:        item.Kind,
			Title:       item.Title,
			DisplayLink: item.DisplayLink,
			Link:        item.Link,
			Image:       image,
			Thumbnails:  thumbnails,
		})
	}

	answerObjectCompleted := activities.Answer{
		ChatID:   chatID,
		Question: prompt,
		Answer:   answerObj.Choices[0].Message.Content,
		Media:    mediaPublicURL,
		Search:   researchResultCompleted,
	}

	// Update Answer in the database
	if err := workflow.ExecuteActivity(ctx, activities.UpdateAnswer, conf, answerObjectCompleted).Get(ctx, &answer); err != nil {
		logger.Error("Activity UpdateAnswer failed.", "Error : ", err)
		return "", err
	}

	// Publish the media Url to the Redis channel
	if err := workflow.ExecuteActivity(ctx, activities.PublisherToChat, conf, redisChanel, answerObjectCompleted).Get(ctx, nil); err != nil {
		logger.Error("Activity PublisherToChat failed.", "Error : ", err)
		return "", err
	}

	// Clear the temporary file
	if err := workflow.ExecuteActivity(ctx, activities.ClearTemp, conf, fileName).Get(ctx, nil); err != nil {
		logger.Error("Activity ClearTemp failed.", "Error : ", err)
		return "", err
	}

	// Return the result
	return answerObj.Choices[0].Message.Content, nil
}
