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

}

type ModelProvider inteface {
    GetModel(ctx context.Context, id string) (*model.Model, error)
}

type FeedbackRepository interface {
	
}