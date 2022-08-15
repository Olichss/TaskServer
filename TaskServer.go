package main

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type Task struct {
	Id   int       `gorm:"primary_key"`
	Text string    `json:"text"`
	Due  time.Time `json:"due"`
}

type TaskServer struct {
	DB *gorm.DB
	Id int
}

func NewTaskServer() (*TaskServer, error) {
	db, err := gorm.Open("mysql", "root:root@/tasklist?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}

	db.Debug().DropTableIfExists(&Task{})
	db.Debug().AutoMigrate(&Task{})

	return &TaskServer{DB: db}, nil
}

func (server *TaskServer) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var task Task
	if err := dec.Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	server.DB.Create(&task)
	result := server.DB.Last(&task)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Value)
}

func (server *TaskServer) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get task at %s\n", r.URL.Path)

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	task, err := server.GetTaskByID(id)
	if err {
		log.WithFields(log.Fields{"Id": id}).Info("Getting TaskItem")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	} else {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"get": false, "error": "Record Not Found"}`)
	}

}

func (server *TaskServer) GetAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get tasks at %s\n", r.URL.Path)

	var tasks []Task
	server.DB.Find(&tasks)
	w.Header().Set("Content-Type", "application/json")
	for _, task := range tasks {
		json.NewEncoder(w).Encode(task)
	}
}

func (server *TaskServer) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get task at %s\n", r.URL.Path)

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	server.DB.Delete(&Task{}, id)
}

func (server *TaskServer) DeleteAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	server.DB.Delete(&Task{})
}

func (server *TaskServer) GetTaskByID(id int) (Task, bool) {
	var task Task
	result := server.DB.First(&task, id)
	if result.Error != nil {
		log.Warn("Task not found in database")
		return task, false
	}
	return task, true
}
