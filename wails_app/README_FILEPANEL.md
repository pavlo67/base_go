# File panel replacement set

Copy these files over `pavlo67/wails_app`.

Start directory selection order:

1. `WAILS_APP_START_DIR` environment variable
2. first CLI argument
3. current working directory

Examples:

```bash
$ WAILS_APP_START_DIR=$HOME wails dev -tags webkit2_41
$ wails dev -tags webkit2_41 -- /tmp
```

After copying:

```bash
$ go mod tidy
$ npm --prefix frontend install
$ npm --prefix frontend run build
$ wails dev -tags webkit2_41
```

The current module is selected in `frontend/src/appConfig.js`.

-------------

Token usage: total=65,337 input=59,799 (+ 509,824 cached) output=5,538 (reasoning 1,177)
To continue this session, run codex resume 019e802f-e032-7301-a033-1c37c3bf8e94
