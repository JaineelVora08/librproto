package moderator

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/JaineelVora08/librproto/models"
	"github.com/google/uuid"
)

func ModeratorResponse(message models.Message) models.ModeratorResponse {
	response := models.ModeratorResponse{}

	response.Mod_id = uuid.New().String()

	randomstat, _ := rand.Int(rand.Reader, big.NewInt(2))
	randomstatus := int(randomstat.Int64())

	if randomstatus == 0 {
		response.Status = "accepted"
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
