package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"bibit.id/challenge/model"
	"github.com/go-redis/redis/v8"
)

const (
	stockSummaryFmt = "stocksummary-%s"
)

func (repo *Repo) GetStockSummary(ctx context.Context, request model.GetStockSummaryRequest) (result []model.Summary, err error) {
	key := fmt.Sprintf(stockSummaryFmt, request.StockCode)

	fromDateUnix := request.FromDate.Unix()
	toDateUnix := request.ToDate.Unix()

	redisResult, err := repo.redisClient.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: strconv.Itoa(int(fromDateUnix)),
		Max: strconv.Itoa(int(toDateUnix)),
	}).Result()
	if err != nil {
		return []model.Summary{}, err
	}

	result = []model.Summary{}
	for _, data := range redisResult {
		summary := model.Summary{}
		err = json.Unmarshal([]byte(data), &summary)
		if err != nil {
			return []model.Summary{}, err
		}

		result = append(result, summary)
	}

	return result, nil
}

func (repo *Repo) UpdateStockSummary(ctx context.Context, summary model.Summary) error {
	key := fmt.Sprintf(stockSummaryFmt, summary.StockCode)
	dateUnix := summary.Date.Unix()

	existingSummary, err := repo.redisClient.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: strconv.Itoa(int(dateUnix)),
		Max: strconv.Itoa(int(dateUnix)),
	}).Result()
	if err != nil {
		return err
	}

	if len(existingSummary) > 0 {
		_, err = repo.redisClient.ZRem(ctx, key, existingSummary).Result()
		if err != nil {
			return err
		}
	}

	value, err := json.Marshal(summary)
	if err != nil {
		return err
	}

	return repo.redisClient.ZAdd(ctx, key, &redis.Z{
		Score:  float64(dateUnix),
		Member: value,
	}).Err()
}
