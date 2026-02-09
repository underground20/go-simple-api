package notification

import (
	"app/internal/config"
	"app/internal/employee/messages"
	kafkaApi "app/lib/kafka"
	formatLog "app/lib/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func ConfigureConsumer(logger *slog.Logger, cfg *config.Config) *kafkaApi.Consumer {
	consumer := kafkaApi.NewConsumer(
		func(m kafka.Message) {
			handleMessage(logger, m, cfg)
		},
		logger,
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		"notification",
	)

	return consumer
}

func handleMessage(logger *slog.Logger, m kafka.Message, cfg *config.Config) {
	logger.Info(
		"Received kafka message",
		"topic", m.Topic,
		"partition", m.Partition,
		"offset", m.Offset,
	)

	var employeeCreatedInfo messages.EmployeeCreated
	err := json.Unmarshal(m.Value, &employeeCreatedInfo)
	if err != nil {
		logger.Error("Failed to unmarshal kafka message", formatLog.Err(err))
		return
	}

	message := fmt.Sprintf("%s\n%s\n%s\n%s", "Employee created",
		"Name: "+employeeCreatedInfo.Name,
		"Age: "+strconv.Itoa(employeeCreatedInfo.Age),
		"Position: "+employeeCreatedInfo.Position,
	)

	if cfg.Telegram.Token == "" || cfg.Telegram.ChatId == "" {
		logger.Info("Telegram token or chat id is empty")
		return
	}

	err = sendTelegramMessage(cfg.Telegram.Token, cfg.Telegram.ChatId, message)
	if err != nil {
		logger.Error("Failed to send telegram message", formatLog.Err(err))
	}
}

func sendTelegramMessage(token, chatID, message string) error {
	url := "https://api.telegram.org/bot" + token + "/sendMessage"
	data := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api error: %d", resp.StatusCode)
	}

	return nil
}
