module github.com/st107853/fast_reading/services

go 1.24.2

replace github.com/st107853/fast_reading/models => ../models

replace github.com/st107853/fast_reading/utils => ../utils

replace github.com/st107853/fast_reading/config => ../config

require (
	github.com/st107853/fast_reading/models v0.0.0-00010101000000-000000000000
	github.com/st107853/fast_reading/utils v0.0.0-00010101000000-000000000000
	gorm.io/gorm v1.31.0
)

require (
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	gorm.io/driver/postgres v1.6.0 // indirect
)
