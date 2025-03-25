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
	SavePrediction(ctx context.Context, p *model.WeatherPrediction) error
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

func (r *Resolver) Query() model.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Mutation() model.MutationResolver {
	return &mutationResolver{r}
}

type queryResolver struct { *Resolver }
type mutationResolver struct { *Resolver }

func (r* queryResolver) GetCurrentPrediction(
	ctx context.Context, 
	input model.PredictionRequest
) (*model.Prediction, error) {
	if err := validateInput(input); err != nil {
		return nil, fmt.Errorf("неверные входные данные: %v", err)
	}

	prediction, err := r.ModelService.Predict(ctx, input)

	if err != nil {
		return nil, fmt.Errorf("ошибка прогнозирования: %v", err)
	}

	if err := r.WeatherRepo.SavePrediction(ctx, prediction); err != nil {
		return nil, fmt.Errorf("ошибка сохранения прогноза: %v", err)
	}

	return prediction, nil
}

func (r* queryResolver) GetHistoricalPredictions (
	ctx context.Context, 
    location *model.GeoPositionn,
	dateForm string,
	dateTo string,
) ([]*model.HistoricalPrediction, error) {
	if location == nil {
        return nil, fmt.Errorf("не указана геопозиция")
    }

    if err := validateDate(dateForm, dateTo); err != nil {
        return nil, fmt.Errorf("неверные даты: %v", err)
    }

    filter := model.HistoryFilter{
        Location: location,
        DateFrom:     dateForm,
        DateTo:       dateTo,
    }

    predictions, err := r.WeatherRepo.GetPredictionsByDate(ctx, filter)

    if err != nil {
        return nil, fmt.Errorf("ошибка получения исторических прогнозов: %v", err)
    }

    return r.WeatherRepo.GetPredictionsByDate(ctx, filter), nil
}

func (r* queryResolver) GetAvailableModels(ctx context.Context) ([]*models.ModelInfo, error) {
	models := r.ModelService.GetModels()

    if len(models) == 0 {
        return nil, fmt.Errorf("нет доступных моделей")
    }

    return models, nil
}

