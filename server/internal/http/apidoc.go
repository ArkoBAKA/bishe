package http

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

type APIDoc struct {
	Name        string
	Method      string
	Path        string
	Auth        string
	Role        string
	Description string
	Query       []FieldDoc
	Body        []FieldDoc
	Data        []FieldDoc
	ErrorCodes  []ErrorCodeDoc
	Rules       []string
}

type FieldDoc struct {
	Name     string
	Type     string
	Required bool
	Default  string
	Desc     string
}

type ErrorCodeDoc struct {
	Code int
	Desc string
	When string
}

var apiDocs = map[string]APIDoc{}

func WithDoc(doc APIDoc, h gin.HandlerFunc) gin.HandlerFunc {
	key := doc.Method + " " + doc.Path
	apiDocs[key] = doc
	return h
}

func DocsMarkdown() string {
	keys := make([]string, 0, len(apiDocs))
	for k := range apiDocs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b bytes.Buffer
	b.WriteString("# API 文档（自动生成）\n\n")
	b.WriteString("统一响应结构：\n\n")
	b.WriteString("```json\n")
	b.WriteString("{\"code\":0,\"message\":\"ok\",\"data\":{}}\n")
	b.WriteString("```\n\n")

	for _, k := range keys {
		doc := apiDocs[k]
		if doc.Name == "" {
			doc.Name = k
		}
		b.WriteString("## ")
		b.WriteString(escapeMD(doc.Name))
		b.WriteString("\n\n")

		writeKV(&b, "Method", "`"+doc.Method+"`")
		writeKV(&b, "Path", "`"+doc.Path+"`")
		if doc.Auth != "" {
			writeKV(&b, "Auth", doc.Auth)
		}
		if doc.Role != "" {
			writeKV(&b, "角色", "`"+doc.Role+"`")
		}
		if doc.Description != "" {
			writeKV(&b, "说明", doc.Description)
		}
		b.WriteString("\n")

		if len(doc.Query) > 0 || len(doc.Body) > 0 {
			b.WriteString("### 请求参数\n\n")
		}
		if len(doc.Query) > 0 {
			b.WriteString("#### Query\n")
			writeFieldTable(&b, doc.Query)
			b.WriteString("\n")
		}
		if len(doc.Body) > 0 {
			b.WriteString("#### Body\n")
			writeFieldTable(&b, doc.Body)
			b.WriteString("\n")
		}

		b.WriteString("### 响应\n\n")
		b.WriteString("#### 成功示例\n")
		b.WriteString("```json\n")
		b.WriteString("{\"code\":0,\"message\":\"ok\",\"data\":{}}\n")
		b.WriteString("```\n\n")

		if len(doc.Data) > 0 {
			b.WriteString("#### data 字段结构\n")
			writeDataTable(&b, doc.Data)
			b.WriteString("\n")
		}

		if len(doc.ErrorCodes) > 0 {
			b.WriteString("### 错误码\n")
			b.WriteString("| code | 含义 | 触发条件 |\n")
			b.WriteString("|---|---|---|\n")
			for _, e := range doc.ErrorCodes {
				b.WriteString("| ")
				b.WriteString(intToStr(e.Code))
				b.WriteString(" | ")
				b.WriteString(escapeMD(e.Desc))
				b.WriteString(" | ")
				b.WriteString(escapeMD(e.When))
				b.WriteString(" |\n")
			}
			b.WriteString("\n")
		}

		if len(doc.Rules) > 0 {
			b.WriteString("### 业务规则（实现必须遵守）\n")
			for _, r := range doc.Rules {
				r = strings.TrimSpace(r)
				if r == "" {
					continue
				}
				b.WriteString("- ")
				b.WriteString(escapeMD(r))
				b.WriteString("\n")
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

func DocsJSON() []byte {
	keys := make([]string, 0, len(apiDocs))
	for k := range apiDocs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]APIDoc, 0, len(keys))
	for _, k := range keys {
		out = append(out, apiDocs[k])
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	return b
}

func DocsSwaggerJSON() []byte {
	keys := make([]string, 0, len(apiDocs))
	for k := range apiDocs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	paths := map[string]any{}
	for _, k := range keys {
		doc := apiDocs[k]
		method := strings.ToLower(strings.TrimSpace(doc.Method))
		if method == "" {
			continue
		}

		path, pathParams := normalizeOpenAPIPath(doc.Path)
		if path == "" {
			continue
		}

		op := map[string]any{}
		if strings.TrimSpace(doc.Name) != "" {
			op["summary"] = doc.Name
		}
		if strings.TrimSpace(doc.Description) != "" {
			op["description"] = doc.Description
		}

		var tag string
		if strings.HasPrefix(path, "/api/v1/") {
			rest := strings.TrimPrefix(path, "/api/v1/")
			seg := rest
			if i := strings.Index(seg, "/"); i >= 0 {
				seg = seg[:i]
			}
			seg = strings.TrimSpace(seg)
			if seg != "" {
				tag = seg
			}
		}
		if tag != "" {
			op["tags"] = []string{tag}
		}

		params := make([]any, 0)
		for _, p := range pathParams {
			params = append(params, map[string]any{
				"name":     p,
				"in":       "path",
				"required": true,
				"schema":   map[string]any{"type": "string"},
			})
		}
		for _, q := range doc.Query {
			schema := schemaForFieldType(q.Type, false)
			if strings.TrimSpace(q.Default) != "" {
				schema["default"] = q.Default
			}
			params = append(params, map[string]any{
				"name":        q.Name,
				"in":          "query",
				"required":    q.Required,
				"description": q.Desc,
				"schema":      schema,
			})
		}
		if len(params) > 0 {
			op["parameters"] = params
		}

		if len(doc.Body) > 0 {
			hasFile := false
			for _, f := range doc.Body {
				if strings.EqualFold(strings.TrimSpace(f.Type), "file") {
					hasFile = true
					break
				}
			}

			contentType := "application/json"
			if hasFile {
				contentType = "multipart/form-data"
			}

			props := map[string]any{}
			required := make([]string, 0)
			for _, f := range doc.Body {
				props[f.Name] = schemaForFieldType(f.Type, true)
				if strings.TrimSpace(f.Desc) != "" {
					props[f.Name].(map[string]any)["description"] = f.Desc
				}
				if f.Required {
					required = append(required, f.Name)
				}
			}

			bodySchema := map[string]any{
				"type":       "object",
				"properties": props,
			}
			if len(required) > 0 {
				bodySchema["required"] = required
			}

			op["requestBody"] = map[string]any{
				"required": len(required) > 0,
				"content": map[string]any{
					contentType: map[string]any{
						"schema": bodySchema,
					},
				},
			}
		}

		if needBearerAuth(doc.Auth) {
			op["security"] = []any{
				map[string]any{
					"bearerAuth": []any{},
				},
			}
		}

		op["responses"] = map[string]any{
			"200": map[string]any{
				"description": "ok",
				"content": map[string]any{
					"application/json": map[string]any{
						"schema": responseSchemaForDocData(doc.Data),
					},
				},
			},
		}

		if docOps, ok := paths[path].(map[string]any); ok {
			docOps[method] = op
			paths[path] = docOps
		} else {
			paths[path] = map[string]any{method: op}
		}
	}

	spec := map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":   "API",
			"version": "1.0.0",
		},
		"servers": []any{
			map[string]any{"url": "/"},
		},
		"paths": paths,
		"components": map[string]any{
			"securitySchemes": map[string]any{
				"bearerAuth": map[string]any{
					"type":         "http",
					"scheme":       "bearer",
					"bearerFormat": "JWT",
				},
			},
		},
	}

	b, _ := json.MarshalIndent(spec, "", "  ")
	return b
}

func writeKV(b *bytes.Buffer, k, v string) {
	b.WriteString("- **")
	b.WriteString(k)
	b.WriteString("**: ")
	b.WriteString(v)
	b.WriteString("\n")
}

func writeFieldTable(b *bytes.Buffer, fields []FieldDoc) {
	b.WriteString("| 字段 | 类型 | 必填 | 默认值 | 说明 |\n")
	b.WriteString("|---|---|---|---|---|\n")
	for _, f := range fields {
		b.WriteString("| ")
		b.WriteString(escapeMD(f.Name))
		b.WriteString(" | ")
		b.WriteString(escapeMD(f.Type))
		b.WriteString(" | ")
		if f.Required {
			b.WriteString("是")
		} else {
			b.WriteString("否")
		}
		b.WriteString(" | ")
		b.WriteString(escapeMD(f.Default))
		b.WriteString(" | ")
		b.WriteString(escapeMD(f.Desc))
		b.WriteString(" |\n")
	}
}

func writeDataTable(b *bytes.Buffer, fields []FieldDoc) {
	b.WriteString("| 字段 | 类型 | 说明 |\n")
	b.WriteString("|---|---|---|\n")
	for _, f := range fields {
		b.WriteString("| ")
		b.WriteString(escapeMD(f.Name))
		b.WriteString(" | ")
		b.WriteString(escapeMD(f.Type))
		b.WriteString(" | ")
		b.WriteString(escapeMD(f.Desc))
		b.WriteString(" |\n")
	}
}

func escapeMD(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "|", "\\|")
	return s
}

func intToStr(n int) string {
	return strconvItoa(n)
}

func strconvItoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [32]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func normalizeOpenAPIPath(p string) (string, []string) {
	p = strings.TrimSpace(p)
	if p == "" {
		return "", nil
	}

	parts := strings.Split(p, "/")
	out := make([]string, 0, len(parts))
	params := make([]string, 0)
	for _, seg := range parts {
		if seg == "" {
			out = append(out, "")
			continue
		}
		if strings.HasPrefix(seg, ":") {
			name := strings.TrimPrefix(seg, ":")
			name = strings.TrimSpace(name)
			if name == "" {
				out = append(out, seg)
				continue
			}
			params = append(params, name)
			out = append(out, "{"+name+"}")
			continue
		}
		if strings.HasPrefix(seg, "*") {
			name := strings.TrimPrefix(seg, "*")
			name = strings.TrimSpace(name)
			if name == "" {
				out = append(out, seg)
				continue
			}
			params = append(params, name)
			out = append(out, "{"+name+"}")
			continue
		}
		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") && len(seg) > 2 {
			name := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(seg, "{"), "}"))
			if name != "" {
				params = append(params, name)
			}
			out = append(out, seg)
			continue
		}
		out = append(out, seg)
	}
	return strings.Join(out, "/"), uniqueStrings(params)
}

func uniqueStrings(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func needBearerAuth(auth string) bool {
	auth = strings.ToLower(strings.TrimSpace(auth))
	if auth == "" {
		return false
	}
	if strings.Contains(auth, "不需要") {
		return false
	}
	if strings.Contains(auth, "public") && strings.Contains(auth, "不需要") {
		return false
	}
	return strings.Contains(auth, "需要")
}

func schemaForFieldType(t string, forBody bool) map[string]any {
	t = strings.ToLower(strings.TrimSpace(t))
	switch t {
	case "file":
		return map[string]any{"type": "string", "format": "binary"}
	case "number", "int", "int64", "uint64", "float":
		return map[string]any{"type": "number"}
	case "bool", "boolean":
		return map[string]any{"type": "boolean"}
	case "array":
		return map[string]any{"type": "array", "items": map[string]any{}}
	case "object":
		return map[string]any{"type": "object"}
	default:
		_ = forBody
		return map[string]any{"type": "string"}
	}
}

func responseSchemaForDocData(data []FieldDoc) map[string]any {
	dataProps := map[string]any{}
	for _, f := range data {
		dataProps[f.Name] = schemaForFieldType(f.Type, false)
		if strings.TrimSpace(f.Desc) != "" {
			dataProps[f.Name].(map[string]any)["description"] = f.Desc
		}
	}
	dataSchema := map[string]any{"type": "object"}
	if len(dataProps) > 0 {
		dataSchema["properties"] = dataProps
	}

	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"code":    map[string]any{"type": "integer"},
			"message": map[string]any{"type": "string"},
			"data":    dataSchema,
		},
		"required": []string{"code", "message", "data"},
	}
}
