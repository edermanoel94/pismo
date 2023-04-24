package data

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/edermanoel94/pismo/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type RepositorySuite struct {
	suite.Suite
	conn *sql.DB
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repo *Repository
	acc  *domain.Account
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (r *RepositorySuite) SetupSuite() {
	var err error

	r.conn, r.mock, err = sqlmock.New()

	assert.NoError(r.T(), err)

	dialector := postgres.New(postgres.Config{
		DriverName:           "postgres",
		DSN:                  "sqlmock_db_0",
		PreferSimpleProtocol: true,
		Conn:                 r.conn,
	})

	r.DB, err = gorm.Open(dialector, &gorm.Config{})

	assert.NoError(r.T(), err)

	r.repo = NewAccountRepository(r.DB)

	assert.IsType(r.T(), &Repository{}, r.repo)

	r.acc = &domain.Account{
		DocumentNumber: "12345678912",
		ID:             1,
	}
}

func (r *RepositorySuite) AfterTest(_, _ string) {
	assert.NoError(r.T(), r.mock.ExpectationsWereMet())
}

func (r *RepositorySuite) TestCreate() {
	r.mock.ExpectBegin()
	r.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "accounts" ("document_number","id") VALUES ($1,$2)`)).
		WithArgs(
			r.acc.DocumentNumber,
			r.acc.ID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	r.mock.ExpectCommit()

	acc, err := r.repo.Create(*r.acc)

	assert.NoError(r.T(), err)

	assert.Equal(r.T(), *r.acc, acc)
}

func (r *RepositorySuite) TestFindById() {

	rows := sqlmock.NewRows([]string{"id", "document_number"}).
		AddRow(
			r.acc.ID,
			r.acc.DocumentNumber,
		)

	r.mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE "accounts"."id" = $1`)).
		WithArgs().
		WillReturnRows(rows)

	acc, err := r.repo.FindById(int(r.acc.ID))

	assert.NoError(r.T(), err)

	assert.Equal(r.T(), *r.acc, acc)
}
