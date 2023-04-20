package processor

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (sp *ScoreProcessor) scoreInsert() error {
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
	isCourse := r.FormValue("isCourse")
	songHash := r.FormValue("songHash")
	courseHash := r.FormValue("courseHash")
	clearType := r.FormValue("clearType")
	lnMode := r.FormValue("lnMode")
	score := r.FormValue("score")
	scoreMax := r.FormValue("scoreMax")
	scoreRate := r.FormValue("scoreRate")
	scoreRank := r.FormValue("scoreRank")
	scorePG := r.FormValue("scorePG")
	scoreGR := r.FormValue("scoreGR")
	scoreGD := r.FormValue("scoreGD")
	scoreBD := r.FormValue("scoreBD")
	scorePR := r.FormValue("scorePR")
	combo := r.FormValue("combo")
	laneOp := r.FormValue("laneOp")
	gaugeOp := r.FormValue("gaugeOp")
	inputType := r.FormValue("inputType")
	updateTime := time.Now().Format("2006-01-02 15:04:05")
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
		IsCourse:       string(isCourse),
		SongHash:       string(songHash),
		CourseHash:     string(courseHash),
		ClearType:      string(clearType),
		LnMode:         string(lnMode),
		Score:          string(score),
		ScoreMax:       string(scoreMax),
		ScoreRate:      string(scoreRate),
		ScoreRank:      string(scoreRank),
		ScorePG:        string(scorePG),
		ScoreGR:        string(scoreGR),
		ScoreGD:        string(scoreGD),
		ScoreBD:        string(scoreBD),
		ScorePR:        string(scorePR),
		Combo:          string(combo),
		LaneOp:         string(laneOp),
		GaugeOp:        string(gaugeOp),
		InputType:      string(inputType),
		UpdateTime:     string(updateTime),
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
	sp.Logger.Info(fmt.Sprintf("start to insert data[%v] to mongo", *sp.data))
	err := sp.scoreInsert()
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
