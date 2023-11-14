// Integration tests.
// We could add "// +build integration" to the top of the file.
package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"goyave.dev/goyave/v5/config"
	"goyave.dev/goyave/v5/database"
	"goyave.dev/goyave/v5/slog"
	"goyave.dev/goyave/v5/util/testutil"
	"goyave.dev/template/database/model"
	"goyave.dev/template/database/repository"

	_ "goyave.dev/goyave/v5/database/dialect/postgres"
)

type logWriter struct{ t *testing.T }

func (lw logWriter) Write(p []byte) (n int, err error) {
	lw.t.Log(string(p))
	return len(p), nil
}

func TestBookSQL_Create(t *testing.T) {
	user1 := model.User{Model: gorm.Model{ID: 1}}

	cfg, err := config.LoadFrom(testutil.FindRootDirectory() + "config.test.json")
	require.NoError(t, err)
	db, err := database.New(cfg, func() *slog.Logger {
		return slog.New(slog.NewHandler(true, logWriter{t}))
	})
	require.NoError(t, err)

	SUT := repository.NewBookSQL(db)
	in := model.Book{Owner: user1, Title: "test"}
	err = SUT.Create(&in)
	require.NoError(t, err)

	assert.NotZero(t, in.ID)
	assert.NotZero(t, in.OwnerID)

	// Would return gorm.ErrRecordNotFound if the book was not created.
	err = db.First(&model.Book{}, "id = ?", in.ID).Error
	require.NoError(t, err)
}
