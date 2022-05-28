package handler

import (
	"fmt"
	"strings"
	"strconv"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-redis/redis"
	"encoding/json"
	"42tokyo-road-to-dojo-go/pkg/server/cache"
)

const (
	nb_get_once int64 = 10
	rank_begin_key string = "start"
)

func HandleRankingGet(db *sql.DB, cli *redis.Client) http.HandlerFunc {
	return func (w http.ResponseWriter, req *http.Request) {
		res, err, status := createRankingGetResponce(db, cli, req)
		if err != nil {
			putError(w, err, status)
			return
		}

		jsondata, err := json.Marshal(res)
		if err != nil {
			putError(w, err, http.StatusInternalServerError)
			return
		}
		w.Write(jsondata)
	}
}

func createRankingGetResponce(db *sql.DB, cli *redis.Client, req *http.Request) (*RankingGetResponce, error, int) {
	rank_begin, err := getRankBegin(req)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}
	rank_end := rank_begin + nb_get_once - 1

	ranking, err := cache.GetRanking(rank_begin, rank_end, cli)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	} else if len(ranking) == 0 {
		return &RankingGetResponce{
			Ranks: make([]*UserRankData, 0),
		}, nil, http.StatusNoContent
	}

	responce, err := setUserName(ranking, db)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return responce, nil, http.StatusOK
}

func getRankBegin(req *http.Request) (int64, error) {
	url_query := req.URL.Query().Get(rank_begin_key)
	rank_begin, err := strconv.ParseInt(url_query, 10, 64)
	if err != nil {
		return 0, err
	} else if (rank_begin <= 0) {
		return 0, fmt.Errorf("Query is 0 or negative")
	}

	return rank_begin, nil
}

func setUserName(ranking map[string]cache.RankData, db *sql.DB) (*RankingGetResponce, error) {
	var (
		ranks []*UserRankData
		userid_buf string
		username_buf string
	)
	rows, err := db.Query(getDBQuery(ranking))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(
			&userid_buf,
			&username_buf,
		)
		if err != nil {
			return nil, err
		}
		ranks = append(ranks, &UserRankData{
			UserId: userid_buf,
			UserName: username_buf,
			Rank: ranking[userid_buf].Rank,
			Score: ranking[userid_buf].Score,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &RankingGetResponce{Ranks: ranks}, err
}

func getDBQuery(ranking map[string]cache.RankData) string {
	ranked_user_id := getRankedUsersId(ranking)
	for i, _ := range ranked_user_id {
		ranked_user_id[i] = strconv.Quote(ranked_user_id[i])
	}

	return fmt.Sprintf(
		"SELECT user_id, user_name FROM users_infos " +
		"WHERE user_id in (%s)", strings.Join(ranked_user_id, ","))
}

func getRankedUsersId(ranking map[string]cache.RankData) []string {
	var keys []string

	for key := range ranking {
		keys = append(keys, key)
	}

	return keys
}

type UserRankData struct {
	UserId string `json:"userId"`
	UserName string `json:"userName"`
	Rank int `json:"rank"`
	Score int `json:"score"`
}

type RankingGetResponce struct {
	Ranks []*UserRankData `json:"ranks"`
}
