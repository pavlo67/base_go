# Wails + Vue + Go hello-world

Мінімальна болванка:

- `frontend/` — Vue + Vite + Bootstrap
- `main.go` — Wails entrypoint + простий Go HTTP backend на стандартних `net/http`, `encoding/json`
- `backend/` — місце для майбутнього рознесення Go-коду
- `build/bin/` — сюди Wails кладе `.exe` / `.app` / linux binary

## Передумови

	$ sudo apt install build-essential pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev nsis
	$ go install github.com/wailsapp/wails/v2/cmd/wails@latest
	$ wails doctor

## Запуск Wails dev

З кореня проєкту:

	$ npm --prefix frontend install
	$ cd frontend
	$ npm install -D @vitejs/plugin-vue
	$ cd ..
	$ go mod tidy
	$ wails dev -tags webkit2_41

Go backend слухає:

**$ http://127.0.0.1:34116/api/hello**

## Окремий Vite-запуск у браузері

Термінал 1:

**$ go run .**

Термінал 2:

**$ npm --prefix frontend run dev**

Потім відкрити:

**$ http://127.0.0.1:5173**

## Build

**$ wails build**

Результат буде в:

**$ build/bin/**
