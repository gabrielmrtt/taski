package projecthttprequests

import (
	"io"
	"mime/multipart"

	"github.com/gabrielmrtt/taski/internal/core"
	projectservice "github.com/gabrielmrtt/taski/internal/project/service"
)

type UpdateProjectDocumentRequest struct {
	Title   *string                `json:"title"`
	Content *string                `json:"content"`
	Version *string                `json:"version"`
	Files   []multipart.FileHeader `form:"files"`
}

func (r *UpdateProjectDocumentRequest) ToInput() projectservice.UpdateProjectDocumentInput {
	var files []core.FileInput = make([]core.FileInput, len(r.Files))
	for i, file := range r.Files {
		f, err := file.Open()
		if err != nil {
			return projectservice.UpdateProjectDocumentInput{}
		}

		f.Seek(0, io.SeekStart)

		content, err := io.ReadAll(f)
		if err != nil {
			return projectservice.UpdateProjectDocumentInput{}
		}

		files[i] = core.FileInput{
			FileName:     file.Filename,
			FileContent:  content,
			FileMimeType: file.Header.Get("Content-Type"),
		}
	}

	return projectservice.UpdateProjectDocumentInput{
		Title:   r.Title,
		Content: r.Content,
		Version: r.Version,
		Files:   files,
	}
}
