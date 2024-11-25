package gamification

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"go.uber.org/zap"
)

func (g *Gamification) getLeaderboard() (*[]*dto.Gamification, error) {
	zap.S().Info("Gamification: Getting leaderboard")

	req := fiber.Get(g.api+"/top4").Set("Accept", "application/json")

	res := new([]*dto.Gamification)
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

	return res, errors.Join(errs...)
}

func downloadAvatar(gam dto.Gamification) (string, error) {
	req := fiber.Get(gam.Avatar)
	status, body, errs := req.Bytes()
	if errs != nil {
		return "", errors.Join(append(errs, errors.New("Gamification: Download avatar request failed"))...)
	}
	if status != fiber.StatusOK {
		return "", fmt.Errorf("Gamification: Download avatar returned bad status code %d", status)
	}

	location := fmt.Sprintf("public/%s.png", gam.Name)
	out, err := os.Create(location)
	if err != nil && err != os.ErrExist {
		return "", err
	}
	defer func() {
		_ = out.Close()
	}()

	_, err = io.Copy(out, bytes.NewReader(body))

	return location, err
}
