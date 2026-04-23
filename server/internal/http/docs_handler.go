/*
 * @Date: 2026-04-23 10:30:13
 * @LastEditors: ArongWang 3312428832@qq.com
 * @LastEditTime: 2026-04-23 11:23:20
 * @FilePath: /server/internal/http/docs_handler.go
 * @Description:
 */
package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func docsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		format := strings.ToLower(strings.TrimSpace(c.Query("format")))
		if format == "" {
			format = "md"
		}

		switch format {
		case "json":
			c.Data(http.StatusOK, "application/json; charset=utf-8", DocsJSON())
		case "swagger", "openapi":
			c.Data(http.StatusOK, "application/json; charset=utf-8", DocsSwaggerJSON())
		default:
			c.Data(http.StatusOK, "text/markdown; charset=utf-8", []byte(DocsMarkdown()))
		}
	}
}
