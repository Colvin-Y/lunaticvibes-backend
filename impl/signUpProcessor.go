package processor

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (sup *SignUpProcessor) signUpDataInsert() error {
	// 设置数据库连接选项
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// 连接到 MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		sup.Logger.Error(err.Error())
		return err
	}

	// 检查连接是否成功
	err = client.Ping(context.Background(), nil)
	if err != nil {
		sup.Logger.Error(err.Error())
		return err
	}

	// 获取指定数据库的集合对象 数据库为 public， 表名为 user
	collection := client.Database("public").Collection("user")

	// 将结构体转换为文档
	doc, err := bson.Marshal(*sup.data)
	if err != nil {
		sup.Logger.Error(err.Error())
		return err
	}

	// 在集合中插入文档
	_, err = collection.InsertOne(context.Background(), doc)
	if err != nil {
		sup.Logger.Error(err.Error())
		return err
	}

	fmt.Println("插入成功！")
	return nil
}

// 生成新的 UserID(unfinished)
func (sup *SignUpProcessor) generateUserID() (int, error) {
	// 设置数据库连接选项
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// 连接到 MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		sup.Logger.Error(err.Error())
		return 0, err
	}

	// 检查连接是否成功
	err = client.Ping(context.Background(), nil)
	if err != nil {
		sup.Logger.Error(err.Error())
		return 0, err
	}

	// 获取指定数据库的集合对象 数据库为 public， 表名为 user
	collection := client.Database("public").Collection("user")

	// 获取计数器文档并更新计数器
	filter := bson.M{"_id": "userCounter"}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var doc bson.M
	err = collection.FindOneAndUpdate(context.Background(), filter, update, options).Decode(&doc)
	if err != nil {
		sup.Logger.Error(err.Error())
		return 0, err
	}
	// 从计数器文档中读取计数器值并返回新的 UserID
	seq, ok := doc["seq"].(int32)
	if !ok {
		sup.Logger.Error(err.Error())
		return 0, err
	}
	return int(seq), nil
}

func doSighUpSyncPrechck(w http.ResponseWriter, r *http.Request, sup *SignUpProcessor) error {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		sup.Logger.Error(fmt.Sprintf("invalid method %v", http.StatusMethodNotAllowed))
		return fmt.Errorf("invalid method")
	}
	return nil
}

func doSighUpSyncDataParse(w http.ResponseWriter, r *http.Request, sup *SignUpProcessor) error {
	userID, err := sup.generateUserID()
	userName := r.FormValue("userName")
	userNickname := r.FormValue("userNickname")
	userEmail := r.FormValue("userEmail")
	userPwd := r.FormValue("userPwd")
	if err != nil {
		sup.Logger.Error(err.Error())
		return err
	}

	data := &UserData{
		UserID:       strconv.Itoa(int(userID)),
		UserName:     string(userName),
		UserNickname: string(userNickname),
		UserEmail:    string(userEmail),
		UserPwd:      string(userPwd),
	}

	err = Validate(*data)
	if err != nil {
		sup.Logger.Error(err.Error())
		return err
	} else {
		sup.data = data
	}

	return nil
}

func doSighUpSyncDataSave(w http.ResponseWriter, r *http.Request, sup *SignUpProcessor) error {
	sup.Logger.Info(fmt.Sprintf("start to insert data[%v] to mongo", *sup.data))
	err := sup.signUpDataInsert()
	if err != nil {
		sup.Logger.Error(err.Error())
		return err
	}
	return nil
}

func doSighUpSyncFinish(w http.ResponseWriter, r *http.Request, sup *SignUpProcessor) error {
	w.Header().Set("Content-Type", "text/plain")
	return nil
}

func (sup *SignUpProcessor) SignUpHandlerSet(w http.ResponseWriter, r *http.Request) {
	handlers := []SignUpProcessorHandlerFunc{
		doSighUpSyncPrechck,
		doSighUpSyncDataParse,
		doSighUpSyncDataSave,
		doSighUpSyncFinish,
	}
	sup.Sync(w, r, handlers...)
}
