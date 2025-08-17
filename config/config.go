package config

import (
	"media-server/utils"
	"strconv"
	"strings"
)

type Config struct {
	UploadDir    string
	Port         string
	PageLimit    int
	FileMaxSize  int64
	FileFilter   string
	AllowedMIMEs map[string]bool
}

func Load() Config {
	pageLimit, _ := strconv.Atoi(utils.GetEnv("PAGE_LIMIT", "20"))
	fileMaxSize, _ := strconv.ParseInt(utils.GetEnv("FILE_MAX_SIZE", "104857600"), 10, 64)

	// читаем список MIME-типа из окружения
	mimesEnv := utils.GetEnv("ALLOWED_MIMES", "image/jpeg,image/png,video/mp4")
	mimesList := strings.Split(mimesEnv, ",")
	mimes := make(map[string]bool)
	for _, m := range mimesList {
		m = strings.TrimSpace(m)
		if m != "" {
			mimes[m] = true
		}
	}

	return Config{
		UploadDir:    utils.GetEnv("UPLOAD_DIR", "./public"),
		Port:         utils.GetEnv("PORT", "8080"),
		PageLimit:    pageLimit,
		FileMaxSize:  fileMaxSize,
		FileFilter:   utils.GetEnv("FILE_FILTER", ".*"),
		AllowedMIMEs: mimes,
	}
}
