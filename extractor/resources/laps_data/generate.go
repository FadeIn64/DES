package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var urlFormat = "https://api.openf1.org/v1/laps?session_key=%d&driver_number=%d"
var sessionKey = 9566

//go:generate go run generate.go
func main() {
	for i := 1; i < 100; i++ {
		response, err := getData(sessionKey, i)
		if err != nil {
			log.Fatal(err)
		}
		if len(response) > 10 {
			err = writeResponse(response, strconv.Itoa(i))
			if err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func writeResponse(body []byte, filePrefix string) error {
	// Запись ответа в файл
	err := os.WriteFile("./"+filePrefix+".json", body, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getData(sessionKey int, driveNumber int) ([]byte, error) {

	url := fmt.Sprintf(urlFormat, sessionKey, driveNumber)

	// Выполнение GET-запроса
	response, err := http.Get(url)
	log.Print("send response to ", url, "\n")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Проверка успешности запроса
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Http response with status: " + response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
