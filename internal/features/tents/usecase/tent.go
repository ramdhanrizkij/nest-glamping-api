package usecase

import (
	"time"

	"github.com/google/uuid"
	tentTypeDomain "github.com/ramdhanrizkij/nest-glamping-api/internal/features/tent-types/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/tents/dto"
	tentDomain "github.com/ramdhanrizkij/nest-glamping-api/internal/features/tents/domain"
	appErr "github.com/ramdhanrizkij/nest-glamping-api/internal/shared/errors"
)

type usecase struct {
	tentRepo     tentDomain.Repository
	tentTypeRepo tentTypeDomain.Repository
}

func NewUsecase(tentRepo tentDomain.Repository, tentTypeRepo tentTypeDomain.Repository) tentDomain.Service {
	return &usecase{tentRepo: tentRepo, tentTypeRepo: tentTypeRepo}
}

func (u *usecase) Create(req dto.CreateTentRequest) (*dto.TentResponse, error) {
	ttUUID, err := uuid.Parse(req.TentTypeID)
	if err != nil {
		return nil, appErr.BadRequest("invalid tent type id")
	}

	if _, err := u.tentTypeRepo.FindByID(ttUUID); err != nil {
		return nil, appErr.NotFound("tent type not found")
	}

	status := "available"
	if req.Status != "" {
		status = req.Status
	}

	tent := &tentDomain.Tent{
		ID:         uuid.New(),
		TentTypeID: ttUUID,
		NameOrNum:  req.NameOrNum,
		Status:     status,
	}

	if err := u.tentRepo.Create(tent); err != nil {
		return nil, appErr.Internal("failed to create tent")
	}

	return &dto.TentResponse{
		ID:         tent.ID.String(),
		TentTypeID: tent.TentTypeID.String(),
		NameOrNum:  tent.NameOrNum,
		Status:     tent.Status,
	}, nil
}

func (u *usecase) List() ([]dto.TentResponse, error) {
	tents, err := u.tentRepo.FindAll()
	if err != nil {
		return nil, appErr.Internal("failed to list tents")
	}

	var result []dto.TentResponse
	for _, t := range tents {
		result = append(result, dto.TentResponse{
			ID:         t.ID.String(),
			TentTypeID: t.TentTypeID.String(),
			NameOrNum:  t.NameOrNum,
			Status:     t.Status,
		})
	}
	return result, nil
}

func (u *usecase) ListByTentTypeID(tentTypeID string) ([]dto.TentResponse, error) {
	ttUUID, err := uuid.Parse(tentTypeID)
	if err != nil {
		return nil, appErr.BadRequest("invalid tent type id")
	}

	tents, err := u.tentRepo.FindByTentTypeID(ttUUID)
	if err != nil {
		return nil, appErr.Internal("failed to list tents")
	}

	var result []dto.TentResponse
	for _, t := range tents {
		result = append(result, dto.TentResponse{
			ID:         t.ID.String(),
			TentTypeID: t.TentTypeID.String(),
			NameOrNum:  t.NameOrNum,
			Status:     t.Status,
		})
	}
	return result, nil
}

func (u *usecase) FindByID(id string) (*dto.TentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, appErr.BadRequest("invalid tent id")
	}

	tent, err := u.tentRepo.FindByID(uid)
	if err != nil {
		return nil, appErr.NotFound("tent not found")
	}

	return &dto.TentResponse{
		ID:         tent.ID.String(),
		TentTypeID: tent.TentTypeID.String(),
		NameOrNum:  tent.NameOrNum,
		Status:     tent.Status,
	}, nil
}

func (u *usecase) Update(id string, req dto.UpdateTentRequest) (*dto.TentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, appErr.BadRequest("invalid tent id")
	}

	tent, err := u.tentRepo.FindByID(uid)
	if err != nil {
		return nil, appErr.NotFound("tent not found")
	}

	if req.NameOrNum != "" {
		tent.NameOrNum = req.NameOrNum
	}
	if req.Status != "" {
		tent.Status = req.Status
	}

	if err := u.tentRepo.Update(tent); err != nil {
		return nil, appErr.Internal("failed to update tent")
	}

	return &dto.TentResponse{
		ID:         tent.ID.String(),
		TentTypeID: tent.TentTypeID.String(),
		NameOrNum:  tent.NameOrNum,
		Status:     tent.Status,
	}, nil
}

func (u *usecase) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return appErr.BadRequest("invalid tent id")
	}

	if _, err := u.tentRepo.FindByID(uid); err != nil {
		return appErr.NotFound("tent not found")
	}

	return u.tentRepo.Delete(uid)
}

func (u *usecase) CheckAvailability(tentTypeID string, checkIn, checkOut time.Time) ([]dto.AvailableTentResponse, error) {
	ttUUID, err := uuid.Parse(tentTypeID)
	if err != nil {
		return nil, appErr.BadRequest("invalid tent type id")
	}

	tentType, err := u.tentTypeRepo.FindByID(ttUUID)
	if err != nil {
		return nil, appErr.NotFound("tent type not found")
	}

	tents, err := u.tentRepo.FindAvailableTents(ttUUID, checkIn, checkOut)
	if err != nil {
		return nil, appErr.Internal("failed to check availability")
	}

	rates, _ := u.tentTypeRepo.FindRatesByTentTypeID(ttUUID)

	var result []dto.AvailableTentResponse
	for _, tent := range tents {
		nights := int(checkOut.Sub(checkIn).Hours() / 24)
		totalPrice := 0.0

		for i := 0; i < nights; i++ {
			currentDate := checkIn.AddDate(0, 0, i+1)
			nightPrice := tentType.BasePrice
			for _, rate := range rates {
				if rate.IsActive && !currentDate.Before(rate.StartDate) && !currentDate.After(rate.EndDate) {
					nightPrice = rate.PricePerNight
					break
				}
			}
			totalPrice += nightPrice
		}

		result = append(result, dto.AvailableTentResponse{
			ID:            tent.ID.String(),
			NameOrNum:     tent.NameOrNum,
			PricePerNight: totalPrice / float64(nights),
		})
	}

	return result, nil
}
