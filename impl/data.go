package processor

import (
	"errors"
	"fmt"
	"reflect"

	logger "github.com/Colvin-Y/lunaticvibes-backend/common/log"
	"gopkg.in/go-playground/validator.v9"
)

// 定义 inerface
type ScoreProcessor struct {
	Logger *logger.Logger
	data   *ScoreData
}

type SignUpProcessor struct {
	Logger *logger.Logger
	// data: 小写
}

// 数据体
type ScoreData struct {
	UserID         string `bson:"userID" validate:"required,string,min=2,max=10"`
	SongHash       string `bson:"SongHash" validate:"required,string,len=32"`
	ClearType      string `bson:"clearType" validate:"required,oneof='FAILED' 'CLEAR' 'ASSIST CLEAR' 'EASY CLEAR' 'HARD CLEAR' 'EX-HARD CLEAR' 'FULLCOMBO'"`
	Score          string `bson:"score" validate:"required,string"`
	ScorePG        string `bson:"scorePG"`
	ScoreGR        string `bson:"scoreGR"`
	ScoreGD        string `bson:"scoreGD"`
	ScoreBD        string `bson:"scoreBD"`
	ScorePR        string `bson:"scorePR"`
	Combo          string `bson:"combo"`
	LaneOp         string `bson:"laneOp"`
	GaugeOp        string `bson:"gaugeOp"`
	InputType      string `bson:"inputType"`
	ReplayFileData string `bson:"fileData"`
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

	// 使用反射方法对输入的 struct 类型进行验证
	value := reflect.ValueOf(data)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		tag := value.Type().Field(i).Tag.Get("validate")
		if tag != "" {
			err = v.Var(field.Interface(), tag)
			if err != nil {
				// 输出验证失败的字段和错误信息
				fieldName := value.Type().Field(i).Name
				errMsg := fmt.Sprintf("Validation error for '%s.%s': %s", value.Type().Name(), fieldName, err.Error())
				return errors.New(errMsg)
			}
		}
	}

	// 验证通过，返回 nil
	return nil
}
