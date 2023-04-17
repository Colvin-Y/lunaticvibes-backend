package processor
import (
	"context"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type MyData struct {
	Score    string `bson:"score"`
    FileData string `bson:"fileData"`
}

func (d *MyData) Insert() {
    // 设置数据库连接选项
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

    // 连接到 MongoDB
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // 检查连接是否成功
    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }

    // 获取指定数据库的集合对象 数据库为 public， 表名为 mytable
    collection := client.Database("public").Collection("mytable")

    // 将结构体转换为文档
    doc, err := bson.Marshal(d)
    if err != nil {
        log.Fatal(err)
    }

    // 在集合中插入文档
    _, err = collection.InsertOne(context.Background(), doc)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("插入成功！")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// 仅处理 POST 请求
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	// 解析表单数据
    score := r.FormValue("score")
    file, _, err := r.FormFile("file")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()
	content, err := ioutil.ReadAll(file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	// 将表单数据写到数据库
	data := &MyData{
		Score: string(score),
		FileData:  string(content),
	}
	data.Insert()

    // http resp
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte(string(score) +  string(content)))
}