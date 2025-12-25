package comparisons

import (
"context"

"cinema.log.server.golang/internal/domain"
"github.com/google/uuid"
)

type Service struct {
ComparisonStore ComparisonStore
}

func NewService(comparisonStore ComparisonStore) *Service {
return &Service{
ComparisonStore: comparisonStore,
}
}

func (s *Service) CreateComparison(ctx context.Context, comparison domain.ComparisonHistory) (*domain.ComparisonHistory, error) {
return s.ComparisonStore.CreateComparison(ctx, comparison)
}

func (s *Service) HasBeenCompared(ctx context.Context, userId, filmAId, filmBId uuid.UUID) (bool, error) {
return s.ComparisonStore.HasBeenCompared(ctx, userId, filmAId, filmBId)
}

func (s *Service) GetComparisonHistory(ctx context.Context, userId uuid.UUID) ([]domain.ComparisonHistory, error) {
return s.ComparisonStore.GetComparisonHistory(ctx, userId)
}
