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

	"crypto/rand"
	"math/big"

	"github.com/JaineelVora08/librproto/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var dbpool *pgxpool.Pool

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbpool, err = pgxpool.New(context.Background(), os.Getenv("connection_string"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	_, err = dbpool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS messages (id UUID PRIMARY KEY, content TEXT NOT NULL, timestamp BIGINT NOT NULL, status VARCHAR(10) NOT NULL)")
}

func addmessage(message models.Message) models.APIResponse {
	message.Message_id = uuid.New().String()
	message.Sent_timestamp = time.Now().Unix()
	message.Status = "pending"

	response := models.ModeratorResponse{}

	datachan := make(chan models.ModeratorResponse)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			response := ModeratorResponse(message)
			datachan <- response
		}()
	}
	wg.Wait()
	close(datachan)

	count := 0
	for response := range datachan {
		if response.Status == "approved" {
			count++
		}
	}

	if count >= 2 {
		message.Status = "approved"
	} else {
		message.Status = "rejected"
	}

	insertData(message)
	final_message := models.APIResponse{}
	final_message = APIresponsefunc(message)
	return final_message
}

func ModeratorResponse(message models.Message) models.ModeratorResponse {
	response := models.ModeratorResponse{}

	response.Mod_id = uuid.New().String()

	randomstat, _ := rand.Int(rand.Reader, big.NewInt(2))
	randomstatus := int(randomstat.Int64())

	if randomstatus == 0 {
		response.Status = "approved"
	} else {
		response.Status = "rejected"
	}

	rawTime, _ := rand.Int(rand.Reader, big.NewInt(3))
	randomtime := rawTime.Int64() + 1
	response.Response_time = int(randomtime)

	time.Sleep(time.Duration(response.Response_time) * time.Second)

	response.Message_id = message.Message_id

	return response
}

func insertData(message models.Message) {
	_, err := dbpool.Exec(context.Background(), "INSERT INTO messages (id, content, timestamp, status) VALUES ($1, $2, $3, $4)", message.Message_id, message.Context, message.Sent_timestamp, message.Status)
	if err != nil {
		log.Fatalf("Unable to insert data: %v\n", err)
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
	rows, _ := dbpool.Query(context.Background(), "SELECT id, content, timestamp, status FROM messages WHERE status=$1 AND timestamp=$2",
		"accepted", timestamp)

	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		_ = rows.Scan(&msg.Message_id, &msg.Context, &msg.Sent_timestamp, &msg.Status)
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
	json.NewEncoder(w).Encode(final_messages)
}

func getallacceptedmessages() []models.Message {
	rows, _ := dbpool.Query(context.Background(), "SELECT * from messages WHERE status=$1", "accepted")
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.Message_id, &msg.Context, &msg.Sent_timestamp, &msg.Status)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages
}

func GetAllAcceptedMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")
	final_messages := getallacceptedmessages()
	json.NewEncoder(w).Encode(final_messages)
}

func APIResponsefunc(message models.Message) models.APIResponse {
	final_response := models.APIResponse{}

	final_response.Message_id = message.Message_id
	final_response.Submission_timestamp = message.Sent_timestamp
	final_response.Status = message.Status

	return final_response
}
