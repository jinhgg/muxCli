package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"muxCli/types"
	"net/http"
	"strconv"
	"time"
)

var log = logrus.New()

func initMongo() (*qmgo.QmgoClient, context.Context, error) {
	ctx := context.Background()
	cli, e := qmgo.Open(
		ctx,
		&qmgo.Config{
			Uri:      viper.GetString("MONGODB.URL"),
			Database: viper.GetString("MONGODB.DATABASE"),
			Coll:     viper.GetString("MONGODB.COLL"),
		},
	)
	return cli, ctx, e
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	cli, ctx, e := initMongo()

	idStr := mux.Vars(r)["id"]
	users := []types.User{}
	if idStr == "" {
		cli.Find(ctx, bson.M{}).All(&users)
	} else {
		id, _ := strconv.ParseInt(idStr, 10, 8)
		cli.Find(ctx, bson.M{"id": id}).All(&users)
	}
	jsons, _ := json.Marshal(users)
	fmt.Fprintf(w, string(jsons))

	//log.Info(r.Method, r.RemoteAddr)
	log.WithFields(logrus.Fields{
		"Method": r.Method,
		"url":    r.RequestURI,
		"time":   time.RFC850,
	}).Info("OK")

	defer func() {
		if e = cli.Close(ctx); e != nil {
			panic(e)
		}
	}()
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	cli, ctx, _ := initMongo()

	user := types.User{}
	err := json.Unmarshal(body, &user)
	if err != nil {
		return
	}
	result, e := cli.InsertOne(ctx, user)

	fmt.Println("result:", result)
	if e != nil {
		fmt.Fprintf(w, "创建失败")
		return
	}
	fmt.Fprintf(w, "OK")
	log.WithFields(logrus.Fields{
		"Method": r.Method,
		"url":    r.RequestURI,
		"time":   time.RFC850,
	}).Info("OK")
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)
	cli, ctx, _ := initMongo()

	e := cli.Remove(ctx, bson.M{"id": id})
	if e != nil {
		fmt.Println(e)
		fmt.Fprintf(w, "删除失败")
		return
	}
	fmt.Fprintf(w, "OK")
	log.WithFields(logrus.Fields{
		"Method": r.Method,
		"url":    r.RequestURI,
		"time":   time.RFC850,
	}).Info("OK")
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	body, _ := ioutil.ReadAll(r.Body)
	user := types.User{}
	json.Unmarshal(body, &user)

	cli, ctx, _ := initMongo()

	e := cli.UpdateOne(ctx,
		bson.M{"id": id},
		bson.M{"$set": bson.M{
			"name": user.Name,
			"age":  user.Age,
		}},
	)
	if e != nil {
		fmt.Println(e)
		fmt.Fprintf(w, "修改错误")
		return
	}
	fmt.Fprintf(w, "OK")
	log.WithFields(logrus.Fields{
		"Method": r.Method,
		"url":    r.RequestURI,
		"time":   time.RFC850,
	}).Info("OK")
}
