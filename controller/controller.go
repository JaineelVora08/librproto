package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/JaineelVora08/librproto/models"
	"github.com/JaineelVora08/librproto/moderator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Dbpool *pgxpool.Pool

func init() {
	fmt.Println("Controller init started")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Dbpool, err = pgxpool.New(context.Background(), os.Getenv("connection_string"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	_, err = Dbpool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS messages (id UUID PRIMARY KEY, content TEXT NOT NULL, timestamp BIGINT NOT NULL, status VARCHAR(10) NOT NULL)")
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}
}

func addmessage(message models.Message) models.APIResponse {
	message.Message_id = uuid.New().String()
	message.Sent_timestamp = time.Now().Unix()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	datachan := make(chan models.ModeratorResponse, 3)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp := moderator.ModeratorResponse(message)
			select {
			case datachan <- resp:
			case <-ctx.Done():
			}
		}()
	}
	wg.Wait()
	close(datachan)

	count := 0
	for response := range datachan {
		if response.Status == "accepted" {
			count++
		}
	}

	if count >= 2 {
		message.Status = "accepted"
	} else {
		message.Status = "rejected"
	}

	insertData(message)
	// final_message := models.APIResponse{}
	final_message := APIResponsefunc(message)
	return final_message
}

func insertData(message models.Message) {
	fmt.Printf("Inserting message: %+v\n", message)

	_, err := Dbpool.Exec(context.Background(),
		"INSERT INTO messages (id, content, timestamp, status) VALUES ($1, $2, $3, $4)",
		message.Message_id, message.Content, message.Sent_timestamp, message.Status)

	if err != nil {
		log.Printf("Insert failed: %v\n", err)
	}
}

func AddNewMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	message := models.Message{}
	_ = json.NewDecoder(r.Body).Decode(&message)
	final_message := addmessage(message)
	json.NewEncoder(w).Encode(final_message)
}

func gettimemessages(timestamp int64) []models.Message {
	rows, _ := Dbpool.Query(context.Background(), "SELECT * FROM messages WHERE status=$1 AND timestamp=$2;",
		"accepted", timestamp)

	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		_ = rows.Scan(&msg.Message_id, &msg.Content, &msg.Sent_timestamp, &msg.Status)
		messages = append(messages, msg)
	}

	return messages
}

func GetTimeMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")
	params := mux.Vars(r)
	timestampStr := params["ts"]
	timestamp, _ := strconv.ParseInt(timestampStr, 10, 64)
	final_messages := gettimemessages(timestamp)

	if len(final_messages) == 0 {
		json.NewEncoder(w).Encode(map[string]string{"message": "no message found"})
	} else {
		json.NewEncoder(w).Encode(final_messages)
	}
}

func getallacceptedmessages() []models.Message {
	query := "SELECT id, content, timestamp, status FROM messages WHERE status=$1"

	rows, err := Dbpool.Query(context.Background(), query, "accepted")
	if err != nil {
		log.Printf("Query error: %v\n", err)
		return nil
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.Message_id, &msg.Content, &msg.Sent_timestamp, &msg.Status)
		if err != nil {
			log.Printf("Row scan error: %v\n", err)
			continue
		}
		log.Printf("Fetched row: %+v\n", msg)
		messages = append(messages, msg)
	}
	return messages
}

func GetAllAcceptedMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")
	final_messages := getallacceptedmessages()

	if len(final_messages) == 0 {
		json.NewEncoder(w).Encode(map[string]string{"message": "no message found"})
	} else {
		json.NewEncoder(w).Encode(final_messages)
	}
}

func APIResponsefunc(message models.Message) models.APIResponse {
	final_response := models.APIResponse{}

	final_response.Message_id = message.Message_id
	final_response.Sent_timestamp = message.Sent_timestamp
	final_response.Status = message.Status

	return final_response
}
