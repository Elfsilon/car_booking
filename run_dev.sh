set -o allexport
source config/dev.env
set +o allexport

go run cmd/main.go