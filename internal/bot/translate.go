package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LibreTranslator struct {
	APIURL string
	APIKey string
}

type Translator interface {
	Translate(text string, fromLang string, toLang string) (string, error)
}

// структура для тела запроса к API LibreTranslator
type translateRequest struct {
	Q      string `json:"q"`
	Source string `json:"source"`
	Target string `json:"target"`
	Format string `json:"format"`
	APIKey string `json:"api_key"`
}

type translateResponse struct {
	TranslatedText string `json:"translatedText"`
}

func (lt *LibreTranslator) Translate(text string, sourceLang string, targetLang string) (string, error) {
	requestBody := translateRequest{
		Q:      text,
		Source: sourceLang,
		Target: targetLang,
		Format: "text",
		APIKey: lt.APIKey,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("Ошибка кодирования JSON: %w", err)
	}

	resp, err := http.Post(lt.APIURL + "/translate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("Ошибка сетевого запроса: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API вернуло ошибку: %s", resp.Status)
	}

	// 4. Прочитай и расшифруй ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var translateResp translateResponse
	err = json.Unmarshal(body, &translateResp)
	if err != nil {
		return "", fmt.Errorf("ошибка расшифровки JSON ответа: %w", err)
	}

	// 5. Верни переведенный текст
	return translateResp.TranslatedText, nil

}

func NewLibreTranslator(apiURL string, apiKey string) *LibreTranslator {
	return &LibreTranslator{
		APIURL: apiURL,
		APIKey: apiKey,
	}
}
