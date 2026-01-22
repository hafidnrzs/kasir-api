Kasir API dibangun dengan bahasa pemrograman Go

## How to Run

Untuk menjalankan aplikasi saat proses development, gunakan perintah berikut:

```bash
go run main.go
```

## Build Binary

Build Standar

```bash
go build -o kasir-api
```

Build Production (Lebih Kecil)

```bash
go build -ldflags="-s -w" -o kasir-api
```

- `-s`: Strip symbol table
- `-w`: Strip debug info

**Cross Compilation**

Build untuk Linux (dari WindowsMac)

```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o kasir-api
```
