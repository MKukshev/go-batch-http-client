package app

import (
	"bufio"
	"bytes"

	"context"
	"encoding/json"
	"fmt"
	. "go-batch-http-client/model"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/time/rate"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func sendRequest(client *http.Client, cfg *Config, jsonStr *string, jsonBytes *[]byte) {
	// Создаем HTTP POST запрос
	req, err := http.NewRequest(cfg.Server.Req, cfg.Server.Url, bytes.NewBuffer(*jsonBytes))
	if err != nil {
		log.Println("Ошибка создания HTTP запроса:", err)
		return
	}

	// Устанавливаем заголовки, если необходимо
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка отправки HTTP запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ сервера и выводим его
	var responseBody bytes.Buffer
	_, err = io.Copy(&responseBody, resp.Body)
	if err != nil {
		log.Println("Ошибка чтения ответа сервера:", err)
		return
	}
	log.Println(fmt.Sprintf("json: %s resp: %s", *jsonStr, resp.Status))

}

func Run(cfg Config) {

	//	apiUrl := "http://192.168.101.235:8080/svetets/orchestrator/speech/payment/offline"
	//	jsonFilesPath := "result/http-38_json.txt"
	// Открываем дирректорию с файлами

	err := filepath.Walk(cfg.Files.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}
		if info.IsDir() {
			return nil // Пропускаем поддиректории
		}
		fileName := info.Name()
		fullPath := filepath.Join(cfg.Files.Path, fileName)
		fmt.Println(fmt.Sprintf("Обработка файла: %s", fileName))
		file, err := os.Open(fullPath)
		if err != nil {
			log.Println("Ошибка открытия файла:", err)
			return nil
		}
		defer file.Close()

		// Создаем HTTP клиента
		client := &http.Client{}
		limiter := rate.NewLimiter(rate.Limit(cfg.Server.Limiter), 1)

		// Читаем файл и отправляем каждую JSON строку как HTTP POST запрос
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			limiter.Wait(context.TODO())
			// Считываем JSON строку
			jsonStr := scanner.Text()

			// Создаем структуру для разбора JSON
			var jsonData interface{}
			err := json.Unmarshal([]byte(jsonStr), &jsonData)
			if err != nil {
				log.Println("Ошибка разбора JSON:", err)
				continue
			}

			// Преобразуем JSON в байты для отправки
			jsonBytes, err := json.Marshal(jsonData)
			if err != nil {
				log.Println("Ошибка преобразования JSON в байты:", err)
				continue
			}

			go sendRequest(client, &cfg, &jsonStr, &jsonBytes)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Ошибка чтения файла:", err)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Ошибка при чтении директории: %s", err)
	}

}
