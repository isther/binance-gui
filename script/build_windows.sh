echo "Start build..."
go build -o binance-gui -ldflags "-s -w -H=windowsgui -extldflags=-static" .
mv binance-gui release/binance-gui
echo "End build..."
