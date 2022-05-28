package handler

import (
	"io"
	"fmt"
	"time"
	"encoding/json"
	"net/http"
	"math/rand"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-redis/redis"
	"42tokyo-road-to-dojo-go/pkg/server/cache"
)

const (
	lack_coins_msg string = "Lack coins error"
	max_times int = 100
)

func HandleGachaDraw(db *sql.DB, cli *redis.Client) http.HandlerFunc {
	return func (w http.ResponseWriter, req *http.Request) {
		res, err, status := createGachaDrawResponce(db, cli, req)
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

func createGachaDrawResponce(db *sql.DB, cli *redis.Client, req *http.Request) (*gachaDrawResponce, error, int) {
	parsed_req, err := parseGachaDrawRequest(db, req)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	results, err := execGachaDraw(db, cli, parsed_req)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	err = postGachaResults(db, parsed_req.UserId, results)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return &gachaDrawResponce{Results: results}, nil, http.StatusOK
}

func parseGachaDrawRequest(db *sql.DB, req *http.Request) (*gachaDrawRequest, error) {
	var parsed_body gachaDrawRequest

	jsonbody, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonbody, &parsed_body)
	if err != nil {
		return nil, err
	}
	parsed_body.UserId = getUserIdFromContext(req)

	if err = validateRequest(db, &parsed_body); err != nil {
		return nil, err
	}

	return &parsed_body, nil
}

func validateRequest(db *sql.DB, req *gachaDrawRequest) error {
	var having_coins int32

	if req.Times < 1 || req.Times > max_times {
		return fmt.Errorf("Invalid times")
	}

	err := db.QueryRow(
		"SELECT having_coins FROM users_infos " +
		"WHERE user_id = ?;",
		req.UserId).Scan(&having_coins)
	if err != nil {
		return err
	}

	gacha_cost, err := getGachaCost()
	if err != nil {
		return err
	}

	if having_coins - (gacha_cost * int32(req.Times)) < 0 {
		return fmt.Errorf(lack_coins_msg)
	}

	return nil
}

func execGachaDraw(db *sql.DB, cli *redis.Client, req *gachaDrawRequest) (*[]gotItem, error) {
	var (
		gacha_base []cache.ItemData
		result cache.ItemData
		results []gotItem
		is_new bool
	)

	items, err := cache.GetItems(cli)
	if err != nil {
		return nil, err
	}

	for i, _ := range *items {
		for j := 0; j < (*items)[i].GachaWeight; j++ {
			gacha_base = append(gacha_base, (*items)[i])
		}
	}

	for i := 0; i < req.Times; i++ {
		rand.Seed(time.Now().UnixNano())
		result = gacha_base[rand.Intn(len(gacha_base))]
		is_new, err = checkIsNew(&result, req.UserId, db)
		if err != nil {
			return nil, err
		}

		results = append(results, gotItem{
			CollectionId: result.Id,
			Name: result.Name,
			Rarity: result.Rarity,
			IsNew: is_new,
		})
	}
	return &results, nil
}

func checkIsNew(item *cache.ItemData, user_id string, db *sql.DB) (bool, error) {
	inventories, err := GetUserInventories(db, user_id)
	if err != nil {
		return false, err
	}

	if strcontains(inventories, item.Id) {
		return false, nil
	}
	return true, nil
}

func postGachaResults(db *sql.DB, user_id string, results *[]gotItem) error {
	gacha_cost, err := getGachaCost()
	if err != nil {
		return err
	}

	if _, err := db.Exec("START TRANSACTION"); err != nil {
		return err
	}

	_, err = db.Exec(createInsertResultQuery(user_id, results))
	if err != nil {
		db.Exec("ROLLBACK")
		return err
	}

	_, err = db.Exec(
		"UPDATE users_infos " +
		"SET having_coins = having_coins - ? " +
		"WHERE user_id = ?;",
		gacha_cost * int32(len(*results)), user_id)
	if err != nil {
		db.Exec("ROLLBACK")
		return err
	}

	db.Exec("COMMIT")
	return nil
}

func getGachaCost() (int32, error) {
	const (
		root_url string = "http://localhost:8080"
		setting_get_url string = "/setting/get"
	)
	var setting_res settingGetResponse

	res, err := http.Get(root_url + setting_get_url)
	if err != nil {
		return 0, err
	}

	jsondata, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(jsondata, &setting_res)
	if err != nil {
		return 0, err
	}

	return setting_res.GachaCoinConsumption, nil
}

func createInsertResultQuery(user_id string, results *[]gotItem) string {
	db_query := "INSERT INTO users_inventories (user_id, item_id) "

	for i := 0; i < len(*results); i++ {
		db_query += fmt.Sprintf("VALUES (\"%s\", \"%s\")", user_id, (*results)[i].CollectionId)
		if i != len(*results) - 1 {
			db_query += ", "
		} else {
			db_query += ";"
		}
	}
	
	return db_query
}

type lackCoinsError error

type gachaDrawRequest struct {
	UserId string
	Times int `json:"times"`
}

type gotItem struct {
	CollectionId string `json:"collectionID"`
	Name string `json:"name"`
	Rarity int `json:"rarity"`
	IsNew bool `json:"isNew"`
}

type gachaDrawResponce struct {
	Results *[]gotItem `json:"results"`
}
