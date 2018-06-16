package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"io/ioutil"
	"os"
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"
)

//表情解码
func UnicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

//表情转换
func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

func FilterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		if value < 65533 {
			new_content += string(value)
		}
		//size := utf8.RuneCountInString(string(value))
		//if size <= 3 {
		//	new_content += string(value)
		//	}
		//}
	}
	return new_content
}

func main() {
	db, err := sql.Open("mysql", "root:znhl2017UP@tcp(tlwl2020.mysql.rds.aliyuncs.com:3686)/fbox?charset=utf8mb4")
	if err != nil {
		panic(err.Error())
	}
	sqlstr := "select id,content,comments from aaa"
	defer db.Close()

	rows, err := db.Query(sqlstr)
	defer rows.Close()
	var rs string
	for rows.Next() {
		var id string
		var content string
		var comments string
		err = rows.Scan(&id, &content, &comments)
		rs += id + "-" + content + "-" + comments + "\n"
	}
	fmt.Printf(FilterEmoji(rs))
	WriteWithIo("data.sql", FilterEmoji(rs))
	err = rows.Err()
	if err != nil {
		panic(err.Error())
	}
}

func WriteWithIo(name,content string) {
	fileObj,err := os.OpenFile(name,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	if err != nil {
		fmt.Println("Failed to open the file",err.Error())
		os.Exit(2)
	}
	if  _,err := io.WriteString(fileObj,content);err == nil {
		fmt.Println("Successful appending to the file with os.OpenFile and io.WriteString.")
	}
}


func WriteWithIoutil(name,content string) {
	data :=  []byte(content)
	if ioutil.WriteFile(name,data,0644) == nil {
		fmt.Println("写入文件成功:",content)
	}
}

func WriteWithFileWrite(name,content string){
	fileObj,err := os.OpenFile(name,os.O_RDWR|os.O_CREATE|os.O_TRUNC,0644)
	if err != nil {
		fmt.Println("Failed to open the file",err.Error())
		os.Exit(2)
	}
	defer fileObj.Close()
	if _,err := fileObj.WriteString(content);err == nil {
		fmt.Println("Successful writing to the file with os.OpenFile and *File.WriteString method.",content)
	}
	contents := []byte(content)
	if _,err := fileObj.Write(contents);err == nil {
		fmt.Println("Successful writing to thr file with os.OpenFile and *File.Write method.",content)
	}
}

func WriteWithBufio(name,content string) {
	if fileObj,err := os.OpenFile(name,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644);err == nil {
		defer fileObj.Close()
		writeObj := bufio.NewWriterSize(fileObj,4096)
		//
		if _,err := writeObj.WriteString(content);err == nil {
			fmt.Println("Successful appending buffer and flush to file with bufio's Writer obj WriteString method",content)
		}

		//使用Write方法,需要使用Writer对象的Flush方法将buffer中的数据刷到磁盘
		buf := []byte(content)
		if _,err := writeObj.Write(buf);err == nil {
			fmt.Println("Successful appending to the buffer with os.OpenFile and bufio's Writer obj Write method.",content)
			if  err := writeObj.Flush(); err != nil {panic(err)}
			fmt.Println("Successful flush the buffer data to file ",content)
		}
	}
}


