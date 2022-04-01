## DB 

Start docker container for local development
```make docker_postgres```, database listen 5441 port with default credentials: ```POSTGRES_DB=boostersnews POSTGRES_USER=boosterdev```

Stop via cmd ```make docker_postgres_stop```

### Migration

Use automigration via [go-migrate](https://github.com/golang-migrate/migrate), store migrations files in .cicd/deploy/migrate and could be configuration in .env via variables ```MIGRATE_DIRECTORY=file://.cicd/deploy/migrate``` 

Run manual migration via cmd: 
```
make migrate_up
```
Create new migration files: 
```
migrate create -ext sql -dir .cicd/deploy/migrate -seq create_posts_table
```



## Environment variables

Variables load via config ./internal/app/config

Sample .env
```dotenv
ENV=dev
POSTGRES_URL=postgres://boosterdev:mysecretpassword@localhost:5442/boostersnews?sslmode=disable
MIGRATE_ENABLE=true
MIGRATE_DIRECTORY=file://.cicd/deploy/migrate
```
