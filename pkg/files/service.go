package files

/*
// Service represents the project service interface

	type Service interface {
		//UploadFile(filename string, data []byte) (*File, error)
		UploadFile(filename string) (*File, error)
		UploadFileStream(r io.Reader, filename string) (*File, error)
		//GetUploadDirectory() string

}

// service implements the Service interface

	type service struct {
		uploadDirectory string
		store           FileStore
	}

// NewService creates a new project service

	func NewService(store FileStore, uploadDirectory string) Service {
		return &service{
			store:           store,
			uploadDirectory: uploadDirectory,
		}
	}

// CreateProject creates a new project

func (s *service) UploadFileStream(r io.Reader, filename string) (*File, error) {
	log.Println("Uploading streams: ", filename)
	// Write to file in chunks

	filePath := fmt.Sprintf("./%s", filename)
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Define the chunk size (1MB)
	const chunkSize = 1 << 20 // 1MB

	// Initialize the total file size
	var totalSize int64

	// Create a buffer to read the file data in chunks
	buf := make([]byte, chunkSize)

	// Read the file data in chunks and write to the file
	for {
		// Read a chunk from the reader
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read from stream: %w", err)
		}
		if n == 0 {
			break
		}

		// Write the chunk to the file
		if _, err := file.Write(buf[:n]); err != nil {
			return nil, fmt.Errorf("failed to write to file: %w", err)
		}

		// Update the total file size
		totalSize += int64(n)
	}

	// Return the file information
	return &File{
		Filename: filename,

	}, nil
}
func (s *service) UploadFile(filename string) (*File, error) {
	log.Println("Uploading file: ", filename)
	return nil, nil
}*/

/*func (s *service) GetUploadDirectory() string {
	return s.uploadDirectory
}*/
