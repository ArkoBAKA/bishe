package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"server/internal/config"
	apphttp "server/internal/http"
	"server/internal/model"
	"server/internal/storage"
)

func main() {
	genDoc := flag.Bool("gen-doc", false, "")
	docOut := flag.String("doc-out", "", "")
	docFormat := flag.String("doc-format", "md", "")
	seedAdmin := flag.Bool("seed-admin", false, "")
	flag.Parse()

	cfg := config.Load()

	if cfg.App.GinMode != "" {
		_ = os.Setenv("GIN_MODE", cfg.App.GinMode)
	}

	if *genDoc {
		engine := apphttp.NewRouter(apphttp.Deps{
			Config: cfg,
		})
		_ = engine

		var out []byte
		switch *docFormat {
		case "json":
			out = apphttp.DocsJSON()
		case "swagger", "openapi":
			out = apphttp.DocsSwaggerJSON()
		default:
			out = []byte(apphttp.DocsMarkdown())
		}

		if *docOut == "" {
			_, _ = os.Stdout.Write(out)
			if len(out) == 0 || out[len(out)-1] != '\n' {
				_, _ = os.Stdout.Write([]byte("\n"))
			}
			return
		}

		outPath := *docOut
		if !filepath.IsAbs(outPath) {
			wd, err := os.Getwd()
			if err == nil {
				outPath = filepath.Join(wd, outPath)
			}
		}
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		if err := os.WriteFile(outPath, out, 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := storage.NewMySQL(cfg)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	if *seedAdmin {
		if err := db.AutoMigrate(&model.User{}); err != nil {
			log.Fatal(err)
		}
		if err := ensureAdminAccount(db, cfg); err != nil {
			log.Fatal(err)
		}
		_ = sqlDB.Close()
		return
	}

	if cfg.App.AutoMigrate {
		if err := db.AutoMigrate(
			&model.User{},
			&model.UploadFile{},
			&model.Forum{},
			&model.Post{},
			&model.Comment{},
			&model.Like{},
			&model.ForumFollow{},
			&model.Follow{},
			&model.Report{},
			&model.Notification{},
		); err != nil {
			log.Fatal(err)
		}
	}

	if err := ensureAdminAccount(db, cfg); err != nil {
		log.Fatal(err)
	}

	rdb, err := storage.NewRedis(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	engine := apphttp.NewRouter(apphttp.Deps{
		DB:     db,
		Redis:  rdb,
		Config: cfg,
	})

	srv := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("http listening on %s", cfg.HTTP.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = rdb.Close()
	_ = srv.Shutdown(shutdownCtx)
	_ = sqlDB.Close()
}

func ensureAdminAccount(db *gorm.DB, cfg config.Config) error {
	account := strings.TrimSpace(cfg.Admin.Account)
	password := cfg.Admin.Password
	if account == "" || password == "" {
		return nil
	}

	nickname := strings.TrimSpace(cfg.Admin.Nickname)
	if nickname == "" {
		nickname = "管理员"
	}

	resetPassword := cfg.Admin.ResetPassword

	log.Printf("admin seed enabled: account=%s resetPassword=%v", account, resetPassword)

	var u model.User
	err := db.Where("account = ?", account).First(&u).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		u = model.User{
			Username:     account,
			Account:      account,
			PasswordHash: string(hash),
			Nickname:     nickname,
			AvatarURL:    "",
			Bio:          "",
			Role:         "admin",
			Status:       "normal",
		}
		if err := db.Create(&u).Error; err != nil {
			return err
		}
		log.Printf("admin seed created: account=%s userId=%d", account, u.ID)
		return nil
	}

	updates := map[string]any{}
	if u.Role != "admin" {
		updates["role"] = "admin"
	}
	if u.Nickname == "" {
		updates["nickname"] = nickname
	}
	if resetPassword {
		updates["password_hash"] = string(hash)
	}
	if len(updates) == 0 {
		log.Printf("admin seed skipped: account=%s userId=%d", account, u.ID)
		return nil
	}
	if err := db.Model(&model.User{}).Where("id = ?", u.ID).Updates(updates).Error; err != nil {
		return err
	}
	log.Printf("admin seed updated: account=%s userId=%d", account, u.ID)
	return nil
}
