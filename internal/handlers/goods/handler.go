package goods

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	gService "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/goods"
)

type Handler struct {
	service gService.Service
}

func NewGoodsHandler(service gService.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		log.WithError(err).Error("invalid limit parameter")
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		log.WithError(err).Error("invalid offset parameter")
		http.Error(w, "Invalid offset", http.StatusBadRequest)
		return
	}

	// Выполнение запроса к базе данных
	resp, err := h.service.GetList(ctx, limit, offset)
	if err != nil {
		log.WithError(err).Error("get records list")
		return
	}

	// Устанавливаем заголовок Content-Type в application/json
	w.Header().Set("Content-Type", "application/json")

	// Кодируем resp в JSON и отправляем его в тело ответа
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithError(err).Error("encode response data")
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ с данными в формате JSON
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр projectId из URL
	projectID, err := strconv.Atoi(r.URL.Query().Get("projectId"))
	if err != nil {
		log.WithError(err).Error("invalid projectId parameter")
		http.Error(w, "Invalid projectId parameter", http.StatusBadRequest)
		return
	}

	// Декодируем JSON из тела запроса
	var postData gService.PostData
	if err := json.NewDecoder(r.Body).Decode(&postData); err != nil {
		log.WithError(err).Error("decode request body")
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	data := &gService.Good{Name: postData.Name, ProjectId: projectID}

	// Создаем новый объект или выполняем другую логику по созданию данных
	// Например, вызываем метод сервиса для создания данных в базе данных
	resp, err := h.service.Create(data)
	if err != nil {
		log.WithError(err).Error("create record")
		http.Error(w, "Failed to create record", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type в application/json
	w.Header().Set("Content-Type", "application/json")

	// Кодируем resp в JSON и отправляем его в тело ответа
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Обработка ошибки при кодировании resp в JSON
		log.WithError(err).Error("encode response data")
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ с данными в формате JSON
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из URL
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.WithError(err).Error("invalid id parameter")
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	projectID, err := strconv.Atoi(r.URL.Query().Get("projectId"))
	if err != nil {
		log.WithError(err).Error("invalid projectId parameter")
		http.Error(w, "Invalid projectId parameter", http.StatusBadRequest)
		return
	}

	// Декодируем JSON из тела запроса
	var patchData gService.PatchData
	if err := json.NewDecoder(r.Body).Decode(&patchData); err != nil {
		log.WithError(err).Error("decode request body")
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Проверяем наличие обязательного поля "name"
	if patchData.Name == "" {
		log.WithError(err).Error("name field is required")
		http.Error(w, "Name field is required", http.StatusBadRequest)
		return
	}

	data := &gService.Good{Name: patchData.Name, Description: patchData.Description, Id: id, ProjectId: projectID}

	// Выполняем обновление записи в базе данных
	resp, err := h.service.Update(data)
	if err != nil {
		log.WithError(err).Error("update record")
		return
	}

	// Устанавливаем заголовок Content-Type в application/json
	w.Header().Set("Content-Type", "application/json")

	// Кодируем resp в JSON и отправляем его в тело ответа
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Обработка ошибки при кодировании resp в JSON
		log.WithError(err).Error("encode response data")
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ с данными в формате JSON
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из URL
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.WithError(err).Error("invalid id parameter")
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	projectID, err := strconv.Atoi(r.URL.Query().Get("projectId"))
	if err != nil {
		log.WithError(err).Error("invalid projectId parameter")
		http.Error(w, "Invalid projectId parameter", http.StatusBadRequest)
		return
	}

	// Вызываем метод Delete сервиса, передавая id и projectID
	resp, err := h.service.Remove(id, projectID)
	if err != nil {
		// Обработка ошибки, если удаление не удалось
		log.WithError(err).Error("delete record")
		http.Error(w, "Failed to delete record", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type в application/json
	w.Header().Set("Content-Type", "application/json")

	// Кодируем resp в JSON и отправляем его в тело ответа
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Обработка ошибки при кодировании resp в JSON
		log.WithError(err).Error("encode response data")
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ с данными в формате JSON
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Reprioritize(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из URL
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.WithError(err).Error("invalid id parameter")
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	projectID, err := strconv.Atoi(r.URL.Query().Get("projectId"))
	if err != nil {
		log.WithError(err).Error("invalid projectId parameter")
		http.Error(w, "Invalid projectId parameter", http.StatusBadRequest)
		return
	}

	// Декодируем JSON из тела запроса
	var patchData struct {
		NewPriority int `json:"newPriority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&patchData); err != nil {
		log.WithError(err).Error("decode request body")
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Проверяем корректность нового приоритета
	if patchData.NewPriority <= 1 {
		log.WithError(err).Error("invalid record new priority")
		http.Error(w, "The new priority must be greater than or equal to one.", http.StatusBadRequest)
		return
	}

	// Вызываем метод UpdatePriority сервиса, передавая данные для обновления приоритета
	resp, err := h.service.Reprioritize(id, projectID, patchData.NewPriority)
	if err != nil {
		// Обработка ошибки, если обновление приоритета не удалось
		log.WithError(err).Error("update priority")
		http.Error(w, "Failed to update priority", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type в application/json
	w.Header().Set("Content-Type", "application/json")

	// Кодируем resp в JSON и отправляем его в тело ответа
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Обработка ошибки при кодировании resp в JSON
		log.WithError(err).Error("encode response data")
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ с данными в формате JSON
	w.WriteHeader(http.StatusOK)
}
