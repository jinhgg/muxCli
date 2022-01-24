/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"muxCli/router"
	"net/http"
)

// startCmd represents the start command
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "启动web服务",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// 绑定路由
		r := router.NewRouter()
		port := viper.GetString("PORT")
		fmt.Printf("running server on http://localhost%s\n", port)
		if e := http.ListenAndServe(port, r); e != nil {
			fmt.Println(e)
		}
	},
}

func init() {
	RootCmd.AddCommand(StartCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
