for GOOS in darwin linux windows
do
    for GOARCH in 386 amd64
    do
        echo "Building $GOOS-$GOARCH"

        export GOOS=$GOOS
        export GOARCH=$GOARCH
        export CGO_ENABLED=1
        
        if [ "$GOOS" == "windows" ]
        then
            go build -o bin/flexo-$GOOS-$GOARCH.exe
        else
            go build -o bin/flexo-$GOOS-$GOARCH
        fi
    done
done
