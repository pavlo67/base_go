# backend/

Поки Go-бек лежить у кореневому `main.go`, бо Wails v2 очікує entrypoint проєкту в корені.

Коли захочеться рознести код культурніше, сюди можна винести HTTP handlers/services/packages, а кореневий `main.go` залишити тільки тонким Wails-entrypoint.
