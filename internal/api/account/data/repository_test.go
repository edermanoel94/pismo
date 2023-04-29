package data

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/edermanoel94/pismo/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type RepositorySuite struct {
	suite.Suite
	conn *sql.DB
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repo    *Repository
	account *domain.Account
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (r *RepositorySuite) SetupSuite() {
	var err error

	r.conn, r.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

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

	r.account = &domain.Account{
		DocumentNumber: "12345678912",
		Balance:        1000.00,
		ID:             1,
	}
}

func (r *RepositorySuite) AfterTest(_, _ string) {
	assert.NoError(r.T(), r.mock.ExpectationsWereMet())
}

func (r *RepositorySuite) TestCreate() {
	r.mock.ExpectBegin()
	r.mock.ExpectQuery(`INSERT INTO "accounts" ("document_number","balance","id") VALUES ($1,$2,$3) RETURNING "id"`).
		WithArgs(
			r.account.DocumentNumber,
			r.account.Balance,
			r.account.ID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	r.mock.ExpectCommit()

	acc, err := r.repo.Create(*r.account)

	assert.NoError(r.T(), err)

	assert.Equal(r.T(), *r.account, acc)
}

func (r *RepositorySuite) TestFindById() {

	rows := sqlmock.NewRows([]string{"id", "balance", "document_number"}).
		AddRow(
			r.account.ID,
			r.account.Balance,
			r.account.DocumentNumber,
		)

	r.mock.ExpectQuery(`SELECT * FROM "accounts" WHERE "accounts"."id" = $1 ORDER BY "accounts"."id" LIMIT 1`).
		WithArgs().
		WillReturnRows(rows)

	acc, err := r.repo.FindById(int(r.account.ID))

	assert.NoError(r.T(), err)

	assert.Equal(r.T(), *r.account, acc)
}

func (r *RepositorySuite) TestUpdateBalance() {

	expectedBalance := 1000.00

	r.mock.ExpectBegin()
	r.mock.ExpectExec(`UPDATE "accounts" SET "balance"=$1 WHERE "id" = $2`).
		WithArgs(
			expectedBalance,
			r.account.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	r.mock.ExpectCommit()

	acc, err := r.repo.UpdateBalance(domain.Account{
		ID:      r.account.ID,
		Balance: r.account.Balance,
	})

	assert.NoError(r.T(), err)

	assert.Equal(r.T(), r.account.Balance, acc.Balance)
	assert.Equal(r.T(), r.account.ID, acc.ID)
}
