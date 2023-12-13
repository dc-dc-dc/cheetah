package util

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type ID ulid.ULID

func EnsureID() ID {
	return ID(ulid.MustNew(ulid.Timestamp(time.Now()), nil))
}

func (i ID) String() string {
	return ulid.ULID(i).String()
}
