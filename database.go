package d

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Orm interface, implement at least the following methods to facilitate internal calls in the devtool library
type InterfaceDatabase interface {
	Init()
}

const (
	ConfigPathDatabaseHost                = "database.host"
	ConfigPathDatabaseName                = "database.name"
	ConfigPathDatabaseUser                = "database.user"
	ConfigPathDatabasePassword            = "database.password"
	ConfigPathTimeoutReconnectionInterval = "database.timeout_reconnection_interval"
)

var (
	database                                   InterfaceDatabase // Global variable, stores the initialized interface, if not initialized, it is nil
	DefaultDatabaseTimeoutReconnectionInterval = 10
)

// ORM library unified access entry
type Database[T InterfaceDatabase] struct{}

// Initialization
func (d Database[T]) Init(conf T) {
	database = conf
}

// Get the initialized interface. If it is not initialized, Turnstile library is used by default.
func (d Database[T]) Get() T {
	if database == nil {
		LibraryGorm{}.Init()
	}
	return database.(T)
}

type LibraryGorm struct {
	*gorm.DB
	OpenDsn string
	Open    func(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error)
}

// Initialization
func (l LibraryGorm) Init() {
	if l.OpenDsn == "" {
		dbHost := Config[InterfaceConfig]{}.Get().GetString(ConfigPathDatabaseHost)
		dbName := Config[InterfaceConfig]{}.Get().GetString(ConfigPathDatabaseName)
		dbUser := Config[InterfaceConfig]{}.Get().GetString(ConfigPathDatabaseUser)
		dbPassword := Config[InterfaceConfig]{}.Get().GetString(ConfigPathDatabasePassword)
		l.OpenDsn = dbUser + ":" + dbPassword + "@tcp(" + dbHost + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	}
	if l.Open == nil {
		l.Open = func(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error) {
			return gorm.Open(dialector, opts...)
		}
	}
	db, err := l.Open(mysql.Open(l.OpenDsn), &gorm.Config{})
	if err != nil {
		// Auto Reconnected
		fmt.Printf("Error encountered while connecting to database: %v, automatically reconnecting after 10 seconds", err.Error())
		tri := Config[InterfaceConfig]{}.Get().GetIntWithDefault(ConfigPathTimeoutReconnectionInterval, DefaultDatabaseTimeoutReconnectionInterval)
		time.Sleep(time.Second * time.Duration(tri))
		l.Init()
		return
	}

	Database[LibraryGorm]{}.Init(LibraryGorm{
		DB: db,
	})
}

// Only support MySQL now
// Generate lazy query parameters based on parameters and value
// Example : GenerateFuzzyQueries(tx, map[string]string{"name": "John", "sex": "female"})
func (l LibraryGorm) GenerateFuzzyQueries(tx *gorm.DB, fields map[string]string) (*gorm.DB, error) {
	whereClause, args, err := MySQL{}.GenerateFuzzyQueries(fields)
	if err != nil {
		return nil, err
	}
	tx = tx.Where(whereClause, args...)
	return tx, nil
}

// Paginate
// https://gorm.io/docs/scopes.html#Pagination
func (l LibraryGorm) Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		q := r.URL.Query()

		page, _ := strconv.Atoi(q.Get(FieldNamePaginationPage))
		if page <= 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get(FieldNamePaginationPageSize))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
