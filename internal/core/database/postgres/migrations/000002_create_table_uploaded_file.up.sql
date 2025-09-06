CREATE TABLE IF NOT EXISTS uploaded_file (
    internal_id UUID NOT NULL PRIMARY KEY,
    public_id VARCHAR(510) UNIQUE NOT NULL,
    file_directory TEXT NOT NULL,
    file_mime_type VARCHAR(100) NOT NULL,
    file_extension VARCHAR(3) NOT NULL,
    user_uploaded_by_internal_id UUID NOT NULL,
    uploaded_at BIGINT NOT NULL,

    CONSTRAINT fk_uploaded_file_user_uploaded_by FOREIGN KEY (user_uploaded_by_internal_id) REFERENCES users(internal_id)
);