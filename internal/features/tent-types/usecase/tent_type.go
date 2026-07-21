package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/tent-types/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/tent-types/dto"
	appErr "github.com/ramdhanrizkij/nest-glamping-api/internal/shared/errors"
)

type usecase struct {
	repo domain.Repository
}

func NewUsecase(repo domain.Repository) domain.Service {
	return &usecase{repo: repo}
}

func (u *usecase) Create(req dto.CreateTentTypeRequest) (*dto.TentTypeResponse, error) {
	tentType := &domain.TentType{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
		BasePrice:   req.BasePrice,
	}

	if err := u.repo.Create(tentType); err != nil {
		return nil, appErr.Internal("failed to create tent type")
	}

	if len(req.AmenityIDs) > 0 {
		amenityIDs := parseUUIDs(req.AmenityIDs)
		if err := u.repo.SetAmenities(tentType.ID, amenityIDs); err != nil {
			return nil, appErr.Internal("failed to set amenities")
		}
	}

	return &dto.TentTypeResponse{
		ID:          tentType.ID.String(),
		Name:        tentType.Name,
		Description: tentType.Description,
		Capacity:    tentType.Capacity,
		BasePrice:   tentType.BasePrice,
	}, nil
}

func (u *usecase) List() ([]dto.TentTypeResponse, error) {
	tentTypes, err := u.repo.FindAll()
	if err != nil {
		return nil, appErr.Internal("failed to list tent types")
	}

	var result []dto.TentTypeResponse
	for _, tt := range tentTypes {
		result = append(result, dto.TentTypeResponse{
			ID:          tt.ID.String(),
			Name:        tt.Name,
			Description: tt.Description,
			Capacity:    tt.Capacity,
			BasePrice:   tt.BasePrice,
		})
	}
	return result, nil
}

func (u *usecase) FindByID(id string) (*dto.TentTypeDetailResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, appErr.BadRequest("invalid tent type id")
	}

	tentType, err := u.repo.FindByID(uid)
	if err != nil {
		return nil, appErr.NotFound("tent type not found")
	}

	images, _ := u.repo.FindImagesByTentTypeID(uid)
	amenityIDs, _ := u.repo.FindAmenitiesByTentTypeID(uid)
	rates, _ := u.repo.FindRatesByTentTypeID(uid)

	var imageResponses []dto.TentTypeImageResponse
	for _, img := range images {
		imageResponses = append(imageResponses, dto.TentTypeImageResponse{
			ID:        img.ID.String(),
			ImageURL:  img.ImageURL,
			IsPrimary: img.IsPrimary,
		})
	}

	var amenityResponses []dto.TentTypeAmenityResponse
	for _, aid := range amenityIDs {
		amenityResponses = append(amenityResponses, dto.TentTypeAmenityResponse{
			AmenityID: aid.String(),
		})
	}

	var rateResponses []dto.TentTypeRateResponse
	for _, r := range rates {
		rateResponses = append(rateResponses, dto.TentTypeRateResponse{
			ID:            r.ID.String(),
			StartDate:     r.StartDate.Format("2006-01-02"),
			EndDate:       r.EndDate.Format("2006-01-02"),
			PricePerNight: r.PricePerNight,
			Description:   r.Description,
			IsActive:      r.IsActive,
		})
	}

	return &dto.TentTypeDetailResponse{
		TentTypeResponse: dto.TentTypeResponse{
			ID:          tentType.ID.String(),
			Name:        tentType.Name,
			Description: tentType.Description,
			Capacity:    tentType.Capacity,
			BasePrice:   tentType.BasePrice,
		},
		Images:    imageResponses,
		Amenities: amenityResponses,
		Rates:     rateResponses,
	}, nil
}

func (u *usecase) Update(id string, req dto.UpdateTentTypeRequest) (*dto.TentTypeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, appErr.BadRequest("invalid tent type id")
	}

	tentType, err := u.repo.FindByID(uid)
	if err != nil {
		return nil, appErr.NotFound("tent type not found")
	}

	if req.Name != "" {
		tentType.Name = req.Name
	}
	if req.Description != "" {
		tentType.Description = req.Description
	}
	if req.Capacity > 0 {
		tentType.Capacity = req.Capacity
	}
	if req.BasePrice > 0 {
		tentType.BasePrice = req.BasePrice
	}

	if err := u.repo.Update(tentType); err != nil {
		return nil, appErr.Internal("failed to update tent type")
	}

	if req.AmenityIDs != nil {
		amenityIDs := parseUUIDs(req.AmenityIDs)
		if err := u.repo.SetAmenities(tentType.ID, amenityIDs); err != nil {
			return nil, appErr.Internal("failed to set amenities")
		}
	}

	return &dto.TentTypeResponse{
		ID:          tentType.ID.String(),
		Name:        tentType.Name,
		Description: tentType.Description,
		Capacity:    tentType.Capacity,
		BasePrice:   tentType.BasePrice,
	}, nil
}

func (u *usecase) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return appErr.BadRequest("invalid tent type id")
	}

	if _, err := u.repo.FindByID(uid); err != nil {
		return appErr.NotFound("tent type not found")
	}

	return u.repo.Delete(uid)
}

// --- Images ---

func (u *usecase) AddImage(tentTypeID string, req dto.AddImageRequest) (*dto.TentTypeImageResponse, error) {
	ttUUID, err := uuid.Parse(tentTypeID)
	if err != nil {
		return nil, appErr.BadRequest("invalid tent type id")
	}

	if _, err := u.repo.FindByID(ttUUID); err != nil {
		return nil, appErr.NotFound("tent type not found")
	}

	if req.IsPrimary {
		if err := u.repo.ClearPrimaryImage(ttUUID); err != nil {
			return nil, appErr.Internal("failed to clear primary images")
		}
	}

	image := &domain.TentTypeImage{
		ID:         uuid.New(),
		TentTypeID: ttUUID,
		ImageURL:   req.ImageURL,
		IsPrimary:  req.IsPrimary,
	}

	if err := u.repo.CreateImage(image); err != nil {
		return nil, appErr.Internal("failed to add image")
	}

	return &dto.TentTypeImageResponse{
		ID:        image.ID.String(),
		ImageURL:  image.ImageURL,
		IsPrimary: image.IsPrimary,
	}, nil
}

func (u *usecase) DeleteImage(imageID string) error {
	imgUUID, err := uuid.Parse(imageID)
	if err != nil {
		return appErr.BadRequest("invalid image id")
	}

	if _, err := u.repo.FindImageByID(imgUUID); err != nil {
		return appErr.NotFound("image not found")
	}

	return u.repo.DeleteImage(imgUUID)
}

func (u *usecase) SetPrimaryImage(tentTypeID string, imageID string) error {
	ttUUID, err := uuid.Parse(tentTypeID)
	if err != nil {
		return appErr.BadRequest("invalid tent type id")
	}

	imgUUID, err := uuid.Parse(imageID)
	if err != nil {
		return appErr.BadRequest("invalid image id")
	}

	if _, err := u.repo.FindImageByID(imgUUID); err != nil {
		return appErr.NotFound("image not found")
	}

	if err := u.repo.ClearPrimaryImage(ttUUID); err != nil {
		return appErr.Internal("failed to clear primary images")
	}

	return u.repo.SetPrimaryImage(ttUUID, imgUUID)
}

// --- Rates ---

func (u *usecase) CreateRate(tentTypeID string, req dto.CreateRateRequest) (*dto.TentTypeRateResponse, error) {
	ttUUID, err := uuid.Parse(tentTypeID)
	if err != nil {
		return nil, appErr.BadRequest("invalid tent type id")
	}

	if _, err := u.repo.FindByID(ttUUID); err != nil {
		return nil, appErr.NotFound("tent type not found")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, appErr.BadRequest("invalid start_date format, use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, appErr.BadRequest("invalid end_date format, use YYYY-MM-DD")
	}

	if !endDate.After(startDate) {
		return nil, appErr.BadRequest("end_date must be after start_date")
	}

	if _, err := u.repo.FindOverlappingRate(ttUUID, startDate, endDate, uuid.Nil); err == nil {
		return nil, appErr.Conflict("rate overlaps with an existing active rate")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	rate := &domain.TentTypeRate{
		ID:            uuid.New(),
		TentTypeID:    ttUUID,
		StartDate:     startDate,
		EndDate:       endDate,
		PricePerNight: req.PricePerNight,
		Description:   req.Description,
		IsActive:      isActive,
	}

	if err := u.repo.CreateRate(rate); err != nil {
		return nil, appErr.Internal("failed to create rate")
	}

	return &dto.TentTypeRateResponse{
		ID:            rate.ID.String(),
		StartDate:     rate.StartDate.Format("2006-01-02"),
		EndDate:       rate.EndDate.Format("2006-01-02"),
		PricePerNight: rate.PricePerNight,
		Description:   rate.Description,
		IsActive:      rate.IsActive,
	}, nil
}

func (u *usecase) ListRates(tentTypeID string) ([]dto.TentTypeRateResponse, error) {
	ttUUID, err := uuid.Parse(tentTypeID)
	if err != nil {
		return nil, appErr.BadRequest("invalid tent type id")
	}

	rates, err := u.repo.FindRatesByTentTypeID(ttUUID)
	if err != nil {
		return nil, appErr.Internal("failed to list rates")
	}

	var result []dto.TentTypeRateResponse
	for _, r := range rates {
		result = append(result, dto.TentTypeRateResponse{
			ID:            r.ID.String(),
			StartDate:     r.StartDate.Format("2006-01-02"),
			EndDate:       r.EndDate.Format("2006-01-02"),
			PricePerNight: r.PricePerNight,
			Description:   r.Description,
			IsActive:      r.IsActive,
		})
	}
	return result, nil
}

func (u *usecase) UpdateRate(rateID string, req dto.UpdateRateRequest) (*dto.TentTypeRateResponse, error) {
	rUUID, err := uuid.Parse(rateID)
	if err != nil {
		return nil, appErr.BadRequest("invalid rate id")
	}

	rate, err := u.repo.FindRateByID(rUUID)
	if err != nil {
		return nil, appErr.NotFound("rate not found")
	}

	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, appErr.BadRequest("invalid start_date format, use YYYY-MM-DD")
		}
		rate.StartDate = startDate
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, appErr.BadRequest("invalid end_date format, use YYYY-MM-DD")
		}
		rate.EndDate = endDate
	}

	if !rate.EndDate.After(rate.StartDate) {
		return nil, appErr.BadRequest("end_date must be after start_date")
	}

	if _, err := u.repo.FindOverlappingRate(rate.TentTypeID, rate.StartDate, rate.EndDate, rUUID); err == nil {
		return nil, appErr.Conflict("rate overlaps with an existing active rate")
	}

	if req.PricePerNight > 0 {
		rate.PricePerNight = req.PricePerNight
	}
	if req.Description != "" {
		rate.Description = req.Description
	}
	if req.IsActive != nil {
		rate.IsActive = *req.IsActive
	}

	if err := u.repo.UpdateRate(rate); err != nil {
		return nil, appErr.Internal("failed to update rate")
	}

	return &dto.TentTypeRateResponse{
		ID:            rate.ID.String(),
		StartDate:     rate.StartDate.Format("2006-01-02"),
		EndDate:       rate.EndDate.Format("2006-01-02"),
		PricePerNight: rate.PricePerNight,
		Description:   rate.Description,
		IsActive:      rate.IsActive,
	}, nil
}

func (u *usecase) DeleteRate(rateID string) error {
	rUUID, err := uuid.Parse(rateID)
	if err != nil {
		return appErr.BadRequest("invalid rate id")
	}

	if _, err := u.repo.FindRateByID(rUUID); err != nil {
		return appErr.NotFound("rate not found")
	}

	return u.repo.DeleteRate(rUUID)
}

func parseUUIDs(ids []string) []uuid.UUID {
	var result []uuid.UUID
	for _, id := range ids {
		if uid, err := uuid.Parse(id); err == nil {
			result = append(result, uid)
		}
	}
	return result
}
