echo "Start build..."
go build -o binance-gui.exe -ldflags "-s -w -H=windowsgui -extldflags=-static" .
# go build -o binance-gui.exe
mv binance-gui.exe release/binance-gui.exe
echo "End build..."
