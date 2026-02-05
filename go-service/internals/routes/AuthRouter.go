package routes

import (
	"go-servie/internals/handler"
	middlewareV1 "go-servie/internals/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func V1Router(handler *handler.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(auth chi.Router) {
			auth.Post("/signup", handler.SignUpHandler)
			auth.Post("/login", handler.LoginHandler)
			auth.Get("/", handler.Test)
			auth.Post("/refresh", handler.Refresh)
			auth.Post("/google", handler.GoogleAuth)
			auth.Post("/test_file", handler.TestAudioDuration)
		})
		r.Route("/user", func(user chi.Router) {
			user.Use(middlewareV1.AuthMiddleware)
			user.Get("/", handler.UserInfo)
			user.Get("/coins", handler.GetUserCoinsHandler)
		})
		r.Route("/note", func(note chi.Router) {
			note.Use(middlewareV1.AuthMiddleware)
			note.Post("/audio", handler.CreateNoteHandler)
			note.Put("/{noteId}", handler.UpdateNoteHandler)
			note.Get("/", handler.UsersNotesHandler)
			note.Get("/notes/{noteId}", handler.NotesById)
			note.Delete("/{noteId}", handler.DeleteNote)
			note.Get("/{noteId}", handler.TestNote)
		})
		r.Route("/search", func(search chi.Router) {
			search.Get("/", handler.SearchHandler)
		})
		r.Route("/coin_packs", func(coins chi.Router) {
			coins.Get("/", handler.GetCoinPacksHandler)
		})
	})
	return r
}
