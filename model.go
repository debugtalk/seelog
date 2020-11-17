package wstailog

type slog struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type logLine struct {
	LogName string `json:"logName"`
	Data    string `json:"data"`
}
