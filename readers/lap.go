package readers

import (
	"DES/models"
	"encoding/json"
	"os"
	"strings"
)

var lapsPath = "./resources/laps_data"

func ReadLapsData() ([]models.Lap, error) {

	var res []models.Lap

	// Получаем список файлов в папке
	files, err := os.ReadDir(lapsPath)
	if err != nil {
		return nil, err
	}

	// Проходим по каждому файлу
	for _, file := range files {
		if file.IsDir() { // Проверяем, что это не папка
			continue
		}
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		filePath := lapsPath + "/" + file.Name()
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		var laps []models.Lap
		err = json.Unmarshal(content, &laps)
		if err != nil {
			return nil, err
		}
		for _, lap := range laps {
			res = append(res, lap)
		}
	}
	return res, nil
}
