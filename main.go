/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"muxCli/types"
	"strconv"
	"time"

	//"muxCli/cmd"
	"net/http"
)

func readConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()

}

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

var log = logrus.New()

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users", getUserHandler).Methods("GET")
	r.HandleFunc("/users/{id}", getUserHandler).Methods("GET")
	r.HandleFunc("/users", createUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", deleteUserHandler).Methods("DELETE")
	r.HandleFunc("/users/{id}", updateUserHandler).Methods("PUT")

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "启动web服务",
		Run: func(cmd *cobra.Command, args []string) {
			readConfig() // 读配置文件
			port := viper.GetString("PORT")
			fmt.Printf("running server on http://localhost%s\n", port)
			e := http.ListenAndServe(port, r)
			if e != nil {
				fmt.Println(e)
			}
		},
	}
	rootCmd := &cobra.Command{Use: "muxCli"}
	rootCmd.AddCommand(startCmd)
	rootCmd.Execute()
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
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

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	cli, ctx, _ := initMongo()

	user := types.User{}
	json.Unmarshal(body, &user)
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

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
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

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
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
