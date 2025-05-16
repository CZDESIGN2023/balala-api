package dwd

import (
	"go-cs/internal/dwh/pkg/model"
	"reflect"
	"testing"
	"time"
)

func Test_structFieldNames(t *testing.T) {
	names := structFieldNames(DwdWitem{})
	t.Log(names)
}

func Test_compareStructFields(t *testing.T) {
	a := &DwdWitem{
		DimModel: model.DimModel{
			GmtCreate: time.Now(),
		},
		ChainModel: model.ChainModel{},
		SpaceId:    1,
	}

	b := DwdWitem{
		DimModel:   model.DimModel{},
		ChainModel: model.ChainModel{},
		SpaceId:    0,
	}

	if compareStructFields(a, b, modelFieldNames) {
		t.Log("same")
	} else {
		t.Log("not same")
	}

}

func Test(t *testing.T) {
	va := reflect.ValueOf(DwdWitem{})
	for i := 0; i < va.NumField(); i++ {
		fieldName := va.Type().Field(i).Name
		t.Log(fieldName)
	}
}
