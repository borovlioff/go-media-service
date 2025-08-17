package handlers

import (
	"github.com/gin-gonic/gin"
	"media-server/config"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func ListFiles(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageStr := c.Query("page")
		page, _ := strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}

		search := strings.ToLower(c.Query("search"))

		files, err := os.ReadDir(cfg.UploadDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read dir"})
			return
		}

		// фильтр по regex + поиск по имени
		re := regexp.MustCompile(cfg.FileFilter)
		filtered := []string{}
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			name := f.Name()
			if re.MatchString(name) && (search == "" || strings.Contains(strings.ToLower(name), search)) {
				filtered = append(filtered, name)
			}
		}

		// сортировка (новые сверху, по имени)
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i] > filtered[j]
		})

		// пагинация
		start := (page - 1) * cfg.PageLimit
		end := start + cfg.PageLimit
		if start > len(filtered) {
			start = len(filtered)
		}
		if end > len(filtered) {
			end = len(filtered)
		}

		list := []string{}
		for _, name := range filtered[start:end] {
			list = append(list, "/public/"+filepath.Base(name))
		}

		c.JSON(http.StatusOK, gin.H{
			"files": list,
			"total": len(filtered),
		})
	}
}
