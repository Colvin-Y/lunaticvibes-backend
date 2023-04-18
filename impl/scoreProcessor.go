package processor

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (sp *ScoreProcessor) Insert() error {
	// 设置数据库连接选项
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// 连接到 MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		sp.Logger.Error(err.Error())
		return err
	}

	// 检查连接是否成功
	err = client.Ping(context.Background(), nil)
	if err != nil {
		sp.Logger.Error(err.Error())
		return err
	}

	// 获取指定数据库的集合对象 数据库为 public， 表名为 mytable
	collection := client.Database("public").Collection("mytable")

	// 将结构体转换为文档
	doc, err := bson.Marshal(*sp.data)
	if err != nil {
		sp.Logger.Error(err.Error())
		return err
	}

	// 在集合中插入文档
	_, err = collection.InsertOne(context.Background(), doc)
	if err != nil {
		sp.Logger.Error(err.Error())
		return err
	}

	fmt.Println("插入成功！")
	return nil
}

func doScoreSyncPrecheck(w http.ResponseWriter, r *http.Request, sp *ScoreProcessor) error {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		sp.Logger.Error(fmt.Sprintf("invalid method %v", http.StatusMethodNotAllowed))
		return fmt.Errorf("invalid method")
	}

	return nil
}

func doScoreSyncDataParse(w http.ResponseWriter, r *http.Request, sp *ScoreProcessor) error {
	// 解析表单数据
	userID := r.FormValue("userID")
	songHash := r.FormValue("songHash")
	clearType := r.FormValue("clearType")
	score := r.FormValue("score")
	scorePG := r.FormValue("scorePG")
	scoreGR := r.FormValue("scorePG")
	scoreGD := r.FormValue("scorePG")
	scoreBD := r.FormValue("scorePG")
	scorePR := r.FormValue("scorePG")
	combo := r.FormValue("combo")
	laneOp := r.FormValue("laneOp")
	gaugeOp := r.FormValue("gaugeOp")
	inputType := r.FormValue("inputType")
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		sp.Logger.Error(err.Error())
		return err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		sp.Logger.Error(err.Error())
		return err
	}

	data := &ScoreData{
		UserID:         string(userID),
		SongHash:       string(songHash),
		ClearType:      string(clearType),
		Score:          string(score),
		ScorePG:        string(scorePG),
		ScoreGR:        string(scoreGR),
		ScoreGD:        string(scoreGD),
		ScoreBD:        string(scoreBD),
		ScorePR:        string(scorePR),
		Combo:          string(combo),
		LaneOp:         string(laneOp),
		GaugeOp:        string(gaugeOp),
		InputType:      string(inputType),
		ReplayFileData: string(content),
	}

	err = Validate(*data)
	if err != nil {
		sp.Logger.Error(err.Error())
		return err
	} else {
		sp.data = data
	}

	return nil
}

func doScoreSyncDataSave(w http.ResponseWriter, r *http.Request, sp *ScoreProcessor) error {
	sp.Logger.Info(fmt.Sprintf("start to inser data[%v] to mongo", *sp.data))
	err := sp.Insert()
	if err != nil {
		sp.Logger.Error(err.Error())
		return err
	}

	return nil
}

func doScoreSyncFinish(w http.ResponseWriter, r *http.Request, sp *ScoreProcessor) error {
	// http resp
	w.Header().Set("Content-Type", "text/plain")
	return nil
}

func (sp *ScoreProcessor) ScoreHandlerSet(w http.ResponseWriter, r *http.Request) {
	handlers := []ScoreProcessorHandlerFunc{
		doScoreSyncPrecheck,
		doScoreSyncDataParse,
		doScoreSyncDataSave,
		doScoreSyncFinish,
	}
	sp.Sync(w, r, handlers...)
}
