package utils

import (
	"bytes"
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct[T any](data []byte) error {
	validate := validator.New()

	// JSONからT型に変換
	var obj T
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&obj)
	if err != nil {
		return err
	}

	// T型でバリデーション
	return validate.Struct(&obj)
}

func ValidateStructTwoWay[oneT any, twoT any](data *oneT) error {
	validate := validator.New()

	// oneT型からJSONに変換
	oneJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// JSONからtwoT型に変換
	var twoObj twoT
	decoder := json.NewDecoder(bytes.NewReader(oneJson))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&twoObj)
	if err != nil {
		return err
	}

	// twoT型でバリデーション
	err = validate.Struct(&twoObj)
	if err != nil {
		return err
	}

	// twoT型からJSONに変換
	twoJson, err := json.Marshal(&twoObj)
	if err != nil {
		return err
	}

	// JSONからoneT型に変換
	var oneObj oneT
	decoder = json.NewDecoder(bytes.NewReader(twoJson))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&oneObj)
	if err != nil {
		return err
	}

	// oneT型でバリデーション
	return validate.Struct(&oneObj)
}
