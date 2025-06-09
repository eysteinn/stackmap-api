package projects

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/justinas/alice"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/contextdata"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/utils"
)

func createDeleteHandler(service ProjectStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		project_name := r.PathValue("project_name")
		//err := service.DeleteProjectByName(project_name)

		w.Write([]byte(project_name))

	}
}

/*
func createDeleteHandler2(service ProjectStore) fiber.Handler {
	return func(c fiber.Ctx) error {
		name := c.Params("projectName")
		if name == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing project name")
		}

		err := service.DeleteProjectByName(name)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return nil
	}
}*/

func createListHandler(service ProjectStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := contextdata.GetContextData(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		projects, err := service.GetProjectsOwnedBy([]int{data.UserID})
		if err != nil {
			slog.Error("Error listing projects", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(projects); err != nil {
			slog.Error("Error encoding response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		/*slog.Debug("Triggering list handler", "user", c.Locals("userID").(string))
		projects, err := service.ListProjects()
		if err != nil {
			log.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(projects)*/
	}
}

/*
func createListHandler3(service ProjectStore) fiber.Handler {
	return func(c fiber.Ctx) error {
		slog.Debug("Triggering list handler", "user", c.Locals("userID").(string))
		projects, err := service.ListProjects()
		if err != nil {
			log.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(projects)
	}
}*/

func createCreateHandler(service ProjectStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type Data struct {
			Name        string `json:"name"`
			Description string `json:"description,omitempty"`
		}

		ctxdata, err := contextdata.GetContextData(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		var data Data
		contentType := utils.ParseContentType(r.Header.Get("Content-Type"))
		switch contentType {
		case "multipart/form-data":
			if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max memory
				slog.Error("Error parsing multipart form", "error", err)
				http.Error(w, "Error parsing multipart form data", http.StatusBadRequest)
				return
			}
			data.Name = r.FormValue("name")
			data.Description = r.FormValue("description")
		case "application/json":
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				slog.Error("Error parsing JSON:", "error", err)
				http.Error(w, "Error parsing JSON data", http.StatusBadRequest)
				return
			}
		default:
			slog.Error("Unsuported media type", "mediatype", contentType)
			http.Error(w, "Supported Media Types are multipart/form-data and application/json", http.StatusUnsupportedMediaType)
			return
		}

		/*
			b, err := json.MarshalIndent(data, "", "   ")
			if err != nil {
				http.Error(w, "Error marshaling data", http.StatusInternalServerError)
				return
			}
			log.Println(string(b))*/

		err = service.CreateProject(data.Name, data.Description, ctxdata.UserID)
		if err != nil {
			slog.Error("error creating project", "project", data.Name, "error", err)
			http.Error(w, "Error creating project", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		ret := map[string]string{"result": "success"}
		if err := json.NewEncoder(w).Encode(ret); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}

	}
}

func AddRoutes(mux *http.ServeMux, chain *alice.Chain, service ProjectStore) {

	//context.WithValue()
	mux.Handle("GET /api/v1/projects", chain.ThenFunc(createListHandler(service)))
	protectProject := chain.Append(CreateAuthProjectAccessMiddleware(service))
	mux.Handle("DELETE /api/v1/projects/{project_name}", protectProject.ThenFunc(createDeleteHandler(service)))
	mux.Handle("POST /api/v1/projects", chain.ThenFunc(createCreateHandler(service)))

}

/*
func AddRoutes3(router fiber.Router, service ProjectStore) {
	//route := &fiber.Registering{app: app, path: path}
	router.Get("/", createListHandler3(service))
	router.Post("/", createCreateHandler2(service))
	//router.Delete("/:projectname", createDeleteHandler2(service))
}

func parseFormValue(data any, c fiber.Ctx) error {
	t := reflect.TypeOf(data).Elem()
	v := reflect.ValueOf(data).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		formTag := field.Tag.Get("json")

		// Get the form value by tag name
		formValue := c.FormValue(formTag)

		// Skip empty form values
		if formValue == "" {
			continue
		}

		// Set the struct field value based on its type
		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(formValue)
		case reflect.Int:
			intValue, err := strconv.Atoi(formValue)
			if err != nil {
				return fmt.Errorf("invalid value for field %s: %s", field.Name, formValue)
			}
			fieldValue.SetInt(int64(intValue))
		// Add other types as needed
		default:
			return fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
		}
	}
	return nil
}
*/
