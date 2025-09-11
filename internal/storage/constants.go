package storage_core

const UploadedFileIdentityPrefix = "fil"

type SupportedImageMimeTypes string

const (
	SupportedImageMimeTypesJPEG SupportedImageMimeTypes = "image/jpeg"
	SupportedImageMimeTypesPNG  SupportedImageMimeTypes = "image/png"
	SupportedImageMimeTypesGIF  SupportedImageMimeTypes = "image/gif"
	SupportedImageMimeTypesWebP SupportedImageMimeTypes = "image/webp"
)

type SupportedVideoMimeTypes string

const (
	SupportedVideoMimeTypesMP4  SupportedVideoMimeTypes = "video/mp4"
	SupportedVideoMimeTypesMPEG SupportedVideoMimeTypes = "video/mpeg"
	SupportedVideoMimeTypesOGG  SupportedVideoMimeTypes = "video/ogg"
	SupportedVideoMimeTypesWebm SupportedVideoMimeTypes = "video/webm"
)

type SupportedPdfMimeTypes string

const (
	SupportedPdfMimeTypesPDF SupportedPdfMimeTypes = "application/pdf"
)
