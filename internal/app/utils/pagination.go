package utils

import (
	"encoding/base64"
	"encoding/json"
)

type Cursors struct {
	NextCursor string `json:"nextCursor"`
	PrevCursor string `json:"prevCursor"`
}

func NewCursors(next *Cursor, prev *Cursor) *Cursors {
	return &Cursors{
		NextCursor: encodeCursor(next),
		PrevCursor: encodeCursor(prev),
	}
}

type Cursor struct {
	Id          int
	IsPointNext bool // for getting the previous part of data in the correct order
}

func NewCursor(id int, isPointNext bool) *Cursor {
	return &Cursor{
		Id:          id,
		IsPointNext: isPointNext,
	}
}

func encodeCursor(cursor *Cursor) string {
	if cursor == nil || cursor.Id == 0 {
		return ""
	}
	serializedCursor, err := json.Marshal(cursor)
	if err != nil {
		return ""
	}
	encodedCursor := base64.StdEncoding.EncodeToString(serializedCursor)
	return encodedCursor
}

func DecodeCursor(cursor string) (*Cursor, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var cur Cursor
	if err := json.Unmarshal(decodedCursor, &cur); err != nil {
		return nil, err
	}
	return &cur, nil
}

func GetPaginationOperator(isPointNext bool, sortOrder string) (string, string) {
	if isPointNext && sortOrder == "asc" {
		return ">", ""
	} else if isPointNext && sortOrder == "desc" {
		return "<", ""
	} else if !isPointNext && sortOrder == "asc" {
		return "<", "DESC"
	} else if !isPointNext && sortOrder == "desc" {
		return ">", "ASC"
	}

	return "", ""
}
