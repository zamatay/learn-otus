package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
)

const event = `{
    "Date": "1708889494",
    "DateInterval": 100000,
    "Description": "Тестовое событие описание",
    "Title": "Тестовое событие",
    "UserID": 1
}`

type Event struct {
	Title        string `db:"title" json:"Title"`
	Date         string `db:"date" json:"Date"`
	DateInterval int32  `db:"date_interval" json:"DateInterval"`
	Description  string `db:"description" json:"Description"`
	UserID       int32  `db:"user_id" json:"UserID"`
}

var (
	errorRead   = errors.New("Error reading response body")
	errorSend   = errors.New("Error sending request")
	errorCreate = errors.New("Error creating request")
)

func main() {
	_, err := addEvent()
	if err != nil {
		os.Exit(1)
	}

	res, err := getEvent(5)
	if err != nil {
		os.Exit(1)
	}

	getEvent := Event{}
	if err := json.Unmarshal(res, &getEvent); err != nil {
		os.Exit(1)
	}

	sendEvent := Event{}
	if err := json.Unmarshal([]byte(event), &sendEvent); err != nil {
		os.Exit(1)
	}

	if !reflect.DeepEqual(sendEvent, getEvent) {
		os.Exit(1)
	}

	os.Exit(0)
}

func addEvent() ([]byte, error) {
	url := "http://0.0.0.0:8080/event/add"

	// Подготовка данных для отправки (в формате JSON, например)
	requestBody := []byte(event)

	// Создание HTTP-запроса типа POST с указанием URL и данных для отправки
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, errorCreate
	}

	// Установка заголовка Content-Type для указания типа содержимого (JSON в данном случае)
	request.Header.Set("Content-Type", "application/json")

	// Создание клиента HTTP
	client := &http.Client{}

	// Отправка запроса
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, errorSend
	}
	defer response.Body.Close()

	// Чтение ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errorRead
	}

	return body, nil
}
func getEvent(id int) ([]byte, error) {
	url := "http://0.0.0.0:8080/event/get/%d"
	client := &http.Client{}

	// Создание запроса
	req, err := http.NewRequest("GET", fmt.Sprintf(url, id), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, errorCreate
	}

	// Выполнение запроса и получение ответа
	resp, err := client.Do(req)
	if err != nil {
		return nil, errorSend
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errorRead
	}

	return body, nil

}
