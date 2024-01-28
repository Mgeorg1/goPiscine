package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Шаг 1: Получение и распаковка JSON-ответа для получения токена
	token, err := requestToken("http://localhost:8800/api/get_token")
	if err != nil {
		fmt.Printf("Error getting token: %s\n", err)
		return
	}

	// Шаг 2: Использование JWT-токена для запроса ресурса
	apiURL := "http://localhost:8800/api/recommend?lat=55.680702&lon=37.608534"
	response, err := makeRequestWithToken(apiURL, token)
	if err != nil {
		fmt.Printf("Error making request: %s\n", err)
		return
	}

	// Вывод ответа
	fmt.Printf("Response from %s:\n%s\n", apiURL, response)
}

// Шаг 1: Запрос и распаковка JSON-ответа для получения токена
func requestToken(tokenURL string) (string, error) {
	response, err := http.Get(tokenURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// Распаковка JSON-ответа
	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return "", err
	}

	// Извлечение токена из JSON
	tokenValue, ok := jsonResponse["token"].(string)
	if !ok {
		return "", fmt.Errorf("Token not found in JSON response")
	}

	return tokenValue, nil
}

// Шаг 2: Использование JWT-токена для запроса ресурса
func makeRequestWithToken(apiURL, token string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	// Добавление JWT-токена в заголовок Authorization
	req.Header.Set("Authorization", "Bearer "+token)

	// Выполнение запроса
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
