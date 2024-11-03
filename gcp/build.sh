go build -ldflags="-s -w" -o uai cmd/main.go
pnpm exec tailwindcss -i input.css -o style.css --minify
upx --best ./uai
