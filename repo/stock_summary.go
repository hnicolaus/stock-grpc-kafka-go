/*
	Hans Nicolaus
	29 Aug 2023
*/

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

// GetStockSummary gets stock summary data for stockCode for the requested date range by performing ZRangeByScore:
// - Key: stockCode
// - Min score: unix value of the requested fromDate
// - Max score: unix value of the requested toDate
// Using ZRangeByScore allows users to retrieve the a stock's summary data over a period of time.
// To get a stock's summary for a single date, specify the same fromDate (inclusive) and toDate (inclusive).
// To get a stock's summary over a period of time, specify a fromDate value that is less than toDate.
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

// UpdateStockSummary upserts stock summary data for a stockCode on a given date by performing the following Redis operations:
// 1. ZRangeByScore to check existing stock summary for stockCode (key) on a given date (same min & max score). Score is unix value of the summary date.
// 2. ZRem to remove any data being returned by ZRangeByScore.
// 3. ZADd to store the stock summary for stockCode (key) for the given date (score).
// This ensures a stockCode to have exactly 1 stock summary per date (score).
func (repo *Repo) UpdateStockSummary(ctx context.Context, summary model.Summary) error {
	key := fmt.Sprintf(stockSummaryFmt, summary.StockCode)

	dateUnix := summary.Date.Unix()
	zRangeMin := strconv.Itoa(int(dateUnix))
	zRangeMax := zRangeMin

	existingSummary, err := repo.redisClient.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: zRangeMin,
		Max: zRangeMax,
	}).Result()
	if err != nil {
		return err
	}

	if len(existingSummary) > 0 {
		err = repo.redisClient.ZRemRangeByScore(ctx, key, zRangeMin, zRangeMax).Err()
		if err != nil {
			return err
		}
	}

	value, err := json.Marshal(summary)
	if err != nil {
		return err
	}

	err = repo.redisClient.ZAdd(ctx, key, &redis.Z{
		Score:  float64(dateUnix),
		Member: value,
	}).Err()
	if err != nil && len(existingSummary) > 0 {
		// Attempt to "recover" (re-add) the previously removed summary
		repo.redisClient.ZAdd(ctx, key, &redis.Z{
			Score:  float64(dateUnix),
			Member: []byte(existingSummary[0]),
		})
	}

	return err
}
