package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"bibit.id/challenge/model"
	"github.com/go-redis/redis/v8"
)

const (
	stockSummaryFmt     = "stocksummary-%s-%s"
	stockSummaryDateFmt = "2006-01-02"
)

func (repo *Repo) GetStockSummary(stockCode string, date time.Time) (summary model.Summary, err error) {
	key := fmt.Sprintf(stockSummaryFmt, stockCode, date.Format(stockSummaryDateFmt))

	value, err := repo.redisClient.Get(context.TODO(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Summary{}, nil
		}

		fmt.Println("Failed to get key:", err)
		return model.Summary{}, err
	}

	err = json.Unmarshal([]byte(value), &summary)
	if err != nil {
		fmt.Println("UNMARSHAL ERROR")
		return model.Summary{}, err
	}

	return summary, nil
}

func (repo *Repo) SetStockSummary(stockCode string, date time.Time, summary model.Summary) (err error) {
	key := fmt.Sprintf(stockSummaryFmt, stockCode, date.Format(stockSummaryDateFmt))

	value, err := json.Marshal(summary)
	if err != nil {
		return err
	}

	return repo.redisClient.Set(context.TODO(), key, value, 0).Err()
}
