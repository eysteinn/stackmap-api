package contextdata

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/user"
)

type contextKey string

const contextDataKey contextKey = "contextData"

type ContextData struct {
	User          user.User
	UsernameClaim string
	UserID        int
	//ProjectID string
	//OtherData any
}

func withContextData(ctx context.Context, data *ContextData) context.Context {
	return context.WithValue(ctx, contextDataKey, data)
}

/*
	func GetOrInitContextData(r *http.Request) *ContextData {
		data := GetContextData(r.Context())
		if data == nil {
			data = &ContextData{}
			ctx := withContextData(r.Context(), data)
			r = r.WithContext(ctx)
		}
		return data
	}
*/
func GetContextData(ctx context.Context) (*ContextData, error) {
	data, ok := ctx.Value(contextDataKey).(*ContextData)
	if ok {
		return data, nil
	}
	slog.Debug("context not found")
	return nil, fmt.Errorf("context data is nil")
}

func AddDataToContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Creating context data")
		ctx := withContextData(r.Context(), &ContextData{})
		next.ServeHTTP(w, r.WithContext(ctx))
		/*
			next.ServeHTTP(w, r.WithContext(withContextData(r.Context(), )))
			next.ServeHTTP(w, r.WithContext(context.WithValue()))
			log.Printf("Started %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
			log.Printf("Completed in %v", time.Since(start))*/
	})
}
