CREATE TABLE IF NOT EXISTS workspace (
    internal_id UUID NOT NULL PRIMARY KEY,
    public_id VARCHAR(510) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(510),
    color VARCHAR(7) NOT NULL,
    status VARCHAR(100) NOT NULL,
    organization_internal_id UUID NOT NULL,
    user_creator_internal_id UUID NOT NULL,
    user_editor_internal_id UUID,
    created_at BIGINT NOT NULL,
    updated_at BIGINT,
    deleted_at BIGINT,

    CONSTRAINT fk_workspace_organization FOREIGN KEY (organization_internal_id) REFERENCES organization(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_workspace_user_creator FOREIGN KEY (user_creator_internal_id) REFERENCES users(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_workspace_user_editor FOREIGN KEY (user_editor_internal_id) REFERENCES users(internal_id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS project (
    internal_id UUID NOT NULL PRIMARY KEY,
    public_id VARCHAR(510) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(510),
    status VARCHAR(100) NOT NULL,
    color VARCHAR(7) NOT NULL,
    priority_level INT NOT NULL,
    start_at BIGINT,
    end_at BIGINT,
    workspace_internal_id UUID NOT NULL,
    user_creator_internal_id UUID NOT NULL,
    user_editor_internal_id UUID,
    created_at BIGINT NOT NULL,
    updated_at BIGINT,
    deleted_at BIGINT,

    CONSTRAINT fk_project_workspace FOREIGN KEY (workspace_internal_id) REFERENCES workspace(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_project_user_creator FOREIGN KEY (user_creator_internal_id) REFERENCES users(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_project_user_editor FOREIGN KEY (user_editor_internal_id) REFERENCES users(internal_id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS workspace_user (
    workspace_internal_id UUID NOT NULL,
    user_internal_id UUID NOT NULL,
    status VARCHAR(100) NOT NULL,

    PRIMARY KEY (workspace_internal_id, user_internal_id),

    CONSTRAINT fk_workspace_user_workspace FOREIGN KEY (workspace_internal_id) REFERENCES workspace(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_workspace_user_user FOREIGN KEY (user_internal_id) REFERENCES users(internal_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS project_user (
    project_internal_id UUID NOT NULL,
    user_internal_id UUID NOT NULL,
    status VARCHAR(100) NOT NULL,

    PRIMARY KEY (project_internal_id, user_internal_id),

    CONSTRAINT fk_project_user_project FOREIGN KEY (project_internal_id) REFERENCES project(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_project_user_user FOREIGN KEY (user_internal_id) REFERENCES users(internal_id) ON DELETE CASCADE
);