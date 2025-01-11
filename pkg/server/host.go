package server

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/justinas/alice"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/auth"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/contextdata"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/files"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/infrastructure"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/projects"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/users"
)

/*
 */

type Server struct {
	//ProjectStore   projects.ProjectStore
	ProjectStore      projects.ProjectStore
	FileStore         files.FileStore
	UserStore         users.UserStore
	RefreshTokenStore auth.RefreshTokenStore
	Infrastructure    infrastructure.Infrastructure
	Port              int
	Logger            *slog.Logger
}

type timeHandler struct {
	format string
}

func (th timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(th.format)
	w.Write([]byte("The time is: " + tm))
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed in %v", time.Since(start))
	})
}

func (server *Server) Serve() error {
	slog.Info("Running server")

	if err := server.Infrastructure.Create(); err != nil {
		return err
	}

	/*secret, err := getJWTSecret()
	if err != nil {
		return err
	}
	fmt.Println("Secret: ", string(secret))
	user := "user1"
	token, err := generateJWT(user)
	if err != nil {
		return err
	}

	println("User: ", user, "   Sample JWT:", token)*/

	mux := http.NewServeMux()
	th := timeHandler{format: time.RFC1123}

	//projectService := projects.NewService(server.ProjectStore)

	//r := fiber.Registering{}
	//projects.AddRoutes(app.Group("/api/v1/projects"), projectService)
	//projects.AddRoutes(v1.Group("/projects"), projectService)
	mux.Handle("POST /login", auth.CreateLoginHandler(server.UserStore, server.RefreshTokenStore))
	mux.Handle("POST /refresh-token", auth.CreateRefreshTokenHandler(server.RefreshTokenStore))
	chain := alice.New(LoggingMiddleware, contextdata.AddDataToContextMiddleware, JWTMiddleware)

	projects.AddRoutes(mux, &chain, server.ProjectStore) // projectService)
	mux.Handle("GET /time", th)
	mux.Handle("POST /register", auth.CreateRegisterUserHandler(server.UserStore))
	//handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.NewLogLogger(server.Logger.Handler(), slog.LevelDebug)
	addr := fmt.Sprintf("0.0.0.0:%v", server.Port)
	slog.Info("Listening...", "address", addr)
	s := &http.Server{Addr: addr, Handler: mux, ErrorLog: logger}
	return s.ListenAndServe()
	//return http.ListenAndServe(addr, mux)
}
