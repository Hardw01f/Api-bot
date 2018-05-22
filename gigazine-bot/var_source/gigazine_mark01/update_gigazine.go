package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"os/exec"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//API叩いてjsonを取得
	resp, err := http.Get("http://gigazine.net/")
	checkErr(err)

	//処理終了後にrespをclose
	defer resp.Body.Close()

	//jsonのbody部を取得
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	content := string(body)

	//fmt.Println(content)
	//collecturls(content)

	//contentをトリミング
	splitedContent := strings.Split(content, "<h2>")
	head := splitedContent[0]
	splitedhead := strings.Split(head, "<title>")

	title := splitedhead[1]
	splitedtitle := strings.Split(title, "</title>")

	hihun := splitedtitle[0]
	//topic名が取れていない時はpanicを起こす
	splitedhihun := strings.Split(hihun,"-")
	if len(splitedhihun) != 2{
			panic("panic occured")
		}

	//fmt.Println(splitedhead)

	//fmt.Println(splitedhihun)
	//fmt.Println(splitedhihun[1])

	topicTitle := splitedhihun[1]

	//fmt.Println(topicTitle)

	//urlを文字列として取得
	url := collecturls(content)


	//dbに接続,deferで処理終了時にclose
	db, err := sql.Open("sqlite3", "/home/tukachan/gigazine_test/gigazinetopic")
	    checkErr(err)

	//古いtopicnameをdbから取得
	var topicname string
	err = db.QueryRow("SELECT topicname FROM topic").Scan(&topicname)
	    checkErr(err)
	defer db.Close()

	//fmt.Println(topicname)

	//古いtopicnameと新しいtopicが異なっていれば新しいtopicnameをdbに更新する
	if topicTitle != topicname{
		stmt, err := db.Prepare("update topic set topicname=? where topicname=?")
		_, err = stmt.Exec(topicTitle, topicname)
		checkErr(err)
		fmt.Println("updated!!!!")
		Message := trimname(topicTitle,url)
		fmt.Println(Message)
		//send(Name)
	}else if topicTitle == topicname{
		fmt.Println("same topic")
	}else{
		fmt.Println("error")
	}

}

//errorチェック用の関数
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

//topicname比較用の関数、未使用
func jadge(topic string, Newtopic string) string {
	if Newtopic == topic {
		return topic
	} else {
		return Newtopic
	}
}


//slackに通知を出すようの関数
func send(topic string){
		out, err := exec.Command("curl", "-XPOST", "-d", "token=xoxp-211441493619-211441493731-249801662402-198194d20ffa667eb2d4515b42166f50", "-d", "channel=#updatetopic", "-d", "text="+topic+" ", "-d", "username=NewTopic", "https://slack.com/api/chat.postMessage","-d", "icon_url=https://pbs.twimg.com/profile_images/876998955627827204/IuoxMaM2_400x400.jpg").Output()
		checkErr(err)
		fmt.Println(string(out))
}


//url取得用の関数
func collecturls(body string) string {
		var splited []string = strings.Split(body,"<div class=\"card\">")
		var splitedthumb []string = strings.Split(splited[1],"<div class=\"thumb\">")
		var splitedhref []string = strings.Split(splitedthumb[1],">")
		var splitedhihun []string = strings.Split(splitedhref[0],"\"")
		var url string = splitedhihun[1]
		//fmt.Println(url)
		return url
}

//トリミング用の関数
func trimname(name string , url string) string {
		var trimedname string = "New topic name :: " + name +" "+ url + " "
		return trimedname
}
