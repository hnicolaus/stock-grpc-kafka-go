/*
	Hans Nicolaus
	29 Aug 2023
*/

package repo

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"testing"
	"time"

	"bibit.id/challenge/model"
	mock "bibit.id/challenge/repo/_mock"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
)

const (
	expectedKey = "stocksummary-BBCA"
)

func Test_Repo_GetStockSummary(t *testing.T) {
	type args struct {
		ctx   context.Context
		input model.GetStockSummaryRequest
	}
	type fields struct {
		redisClient func(ctrl *gomock.Controller) RedisClient
	}
	tests := []struct {
		name   string
		args   args
		fields fields

		wantResponse []model.Summary
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				input: model.GetStockSummaryRequest{
					StockCode: "BBCA",
					FromDate:  time.Time{}.AddDate(0, 0, 1),
					ToDate:    time.Time{}.AddDate(0, 0, 2),
				},
			},
			fields: fields{
				redisClient: func(ctrl *gomock.Controller) RedisClient {
					m := mock.NewMockRedisClient(ctrl)

					fromDate := time.Time{}.AddDate(0, 0, 1)
					toDate := time.Time{}.AddDate(0, 0, 2)

					expectedSummaryOne := model.Summary{
						StockCode: "BBCA",
						Date:      fromDate,
						Prev:      8000,
						Open:      8050,
						High:      8100,
						Low:       7950,
						Close:     8100,
						Volume:    900,
						Value:     7210000,
						Average:   8011,
					}
					expectedSummaryOneJSON, _ := json.Marshal(expectedSummaryOne)

					expectedSummaryTwo := model.Summary{
						StockCode: "BBCA",
						Date:      toDate,
						Prev:      8000,
						Open:      8050,
						High:      8100,
						Low:       7950,
						Close:     8100,
						Volume:    900,
						Value:     7210000,
						Average:   8011,
					}
					expectedSummaryTwoJSON, _ := json.Marshal(expectedSummaryTwo)

					expectedResult := []string{string(expectedSummaryOneJSON), string(expectedSummaryTwoJSON)}

					m.EXPECT().ZRangeByScore(gomock.Any(), expectedKey, &redis.ZRangeBy{
						Min: strconv.Itoa(int(fromDate.Unix())),
						Max: strconv.Itoa(int(toDate.Unix())),
					}).Return(redis.NewStringSliceResult(expectedResult, nil))

					return m
				},
			},
			wantResponse: []model.Summary{
				{
					StockCode: "BBCA",
					Date:      time.Time{}.AddDate(0, 0, 1),
					Prev:      8000,
					Open:      8050,
					High:      8100,
					Low:       7950,
					Close:     8100,
					Volume:    900,
					Value:     7210000,
					Average:   8011,
				},
				{
					StockCode: "BBCA",
					Date:      time.Time{}.AddDate(0, 0, 2),
					Prev:      8000,
					Open:      8050,
					High:      8100,
					Low:       7950,
					Close:     8100,
					Volume:    900,
					Value:     7210000,
					Average:   8011,
				},
			},
		},
		{
			name: "error-redis",
			args: args{
				ctx: context.Background(),
				input: model.GetStockSummaryRequest{
					StockCode: "BBCA",
					FromDate:  time.Time{}.AddDate(0, 0, 1),
					ToDate:    time.Time{}.AddDate(0, 0, 2),
				},
			},
			fields: fields{
				redisClient: func(ctrl *gomock.Controller) RedisClient {
					m := mock.NewMockRedisClient(ctrl)

					fromDate := time.Time{}.AddDate(0, 0, 1)
					toDate := time.Time{}.AddDate(0, 0, 2)

					m.EXPECT().ZRangeByScore(gomock.Any(), expectedKey, &redis.ZRangeBy{
						Min: strconv.Itoa(int(fromDate.Unix())),
						Max: strconv.Itoa(int(toDate.Unix())),
					}).Return(redis.NewStringSliceResult([]string{}, errors.New("error-redis")))

					return m
				},
			},
			wantResponse: []model.Summary{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			repo := &Repo{
				redisClient: tt.fields.redisClient(ctrl),
			}

			gotResponse, err := repo.GetStockSummary(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("repo.GetStockSummary() err = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("repo.GetStockSummary() gotResponse = %v, wantResponse %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func Test_Repo_UpdateStockSummary(t *testing.T) {
	type args struct {
		ctx   context.Context
		input model.Summary
	}
	type fields struct {
		redisClient func(ctrl *gomock.Controller) RedisClient
	}
	tests := []struct {
		name   string
		args   args
		fields fields

		wantErr bool
	}{
		{
			name: "success-existing-summary",
			args: args{
				ctx: context.Background(),
				input: model.Summary{
					StockCode: "BBCA",
					Date:      time.Time{}.AddDate(0, 0, 1),
					Prev:      9999,
					Open:      9999,
					High:      9999,
					Low:       9999,
					Close:     9999,
					Volume:    9999,
					Value:     99999999,
					Average:   9999,
				},
			},
			fields: fields{
				redisClient: func(ctrl *gomock.Controller) RedisClient {
					m := mock.NewMockRedisClient(ctrl)

					expectedDate := time.Time{}.AddDate(0, 0, 1)

					expectedExistingSummary := model.Summary{
						StockCode: "BBCA",
						Date:      expectedDate,
						Prev:      8000,
						Open:      8050,
						High:      8100,
						Low:       7950,
						Close:     8100,
						Volume:    900,
						Value:     7210000,
						Average:   8011,
					}
					expectedExistingSummaryJSON, _ := json.Marshal(expectedExistingSummary)
					expectedExistingSummaryList := []string{string(expectedExistingSummaryJSON)}

					m.EXPECT().ZRangeByScore(gomock.Any(), expectedKey, &redis.ZRangeBy{
						Min: strconv.Itoa(int(expectedDate.Unix())),
						Max: strconv.Itoa(int(expectedDate.Unix())),
					}).Return(redis.NewStringSliceResult(expectedExistingSummaryList, nil))

					m.EXPECT().ZRem(gomock.Any(), expectedKey, expectedExistingSummaryList).Return(redis.NewIntCmd(context.Background(), int64(1), nil))

					expectedNewSummary := model.Summary{
						StockCode: "BBCA",
						Date:      expectedDate,
						Prev:      9999,
						Open:      9999,
						High:      9999,
						Low:       9999,
						Close:     9999,
						Volume:    9999,
						Value:     99999999,
						Average:   9999,
					}
					expectedNewSummaryJSON, _ := json.Marshal(expectedNewSummary)
					expectedMembers := []*redis.Z{
						{
							Score:  float64(expectedDate.Unix()),
							Member: expectedNewSummaryJSON,
						},
					}

					m.EXPECT().ZAdd(gomock.Any(), expectedKey, expectedMembers).Return(redis.NewIntCmd(context.Background(), int64(1), nil))

					return m
				},
			},
		},
		{
			name: "success-no-existing-summary",
			args: args{
				ctx: context.Background(),
				input: model.Summary{
					StockCode: "BBCA",
					Date:      time.Time{}.AddDate(0, 0, 1),
					Prev:      9999,
					Open:      9999,
					High:      9999,
					Low:       9999,
					Close:     9999,
					Volume:    9999,
					Value:     99999999,
					Average:   9999,
				},
			},
			fields: fields{
				redisClient: func(ctrl *gomock.Controller) RedisClient {
					m := mock.NewMockRedisClient(ctrl)

					expectedDate := time.Time{}.AddDate(0, 0, 1)

					m.EXPECT().ZRangeByScore(gomock.Any(), expectedKey, &redis.ZRangeBy{
						Min: strconv.Itoa(int(expectedDate.Unix())),
						Max: strconv.Itoa(int(expectedDate.Unix())),
					}).Return(redis.NewStringSliceResult([]string{}, nil))

					expectedNewSummary := model.Summary{
						StockCode: "BBCA",
						Date:      expectedDate,
						Prev:      9999,
						Open:      9999,
						High:      9999,
						Low:       9999,
						Close:     9999,
						Volume:    9999,
						Value:     99999999,
						Average:   9999,
					}
					expectedNewSummaryJSON, _ := json.Marshal(expectedNewSummary)
					expectedMembers := []*redis.Z{
						{
							Score:  float64(expectedDate.Unix()),
							Member: expectedNewSummaryJSON,
						},
					}

					m.EXPECT().ZAdd(gomock.Any(), expectedKey, expectedMembers).Return(redis.NewIntCmd(context.Background(), int64(1), nil))

					return m
				},
			},
		},
		{
			name: "error-zrangebyscore",
			args: args{
				ctx: context.Background(),
				input: model.Summary{
					StockCode: "BBCA",
					Date:      time.Time{}.AddDate(0, 0, 1),
					Prev:      9999,
					Open:      9999,
					High:      9999,
					Low:       9999,
					Close:     9999,
					Volume:    9999,
					Value:     99999999,
					Average:   9999,
				},
			},
			fields: fields{
				redisClient: func(ctrl *gomock.Controller) RedisClient {
					m := mock.NewMockRedisClient(ctrl)

					expectedDate := time.Time{}.AddDate(0, 0, 1)

					m.EXPECT().ZRangeByScore(gomock.Any(), expectedKey, &redis.ZRangeBy{
						Min: strconv.Itoa(int(expectedDate.Unix())),
						Max: strconv.Itoa(int(expectedDate.Unix())),
					}).Return(redis.NewStringSliceResult([]string{}, errors.New("error-zrangebyscore")))

					return m
				},
			},
			wantErr: true,
		},
		{
			name: "error-zrem",
			args: args{
				ctx: context.Background(),
				input: model.Summary{
					StockCode: "BBCA",
					Date:      time.Time{}.AddDate(0, 0, 1),
					Prev:      9999,
					Open:      9999,
					High:      9999,
					Low:       9999,
					Close:     9999,
					Volume:    9999,
					Value:     99999999,
					Average:   9999,
				},
			},
			fields: fields{
				redisClient: func(ctrl *gomock.Controller) RedisClient {
					m := mock.NewMockRedisClient(ctrl)

					expectedDate := time.Time{}.AddDate(0, 0, 1)

					expectedExistingSummary := model.Summary{
						StockCode: "BBCA",
						Date:      expectedDate,
						Prev:      8000,
						Open:      8050,
						High:      8100,
						Low:       7950,
						Close:     8100,
						Volume:    900,
						Value:     7210000,
						Average:   8011,
					}
					expectedExistingSummaryJSON, _ := json.Marshal(expectedExistingSummary)
					expectedExistingSummaryList := []string{string(expectedExistingSummaryJSON)}

					m.EXPECT().ZRangeByScore(gomock.Any(), expectedKey, &redis.ZRangeBy{
						Min: strconv.Itoa(int(expectedDate.Unix())),
						Max: strconv.Itoa(int(expectedDate.Unix())),
					}).Return(redis.NewStringSliceResult(expectedExistingSummaryList, nil))

					m.EXPECT().ZRem(gomock.Any(), expectedKey, expectedExistingSummaryList).Return(redis.NewIntResult(int64(0), errors.New("error-zrem")))

					return m
				},
			},
			wantErr: true,
		},
		{
			name: "error-zadd",
			args: args{
				ctx: context.Background(),
				input: model.Summary{
					StockCode: "BBCA",
					Date:      time.Time{}.AddDate(0, 0, 1),
					Prev:      9999,
					Open:      9999,
					High:      9999,
					Low:       9999,
					Close:     9999,
					Volume:    9999,
					Value:     99999999,
					Average:   9999,
				},
			},
			fields: fields{
				redisClient: func(ctrl *gomock.Controller) RedisClient {
					m := mock.NewMockRedisClient(ctrl)

					expectedDate := time.Time{}.AddDate(0, 0, 1)

					expectedExistingSummary := model.Summary{
						StockCode: "BBCA",
						Date:      expectedDate,
						Prev:      8000,
						Open:      8050,
						High:      8100,
						Low:       7950,
						Close:     8100,
						Volume:    900,
						Value:     7210000,
						Average:   8011,
					}
					expectedExistingSummaryJSON, _ := json.Marshal(expectedExistingSummary)
					expectedExistingSummaryList := []string{string(expectedExistingSummaryJSON)}

					m.EXPECT().ZRangeByScore(gomock.Any(), expectedKey, &redis.ZRangeBy{
						Min: strconv.Itoa(int(expectedDate.Unix())),
						Max: strconv.Itoa(int(expectedDate.Unix())),
					}).Return(redis.NewStringSliceResult(expectedExistingSummaryList, nil))

					m.EXPECT().ZRem(gomock.Any(), expectedKey, expectedExistingSummaryList).Return(redis.NewIntCmd(context.Background(), int64(1), nil))

					expectedNewSummary := model.Summary{
						StockCode: "BBCA",
						Date:      expectedDate,
						Prev:      9999,
						Open:      9999,
						High:      9999,
						Low:       9999,
						Close:     9999,
						Volume:    9999,
						Value:     99999999,
						Average:   9999,
					}
					expectedNewSummaryJSON, _ := json.Marshal(expectedNewSummary)
					expectedMembers := []*redis.Z{
						{
							Score:  float64(expectedDate.Unix()),
							Member: expectedNewSummaryJSON,
						},
					}
					m.EXPECT().ZAdd(gomock.Any(), expectedKey, expectedMembers).Return(redis.NewIntResult(int64(0), errors.New("error-zadd")))

					expectedMembers = []*redis.Z{
						{
							Score:  float64(expectedDate.Unix()),
							Member: expectedExistingSummaryJSON,
						},
					}
					m.EXPECT().ZAdd(gomock.Any(), expectedKey, expectedMembers).Return(redis.NewIntResult(int64(1), nil))

					return m
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			repo := &Repo{
				redisClient: tt.fields.redisClient(ctrl),
			}

			err := repo.UpdateStockSummary(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("repo.GetStockSummary() err = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
