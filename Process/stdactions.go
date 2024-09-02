package stdactions

import (
	"encoding/json"
	"net/http"
	structs "usr/local/go/bin/Structs"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func AddStudent(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		response := map[string]interface{}{
			"status": "false",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}
	}
	db, err := gorm.Open(sqlite.Open("Data.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Convert map values to the expected types and create the Student
	student := structs.Student{
		Name:           data["name"].(string),
		Code:           int(data["code"].(float64)),
		AttendanceRate: int(data["ar"].(float64)),
		Rank:           int(data["rank"].(float64)),
		Messages:       []structs.Messages{},
		Degrees:        []structs.Degrees{},
	}

	var existingStudent structs.Student
	if err := db.Where("id = ?", int(data["id"].(float64))).First(&existingStudent).Error; err == nil {
		db.Model(&existingStudent).Updates(map[string]interface{}{"attendance_rate": int(data["ar"].(float64)), "rank": int(data["rank"].(float64))})
		return
	}
	db.Create(&student)

	response := map[string]interface{}{
		"status": "success",
	}

	// Set the content type to application/json

	// Encode the map to JSON and write it to the response writer
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}
func SearchStudent(w http.ResponseWriter, r *http.Request) {
	// Ensure the request body is closed after reading
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	// Decode the JSON body into a map
	data := map[string]interface{}{
		"id": r.URL.Query().Get("id"),
	}

	// Open the database
	db, err := gorm.Open(sqlite.Open("Data.db"), &gorm.Config{})
	if err != nil {
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}

	// Find the student by ID and preload related Messages and Degrees
	var student structs.Student
	if err := db.Preload("Messages").Preload("Degrees").First(&student, data["id"]).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Process Messages
	messages := []interface{}{}
	for _, message := range student.Messages {
		messages = append(messages, map[string]interface{}{
			"id":      message.ID,
			"content": message.Content,
		})
	}

	// Process Degrees
	degrees := []interface{}{}
	// Prepare the response
	response := map[string]interface{}{
		"std-id":   student.ID,
		"std-ar":   student.AttendanceRate,
		"std-rank": student.Rank,
		"std-name": student.Name,
		"messages": messages,
		"degrees":  degrees,
	}

	// Encode the response as JSON and send it
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}
}
func AddMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	// Decode the JSON body into a map
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error decoding JSON payload", http.StatusBadRequest)
		return
	}
	// Type assertion with error handling
	studentID, ok1 := data["id"].(float64) // JSON numbers are decoded as float64
	content, ok2 := data["content"].(string)
	if !ok1 || !ok2 {
		http.Error(w, "Invalid data types", http.StatusBadRequest)
		return
	}

	// Create new message
	newMessage := structs.Messages{
		StudentID: int(studentID), // Convert float64 to int
		Content:   content,
		Type:      data["type"].(string),
	}

	// Open database connection
	db, err := gorm.Open(sqlite.Open("Data.db"), &gorm.Config{})
	if err != nil {
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}

	// Create new message record
	db.Create(&newMessage)

	// Respond with success
	response := map[string]string{"status": "success"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
func Post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newpost := structs.Posts{
		Genre:        data["genre"].(string),
		Title:        data["title"].(string),
		Description:  data["description"].(string),
		EmbededLinks: data["link"].(string),
	}

	db, err := gorm.Open(sqlite.Open("Data.db"), &gorm.Config{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	db.Create(&newpost)

	response := map[string]string{"status": "success"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func Posts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	db, err := gorm.Open(sqlite.Open("Data.db"), &gorm.Config{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var posts []structs.Posts
	db.Find(&posts)

	var postJs []interface{}
	for _, post := range posts {
		postJs = append(postJs, map[string]interface{}{
			"title": post.Title,
			"desc":  post.Description,
			"link":  post.EmbededLinks,
			"Genre": post.Genre,
		})
	}

	response := postJs
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	// Decode the JSON body into a map
	data := map[string]interface{}{
		"id": r.URL.Query().Get("id"),
	}
	db, err := gorm.Open(sqlite.Open("Data.db"), &gorm.Config{})
	if err != nil {
		return
	}
	var student structs.Student
	if err := db.Preload("Messages").Preload("Degrees").First(&student, data["id"]).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Process Messages
	messages := []interface{}{}
	for _, message := range student.Messages {
		messages = append(messages, map[string]interface{}{
			"id":      message.ID,
			"content": message.Content,
			"type":    message.Type,
		})
	}
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}

}
