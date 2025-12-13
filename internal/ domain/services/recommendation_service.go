package services

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/r6lik/recommendation_service/internal/ domain/entities"
	"github.com/r6lik/recommendation_service/internal/ domain/ports/output"
)

type RecommendationService struct {
	productRepo     output.ProductRepository
	interactionRepo output.InteractionRepository
	recommendCache  output.RecommendationCache
}

func NewRecommendationService(
	productRepo output.ProductRepository,
	interactionRepo output.InteractionRepository,
	recommendCache output.RecommendationCache,
) *RecommendationService {
	return &RecommendationService{
		productRepo:     productRepo,
		interactionRepo: interactionRepo,
		recommendCache:  recommendCache,
	}
}

func (s *RecommendationService) GetRecommendations(
	ctx context.Context,
	userID int64,
	device entities.DeviceType,
) ([]entities.Recommendation, error) {
	if cached, err := s.recommendCache.GetRecommendations(ctx, userID); err != nil && cached != nil {
		return cached, nil
	}

	currentContext := s.buildContext(device)

	recommendations := s.generateRecommendations(ctx, userID, currentContext)

	ttl := s.calculateCacheTTL(currentContext)

	err := s.recommendCache.SetRecommendations(ctx, userID, recommendations, ttl)
	if err != nil {
		return nil, err
	}

	return recommendations, nil
}

func (s *RecommendationService) generateRecommendations(
	ctx context.Context,
	userID int64,
	context *entities.Context,
) []entities.Recommendation {
	var recommendations []entities.Recommendation

	seasonalRecs := s.getSeasonalRecommendations(ctx, context.Season)
	recommendations = append(recommendations, seasonalRecs...)

	timeRecs := s.getTimeBasedRecommendations(ctx, context.TimeOfDay)
	recommendations = append(recommendations, timeRecs...)

	popularProducts, _ := s.productRepo.GetPopularProducts(ctx, context, 5)
	for _, product := range popularProducts {
		recommendations = append(recommendations, entities.Recommendation{
			ProductID: product.ID,
			Score:     85.0,       // !TODO()
			Reason:    "trending", // !TODO()
			ReasonDetails: map[string]any{
				"category": product.CategoryID,
			},
		})
	}

	recommendations = s.deduplicateAndSort(recommendations)

	if len(recommendations) > 10 {
		recommendations = recommendations[:10]
	}

	return recommendations
}

func (s *RecommendationService) getSeasonalRecommendations(
	ctx context.Context,
	season entities.Season,
) []entities.Recommendation {
	var recommendations []entities.Recommendation

	seasonalCategories := map[entities.Season]uuid.UUID{
		entities.SeasonSummer: uuid.UUID{}, // !TODO(): Get from services
		entities.SeasonWinter: uuid.UUID{},
		entities.SeasonSpring: uuid.UUID{},
		entities.SeasonFall:   uuid.UUID{},
	}

	category := seasonalCategories[season]
	if category == uuid.Nil {
		return recommendations
	}

	products, _ := s.productRepo.GetProductsByCategory(ctx, category)
	for _, product := range products {
		recommendations = append(recommendations, entities.Recommendation{
			ProductID: product.ID,
			Score:     75.0, // !TODO()
		})
	}

	return recommendations
}

func (s *RecommendationService) getTimeBasedRecommendations(ctx context.Context, timeOfDay entities.TimeOfDay) []entities.Recommendation {
	var recommendations []entities.Recommendation

	// TODO(): Pull date from shop
	timeCategories := map[entities.TimeOfDay][]uuid.UUID{
		entities.TimeOfDayEvening:   {},
		entities.TimeOfDayMorning:   {},
		entities.TimeOfDayAfternoon: {},
		entities.TimeOfDayNight:     {},
	}

	categories := timeCategories[timeOfDay]
	for _, category := range categories {
		products, _ := s.productRepo.GetProductsByCategory(ctx, category)
		for _, product := range products {
			recommendations = append(recommendations, entities.Recommendation{
				ProductID: product.ID,
				Score:     70.0,         // !TODO(): calculate dynamically
				Reason:    "time_based", // !TODO()
				ReasonDetails: map[string]any{
					"time_of_day": timeOfDay,
				},
			})
		}
	}

	return recommendations
}

func (s *RecommendationService) RecordEvent(ctx context.Context, interaction *entities.UserInteraction) error {
	if err := s.interactionRepo.SaveInteraction(ctx, interaction); err != nil {
		return err
	}

	err := s.recommendCache.InvalidateRecommendations(ctx, interaction.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *RecommendationService) buildContext(device entities.DeviceType) *entities.Context {
	now := time.Now()
	season := s.getSeason(now)
	timeOfDay := s.getTimeOfDay(now)

	return &entities.Context{
		Timestamp: now,
		Season:    season,
		TimeOfDay: timeOfDay,
		DayOfWeek: now.Weekday().String(),
		Device:    device,
		Region:    "RU",
	}
}

func (s *RecommendationService) deduplicateAndSort(
	recommendations []entities.Recommendation,
) []entities.Recommendation {
	seen := make(map[uuid.UUID]bool)
	var unique []entities.Recommendation

	for _, recommendations := range recommendations {
		if !seen[recommendations.ProductID] {
			seen[recommendations.ProductID] = true
			unique = append(unique, recommendations)
		}
	}

	sort.Slice(unique, func(i, j int) bool {
		return unique[i].Score > unique[j].Score
	})

	return unique
}

func (s *RecommendationService) calculateCacheTTL(context *entities.Context) int {
	baseTTL := 3600 // 1h

	switch context.TimeOfDay {
	case entities.TimeOfDayEvening:
		baseTTL = 1800 // 30m
	case entities.TimeOfDayNight:
		baseTTL = 21600 // 6h
	}

	return baseTTL
}

func (s *RecommendationService) getSeason(t time.Time) entities.Season {
	month := t.Month()
	switch {
	case month >= 3 && month <= 5:
		return entities.SeasonSpring
	case month >= 6 && month <= 8:
		return entities.SeasonSummer
	case month >= 9 && month <= 11:
		return entities.SeasonFall
	default:
		return entities.SeasonWinter
	}
}

func (s *RecommendationService) getTimeOfDay(t time.Time) entities.TimeOfDay {
	hour := t.Hour()
	switch {
	case hour >= 5 && hour < 12:
		return entities.TimeOfDayMorning
	case hour >= 12 && hour < 17:
		return entities.TimeOfDayAfternoon
	case hour >= 17 && hour < 21:
		return entities.TimeOfDayAfternoon
	default:
		return entities.TimeOfDayAfternoon
	}
}
