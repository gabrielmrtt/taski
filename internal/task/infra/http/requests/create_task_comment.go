package taskhttprequests

import (
	"io"
	"mime/multipart"

	"github.com/gabrielmrtt/taski/internal/core"
	taskservice "github.com/gabrielmrtt/taski/internal/task/services"
)

type CreateTaskCommentRequest struct {
	Content string                 `json:"content"`
	Files   []multipart.FileHeader `json:"files"`
}

func (r *CreateTaskCommentRequest) ToInput() taskservice.CreateTaskCommentInput {
	var files []core.FileInput = make([]core.FileInput, len(r.Files))
	for i, file := range r.Files {
		f, err := file.Open()
		if err != nil {
			return taskservice.CreateTaskCommentInput{}
		}

		f.Seek(0, io.SeekStart)

		content, err := io.ReadAll(f)
		if err != nil {
			return taskservice.CreateTaskCommentInput{}
		}

		files[i] = core.FileInput{
			FileName:     file.Filename,
			FileContent:  content,
			FileMimeType: file.Header.Get("Content-Type"),
		}
	}

	return taskservice.CreateTaskCommentInput{
		Content: r.Content,
		Files:   files,
	}
}
