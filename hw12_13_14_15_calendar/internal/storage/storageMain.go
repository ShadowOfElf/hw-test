package storage

import (
	memorystorage "github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage/unityres"
)

func NewStorage(isDB bool) unityres.UnityStorageInterface {
	if !isDB {
		return memorystorage.New()
	}
	return sqlstorage.New()
}
