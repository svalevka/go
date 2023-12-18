package api

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/laserdigital/platform/go/pkg/encoding"
)

// None is a placeholder Request Body used to identify requests that are not
// expecting a request to have a body.
type None struct{}

// Request wraps the raw HTTP Request the the automatically unmarshaled
// request body.
type Request[T any] struct {
	*http.Request

	Body *T
}

// Response controls what and how the response is written to the client,
// including marshaling the response body.
type Response[T any] struct {
	// StatusCode overrides the default HTTP 200 OK Status Code given to tbe
	// client before the response body is marshaled.
	StatusCode int

	// Headers will append additional HTP Headers to the response.
	Headers http.Header

	// Body, when not nil, will be marshaled to the client using the
	// configured encoding.
	Body *T
}

// App is an abstraction on a Chi Router that carries configuration context
// between the top-level HTTP Method calls in this package, such as for
// encoding, error handling and logging.
type App struct {
	// Encoding configures how request bodies are read from the wire and
	// response bodies are written to the wire.
	Encoding encoding.Encoding

	// ErrorHandler is optionally invoked to handle unexpected errors, this can
	// be used to return a custom error type. If not set or returns nil, an
	// HTTP 500 Internal Server Error is returned.
	ErrorHandler func(error) any

	// Logger is the optional destination for unexpected errors and debug
	// information.
	Logger *slog.Logger

	router chi.Router
}

// New initializes a new App from a fresh router.
func New(enc encoding.Encoding, log *slog.Logger) *App {
	return From(enc, log, chi.NewRouter())
}

// From initializes an App from an existing router.
func From(enc encoding.Encoding, log *slog.Logger, r chi.Router) *App {
	return &App{
		Encoding: enc,
		Logger:   log,
		router:   r,
	}
}

// Use attaches one or more HTTP Middleware functions to the routing stack to
// be called before the request is executed.
func (a *App) Use(next ...func(http.Handler) http.Handler) {
	a.router.Use(next...)
}

// Route initializes a sub-router of the given path with an instance of App
// where further methods can be attached.
func (a *App) Route(path string, setup func(*App)) {
	a.router.Route(path, func(r chi.Router) {
		setup(&App{
			ErrorHandler: a.ErrorHandler,
			Encoding:     a.Encoding,
			Logger:       a.Logger,
			router:       r,
		})
	})
}

// NotFound configures the HTTP Handler for requests to paths that don't exist.
func (a *App) NotFound(body any) {
	a.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		a.writeResponse(w, http.StatusNotFound, body)
	})
}

// MethodNotAllowed configures the HTTP Handler for requests to paths that
// exist but are not expecting the requested HTTP Method.
func (a *App) MethodNotAllowed(body any) {
	a.router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		a.writeResponse(w, http.StatusMethodNotAllowed, body)
	})
}

func (a *App) readRequest(r *http.Request, dst any) error {
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		return nil
	}

	if _, ok := dst.(*None); ok {
		return nil
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = a.Encoding.Decode(bytes, dst)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) writeResponse(w http.ResponseWriter, status int, src any) {
	if sc, ok := src.(interface {
		StatusCode() int
	}); ok {
		status = sc.StatusCode()
	}

	if status < 1 {
		status = http.StatusOK
	}

	// only write the body if the response is not HTTP 204 No Content or
	// explicitly configured with the None type.
	body := status != http.StatusNoContent
	if _, ok := src.(*None); ok {
		body = false
	}

	if body {
		w.Header().Set("Content-Type", a.Encoding.ContentType()+"; charset=utf-8")
	}

	w.WriteHeader(status)

	if body {
		bytes, err := a.Encoding.Encode(src)
		if err != nil {
			a.Logger.Error("error marshaling response body", slog.String("error", err.Error()))
			return
		}

		_, err = w.Write(bytes)
		if err != nil {
			a.Logger.Error("error writing response body to client", slog.String("error", err.Error()))
			return
		}
	}
}

// URLParam retrieves the value of a parameterized request path, or an empty
// string if it was not set.
func URLParam(ctx context.Context, key string) string {
	return chi.URLParamFromCtx(ctx, key)
}

// Get registers an HTTP Method GET request with the Application.
func Get[REQ, RES any](app *App, path string, fn func(context.Context, *Request[REQ]) (*Response[RES], error)) {
	app.router.Get(path, handle(app, fn))
}

// Post registers an HTTP Method POST request with the Application.
func Post[REQ, RES any](app *App, path string, fn func(context.Context, *Request[REQ]) (*Response[RES], error)) {
	app.router.Post(path, handle(app, fn))
}

// Put registers an HTTP Method PUT request with the Application.
func Put[REQ, RES any](app *App, path string, fn func(context.Context, *Request[REQ]) (*Response[RES], error)) {
	app.router.Put(path, handle(app, fn))
}

// Patch registers an HTTP Method PATCH request with the Application.
func Patch[REQ, RES any](app *App, path string, fn func(context.Context, *Request[REQ]) (*Response[RES], error)) {
	app.router.Patch(path, handle(app, fn))
}

// Delete registers an HTTP Method DELETE request with the Application.
func Delete[REQ, RES any](app *App, path string, fn func(context.Context, *Request[REQ]) (*Response[RES], error)) {
	app.router.Delete(path, handle(app, fn))
}

func handle[REQ, RES any](app *App, fn func(context.Context, *Request[REQ]) (*Response[RES], error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &Request[REQ]{
			Request: r,
			Body:    new(REQ),
		}

		err := app.readRequest(r, req.Body)
		if err != nil {
			if app.ErrorHandler != nil {
				app.writeResponse(w, http.StatusInternalServerError, app.ErrorHandler(err))
				return
			}

			http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
			return
		}

		res, err := fn(r.Context(), req)
		if err != nil {
			if app.ErrorHandler != nil {
				app.writeResponse(w, http.StatusInternalServerError, app.ErrorHandler(err))
				return
			}

			http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
			return
		}

		app.writeResponse(w, res.StatusCode, res.Body)
	}
}
