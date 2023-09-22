package initutil

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/stellar-payment/sp-payment/internal/config"
)

func InitDirectory() {
	conf := config.Get()
	logDir := filepath.Join(conf.FilePath, "logs")

	_, err := os.Stat(logDir)
	if os.IsNotExist(err) {
		os.Mkdir(logDir, fs.ModeDir)
	}
}
