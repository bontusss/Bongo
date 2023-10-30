package bongo

import (
	"fmt"
	"github.com/bontusss/bongo/session"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
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
)

var ProjectFolders = []string{"models", "templates", "handlers", "settings", "migrations", "tmp", "public", "logs"}

type Bongo struct {
	AppName  string
	Version  string
	Debug    bool
	Logger   *zap.Logger
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
		logger:         os.Getenv("LOGGER"),
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

	// initialize logger
	logger := b.createLogger()
	//if b.Debug == true {
	//	logger = zap.Must(zap.NewDevelopment())
	//}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Println("Logger", err)
		}
	}(logger)
	b.Logger = logger

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
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     log.Default(),
		Handler:      b.Router,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	fmt.Println("Server is running on port", b.config.port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Error starting server")
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

//	configure loggers
//
// betterstack.com/community/guides/logging/go/zap
// Ask chatgpt to create a custom zap logger
func (b *Bongo) createLogger() *zap.Logger {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, //megabytes
		MaxAge:     7,
		MaxBackups: 3,
		LocalTime:  false,
		Compress:   false,
	})

	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(zapcore.NewCore(consoleEncoder, stdout, level), zapcore.NewCore(fileEncoder, file, level))

	return zap.New(core)
}
