package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"bibit.id/challenge/model"
	"github.com/go-redis/redis/v8"
)

const (
	stockSummaryFmt     = "stocksummary-%s-%s"
	stockSummaryDateFmt = "2006-01-02"
)

func (repo *Repo) GetStockSummary(ctx context.Context, request model.GetStockSummaryRequest) (summary model.Summary, err error) {
	key := fmt.Sprintf(stockSummaryFmt, request.StockCode, request.Date.Format(stockSummaryDateFmt))

	value, err := repo.redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Summary{}, nil
		}

		return model.Summary{}, err
	}

	err = json.Unmarshal([]byte(value), &summary)
	if err != nil {
		return model.Summary{}, err
	}

	return summary, nil
}

func (repo *Repo) SetStockSummary(ctx context.Context, summary model.Summary) (err error) {
	key := fmt.Sprintf(stockSummaryFmt, summary.StockCode, summary.Date.Format(stockSummaryDateFmt))

	value, err := json.Marshal(summary)
	if err != nil {
		return err
	}

	return repo.redisClient.Set(ctx, key, value, 0).Err()
}
