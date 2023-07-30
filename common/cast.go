package common

import (
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

func StrToUUID(uuidStr string) (pgtype.UUID, error) {
	if len(uuidStr) != 36 {
		return pgtype.UUID{}, errors.New("invalid UUID")
	}
	index := [16]int{0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34}
	var id [16]byte
	for i := 0; i < 16; i++ {
		beg := index[i]
		b, err := strconv.ParseInt(uuidStr[beg:beg+2], 16, 64)
		if err != nil {
			return pgtype.UUID{}, err
		}
		id[i] = byte(b)
	}
	return pgtype.UUID{Bytes: id, Valid: true}, nil
}
