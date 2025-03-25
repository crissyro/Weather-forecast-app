package graph

import (
	"context"
    "fmt"
	"time"

    "github.com/crissyro/weatherapi/graph/model"
)

type Resolver struct {
	WeatherRepo WeatherRepository
	ModelService ModelProvider
	FeedbackStore FeedbackRepository
}

type WeatherRepository inteface {
	SavePrediction(ctx context.Context, p* model.WeatherPrediction) errorconst
	GetPredictionsByDate(ctx context.Context, filter model.HistoryFilter) ([]*model.HistoricalPrediction, error)
}

type ModelProvider interface {
	GetModels() []*model.ModelInfo
	Predict(ctx context.Context, req model.PredictionRequest) (*model.WeatherPrediction, error)
}

type FeedbackRepository interface {
	AddFeedback(ctx context.Context, feedback *model.FeedbackResult) error
	CalculateAccuracy(ctx context.Context, predictionID string) (float64, error)
}