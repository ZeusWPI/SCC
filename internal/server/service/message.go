package service

import (
	"context"
	"math"
	"slices"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/buzzer"
	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/internal/ledstrip"
	"github.com/zeusWPI/scc/internal/server/dto"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

type Message struct {
	service Service

	message repository.Message
	reply   repository.Reply

	mu         sync.RWMutex
	idToClient map[int]*websocket.Conn

	buzzer    buzzer.Client
	blacklist []string
	ledstrip  ledstrip.Client
}

var messageSingleton *Message

func (s *Service) NewMessage() *Message {
	if messageSingleton == nil {
		messageSingleton = &Message{
			service:    *s,
			message:    *s.repo.NewMessage(),
			reply:      *s.repo.NewReply(),
			idToClient: map[int]*websocket.Conn{},
			buzzer:     *buzzer.New(),
			blacklist:  config.GetDefaultStringSlice("cammie.blacklist", []string{}),
			ledstrip:   *ledstrip.New(),
		}
	}

	return messageSingleton
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
	for _, msg := range messagesDB {
		if msg.ID < minMsgID {
			minMsgID = msg.ID
		}
		messages = append(messages, dto.MessageDTO(msg))
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

	for _, msg := range messages {
		// Already a cluster?
		if idx, ok := lastClusterIdxByAuthor[msg.Name]; ok {
			c := clusters[idx]
			lastMsg := c.Messages[len(c.Messages)-1]

			delta := msg.SendAt.Sub(lastMsg.SendAt)
			if delta >= 0 && delta <= gap {
				c.Messages = append(c.Messages, msg)
				if _, ok := m.idToClient[msg.ID]; ok {
					c.Connected = true
				}

				clusters[idx] = c
				continue
			}
		}

		// New cluster
		cluster := dto.MessageCluster{
			Messages: []dto.Message{msg},
		}
		if _, ok := m.idToClient[msg.ID]; ok {
			cluster.Connected = true
		}

		clusters = append(clusters, cluster)
		lastClusterIdxByAuthor[msg.Name] = len(clusters) - 1
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

	// Reverse the days
	slices.Reverse(days)
	// Reverse the clusters in days
	for idx := range days {
		slices.Reverse(days[idx].Clusters)
	}

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

func (m *Message) Create(ctx context.Context, msgSave dto.MessageSave, conn *websocket.Conn) (dto.Message, error) {
	msg := msgSave.ToModel()
	if err := m.message.Create(ctx, msg); err != nil {
		zap.S().Error(err)
		return dto.Message{}, fiber.ErrInternalServerError
	}

	if conn != nil {
		m.mu.Lock()
		m.idToClient[msg.ID] = conn
		m.mu.Unlock()
	}

	if !slices.Contains(m.blacklist, msg.Name) {
		go m.buzzer.Play()
	}

	_ = m.ledstrip.Flash(*msg)

	return dto.MessageDTO(msg), nil
}

func (m *Message) Reply(ctx context.Context, replySave dto.ReplySave) (dto.Reply, error) {
	msg, err := m.message.Get(ctx, replySave.MessageID)
	if err != nil {
		zap.S().Error(err)
		return dto.Reply{}, fiber.ErrInternalServerError
	}
	if msg == nil {
		return dto.Reply{}, fiber.ErrNotFound
	}

	m.mu.Lock()
	conn, ok := m.idToClient[msg.ID]
	m.mu.Unlock()
	if !ok {
		return dto.Reply{}, fiber.ErrGone
	}

	reply := replySave.ToModel()
	if err := m.reply.Create(ctx, reply); err != nil {
		zap.S().Error(err)
		return dto.Reply{}, fiber.ErrInternalServerError
	}

	if conn != nil {
		_ = conn.WriteJSON(dto.WSFrame{
			Event: "replymessage",
			Data: map[string]any{
				"name":       reply.Name,
				"message":    reply.Message,
				"id":         reply.ID,
				"message_id": msg.ID,
			},
		})
	}

	return dto.ReplyDTO(reply), nil
}

func (m *Message) ListenerRemove(conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range m.idToClient {
		if v == conn {
			delete(m.idToClient, k)
		}
	}
}
