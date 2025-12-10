package output

import (
	"context"

	"github.com/google/uuid"
	"github.com/r6lik/recommendation_service/internal/ domain/entities"
)

type ProductRepository interface {
	GetProductByID(ctx context.Context, uuid uuid.UUID) (*entities.Product, error)
	GetProductsByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entities.Product, error)
	GetPopularProducts(ctx context.Context, context *entities.Context, limit int) ([]entities.Product, error)
	SaveProduct(ctx context.Context, product *entities.Product) error
	GetAllProducts(ctx context.Context) ([]entities.Product, error)
}

type InteractionRepository interface {
	SaveInteraction(ctx context.Context, interaction *entities.UserInteraction) error
	GetUserInteractions(ctx context.Context, userID int64, limit int) ([]entities.UserInteraction, error)
	GetProductAssociations(ctx context.Context, productID uuid.UUID, limit int) (map[int64]int, error)
}

type UserProfileRepository interface {
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*entities.UserProfile, error)
	SaveUserProfile(ctx context.Context, profile *entities.UserProfile) error
	UpdateUserSegment(ctx context.Context, userID int64, segment entities.UserSegment) error
}

type RecommendationCache interface {
	GetRecommendations(ctx context.Context, userID int64) ([]entities.Recommendation, error)
	SetRecommendations(ctx context.Context, userID int64, recs []entities.Recommendation, ttl int) error
	InvalidateRecommendations(ctx context.Context, userID int64) error
}
