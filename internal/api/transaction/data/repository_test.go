package data

import (
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/edermanoel94/pismo/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type RepositorySuite struct {
	suite.Suite
	conn *sql.DB
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repo        *Repository
	transaction *domain.Transaction
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

	r.repo = NewTransactionRepository(r.DB)

	assert.IsType(r.T(), &Repository{}, r.repo)

	r.transaction = &domain.Transaction{
		ID:              1,
		Amount:          45.00,
		AccountID:       1,
		OperationTypeID: 4,
		EventDate:       time.Now(),
	}
}

func (r *RepositorySuite) AfterTest(_, _ string) {
	assert.NoError(r.T(), r.mock.ExpectationsWereMet())
}

func (r *RepositorySuite) TestCreate() {
	r.mock.ExpectBegin()
	r.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "transactions" ("amount","event_date","account_id","operation_type_id","id")`)).
		WithArgs(
			r.transaction.Amount,
			AnyTime{},
			r.transaction.AccountID,
			r.transaction.OperationTypeID,
			r.transaction.ID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	r.mock.ExpectCommit()

	acc, err := r.repo.Create(*r.transaction)

	assert.NoError(r.T(), err)

	assert.Equal(r.T(), *r.transaction, acc)
}
