GOOS=windows GOARCH=386 go build -v -o ./../dst/tsfp.exe
GOOS=linux GOARCH=386 go build -v -o ./../dst/tsfp.linux
GOARCH=386 go build -v -o ./../dst/tsfp.mac
