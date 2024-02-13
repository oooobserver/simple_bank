package util

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Convert google uuid to pgtype.UUID
func ConvertGU(gu uuid.UUID) pgtype.UUID {
	var tmp pgtype.UUID
	tmp.Valid = true
	tmp.Bytes = gu
	return tmp
}

func ConverTime(t time.Time) pgtype.Timestamptz {
	var res pgtype.Timestamptz
	res.Valid = true
	res.Time = t
	return res
}

func ConvertString(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}
