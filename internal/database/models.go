// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type StatusEnum string

const (
	StatusEnumTodo    StatusEnum = "Todo"
	StatusEnumProcess StatusEnum = "Process"
	StatusEnumDone    StatusEnum = "Done"
)

func (e *StatusEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = StatusEnum(s)
	case string:
		*e = StatusEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for StatusEnum: %T", src)
	}
	return nil
}

type NullStatusEnum struct {
	StatusEnum StatusEnum
	Valid      bool // Valid is true if StatusEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullStatusEnum) Scan(value interface{}) error {
	if value == nil {
		ns.StatusEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.StatusEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullStatusEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.StatusEnum), nil
}

type Diary struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DatetimeExc time.Time
	UserID      uuid.UUID
	Title       string
	Description sql.NullString
	Images      []string
	Urls        []string
	Icon        sql.NullString
	Background  sql.NullString
}

type FollowHabit struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	UserID         uuid.UUID
	HabitID        uuid.UUID
	PromiseFrom    time.Time
	PromiseEnd     time.Time
	Processes      []time.Time
	IncludedReport bool
}

type FollowProject struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	ProjectID uuid.UUID
}

type Habit struct {
	ID                uuid.UUID
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Title             string
	Creator           uuid.UUID
	Purpose           sql.NullString
	Description       sql.NullString
	Icon              sql.NullString
	Background        sql.NullString
	Images            []string
	Urls              []string
	TimeInDay         time.Time
	LoopWeek          []int32
	LoopMonth         []int32
	RecommendDuration int32
}

type Project struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Creator     uuid.UUID
	Purpose     sql.NullString
	Description sql.NullString
	Background  sql.NullString
}

type Task struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DatetimeExc time.Time
	Title       string
	Purpose     sql.NullString
	Description sql.NullString
	Images      []string
	Urls        []string
	Status      StatusEnum
	UserID      uuid.UUID
	ProjectID   uuid.NullUUID
	ParentID    uuid.NullUUID
}

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Username       string
	Fullname       string
	HashedPassword string
	Bio            sql.NullString
	Avatar         sql.NullString
	Email          string
	ApiKey         string
}
