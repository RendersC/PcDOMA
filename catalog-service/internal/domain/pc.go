package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PCStatus string

const (
	StatusAvailable   PCStatus = "available"
	StatusBooked      PCStatus = "booked"
	StatusMaintenance PCStatus = "maintenance"
)

type Location struct {
	District    string  `bson:"district" json:"district"`
	Address     string  `bson:"address" json:"address"`
	Lat         float64 `bson:"lat" json:"lat"`
	Lng         float64 `bson:"lng" json:"lng"`
}

type PCSpecs struct {
	CPU     string `bson:"cpu" json:"cpu"`
	GPU     string `bson:"gpu" json:"gpu"`
	RAM     string `bson:"ram" json:"ram"`
	Storage string `bson:"storage" json:"storage"`
}

type PC struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Slug         string             `bson:"slug" json:"slug"`
	Specs        PCSpecs            `bson:"specs" json:"specs"`
	PricePerHour float64            `bson:"price_per_hour" json:"price_per_hour"`
	PricePerDay  float64            `bson:"price_per_day" json:"price_per_day"`
	Currency     string             `bson:"currency" json:"currency"`
	Location     Location           `bson:"location" json:"location"`
	Status       PCStatus           `bson:"status" json:"status"`
	Images       []string           `bson:"images" json:"images"`
	RatingAvg    float64            `bson:"rating_avg" json:"rating_avg"`
	RatingCount  int                `bson:"rating_count" json:"rating_count"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type PeripheralType string

const (
	PeripheralMouse      PeripheralType = "mouse"
	PeripheralHeadphones PeripheralType = "headphones"
	PeripheralHeadset    PeripheralType = "headset"
	PeripheralMonitor    PeripheralType = "monitor"
	PeripheralKeyboard   PeripheralType = "keyboard"
)

type Peripheral struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Type         PeripheralType     `bson:"type" json:"type"`
	Description  string             `bson:"description" json:"description"`
	Specs        map[string]string  `bson:"specs" json:"specs"`
	PricePerHour float64            `bson:"price_per_hour" json:"price_per_hour"`
	PricePerDay  float64            `bson:"price_per_day" json:"price_per_day"`
	Currency     string             `bson:"currency" json:"currency"`
	LocationID   string             `bson:"location_id" json:"location_id"`
	Status       PCStatus           `bson:"status" json:"status"`
	Images       []string           `bson:"images" json:"images"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type SetupCategory string

const (
	SetupGaming    SetupCategory = "gaming"
	SetupWork      SetupCategory = "work"
	SetupDesign    SetupCategory = "design"
	SetupStreaming  SetupCategory = "streaming"
)

type Setup struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Category         SetupCategory      `bson:"category" json:"category"`
	PCID             string             `bson:"pc_id" json:"pc_id"`
	PeripheralIDs    []string           `bson:"peripheral_ids" json:"peripheral_ids"`
	TotalPricePerHour float64           `bson:"total_price_per_hour" json:"total_price_per_hour"`
	TotalPricePerDay  float64           `bson:"total_price_per_day" json:"total_price_per_day"`
	DiscountPercent  float64            `bson:"discount_percent" json:"discount_percent"`
	Description      string             `bson:"description" json:"description"`
	Images           []string           `bson:"images" json:"images"`
	RatingAvg        float64            `bson:"rating_avg" json:"rating_avg"`
	RatingCount      int                `bson:"rating_count" json:"rating_count"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

// Request/Response types

type CreatePCRequest struct {
	Name         string   `json:"name" binding:"required"`
	Specs        PCSpecs  `json:"specs" binding:"required"`
	PricePerHour float64  `json:"price_per_hour" binding:"required,gt=0"`
	PricePerDay  float64  `json:"price_per_day" binding:"required,gt=0"`
	Location     Location `json:"location" binding:"required"`
	Images       []string `json:"images"`
}

type UpdatePCRequest struct {
	Name         string   `json:"name"`
	Specs        *PCSpecs `json:"specs"`
	PricePerHour float64  `json:"price_per_hour"`
	PricePerDay  float64  `json:"price_per_day"`
	Location     *Location `json:"location"`
	Images       []string `json:"images"`
}

type PCFilter struct {
	Status   string
	Location string
	MinPrice float64
	MaxPrice float64
	Sort     string
	Limit    int
	Offset   int
}

type CreatePeripheralRequest struct {
	Name         string            `json:"name" binding:"required"`
	Type         PeripheralType    `json:"type" binding:"required"`
	Description  string            `json:"description"`
	Specs        map[string]string `json:"specs"`
	PricePerHour float64           `json:"price_per_hour" binding:"required,gt=0"`
	PricePerDay  float64           `json:"price_per_day" binding:"required,gt=0"`
	LocationID   string            `json:"location_id" binding:"required"`
	Images       []string          `json:"images"`
}

type CreateSetupRequest struct {
	Name            string        `json:"name" binding:"required"`
	Category        SetupCategory `json:"category" binding:"required"`
	PCID            string        `json:"pc_id" binding:"required"`
	PeripheralIDs   []string      `json:"peripheral_ids"`
	DiscountPercent float64       `json:"discount_percent"`
	Description     string        `json:"description"`
	Images          []string      `json:"images"`
}

type AvailabilityResponse struct {
	Available    bool    `json:"available"`
	PricePerHour float64 `json:"price_per_hour"`
	PricePerDay  float64 `json:"price_per_day"`
}
