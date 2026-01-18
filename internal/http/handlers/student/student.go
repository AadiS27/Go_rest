package student

//slog  is used for logging and is different from log as it is more efficient and faster
import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AadiS27/Go_rest/internal/storage"
	"github.com/AadiS27/Go_rest/internal/types"
	"github.com/AadiS27/Go_rest/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
			return
		}
		//validation of request

		if err := validator.New().Struct(student); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}
		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create student"})
			return
		}
		student.Id = int(lastId)
		slog.Info("student created", slog.Int("id", student.Id))

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid student ID"})
			return
		}
		student, err := storage.GetStudent(id)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get student"})
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}
