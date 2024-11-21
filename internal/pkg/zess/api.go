package zess

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"go.uber.org/zap"
)

func (z *Zess) getSeasons() (*[]*dto.Season, error) {
	zap.S().Info("Zess: Getting seasons")

	req := fiber.Get(z.api + "/seasons")

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

	req := fiber.Get(z.api + "/recent_scans")

	res := new([]*dto.Scan)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return nil, errors.Join(append([]error{errors.New("Zess: Scan API request failed")}, errs...)...)
	}
	if status != fiber.StatusOK {
		return nil, errors.New("error getting scans")
	}

	errs = make([]error, 0)
	for _, scan := range *res {
		if err := dto.Validate.Struct(scan); err != nil {
			errs = append(errs, err)
		}
	}

	return res, errors.Join(errs...)
}
