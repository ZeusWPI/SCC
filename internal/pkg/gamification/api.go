package gamification

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"go.uber.org/zap"
)

type gamificationItem struct {
	ID        int32  `json:"id"`
	Name      string `json:"github_name"`
	Score     int32  `json:"score"`
	AvatarURL string `json:"avatar_url"`
}

func (g *Gamification) getLeaderboard() ([]dto.Gamification, error) {
	zap.S().Info("Gamification: Getting leaderboard")

	req := fiber.Get(g.api+"/top4").Set("Accept", "application/json")

	res := new([]gamificationItem)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return nil, errors.Join(append(errs, errors.New("Gamification: Leaderboard API request failed"))...)
	}
	if status != fiber.StatusOK {
		return nil, fmt.Errorf("Gamification: Leaderboard API request returned bad status code %d", status)
	}

	errs = make([]error, 0)
	for _, gam := range *res {
		if err := dto.Validate.Struct(gam); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	gams := make([]dto.Gamification, 0, 4)
	for _, item := range *res {
		gam, err := downloadAvatar(item)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		gams = append(gams, gam)
	}

	return gams, errors.Join(errs...)
}

func downloadAvatar(gam gamificationItem) (dto.Gamification, error) {
	zap.S().Info("Gamification: Getting avatar for ", gam.Name)

	req := fiber.Get(gam.AvatarURL)
	status, body, errs := req.Bytes()
	if len(errs) != 0 {
		return dto.Gamification{}, errors.Join(append(errs, errors.New("Gamification: Download avatar request failed"))...)
	}
	if status != fiber.StatusOK {
		return dto.Gamification{}, fmt.Errorf("Gamification: Download avatar returned bad status code %d", status)
	}

	g := dto.Gamification{
		ID:     gam.ID,
		Name:   gam.Name,
		Score:  gam.Score,
		Avatar: body,
	}

	return g, nil
}
