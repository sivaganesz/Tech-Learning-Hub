package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const uploadDir = "./uploads"

func ensureUploadDir() error {
	return os.MkdirAll(uploadDir, 0755)
}

func uploadSingle(c *gin.Context) {
	if err := ensureUploadDir(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create upload dir"})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	dst := filepath.Join(uploadDir, filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"filename": file.Filename})
}

func uploadMultiple(c *gin.Context) {
	if err := ensureUploadDir(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create upload dir"})
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad multipart form"})
		return
	}
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files provided"})
		return
	}
	saved := []string{}
	for _, f := range files {
		dst := filepath.Join(uploadDir, filepath.Base(f.Filename))
		if err := c.SaveUploadedFile(f, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		saved = append(saved, f.Filename)
	}
	c.JSON(http.StatusCreated, gin.H{"files": saved})
}

func listFiles(c *gin.Context) {
	if err := ensureUploadDir(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot access upload dir"})
		return
	}
	dirEntries, err := os.ReadDir(uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	names := []string{}
	for _, e := range dirEntries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	c.JSON(http.StatusOK, gin.H{"files": names})
}

func downloadFile(c *gin.Context) {
	name := c.Param("name")
	path := filepath.Join(uploadDir, filepath.Base(name))
	// simple existence check
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	c.File(path)
}

func main() {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MB

	router.POST("/upload", uploadSingle)
	router.POST("/upload/multi", uploadMultiple)
	router.GET("/files", listFiles)
	router.GET("/files/:name", downloadFile)

	// allow static access too if desired:
	// router.Static("/uploads", uploadDir)

	router.Run(":8080")
}
