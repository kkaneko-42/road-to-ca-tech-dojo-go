package cache

import (
	"fmt"
	"github.com/go-redis/redis"
)

const (
	ranking_key string = "ranking"
	redis_nil_err string = "redis: nil"
)

func PushScore(user_id string, score int, cli *redis.Client) error {
	prev_score_z := cli.ZScore(ranking_key, user_id)

	// すでにスコアが記録されていた場合
	if prev_score_z != nil {
		if int(prev_score_z.Val()) >= score {
			return nil
		} else {
			if err := cli.ZRem(ranking_key, user_id).Err(); err != nil {
				return err
			}
		}
	}

	// redis serverへ追加
	err := cli.ZAdd(ranking_key, redis.Z{
			Score: float64(score),
			Member: user_id,
		}).Err()
	if err != nil {
		return err
	}

	// 永続化
	err = cli.Save().Err()
	if err != nil {
		return err
	}

	return nil
}

func GetUserHighScore(user_id string, cli *redis.Client) (int, error) {
	high_score, err := cli.ZScore(ranking_key, user_id).Result()
	if err != nil {
		if fmt.Sprintf("%s", err) == redis_nil_err {
			return 0, nil
		} else {
			return 0, err
		}
	}

	return int(high_score), nil
}

func GetRanking(start, end int64, cli *redis.Client) (map[string]RankData, error) {
	var ranking map[string]RankData = map[string]RankData{}

	ranked_users := cli.ZRevRangeWithScores(ranking_key, start - 1, end - 1)
	if err := ranked_users.Err(); err != nil {
		return nil, err
	}

	for i, user := range ranked_users.Val() {
		ranking[user.Member.(string)] = RankData{
			Score: int(user.Score),
			Rank: int(start) + i,
		}
	}
	return ranking, nil
}

type RankData struct {
	Score int
	Rank int
}
