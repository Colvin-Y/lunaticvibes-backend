package processor

import (
	"errors"
	"fmt"

	logger "github.com/Colvin-Y/lunaticvibes-backend/common/log"
	"github.com/go-playground/validator/v10"
)

// 定义 interface
type ScoreProcessor struct {
	Logger *logger.Logger
	data   *ScoreData
}

type SignUpProcessor struct {
	Logger *logger.Logger
	data   *UserData
}

// 数据体
type SongInfo struct {
	SongHash string `bson:"SongHash,omitempty" validate:"required,string,len=32"`
	//SongID     string
	SongName   string `bson:"songName" validate:"required,string"`
	SongArtist string `bson:"songArtist" validate:"required,string"`
	MaxBpm     string `bson:"maxBpm" validate:"required,string"`
	MinBpm     string `bson:"minBpm" validate:"required,string"`
	SongTotal  string `bson:"songTotal" validate:""`
	TotalNotes string `bson:"totalNotes" validate:"required,string"`
	SongLevel  string `bson:"songLevel" validate:"required,string"`
	Keys       string `bson:"keys" validate:"required,string"`
	JudgeRank  string `bson:"judgeRank" validate:"required,string"`
}

type ScoreData struct {
	UserID         string `bson:"userID" validate:"required,string,min=1,max=10"`
	IsCourse       string `bson:"isCourse" validate:"required,oneof=0 1"`
	SongHash       string `bson:"SongHash,omitempty" validate:"required_if=IsCourse 0|md5"`
	CourseHash     string `bson:"CourseHash,omitempty" validate:"required_if=IsCourse 1|md5"`
	ClearType      string `bson:"clearType" validate:"required,string"`
	LnMode         string `bson:"lnMode,omitempty" `
	Score          string `bson:"score" validate:"required,string"`
	ScoreMax       string `bson:"scoreMax" validate:"required,string"`
	ScoreRate      string `bson:"scoreRate" validate:"required,string"`
	ScoreRank      string `bson:"scoreRank" validate:"required,string"`
	ScorePG        string `bson:"scorePG" validate:"required,string,min=0"`
	ScoreGR        string `bson:"scoreGR" validate:"required,string,min=0"`
	ScoreGD        string `bson:"scoreGD" validate:"required,string,min=0"`
	ScoreBD        string `bson:"scoreBD" validate:"required,string,min=0"`
	ScorePR        string `bson:"scorePR" validate:"required,string,min=0"`
	Combo          string `bson:"combo" validate:"required,string,min=0"`
	LaneOp         string `bson:"laneOp" validate:"required,string"`
	GaugeOp        string `bson:"gaugeOp" validate:"required,string"`
	InputType      string `bson:"inputType" validate:"required,string"`
	UpdateTime     string `bson:"updateTime"`
	ReplayFileData string `bson:"fileData"`
}

type CourseData struct {
	CourseHash string `bson:"courseHash" validate:"required,string,len=32"`
	//CourseID   string
	Songs []string
}

type UserData struct {
	UserID       string `bson:"userID" validate:"required,string"`
	UserName     string `bson:"userName" validate:"required,string"`
	UserNickname string `bson:"userNickname" validate:"required,string"`
	UserEmail    string `bson:"userEmail" validate:"required,email"`
	UserPwd      string `bson:"userPwd" validate:"required,string"`
}

func Validate(data interface{}) error {
	// 新建一个 validator 实例
	v := validator.New()

	// 注册 "string" tag 的验证器
	err := v.RegisterValidation("string", func(fl validator.FieldLevel) bool {
		_, ok := fl.Field().Interface().(string)
		return ok
	})
	if err != nil {
		return err
	}

	/* 使用反射方法对输入的 struct 类型进行验证
	value := reflect.ValueOf(data)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		tag := value.Type().Field(i).Tag.Get("validate")
		if tag != "" {
			err = v.Var(field.Interface(), tag) //会出现invalid field namespace
			if err != nil {
				// 输出验证失败的字段和错误信息
				fieldName := value.Type().Field(i).Name
				errMsg := fmt.Sprintf("Validation error for '%s.%s': %s", value.Type().Name(), fieldName, err.Error())
				return errors.New(errMsg)
			}
		}
	}*/

	//直接对struct进行验证
	err = v.Struct(data)
	if err != nil {
		//输出失败信息
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range ve {
				errMsg := fmt.Sprintf("Validation error for '%s.%s': %s", fieldError.Field(), fieldError.Tag(), fieldError.Error())
				return errors.New(errMsg)
			}
		}
	}
	// 验证通过，返回 nil
	return nil
}
