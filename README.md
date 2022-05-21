# go cos disk

## run ftp

```
go run ftp/main/main.go
```

## run server

```
go run main.go
```

## start tikv and tidb

```
tiup clean --all
tiup playground
```

## fuse

```
mkdir ./cmd/test
go run cmd/fuse/main.go ./cmd/test
```