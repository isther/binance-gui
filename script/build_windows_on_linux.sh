echo "Start build..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ HOST=x86_64-w64-mingw32 go build -ldflags "-s -w -H=windowsgui -extldflags=-static" -p 4 -v -o binance-gui.exe
mv binance-gui.exe release/binance-gui.exe
echo "End build..."
