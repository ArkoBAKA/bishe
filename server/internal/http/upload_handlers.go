package http

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"server/internal/model"
)

func uploadHandler(deps Deps) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(deps.Config.Upload.AllowedExt))
	for _, e := range deps.Config.Upload.AllowedExt {
		e = strings.ToLower(strings.TrimSpace(e))
		e = strings.TrimPrefix(e, ".")
		if e == "" {
			continue
		}
		allowed[e] = struct{}{}
	}
	maxFileBytes := int64(deps.Config.Upload.MaxFileMB) * 1024 * 1024
	if maxFileBytes <= 0 {
		maxFileBytes = 50 * 1024 * 1024
	}
	maxReqBytes := int64(deps.Config.Upload.MaxRequestMB) * 1024 * 1024
	if maxReqBytes <= 0 {
		maxReqBytes = 100 * 1024 * 1024
	}

	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxReqBytes)
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			RespondFail(c, http.StatusBadRequest, 40001, "invalid multipart form")
			return
		}

		bucket := strings.TrimSpace(c.PostForm("bucket"))
		if bucket == "" {
			bucket = "public"
		}
		if bucket != "public" && bucket != "private" {
			RespondFail(c, http.StatusBadRequest, 40002, "invalid bucket")
			return
		}

		scene := strings.TrimSpace(c.PostForm("scene"))
		if len(scene) > 32 {
			scene = scene[:32]
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			RespondFail(c, http.StatusBadRequest, 40003, "file is required")
			return
		}

		if fileHeader.Size <= 0 || fileHeader.Size > maxFileBytes {
			RespondFail(c, http.StatusBadRequest, 40004, "file too large")
			return
		}

		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fileHeader.Filename), "."))
		if ext == "" {
			RespondFail(c, http.StatusBadRequest, 40005, "file extension not allowed")
			return
		}
		if _, ok := allowed[ext]; !ok {
			RespondFail(c, http.StatusBadRequest, 40005, "file extension not allowed")
			return
		}

		mimeType := fileHeader.Header.Get("Content-Type")
		if mimeType == "" {
			mimeType = mime.TypeByExtension("." + ext)
		}
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		now := time.Now()
		year, month, day := now.Date()
		storedName := newUUIDv4() + "." + ext

		relPath := filepath.Join(bucket, fmt.Sprintf("%04d", year), fmt.Sprintf("%02d", int(month)), fmt.Sprintf("%02d", day), storedName)
		baseDir := deps.Config.Upload.BaseDir
		if baseDir == "" {
			baseDir = "data/uploads"
		}

		tmpDir := filepath.Join(baseDir, "tmp")
		finalPath := filepath.Join(baseDir, relPath)
		if err := os.MkdirAll(filepath.Dir(finalPath), 0o755); err != nil {
			RespondFail(c, http.StatusInternalServerError, 50001, "create dir failed")
			return
		}
		if err := os.MkdirAll(tmpDir, 0o755); err != nil {
			RespondFail(c, http.StatusInternalServerError, 50001, "create dir failed")
			return
		}

		src, err := fileHeader.Open()
		if err != nil {
			RespondFail(c, http.StatusBadRequest, 40006, "open file failed")
			return
		}
		defer src.Close()

		tmpFile, err := os.CreateTemp(tmpDir, "upload-*."+ext)
		if err != nil {
			RespondFail(c, http.StatusInternalServerError, 50002, "create temp failed")
			return
		}
		tmpName := tmpFile.Name()

		limited := io.LimitReader(src, maxFileBytes+1)
		written, copyErr := io.Copy(tmpFile, limited)
		closeErr := tmpFile.Close()
		if copyErr != nil || closeErr != nil {
			_ = os.Remove(tmpName)
			RespondFail(c, http.StatusInternalServerError, 50003, "save file failed")
			return
		}
		if written > maxFileBytes {
			_ = os.Remove(tmpName)
			RespondFail(c, http.StatusBadRequest, 40004, "file too large")
			return
		}

		if err := os.Rename(tmpName, finalPath); err != nil {
			_ = os.Remove(tmpName)
			RespondFail(c, http.StatusInternalServerError, 50004, "finalize file failed")
			return
		}

		url := fmt.Sprintf("/api/v1/static/%s/%s", bucket, filepath.ToSlash(filepath.Dir(relPath[len(bucket)+1:]))+"/"+storedName)
		if strings.Contains(url, "//") {
			url = strings.ReplaceAll(url, "//", "/")
		}

		rec := model.UploadFile{
			Bucket:       bucket,
			Scene:        scene,
			OriginalName: fileHeader.Filename,
			StoredName:   storedName,
			Ext:          ext,
			Size:         written,
			MimeType:     mimeType,
			RelPath:      filepath.ToSlash(relPath),
			URL:          url,
			UploaderID:   ai.UserID,
		}
		if err := deps.DB.WithContext(c.Request.Context()).Create(&rec).Error; err != nil {
			_ = os.Remove(finalPath)
			RespondFail(c, http.StatusInternalServerError, 50005, "db error")
			return
		}

		RespondOK(c, map[string]any{
			"fileId": rec.ID,
			"url":    rec.URL,
		})
	}
}

func staticHandler(deps Deps) gin.HandlerFunc {
	baseDir := deps.Config.Upload.BaseDir
	if baseDir == "" {
		baseDir = "data/uploads"
	}

	return func(c *gin.Context) {
		bucket := c.Param("bucket")
		if bucket != "public" && bucket != "private" {
			c.Status(http.StatusNotFound)
			return
		}
		if bucket == "private" {
			if _, ok := RequireAuth(c, deps); !ok {
				return
			}
		}

		fp := strings.TrimPrefix(c.Param("filepath"), "/")
		fp = filepath.Clean(fp)
		if fp == "." || strings.HasPrefix(fp, "..") || strings.Contains(fp, `\`) {
			c.Status(http.StatusNotFound)
			return
		}

		full := filepath.Join(baseDir, bucket, fp)
		st, err := os.Stat(full)
		if err != nil || st.IsDir() {
			c.Status(http.StatusNotFound)
			return
		}

		etag := fmt.Sprintf("W/\"%d-%d\"", st.Size(), st.ModTime().UnixNano())
		if inm := c.GetHeader("If-None-Match"); inm != "" && inm == etag {
			c.Status(http.StatusNotModified)
			return
		}

		c.Header("ETag", etag)
		c.Header("Last-Modified", st.ModTime().UTC().Format(http.TimeFormat))
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
		c.File(full)
	}
}

func getUploadMetaHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := RequireAuth(c, deps); !ok {
			return
		}
		id, err := strconv.ParseUint(c.Param("fileId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusBadRequest, 40010, "invalid fileId")
			return
		}

		var rec model.UploadFile
		if err := deps.DB.WithContext(c.Request.Context()).First(&rec, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondFail(c, http.StatusNotFound, 40400, "not found")
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50010, "db error")
			return
		}

		RespondOK(c, map[string]any{
			"id":           rec.ID,
			"bucket":       rec.Bucket,
			"scene":        rec.Scene,
			"originalName": rec.OriginalName,
			"storedName":   rec.StoredName,
			"ext":          rec.Ext,
			"size":         rec.Size,
			"mimeType":     rec.MimeType,
			"url":          rec.URL,
			"uploaderId":   rec.UploaderID,
			"createdAt":    rec.CreatedAt.Format(time.RFC3339),
		})
	}
}

func deleteUploadHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ai, ok := RequireAuth(c, deps)
		if !ok {
			return
		}

		id, err := strconv.ParseUint(c.Param("fileId"), 10, 64)
		if err != nil {
			RespondFail(c, http.StatusBadRequest, 40010, "invalid fileId")
			return
		}

		var rec model.UploadFile
		if err := deps.DB.WithContext(c.Request.Context()).First(&rec, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				RespondFail(c, http.StatusNotFound, 40400, "not found")
				return
			}
			RespondFail(c, http.StatusInternalServerError, 50010, "db error")
			return
		}

		if ai.Role != "admin" && ai.UserID != rec.UploaderID {
			RespondFail(c, http.StatusForbidden, 40300, "forbidden")
			return
		}

		if err := deps.DB.WithContext(c.Request.Context()).Delete(&rec).Error; err != nil {
			RespondFail(c, http.StatusInternalServerError, 50010, "db error")
			return
		}

		RespondOK(c, map[string]any{})
	}
}

func newUUIDv4() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	var out [36]byte
	hex.Encode(out[0:8], b[0:4])
	out[8] = '-'
	hex.Encode(out[9:13], b[4:6])
	out[13] = '-'
	hex.Encode(out[14:18], b[6:8])
	out[18] = '-'
	hex.Encode(out[19:23], b[8:10])
	out[23] = '-'
	hex.Encode(out[24:36], b[10:16])
	return string(out[:])
}
