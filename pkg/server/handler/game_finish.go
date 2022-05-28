package handler

import (
	"io"
	"fmt"
	"net/http"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-redis/redis"
	"42tokyo-road-to-dojo-go/pkg/server/cache"
)

func HandleGameFinish(db *sql.DB, cli *redis.Client) http.HandlerFunc {
	return func (w http.ResponseWriter, req *http.Request) {
		res, err := createGameFinishResponce(db, cli, req)
		if err != nil {
			putError(w, err)
			return
		}

		jsondata, err := json.Marshal(res)
		if err != nil {
			putError(w, err)
			return
		}
		w.Write(jsondata)
	}
}

func createGameFinishResponce(db *sql.DB, cli *redis.Client, req *http.Request) (*gameFinishResponce, error) {
	req_body, err := parseGameFinishRequest(req)
	if err != nil {
		return nil, err
	}

	reward, err := postGameFinish(db, cli, req_body)
	if err != nil {
		return nil, err
	}
	return &gameFinishResponce{Coin: reward}, nil
}

func parseGameFinishRequest(req *http.Request) (*gameFinishRequest, error){
	jsonbody, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var parsed_body gameFinishRequest
	parsed_body.UserId = getUserIdFromContext(req)
	err = json.Unmarshal(jsonbody, &parsed_body)
	if err != nil {
		return nil, err
	} else if parsed_body.Score < 0 {
		return nil, fmt.Errorf("Invalid score")
	}

	return &parsed_body, nil
}

func postGameFinish(db *sql.DB, cli *redis.Client, req *gameFinishRequest) (int, error) {
	if err := cache.PushScore(req.UserId, req.Score, cli); err != nil {
		return 0, err
	}
	reward, err := updateCoin(db, req)
	if err != nil {
		return 0, err
	}

	return reward, nil
}

func updateCoin(db *sql.DB, req *gameFinishRequest) (int, error) {
	reward := calcReward(req.Score)

	_, err := db.Exec(
		"UPDATE users_infos " +
		"SET having_coins = having_coins + ? " +
		"WHERE user_id = ?",
		reward, req.UserId)
	if err != nil {
		return 0, err
	}

	return reward, nil
}

func calcReward(score int) int {
	return score / 100
}

type gameFinishRequest struct {
	UserId string
	Score int
}

type gameFinishResponce struct {
	Coin int `json:"coin"`
}
