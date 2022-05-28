package server

import (
	"log"
	"net/http"
	"42tokyo-road-to-dojo-go/pkg/server/handler"
	"42tokyo-road-to-dojo-go/pkg/server/cache"
	"42tokyo-road-to-dojo-go/pkg/http/middleware"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-redis/redis"
)

// Serve HTTPサーバを起動する
func Serve(addr string) {

	/* ==== DBのOpen ==== */
	db, err := sql.Open("mysql", "root:ca-tech-dojo@(127.0.0.1:3306)/road_to_ca")
	if err != nil {
		log.Fatal("DB connection failed: ", err)
		return;
	}
	defer db.Close()

	/* ==== Master dataのキャッシュ ==== */
	cli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	if err = cli.Del(cache.Item_key).Err(); err != nil {
		log.Fatal("Cache Initialization failed: ", err)
	}

	err = cacheMasterData(db, cli)
	if err != nil {
		log.Fatal("Caching failed: ", err)
	}

	/* ===== URLマッピングを行う ===== */
	mappingURL(db, cli)

	/* ===== サーバの起動 ===== */
	log.Println("Server running...")
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}

func cacheMasterData(db *sql.DB, cli *redis.Client) error {
	err := cache.CacheItems(db, cli)
	if err != nil {
		return err
	}
	return nil
}

func mappingURL(db *sql.DB, cli *redis.Client) {
	http.HandleFunc("/setting/get", get(handler.HandleSettingGet()))
	http.HandleFunc("/user/create", post(handler.HandleUserCreate(db)))
	http.HandleFunc("/user/get",
		get(middleware.Authenticate(handler.HandleUserGet(db, cli))))
	http.HandleFunc("/user/update",
		post(middleware.Authenticate(handler.HandleUserUpdate(db))))
	http.HandleFunc("/collection/list",
		get(middleware.Authenticate(handler.HandleCollectionListGet(db, cli))))
	http.HandleFunc("/ranking/list",
		get(middleware.Authenticate(handler.HandleRankingGet(db, cli))))
	http.HandleFunc("/game/finish",
		post(middleware.Authenticate(handler.HandleGameFinish(db, cli))))
	http.HandleFunc("/gacha/draw",
		post(middleware.Authenticate(handler.HandleGachaDraw(db, cli))))
}

// get GETリクエストを処理する
func get(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodGet)
}

// post POSTリクエストを処理する
func post(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodPost)
}

// httpMethod 指定したHTTPメソッドでAPIの処理を実行する
func httpMethod(apiFunc http.HandlerFunc, method string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// CORS対応
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Accept,Origin,x-token")

		// プリフライトリクエストは処理を通さない
		if request.Method == http.MethodOptions {
			return
		}
		// 指定のHTTPメソッドでない場合はエラー
		if request.Method != method {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method Not Allowed"))
			return
		}

		// 共通のレスポンスヘッダを設定
		writer.Header().Add("Content-Type", "application/json")
		apiFunc(writer, request)
	}
}
