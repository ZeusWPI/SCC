package service

import (
	"context"
	"math"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/buzzer"
	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/internal/server/dto"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

type Message struct {
	message repository.Message
	reply   repository.Reply

	buzzer    buzzer.Client
	blacklist []string
}

func (s *Service) NewMessage() *Message {
	return &Message{
		message:   *s.repo.NewMessage(),
		reply:     *s.repo.NewReply(),
		buzzer:    *buzzer.New(),
		blacklist: config.GetDefaultStringSlice("cammie.blacklist", []string{}),
	}
}

func (m *Message) Get(ctx context.Context, sinceID int, dayLimit int) ([]dto.MessageDayGroup, error) {
	messagesDB, err := m.message.GetSinceID(ctx, sinceID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if len(messagesDB) == 0 {
		return nil, nil
	}

	minMsgID := math.MaxInt
	messages := make([]dto.Message, 0, len(messagesDB))
	for _, m := range messagesDB {
		if m.ID < minMsgID {
			minMsgID = m.ID
		}
		messages = append(messages, dto.MessageDTO(m))
	}

	repliesDB, err := m.reply.GetSinceMessageID(ctx, minMsgID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	byID := make(map[int]*dto.Message, len(messages))
	for i := range messages {
		byID[messages[i].ID] = &messages[i]
	}
	for _, r := range repliesDB {
		if msg := byID[r.MessageID]; msg != nil {
			msg.Replies = append(msg.Replies, dto.ReplyDTO(r))
		}
	}

	slices.SortFunc(messages, func(a, b dto.Message) int { return a.SendAt.Compare(b.SendAt) })

	const gap = 10 * time.Minute

	clusters := []dto.MessageCluster{}
	lastClusterIdxByAuthor := map[string]int{}

	for _, m := range messages {
		// Already a cluster?
		if idx, ok := lastClusterIdxByAuthor[m.Name]; ok {
			c := clusters[idx]
			lastMsg := c.Messages[len(c.Messages)-1]

			delta := m.SendAt.Sub(lastMsg.SendAt)
			if delta >= 0 && delta <= gap {
				c.Messages = append(c.Messages, m)
				clusters[idx] = c
				continue
			}
		}

		// New cluster
		clusters = append(clusters, dto.MessageCluster{Messages: []dto.Message{m}})
		lastClusterIdxByAuthor[m.Name] = len(clusters) - 1
	}

	days := []dto.MessageDayGroup{}
	dayIdx := map[string]int{}

	for _, c := range clusters {
		first := c.Messages[0]
		key := first.SendAt.Format("2006-01-02")

		if i, ok := dayIdx[key]; ok {
			days[i].Clusters = append(days[i].Clusters, c)
			continue
		}

		dayIdx[key] = len(days)
		days = append(days, dto.MessageDayGroup{
			DateKey:   key,
			DateLabel: first.SendAt.Format("Mon, 02 Jan 2006"),
			Clusters:  []dto.MessageCluster{c},
		})
	}

	if dayLimit > 0 && len(days) > dayLimit {
		days = days[len(days)-dayLimit:]
	}

	slices.Reverse(days)

	return days, nil
}

func (m *Message) GetLast(ctx context.Context) (dto.Message, error) {
	msg, err := m.message.GetLast(ctx)
	if err != nil {
		zap.S().Error(err)
		return dto.Message{}, fiber.ErrInternalServerError
	}
	if msg == nil {
		return dto.Message{}, fiber.ErrNotFound
	}

	return dto.MessageDTO(msg), nil
}

func (m *Message) Create(ctx context.Context, msgSave dto.MessageSave) (dto.Message, error) {
	msg := msgSave.ToModel()
	if err := m.message.Create(ctx, msg); err != nil {
		zap.S().Error(err)
		return dto.Message{}, fiber.ErrInternalServerError
	}

	if !slices.Contains(m.blacklist, msg.Name) {
		m.buzzer.Play()
	}

	return dto.MessageDTO(msg), nil
}
