# trainingFinder
.env можно задать свой:

DB_HOST=localhost
HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin
DB_NAME=godb
DB_SSLMODE=disable
GRPC_PORT=:9090
HTTP_PORT=:8080
SECRET=pofig_kakoi
ACCESS_TOKEN_DURATION=1h
AUTH_CONFIG_PATH=./config.yml

бд поднимается по start-db

переменные окружения загружаются по make export

затем make bin-deps, make migration-up, make start

эндпоинты можно посмотреть в сваггере http://localhost:8080/docs/# 
авторизация там тоже рабочая