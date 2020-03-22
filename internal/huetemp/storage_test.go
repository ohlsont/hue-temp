package huetemp

import (
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMigrate(t *testing.T) {
	db, err := bolt.Open("my.db", 0600, nil)
	require.NoError(t, err)
	defer db.Close()
	s := &Service{db:db}
	err = s.migrateData()
	require.NoError(t, err)
}
