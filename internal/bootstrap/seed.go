package bootstrap

import (
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"github.com/ramdhanrizkij/nest-glamping-api/internal/bootstrap/migrations"
	"github.com/ramdhanrizkij/nest-glamping-api/pkg/hash"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	log.Info("Seeding database...")

	var roleAdmin migrations.Role
	db.Where("name = ?", "admin").First(&roleAdmin)

	var roleCustomer migrations.Role
	db.Where("name = ?", "customer").First(&roleCustomer)

	// --- Users ---
	users := []migrations.User{
		{
			ID:           uuid.New(),
			RoleID:       roleAdmin.ID,
			Name:         "Admin Utama",
			Email:        "admin@glamping.com",
			PasswordHash: mustHash("admin123"),
			PhoneNumber:  "081234567890",
		},
		{
			ID:           uuid.New(),
			RoleID:       roleCustomer.ID,
			Name:         "Budi Santoso",
			Email:        "budi@example.com",
			PasswordHash: mustHash("password123"),
			PhoneNumber:  "081234567891",
		},
		{
			ID:           uuid.New(),
			RoleID:       roleCustomer.ID,
			Name:         "Siti Rahayu",
			Email:        "siti@example.com",
			PasswordHash: mustHash("password123"),
			PhoneNumber:  "081234567892",
		},
		{
			ID:           uuid.New(),
			RoleID:       roleCustomer.ID,
			Name:         "Andi Wijaya",
			Email:        "andi@example.com",
			PasswordHash: mustHash("password123"),
			PhoneNumber:  "081234567893",
		},
	}
	seedUsers(db, users)

	// --- Amenities ---
	amenities := []migrations.Amenity{
		{Name: "WiFi Gratis", IconURL: "/icons/wifi.png", Description: "Akses internet nirkabel di area umum"},
		{Name: "AC", IconURL: "/icons/ac.png", Description: "Pendingin ruangan di dalam tenda"},
		{Name: "Kamar Mandi Dalam", IconURL: "/icons/bathroom.png", Description: "Kamar mandi pribadi dengan air panas"},
		{Name: "BBQ Area", IconURL: "/icons/bbq.png", Description: "Area barbekyu bersama"},
		{Name: "Spot Foto", IconURL: "/icons/camera.png", Description: "Area foto dengan view terbaik"},
		{Name: "Breakfast", IconURL: "/icons/breakfast.png", Description: "Sarapan pagi untuk 2 orang"},
		{Name: "Bonfire", IconURL: "/icons/fire.png", Description: "Api unggun malam hari"},
		{Name: "Parking", IconURL: "/icons/parking.png", Description: "Parkir kendaraan pribadi"},
	}
	amenityIDs := seedAmenities(db, amenities)

	// --- Tent Types ---
	tentTypes := []migrations.TentType{
		{Name: "Safari Tent", Description: "Tenda safari premium dengan view pegunungan. Cocok untuk pasangan atau keluarga kecil.", Capacity: 2, BasePrice: 750000},
		{Name: "Deluxe Camp", Description: "Tenda deluxe dengan fasilitas lengkap. Luas dan nyaman untuk liburan.", Capacity: 4, BasePrice: 1200000},
		{Name: "Family Glamp", Description: "Tenda besar untuk keluarga. Ruang luas dengan 2 kamar tidur.", Capacity: 6, BasePrice: 1800000},
		{Name: "Couple's Nest", Description: "Tenda romantis untuk honeymoon. Private dengan jacuzzi.", Capacity: 2, BasePrice: 1500000},
	}
	ttIDs := seedTentTypes(db, tentTypes)

	// --- Tent Type Amenities ---
	seedTentTypeAmenities(db, ttIDs, amenityIDs)

	// --- Tent Type Images ---
	images := []migrations.TentTypeImage{
		{TentTypeID: ttIDs[0], ImageURL: "/images/safari-1.jpg", IsPrimary: true},
		{TentTypeID: ttIDs[0], ImageURL: "/images/safari-2.jpg", IsPrimary: false},
		{TentTypeID: ttIDs[1], ImageURL: "/images/deluxe-1.jpg", IsPrimary: true},
		{TentTypeID: ttIDs[1], ImageURL: "/images/deluxe-2.jpg", IsPrimary: false},
		{TentTypeID: ttIDs[2], ImageURL: "/images/family-1.jpg", IsPrimary: true},
		{TentTypeID: ttIDs[3], ImageURL: "/images/couple-1.jpg", IsPrimary: true},
	}
	seedImages(db, images)

	// --- Tent Type Rates (Dynamic Pricing) ---
	now := time.Now()
	rates := []migrations.TentTypeRate{
		// High season Lebaran (misalnya Juni 2025)
		{TentTypeID: ttIDs[0], StartDate: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC), PricePerNight: 950000, Description: "High Season Lebaran", IsActive: true},
		{TentTypeID: ttIDs[1], StartDate: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC), PricePerNight: 1500000, Description: "High Season Lebaran", IsActive: true},
		// Promo Natal (Desember)
		{TentTypeID: ttIDs[0], StartDate: time.Date(2025, 12, 20, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC), PricePerNight: 650000, Description: "Promo Tahun Baru", IsActive: true},
		// Weekend surcharge (hanya contoh)
		{TentTypeID: ttIDs[2], StartDate: now.AddDate(0, 0, 30), EndDate: now.AddDate(0, 0, 32), PricePerNight: 2000000, Description: "Weekend Spesial", IsActive: true},
	}
	seedRates(db, rates)

	// --- Tents (Units) ---
	tents := []migrations.Tent{
		// Safari Tent (3 units)
		{TentTypeID: ttIDs[0], NameOrNum: "Safari-01", Status: "available"},
		{TentTypeID: ttIDs[0], NameOrNum: "Safari-02", Status: "available"},
		{TentTypeID: ttIDs[0], NameOrNum: "Safari-03", Status: "maintenance"},
		// Deluxe Camp (4 units)
		{TentTypeID: ttIDs[1], NameOrNum: "Deluxe-01", Status: "available"},
		{TentTypeID: ttIDs[1], NameOrNum: "Deluxe-02", Status: "available"},
		{TentTypeID: ttIDs[1], NameOrNum: "Deluxe-03", Status: "available"},
		{TentTypeID: ttIDs[1], NameOrNum: "Deluxe-04", Status: "available"},
		// Family Glamp (2 units)
		{TentTypeID: ttIDs[2], NameOrNum: "Family-01", Status: "available"},
		{TentTypeID: ttIDs[2], NameOrNum: "Family-02", Status: "available"},
		// Couple's Nest (2 units)
		{TentTypeID: ttIDs[3], NameOrNum: "Couple-01", Status: "available"},
		{TentTypeID: ttIDs[3], NameOrNum: "Couple-02", Status: "available"},
	}
	seedTents(db, tents)

	log.Info("Seeding completed")
	return nil
}

func mustHash(password string) string {
	h, _ := hash.HashPassword(password)
	return h
}

func seedUsers(db *gorm.DB, users []migrations.User) {
	for _, u := range users {
		var existing migrations.User
		if err := db.Where("email = ?", u.Email).First(&existing).Error; err != nil {
			db.Create(&u)
		}
	}
}

func seedAmenities(db *gorm.DB, amenities []migrations.Amenity) []uuid.UUID {
	var ids []uuid.UUID
	for _, a := range amenities {
		var existing migrations.Amenity
		if err := db.Where("name = ?", a.Name).First(&existing).Error; err != nil {
			db.Create(&a)
			ids = append(ids, a.ID)
		} else {
			ids = append(ids, existing.ID)
		}
	}
	return ids
}

func seedTentTypes(db *gorm.DB, tentTypes []migrations.TentType) []uuid.UUID {
	var ids []uuid.UUID
	for _, tt := range tentTypes {
		var existing migrations.TentType
		if err := db.Where("name = ?", tt.Name).First(&existing).Error; err != nil {
			db.Create(&tt)
			ids = append(ids, tt.ID)
		} else {
			ids = append(ids, existing.ID)
		}
	}
	return ids
}

func seedTentTypeAmenities(db *gorm.DB, ttIDs, amenityIDs []uuid.UUID) {
	for i, ttID := range ttIDs {
		for _, aID := range amenityIDs {
			var existing migrations.TentTypeAmenity
			if err := db.Where("tent_type_id = ? AND amenity_id = ?", ttID, aID).First(&existing).Error; err != nil {
				// Setiap tipe punya semua amenities (sederhana untuk dummy)
				db.Create(&migrations.TentTypeAmenity{
					ID:         uuid.New(),
					TentTypeID: ttID,
					AmenityID:  aID,
				})
			}
		}
		_ = i
	}
}

func seedImages(db *gorm.DB, images []migrations.TentTypeImage) {
	for _, img := range images {
		var existing migrations.TentTypeImage
		if err := db.Where("tent_type_id = ? AND image_url = ?", img.TentTypeID, img.ImageURL).First(&existing).Error; err != nil {
			db.Create(&img)
		}
	}
}

func seedRates(db *gorm.DB, rates []migrations.TentTypeRate) {
	for _, r := range rates {
		var existing migrations.TentTypeRate
		if err := db.Where("tent_type_id = ? AND start_date = ? AND end_date = ?", r.TentTypeID, r.StartDate, r.EndDate).First(&existing).Error; err != nil {
			db.Create(&r)
		}
	}
}

func seedTents(db *gorm.DB, tents []migrations.Tent) {
	for _, t := range tents {
		var existing migrations.Tent
		if err := db.Where("tent_type_id = ? AND name_or_number = ?", t.TentTypeID, t.NameOrNum).First(&existing).Error; err != nil {
			db.Create(&t)
		}
	}
}
