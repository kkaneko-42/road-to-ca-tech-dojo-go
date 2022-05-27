package cache

import (
	"fmt"
	"github.com/go-redis/redis"
)

const (
	ranking_key string = "ranking"
)

func PushScore(user_id string, score int, cli *redis.Client) error {
	err := cli.ZAdd(ranking_key, redis.Z{float64(score), user_id}).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetUserHighScore(user_id string, cli *redis.Client) (int, error) {
	high_score, err := cli.ZScore(ranking_key, user_id).Result()
	if err != nil {
		return 0, err
	}

	return int(high_score), nil
}

func GetRanking(start, end int64, cli *redis.Client) (*[]RankData, error) {
	var ranking []RankData

	ranked_users := cli.ZRangeWithScores(ranking_key, start, end).Val()
	if isEmpty(ranked_users) {
		return nil, fmt.Errorf("Ranking is empty")
	}

	for i, user := range ranked_users {
		ranking = append(ranking, RankData{
			UserId: user.Member.(string),
			Score: int(user.Score),
			Rank: i + 1,
		})
	}
	return &ranking, nil
}

func isEmpty(ranking []redis.Z) bool {
	if len(ranking) == 0 {
		return true
	}
	return false
}

type RankData struct {
	UserId string
	Score int
	Rank int
}
