package zess

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

func (z *Zess) getSeasons() (*[]*dto.Season, error) {
	zap.S().Info("Zess: Getting seasons")

	api := config.GetDefaultString("backend.zess.api", "https://zess.zeus.gent")
	req := fiber.Get(api + "/seasons")

	res := new([]*dto.Season)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return nil, errors.Join(append([]error{errors.New("Zess: Season API request failed")}, errs...)...)
	}
	if status != fiber.StatusOK {
		return nil, errors.New("error getting seasons")
	}

	errs = make([]error, 0)
	for _, season := range *res {
		if err := dto.Validate.Struct(season); err != nil {
			errs = append(errs, err)
		}
	}

	return res, errors.Join(errs...)
}

func (z *Zess) getScans() (*[]*dto.Scan, error) {
	zap.S().Info("Zess: Getting scans")

	api := config.GetDefaultString("backend.zess.api", "https://zess.zeus.gent")
	req := fiber.Get(api + "/recent_scans")

	res := new([]*dto.Scan)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return nil, errors.Join(append(errs, errors.New("Zess: Scan API request failed"))...)
	}
	if status != fiber.StatusOK {
		return nil, fmt.Errorf("Zess: Scan API returned bad status code %d", status)
	}

	errs = make([]error, 0)
	for _, scan := range *res {
		if err := dto.Validate.Struct(scan); err != nil {
			errs = append(errs, err)
		}
	}

	return res, errors.Join(errs...)
}
