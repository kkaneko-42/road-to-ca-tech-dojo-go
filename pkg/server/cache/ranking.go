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

	if prev_score_z != nil && int(prev_score_z.Val()) >= score {
		return nil
	}
	err := cli.ZAdd(ranking_key, redis.Z{
			Score: float64(score),
			Member: user_id,
		}).Err()
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

	ranked_users := cli.ZRangeWithScores(ranking_key, start - 1, end - 1).Val()
	if isEmpty(ranked_users) {
		return nil, fmt.Errorf("Ranking is empty")
	}

	for i, user := range ranked_users {
		ranking[user.Member.(string)] = RankData{
			Score: int(user.Score),
			Rank: i + 1,
		}
	}
	return ranking, nil
}

func isEmpty(ranking []redis.Z) bool {
	if len(ranking) == 0 {
		return true
	}
	return false
}

type RankData struct {
	Score int
	Rank int
}
