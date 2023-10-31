package bongo

import (
	"fmt"
	"github.com/bontusss/bongo/session"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/bontusss/bongo/template"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const (
	VERSION = "1.0.0"
	PORT    = "4000"
)

var goh = "\n  ____                          \n |  _ \\                         \n | |_) | ___  _ __   __ _  ___  \n |  _ < / _ \\| '_ \\ / _` |/ _ \\ \n | |_) | (_) | | | | (_| | (_) |\n |____/ \\___/|_| |_|\\__, |\\___/ \n                     __/ |      \n                    |___/       \n"

var ProjectFolders = []string{"models", "templates", "handlers", "settings", "migrations", "tmp", "public", "logs"}

type Bongo struct {
	AppName  string
	Version  string
	Debug    bool
	Logger   *slog.Logger
	config   config
	RootPath string
	Router   *chi.Mux
	Template *template.Template
	JetViews *jet.Set
	Session  *scs.SessionManager
}

type config struct {
	logger         string
	port           string
	templateEngine string
	cookie         cookieConfig
	sessionType    string
}

type cookieConfig struct {
	name     string
	lifetime string
	persist  string
	secure   string
	domain   string
}

// New creates a new instance of Bongo app.
func (b *Bongo) New(rootPath string) error {
	err := b.initBongo(rootPath, ProjectFolders)
	if err != nil {
		return err
	}

	// check that .env file exists
	if err := b.confirmEnvExists(rootPath); err != nil {
		return err
	}

	// read .env
	if err = godotenv.Load(); err != nil {
		return err
	}

	// config
	b.config = config{
		port:           os.Getenv("PORT"),
		templateEngine: strings.ToLower(os.Getenv("TEMPLATE")),

		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSIST"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
	}

	// create sessions
	sess := &session.Session{
		CookieSecure:   strings.ToLower(b.config.cookie.secure),
		CookieLifetime: strings.ToLower(b.config.cookie.lifetime),
		SessionType:    strings.ToLower(b.config.sessionType),
		CookieDomain:   strings.ToLower(b.config.cookie.domain),
		CookieName:     strings.ToLower(b.config.cookie.name),
		CookiePersist:  strings.ToLower(b.config.cookie.persist),
	}
	b.Session = sess.InitSession()

	// initialize logger.
	b.Logger = b.bLogger()

	b.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	b.Version = VERSION
	b.RootPath = rootPath
	b.Router = b.router().(*chi.Mux)

	// init Jet template
	var views = jet.NewSet(jet.NewOSFileSystemLoader(fmt.Sprintf("%s/templates", rootPath)), jet.InDevelopmentMode()) // set inDevelopmentMode when debug is true
	b.JetViews = views

	// init app template renderer engine
	b.Template = b.initTemplateEngine()

	return nil
}

// initBongo creates Bongo project folders.
func (b *Bongo) initBongo(rootPath string, folders []string) error {
	for _, folder := range folders {
		fullPath := filepath.Join(rootPath, folder)

		_, err := os.Stat(fullPath)
		if os.IsNotExist(err) {
			err := os.MkdirAll(fullPath, os.ModePerm) // check if os.ModePerm is the best permission to use.
			if err != nil {
				return err
			}
			fmt.Println("Created folder", fullPath)
		} else if err != nil {
			return err
		}
	}

	return nil

}

func (b *Bongo) Serve() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", b.config.port),
		ErrorLog:     log.Default(),
		Handler:      b.Router,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}
	if b.config.port == "" {
		fmt.Printf("PORT is not configured, defaulting to %s\n", PORT)
		b.config.port = PORT
	}
	fmt.Println(goh)
	fmt.Printf("%s is running on port %s\n", b.AppName, b.config.port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Ouch... %v", err)
	}

}

// Check if settings file exists in the settings folder on the root directory
func (b *Bongo) confirmEnvExists(rootPath string) error {
	_, err := os.Stat(fmt.Sprintf("%s/.env", rootPath))
	if err != nil {
		return err
	}
	return nil
}

//func (b *Bongo) startLoggers(loggerType string, ctx context.Context) *logger.LoggerWrapper {
//	wrapper := b.Logger.NewLoggerWrapper(loggerType, ctx)
//	return wrapper
//}

func (b *Bongo) initTemplateEngine() *template.Template {
	return &template.Template{
		Engine:   b.config.templateEngine,
		RootPath: b.RootPath,
		Port:     b.config.port,
		JetViews: b.JetViews,
	}
}

// configure logger
func (b *Bongo) bLogger() *slog.Logger {
	//var logger *slog.Logger

	//for development
	// Set Debug as default log level
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)

	if b.Debug == false {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	}
	logger := slog.New(handler)
	return logger
}
