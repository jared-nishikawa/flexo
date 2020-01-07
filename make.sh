# needed sudo apt install mingw-w64
BIN="bin"
if [ -z $1 ]
then
    echo Need arg
    exit 0
fi
if [ $1 == "clean" ]
then
    rm $BIN/*
    exit
fi
#GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o $BIN/${1%.*}_x64.exe $1
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o $BIN/windows_x64.exe $1
#go build -o $BIN/${1%.*}_x64 $1
go build -o $BIN/linux_x64 $1
