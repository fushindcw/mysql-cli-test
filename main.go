package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"os"
)

var (
	host        string
	port        string
	username    string
	password    string
	showVersion bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "host")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "P", "3306", "port")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "V", false, "app version")
	rootCmd.PersistentFlags().StringVarP(&username, "user", "U", "", "login username")
}

const maxFailedTimes = 5

var rootCmd = &cobra.Command{
	Use:   "mysql-cli-test",
	Short: "demo04_short",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println("Lua 5.4.4  Copyright (C) 1994-2022 Lua.org, PUC-Rio")
			return nil
		}
		createJh()
		return nil
	},
}

var db *sql.DB
var tx *sql.Tx

var templates = &promptui.PromptTemplates{
	Prompt:  "{{ . }} ",
	Valid:   "{{ . | green }} ",
	Invalid: "{{ . | red }} ",
	Success: "{{ . | bold }} ",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func createDb() bool {
	dbConnectStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, host, port)
	fmt.Println(dbConnectStr)
	var err error
	db, err = sql.Open("mysql", dbConnectStr)
	if err != nil {
		fmt.Println(err)
		fmt.Println("数据库连接失败")
		return false
	} else {
		err = db.Ping()
		if err != nil {
			fmt.Println(err)
			fmt.Println("数据库连接失败")
			return false
		} else {
			return true
		}
	}
}

func createJh() {
	//判断用户名是否输入
	if username == "" {
		prompt := promptui.Prompt{
			Label:     "请输入用户名: ",
			Templates: templates,
		}
		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}
		username = result
	}
	failedTimes := 0
	//输入密码
	for {
		if failedTimes == maxFailedTimes {
			fmt.Printf("密码错误次数到达%d次\n", maxFailedTimes)
			os.Exit(1)
		}
		prompt := promptui.Prompt{
			Label:     "请输入密码: ",
			Templates: templates,
		}
		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}
		password = result
		if createDb() {
			break
		} else {
			failedTimes++
		}
	}
	for {
		prompt := promptui.Prompt{
			Label:     ">",
			Templates: templates,
		}
		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}
		if result == "exit" {
			db.Close()
			fmt.Println("bye!")
			os.Exit(1)
		} else {
			execute(result)
		}
	}
}

func execute(sql string) {
	rs, err := db.Exec(sql)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(&rs)
}
