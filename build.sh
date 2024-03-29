echo Building windows/amd64
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o bin/flexo-windows-amd64.exe
echo Building windows/386
GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc go build -o bin/flexo-windows-386.exe

echo Building darwin/amd64
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-apple-darwin15-clang go build -o bin/flexo-darwin-amd64
echo Buliding darwin/386
GOOS=darwin GOARCH=386 CGO_ENABLED=1 CC=i386-apple-darwin15-clang go build -o bin/flexo-darwin-386

echo Building linux/amd64
go build -o bin/flexo-linux-amd64
echo Building linux/386
GOOS=linux GOARCH=386 CGO_ENABLED=1 go build -o bin/flexo-linux-386
