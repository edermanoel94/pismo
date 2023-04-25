package data

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
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

	r.repo = NewOperationTypeRepository(r.DB)

	assert.IsType(r.T(), &Repository{}, r.repo)
}

func (r *RepositorySuite) AfterTest(_, _ string) {
	assert.NoError(r.T(), r.mock.ExpectationsWereMet())
}

func (r *RepositorySuite) TestList() {
	r.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "operation_types"`)).
		WillReturnRows(sqlmock.NewRows(nil))

	operationTypes, err := r.repo.List()

	assert.NoError(r.T(), err)

	assert.Empty(r.T(), operationTypes)
}
