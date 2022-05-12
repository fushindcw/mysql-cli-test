package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	host        string
	port        string
	username    string
	password    string
	showVersion bool
	database    string
)

var languageMatch = make(map[string]func(string))

func init() {
	languageMatch["select"] = query
	languageMatch["help"] = help
	languageMatch["h"] = help
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "host")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "P", "3306", "port")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "V", false, "app version")
	rootCmd.PersistentFlags().StringVarP(&username, "user", "U", "", "login username")
	rootCmd.PersistentFlags().StringVarP(&password, "pwd", "", "Tv=9O9k:NlmB.s3+", "login password")
	rootCmd.PersistentFlags().StringVarP(&database, "db", "D", "xiaodu_v2_dev", "database")
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
	dbConnectStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
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
	if password == "" {
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
	} else {
		createDb()
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
			if result == "" {
				help(result)
				continue
			}
			startStr := strings.Split(result, " ")[0]
			startStr = strings.ToLower(startStr)
			fnc := languageMatch[startStr]
			if fnc != nil {
				fnc(result)
			} else {
				execute(result)
			}
		}
	}
}

func execute(sql string) {
	rs, err := db.Exec(sql)
	if err != nil {
		fmt.Println(err)
		return
	}
	effectRow, err := rs.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%d rows effect", effectRow)
}

func query(sql string) {
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println(err)
		return
	}
	columnTypes, _ := rows.ColumnTypes()
	for _, v := range columnTypes {
		fmt.Printf("%s\t", v.Name())
	}
	fmt.Println()
	for rows.Next() {
		s := make([]interface{}, len(columnTypes))
		for i, _ := range columnTypes {
			s[i] = new(string)
		}
		rows.Scan(s...)
		for _, v := range s {
			fmt.Printf("%s\t", *(v.(*string)))
		}
		fmt.Println()
	}
}

func help(string) {
	fmt.Printf("%s\t%s\n", "help", "help options")
	fmt.Printf("%s\t%s\n", "update", "update options")
	fmt.Printf("%s\t%s\n", "delete", "delete options")
	fmt.Printf("%s\t%s\n", "insert", "insert options")
	fmt.Printf("%s\t%s\n", "select", "select options")
	fmt.Printf("%s\t%s\n", "exit", "exit app")
}
