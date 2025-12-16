CREATE TABLE project_document_version_manager (
    internal_id UUID NOT NULL PRIMARY KEY,
    public_id VARCHAR(510) UNIQUE NOT NULL,
    project_internal_id UUID NOT NULL,

    CONSTRAINT fk_project_document_version_manager_project FOREIGN KEY (project_internal_id) REFERENCES project(internal_id) ON DELETE CASCADE
);

CREATE TABLE project_document_version (
    internal_id UUID NOT NULL PRIMARY KEY,
    public_id VARCHAR(510) UNIQUE NOT NULL,
    project_document_version_manager_internal_id UUID NOT NULL,
    version VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    user_creator_internal_id UUID NOT NULL,
    user_editor_internal_id UUID,
    latest BOOLEAN NOT NULL DEFAULT FALSE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT,

    CONSTRAINT fk_project_document_version_project_document_version_manager FOREIGN KEY (project_document_version_manager_internal_id) REFERENCES project_document_version_manager(internal_id) ON DELETE CASCADE
);

CREATE TABLE project_document_file (
    internal_id UUID NOT NULL PRIMARY KEY,
    project_document_version_internal_id UUID NOT NULL,
    file_internal_id UUID NOT NULL,

    CONSTRAINT fk_project_document_file_file FOREIGN KEY (file_internal_id) REFERENCES uploaded_file(internal_id) ON DELETE CASCADE
);