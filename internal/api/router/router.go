package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gitlab.com/EysteinnSig/stackmap-api/internal/api/pkg/database"
	"gitlab.com/EysteinnSig/stackmap-api/internal/api/pkg/files"
	"gitlab.com/EysteinnSig/stackmap-api/internal/api/pkg/projects"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		next.ServeHTTP(w, r)
	})
}

/*
	func times() *chi.Mux {
		router := chi.NewRouter()
		router.Get("/times", func(w http.ResponseWriter, r *http.Request) {
*/
func timesRoute(w http.ResponseWriter, r *http.Request) {
	product := r.URL.Query().Get("product")
	project := r.URL.Query().Get("project")
	layers, err := database.GetAvailableTimes(project, product)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprint(err)))
		return
	}

	response, _ := json.Marshal(layers)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}

/*})
	return router
}*/

func products() *chi.Mux {
	log.Println("Setting up routing.")
	router := chi.NewRouter()
	type Layer struct {
		product string
	}

	router.Post("/projects", projects.PostHandler)
	/*router.Get("/projects", func(w http.ResponseWriter, r *http.Request) {
		projects, err := database.GetUniqueProjects()
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}

		response, _ := json.Marshal(projects)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(response)
	})*/
	router.Delete("/projects/{project}", projects.DeleteHandler)
	router.Get("/projects", projects.GetHandler)

	//router.Route("/projects/{project}", func(router chi.Router) {
	//router.Get("/products", func(w http.ResponseWriter, r *http.Request) {
	router.Get("/projects/{project}/products", func(w http.ResponseWriter, r *http.Request) {

		project := chi.URLParam(r, "project")
		layers, err := database.GetUniqueProducts(project)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}

		response, _ := json.Marshal(layers)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(response)
	})

	router.Get("/projects/{project}/files", files.GetHandler)
	//router.Route("/products/{product}/", func(router chi.Router) {
	router.Get("/projects/{project}/products/{product}/times", func(w http.ResponseWriter, r *http.Request) {
		/*w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(chi.URLParam(r, "product") + "   " + chi.URLParam(r, "project")))*/
		project := chi.URLParam(r, "project")
		product := chi.URLParam(r, "product")
		layers, err := database.GetAvailableTimes(project, product)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}

		response, _ := json.Marshal(layers)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(response)
	})
	/*},
	)*/
	//})

	/*router.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		// swagger:route GET /products pets users uniqueLayers
		//
		// Lists pets filtered by some parameters.
		//
		//     Parameters:
		//       + name: limit
		//         in: query
		//         description: maximum numnber of results to return
		//         required: false
		//         type: integer
		//         format: int32

		layers, err := database.GetUniqueProducts()
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}

		response, _ := json.Marshal(layers)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(response)
	})*/
	return router
}
func Setup() *chi.Mux {
	router := chi.NewRouter()

	/*router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	  }))*/
	router.Use(Cors)
	router.Use(middleware.Logger)
	router.Mount("/api/v1/", products())

	router.Get("/api/v1/times", timesRoute)
	//router.HandleFunc("/*", http.NotFoundHandler().ServeHTTP)

	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Problem with routing")
		http.Error(w, "404 page not found, problem with routing", http.StatusNotFound)
	})
	/*router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		db := database.GetDB()
		fmt.Println(db)
		var geom database.Raster_geoms
		db.First(&geom)
		fmt.Println(geom)
		fmt.Println("Product: ", geom.Product)

		response, _ := json.Marshal(geom)
		//fmt.Println(payload)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(response)
		//w.Write([]byte("welcome"))
	})*/
	return router
}
