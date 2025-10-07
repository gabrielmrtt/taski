CREATE TABLE IF NOT EXISTS project_task_status (
    internal_id UUID NOT NULL PRIMARY KEY,
    public_id VARCHAR(510) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    color VARCHAR(7) NOT NULL,
    status_order INT,
    should_set_task_to_completed BOOLEAN NOT NULL DEFAULT FALSE,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    project_internal_id UUID NOT NULL,

    CONSTRAINT fk_project_task_status_project FOREIGN KEY (project_internal_id) REFERENCES project(internal_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS project_task_category (
    internal_id UUID NOT NULL PRIMARY KEY,
    public_id VARCHAR(510) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    color VARCHAR(7) NOT NULL,
    project_internal_id UUID NOT NULL,

    CONSTRAINT fk_project_task_category_project FOREIGN KEY (project_internal_id) REFERENCES project(internal_id) ON DELETE CASCADE
);