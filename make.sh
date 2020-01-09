
# needed sudo apt install mingw-w64
BIN="bin"
if [ "$1" == "clean" ]
then
    rm $BIN/*
    exit
fi

#docker run -v $PWD:/flexo -v ~/go:/go liamg/golang-opengl bash -c 'cd /flexo && bash build.sh'
docker run -v $PWD:/flexo -v ~/go:/go golang-cross bash -c 'cd /flexo && bash build.sh'

#GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o $BIN/${1%.*}_x64.exe $1
#GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o $BIN/flexo-windows-x64.exe
#go build -o $BIN/${1%.*}_x64 $1
#go build -o $BIN/flexo-linux-x64
