# Robot Apocalypse

# Configuratino

```apocalypse.yaml```

## Tests
```
go test ./...
```

## Make the API documentation

```
make
```

## Documentation available at
```
http://localhost:8080/docs
```

## Sample requests

```
curl -X POST localhost:8080/survivors -d @sample1.json
curl -X POST localhost:8080/survivors -d @sample2.json
curl -X GET localhost:8080/survivors
curl -X PUT localhost:8080/survivors/infected -d '{"id": "HD138VOP34219" }'
curl -X GET localhost:8080/survivors/infected?status=true
curl -X GET localhost:8080/survivors/infected?status=false
curl -X PUT localhost:8080/survivors/location -d '{"id": "HD138VOP34219", "Latitude": 1, "Longitude": 2 }'
curl -X GET localhost:8080/survivors/stats
```

## Visit `http://localhost:8080/reportweb` to view the records of survivors from the web


## for a reactjs web app that uses the api
```
cd web/apocalypseweb
npm install
npm start
```
Visit `http://localhost:3000`