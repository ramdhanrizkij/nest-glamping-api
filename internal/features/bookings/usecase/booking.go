package usecase

import (
	"time"

	"github.com/google/uuid"
	bookingDomain "github.com/ramdhanrizkij/nest-glamping-api/internal/features/bookings/domain"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/features/bookings/dto"
	tentTypeDomain "github.com/ramdhanrizkij/nest-glamping-api/internal/features/tent-types/domain"
	tentDomain "github.com/ramdhanrizkij/nest-glamping-api/internal/features/tents/domain"
	appErr "github.com/ramdhanrizkij/nest-glamping-api/internal/shared/errors"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/shared/utils"
)

type usecase struct {
	bookingRepo bookingDomain.Repository
	tentRepo    tentDomain.Repository
	tentTypeRepo tentTypeDomain.Repository
}

func NewUsecase(
	bookingRepo bookingDomain.Repository,
	tentRepo tentDomain.Repository,
	tentTypeRepo tentTypeDomain.Repository,
) bookingDomain.Service {
	return &usecase{
		bookingRepo:  bookingRepo,
		tentRepo:     tentRepo,
		tentTypeRepo: tentTypeRepo,
	}
}

func (u *usecase) CreateBooking(userID string, req dto.CreateBookingRequest) (*dto.BookingResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, appErr.BadRequest("invalid user id")
	}

	tentTypeUUID, err := uuid.Parse(req.TentTypeID)
	if err != nil {
		return nil, appErr.BadRequest("invalid tent type id")
	}

	checkIn, err := time.Parse("2006-01-02", req.CheckInDate)
	if err != nil {
		return nil, appErr.BadRequest("invalid check_in_date format, use YYYY-MM-DD")
	}

	checkOut, err := time.Parse("2006-01-02", req.CheckOutDate)
	if err != nil {
		return nil, appErr.BadRequest("invalid check_out_date format, use YYYY-MM-DD")
	}

	if !checkOut.After(checkIn) {
		return nil, appErr.BadRequest("check_out_date must be after check_in_date")
	}

	if req.GuestCount < 1 {
		return nil, appErr.BadRequest("guest_count must be at least 1")
	}

	tentType, err := u.tentTypeRepo.FindByID(tentTypeUUID)
	if err != nil {
		return nil, appErr.NotFound("tent type not found")
	}

	if req.GuestCount > tentType.Capacity {
		return nil, appErr.BadRequest("guest_count exceeds tent capacity")
	}

	// Validate all tent IDs belong to this tent type and are available
	var tentIDs []uuid.UUID
	for _, tid := range req.TentIDs {
		tUUID, err := uuid.Parse(tid)
		if err != nil {
			return nil, appErr.BadRequest("invalid tent id: " + tid)
		}
		tent, err := u.tentRepo.FindByID(tUUID)
		if err != nil {
			return nil, appErr.NotFound("tent not found: " + tid)
		}
		if tent.TentTypeID != tentTypeUUID {
			return nil, appErr.BadRequest("tent " + tid + " does not belong to this tent type")
		}
		if tent.Status != "available" {
			return nil, appErr.BadRequest("tent " + tid + " is not available")
		}
		tentIDs = append(tentIDs, tUUID)
	}

	// Re-check availability (prevent race condition)
	availableTents, err := u.tentRepo.FindAvailableTents(tentTypeUUID, checkIn, checkOut)
	if err != nil {
		return nil, appErr.Internal("failed to check availability")
	}

	availableMap := make(map[uuid.UUID]bool)
	for _, t := range availableTents {
		availableMap[t.ID] = true
	}

	for _, tid := range tentIDs {
		if !availableMap[tid] {
			return nil, appErr.Conflict("tent is no longer available for the selected dates")
		}
	}

	// Calculate price per night for each tent (dynamic pricing)
	rates, _ := u.tentTypeRepo.FindRatesByTentTypeID(tentTypeUUID)
	nights := int(checkOut.Sub(checkIn).Hours() / 24)

	var bookingTents []bookingDomain.BookingTent
	var totalAmount float64

	for _, tentID := range tentIDs {
		var tentTotal float64

		for i := 0; i < nights; i++ {
			currentDate := checkIn.AddDate(0, 0, i+1)
			nightPrice := tentType.BasePrice
			for _, rate := range rates {
				if rate.IsActive && !currentDate.Before(rate.StartDate) && !currentDate.After(rate.EndDate) {
					nightPrice = rate.PricePerNight
					break
				}
			}
			tentTotal += nightPrice
		}

		// Snapshot: price_per_night = total for this tent / nights (averaged)
		snapshotPrice := tentTotal / float64(nights)

		bookingTents = append(bookingTents, bookingDomain.BookingTent{
			ID:            uuid.New(),
			TentID:        tentID,
			PricePerNight: snapshotPrice,
		})

		totalAmount += tentTotal
	}

	bookingCode := utils.GenerateBookingCode()

	booking := &bookingDomain.Booking{
		ID:              uuid.New(),
		UserID:          userUUID,
		BookingCode:     bookingCode,
		CheckInDate:     checkIn,
		CheckOutDate:    checkOut,
		TotalAmount:     totalAmount,
		Status:          "pending",
		GuestCount:      req.GuestCount,
		SpecialRequests: req.SpecialRequests,
	}

	// Set booking_id on booking tents
	for i := range bookingTents {
		bookingTents[i].BookingID = booking.ID
	}

	if err := u.bookingRepo.CreateBooking(booking, bookingTents); err != nil {
		return nil, appErr.Internal("failed to create booking")
	}

	return u.toBookingResponse(booking, bookingTents), nil
}

func (u *usecase) ListMyBookings(userID string, page, perPage int) (*dto.BookingListResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, appErr.BadRequest("invalid user id")
	}

	bookings, err := u.bookingRepo.FindUserBookingsWithTents(userUUID)
	if err != nil {
		return nil, appErr.Internal("failed to list bookings")
	}

	total := int64(len(bookings))
	offset := (page - 1) * perPage
	if offset >= len(bookings) {
		offset = len(bookings)
	}
	end := offset + perPage
	if end > len(bookings) {
		end = len(bookings)
	}
	pageBookings := bookings[offset:end]

	var result []dto.BookingResponse
	for _, b := range pageBookings {
		tents, _ := u.bookingRepo.FindBookingTentsByBookingID(b.ID)
		result = append(result, *u.toBookingResponse(&b, tents))
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return &dto.BookingListResponse{
		Data:       result,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (u *usecase) GetBookingDetail(bookingID, userID string, isAdmin bool) (*dto.BookingDetailResponse, error) {
	bookingUUID, err := uuid.Parse(bookingID)
	if err != nil {
		return nil, appErr.BadRequest("invalid booking id")
	}

	booking, bookingTents, err := u.bookingRepo.FindBookingByIDWithTents(bookingUUID)
	if err != nil {
		return nil, appErr.NotFound("booking not found")
	}

	if !isAdmin {
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return nil, appErr.BadRequest("invalid user id")
		}
		if booking.UserID != userUUID {
			return nil, appErr.NotFound("booking not found")
		}
	}

	resp := u.toBookingResponse(booking, bookingTents)
	return &dto.BookingDetailResponse{
		BookingResponse: *resp,
		CreatedAt:       booking.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (u *usecase) CancelBooking(bookingID, userID string) error {
	bookingUUID, err := uuid.Parse(bookingID)
	if err != nil {
		return appErr.BadRequest("invalid booking id")
	}

	booking, err := u.bookingRepo.FindBookingByID(bookingUUID)
	if err != nil {
		return appErr.NotFound("booking not found")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return appErr.BadRequest("invalid user id")
	}

	if booking.UserID != userUUID {
		return appErr.NotFound("booking not found")
	}

	if booking.Status != "pending" {
		return appErr.BadRequest("can only cancel pending bookings")
	}

	return u.bookingRepo.UpdateBookingStatus(bookingUUID, "cancelled")
}

func (u *usecase) ListAllBookings(page, perPage int) (*dto.BookingListResponse, error) {
	bookings, err := u.bookingRepo.FindAllBookingsWithTents()
	if err != nil {
		return nil, appErr.Internal("failed to list bookings")
	}

	total := int64(len(bookings))
	offset := (page - 1) * perPage
	if offset >= len(bookings) {
		offset = len(bookings)
	}
	end := offset + perPage
	if end > len(bookings) {
		end = len(bookings)
	}
	pageBookings := bookings[offset:end]

	var result []dto.BookingResponse
	for _, b := range pageBookings {
		tents, _ := u.bookingRepo.FindBookingTentsByBookingID(b.ID)
		result = append(result, *u.toBookingResponse(&b, tents))
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return &dto.BookingListResponse{
		Data:       result,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (u *usecase) ConfirmBooking(bookingID string) error {
	bookingUUID, err := uuid.Parse(bookingID)
	if err != nil {
		return appErr.BadRequest("invalid booking id")
	}

	booking, err := u.bookingRepo.FindBookingByID(bookingUUID)
	if err != nil {
		return appErr.NotFound("booking not found")
	}

	if booking.Status != "pending" {
		return appErr.BadRequest("can only confirm pending bookings")
	}

	return u.bookingRepo.UpdateBookingStatus(bookingUUID, "confirmed")
}

func (u *usecase) toBookingResponse(booking *bookingDomain.Booking, bookingTents []bookingDomain.BookingTent) *dto.BookingResponse {
	var tents []dto.BookingTentItem
	for _, bt := range bookingTents {
		tent, _ := u.tentRepo.FindByID(bt.TentID)
		tentName := ""
		if tent != nil {
			tentName = tent.NameOrNum
		}
		tents = append(tents, dto.BookingTentItem{
			ID:            bt.ID.String(),
			TentID:        bt.TentID.String(),
			TentName:      tentName,
			PricePerNight: bt.PricePerNight,
		})
	}

	return &dto.BookingResponse{
		ID:              booking.ID.String(),
		BookingCode:     booking.BookingCode,
		CheckInDate:     booking.CheckInDate.Format("2006-01-02"),
		CheckOutDate:    booking.CheckOutDate.Format("2006-01-02"),
		TotalAmount:     booking.TotalAmount,
		Status:          booking.Status,
		GuestCount:      booking.GuestCount,
		SpecialRequests: booking.SpecialRequests,
		Tents:           tents,
	}
}
