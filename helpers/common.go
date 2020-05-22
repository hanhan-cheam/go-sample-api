package helpers

import (
	"strconv"
)

func ToInt(value string) int {
	integer, _ := strconv.Atoi(value)
	return integer
}
func ToString(value uint) string {
	return strconv.FormatUint(uint64(value), 10)
}

func PadLeft(id string, pad string, length int) string {
	for {
		id = pad + id
		if len(id) >= length {
			return id[0:length]
		}
	}
}

// func AcceptToken(target_id uint, station_id uint) string {
// 	token := jwt.New(jwt.SigningMethodHS256)
// 	stringTargetID := strconv.FormatUint(uint64(target_id), 10)
// 	stringStationID := strconv.FormatUint(uint64(station_id), 10)
// 	fmt.Println("stringTargetID", stringTargetID)
// 	fmt.Println("stringStationID", stringStationID)
// 	token.Claims = &models.Tokens{
// 		TargetId:  stringTargetID,
// 		StationId: stringStationID,
// 	}
// 	tokenString, _ := token.SignedString([]byte("secret"))
// 	// type tempToken struct{ Token string }
// 	return tokenString
// }

func IsEmpty(data string) bool {
	if len(data) == 0 {
		return true
	} else {
		return false
	}
}
