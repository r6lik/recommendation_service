package entities

import (
	"time"

	"github.com/google/uuid"
)

type ActionType string
type Season string
type TimeOfDay string
type DeviceType string
type UserSegment string

const (
	ActionView     ActionType = "view"
	ActionWishlist ActionType = "wishlist"
	ActionCart     ActionType = "cart"
	ActionPurchase ActionType = "purchase"
	ActionReview   ActionType = "review"
)

const (
	SeasonSpring Season = "SPRING"
	SeasonSummer Season = "SUMMER"
	SeasonFall   Season = "FALL"
	SeasonWinter Season = "WINTER"
)

const (
	DeviceMobile  DeviceType = "MOBILE"
	DeviceDesktop DeviceType = "DESKTOP"
	DeviceTablet  DeviceType = "TABLET"
)

const (
	SegmentHighValue         UserSegment = "HIGH"
	SegmentBrowser           UserSegment = "BROWSER"
	SegmentWishlistCollector UserSegment = "WISHLIST_COLLECTOR"
	SegmentNewUser           UserSegment = "NEW_USER"
	SegmentVIP               UserSegment = "VIP"
)

type Product struct {
	ID            uuid.UUID `db:"uuid" json:"uuid"`
	Name          string    `db:"name" json:"name"`
	Price         int       `db:"price" json:"price"`
	DiscountPrice int       `db:"discount_price" json:"discount_price"`
	CategoryID    uuid.UUID `db:"category_id" json:"category_id"`
	Tags          []string  `db:"tags" json:"tags"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type Recommendation struct {
	ProductID     uuid.UUID      `db:"product_id" json:"product_id"`
	Score         float64        `db:"score" json:"score"`
	Reason        string         `db:"reason" json:"reason"`
	ReasonDetails map[string]any `db:"reason_details" json:"reason_details,omitempty"`
}

type UserInteraction struct {
	ID         int64      `db:"id" json:"id"`
	UserID     int64      `db:"user_id" json:"user_id"`
	ProductID  uuid.UUID  `db:"product_id" json:"product_id"`
	ActionType ActionType `db:"action_type" json:"action_type"`
	Weight     int64      `db:"weight" json:"weight"`
	Timestamp  time.Time  `db:"timestamp" json:"timestamp"`
}

type Context struct {
	Timestamp   time.Time
	Season      Season
	TimeOfDay   TimeOfDay
	DayOfWeek   string
	Region      string
	DeviceType  DeviceType
	UserSegment UserSegment
	IsHoliday   bool
}

type UserProfile struct {
	UserID             int64       `db:"user_id" json:"user_id"`
	TotalPurchases     int         `db:"total_purchases" json:"total_purchases"`
	TotalSpent         float64     `db:"total_spent" json:"total_spent"`
	FavoriteCategories []string    `db:"favorite_categories" json:"favorite_categories"`
	Segment            UserSegment `db:"segment" json:"segment"`
	LastInteractionAt  time.Time   `db:"last_interaction_at" json:"last_interaction_at"`
}

func (a ActionType) GetWeight() int64 {
	weight := map[ActionType]int64{
		ActionView:     1,
		ActionWishlist: 10,
		ActionCart:     5,
		ActionPurchase: 50,
		ActionReview:   15,
	}
	if weight, ok := weight[a]; ok {
		return weight
	}
	return 1
}
