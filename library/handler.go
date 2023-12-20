package library

import (
	"encoding/json"
	"golibrary/utils"
	"net/http"
)

func (lf *LibraryFacade) GetAuthorsHandler(w http.ResponseWriter, r *http.Request) {
	authors, err := utils.GetAuthors(lf.DB)
	if err != nil {
		http.Error(w, "Ошибка при получении информации об авторах", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Явное указание успешного кода состояния

	err = json.NewEncoder(w).Encode(authors)
	if err != nil {
		// Логирование ошибки при кодировании JSON
		http.Error(w, "Ошибка при кодировании ответа", http.StatusInternalServerError)
	}
}
