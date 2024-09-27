### Challenge
- Pre-existing Golang knowledge: zero
- Golang theory: 2 days
- Challenge: 10 days
- Help: 
  - chatgpt-4o
  - claude-sonnet-3.5

### Usage
A simple .env file has been provided.

```
docker compose up --build
```

Check all is good through Grafana at `http://localhost:3000` with user-password both `admin`.


If you restart the container, producer will fail because it will try to create tasks in the database that already exist from the previous run. So, use `-v` flag to remove the volume of the db before relaunching:

```
docker compose down -v
```

### Tests

```
go test --cover -v ./...
```

Most tests mock the db, but some need a real db instance. Those will fail. 


Simply run:

```
docker compose -f docker-compose-test.yaml up -d --build
```

Now, you have a database instance. You just need to export the DB_URL to your terminal session:
```
export DB_URL=postgres://user:password@localhost:5432/golang_db_test
```

Now all tests should pass with `go test --cover -v ./...`

### Design overview

![Design image](design.png)

### ToDo
1. Better tests, more tests.
2. Share a bit more code between the services
3. Rethink extendability into design

