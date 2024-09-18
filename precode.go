package main

import (
    "encoding/json"
	"fmt"
	"net/http"
	"bytes"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

//Обработчик для получения всех задач
func getTasks(w http.ResponseWriter, r *http.Request) {
    // сериализуем данные из слайса tasks
    resp, err := json.Marshal(tasks)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(resp)
}

//Обработчик для отправки задачи на сервер
func postTasks(w http.ResponseWriter, r *http.Request) {
    var task Task
    var buf bytes.Buffer

    _, err := buf.ReadFrom(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    tasks[task.ID] = task

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
}

//Обработчик для получения задачи по ID
func getTask(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    task, ok := tasks[id]
    if !ok {
        http.Error(w, "Артист не найден", http.StatusNoContent)
        return
    }

    resp, err := json.Marshal(task)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
}

func dropTask(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    task, ok := tasks[id]
    if !ok {
        http.Error(w, "Артист не найден", http.StatusNoContent)
        return
    }

    resp, err := json.Marshal(task)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	delete(tasks, id)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(resp)
}

func main() {
	r := chi.NewRouter()

	//Маршрутизатор для получения всех задач
	r.Get("/tasks", getTasks)
	//Маршрутизатор для добавления задачи
	r.Post("/tasks", postTasks)
	//Маршрутизатор для получения задачи по id
	r.Get("/task/{id}", getTask)
	//Маршрутизатор для удаления задачи по id
	r.Delete("/delete-task/{id}", dropTask)


	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
