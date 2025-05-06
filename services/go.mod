module github.com/st107853/fast_reading/services

go 1.24.2

replace github.com/st107853/fast_reading/models => ../models

replace github.com/st107853/fast_reading/utils => ../utils

replace github.com/st107853/fast_reading/config => ../config

require (
	github.com/st107853/fast_reading/models v0.0.0-00010101000000-000000000000
	github.com/st107853/fast_reading/utils v0.0.0-00010101000000-000000000000
	go.mongodb.org/mongo-driver v1.17.3
)

require (
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/st107853/fast_reading/config v0.0.0-00010101000000-000000000000 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
