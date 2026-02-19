package ledstrip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

type Client struct {
	url string
}

const morseTopic = "morsemessage"

type ledMessage struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

func New() *Client {
	url := config.GetDefaultString("backend.ledstrip.url", "")

	if url == "" {
		zap.S().Info("No ledstrip url configured.\nMock messages will be used insteadd")
	}

	return &Client{
		url: url,
	}
}

func (c *Client) Flash(message model.Message) error {
	if c.url == "" {
		zap.S().Info("FLASH FLASH FLAAAAAAASH")
		return nil
	}

	msg := ledMessage{
		Topic:   morseTopic,
		Message: fmt.Sprintf("<%s> %s", message.Name, message.Message),
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(msg); err != nil {
		return fmt.Errorf("encode msg %w", err)
	}

	req, err := http.NewRequest("PUT", c.url, &buf)
	if err != nil {
		return fmt.Errorf("make request %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do req %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %s", resp.Status)
	}

	return nil
}
