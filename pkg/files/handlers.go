package files

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/justinas/alice"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/utils"
)

func createUploadHandler(service FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := utils.ParseContentType(r.Header.Get("Content-Type"))
		log.Println("ContentType: ", contentType)

		if contentType != "multipart/form-data" {
			slog.Error("Unsupported Media Type", "contentType", contentType)
			http.Error(w, "Supported Media Types are multipart/form-data", http.StatusUnsupportedMediaType)
			return
		}

		err := r.ParseMultipartForm(10 << 20) // 10MB max memory
		if err != nil {
			slog.Error("Failed to parse form", "error", err)
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			slog.Error("Failed to read file", "error", err)
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		filePath := "/tmp/" + fileHeader.Filename
		dst, err := os.Create(filePath)
		if err != nil {
			slog.Error("Failed to save file", "error", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		defer dst.Close()

		_, err = dst.ReadFrom(file)
		if err != nil {
			slog.Error("Failed to save file data", "error", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		/*
			service.UploadFile(filePath)
			service.UploadFileStream(file, fileHeader.Filename)
		*/

		//service.Upload("someproj", nil, nil)
		//fmt.Fprintln(w, "File uploaded successfully")
		slog.Debug("File uploaded successfully", "filename", fileHeader.Filename)
		w.Header().Set("Content-Type", "application/json")
		ret := map[string]string{"result": "success"}
		if err := json.NewEncoder(w).Encode(ret); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

func AddRoutes(mux *http.ServeMux, chain *alice.Chain, service FileStore) {

	slog.Info("Adding route POST /api/v1/projects/{project_name}/files")
	mux.Handle("POST /api/v1/projects/{project_name}/files", chain.ThenFunc(createUploadHandler(service)))
	//mux.Handle("POST /api/v1/projects/{project_name}/files", createUploadHandler(service))

}

/*
func AddRoutes(mux *http.ServeMux, chain *alice.Chain, service ProjectStore) {

	//context.WithValue()
	mux.Handle("GET /api/v1/projects", chain.ThenFunc(createListHandler(service)))
	mux.Handle("DELETE /api/v1/projects/{project_name}", chain.ThenFunc(createDeleteHandler(service)))
	mux.Handle("POST /api/v1/projects", chain.ThenFunc(createCreateHandler(service)))

}
*/
/*

func createGetHandler(service Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		name := c.Params("project_id")

		if name == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing project name")
		}

		fmt.Println("GetProject triggered")
		return nil
	}
}

func createUploadHandler(service Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		contentType := utils.ParseContentType(c.Get("Content-Type"))
		log.Println("ContentType: ", contentType)

		switch contentType {
		case "multipart/form-data":
			fileHeader, err := c.FormFile("file")
			if err != nil {
				return c.Status(fiber.StatusBadRequest).SendString("Failed to read file")
			}
			c.SaveFile(fileHeader, "/tmp/"+fileHeader.Filename)
			service.UploadFile("/tmp/" + fileHeader.Filename)

			file, err := fileHeader.Open()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Failed to open file")
			}

			service.UploadFileStream(file, fileHeader.Filename)

			defer file.Close()
			/*
				// Read the file data into a byte slice
				fileData, err := io.ReadAll(file)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).SendString("Failed to read file data")
				}

				/*filePath := fmt.Sprintf("./%s", fileHeader.Filename)
				if err := c.SaveFile(fileHeader, filePath); err != nil {
					return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file")
				}*/

/*log.Println("Received file "+fileHeader.Filename+" with length", fileHeader.Size)
service.UploadFile(fileHeader.Filename, fileData)*/

//data.Description = c.Format()
//reflect.TypeOf(data).Elem(
/*
		default:
			return c.Status(fiber.StatusUnsupportedMediaType).SendString("Supported Media Types are multipart/form-data")
		}

		return nil
		//return c.JSON(project)
	}
}

func AddRoutes(router fiber.Router, service Service) {

	router.Post("/", createUploadHandler(service))
	router.Get("/", createGetHandler(service))
	//router.Delete("/:projectname", createDeleteHandler(service))
}
*/
