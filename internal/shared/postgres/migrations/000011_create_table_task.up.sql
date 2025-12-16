CREATE TABLE task (
    internal_id UUID NOT NULL PRIMARY KEY,
    public_id VARCHAR(510) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(510),
    estimated_minutes INT NOT NULL,
    priority_level INT NOT NULL,
    due_date BIGINT,
    completed_at BIGINT,
    type VARCHAR(100) NOT NULL,
    project_task_status_internal_id UUID NOT NULL,
    project_task_category_internal_id UUID NOT NULL,
    parent_task_internal_id UUID,
    project_internal_id UUID NOT NULL,
    user_completed_internal_id UUID,
    user_creator_internal_id UUID NOT NULL,
    user_editor_internal_id UUID,
    created_at BIGINT NOT NULL,
    updated_at BIGINT,
    deleted_at BIGINT,

    CONSTRAINT fk_task_project_task_status FOREIGN KEY (project_task_status_internal_id) REFERENCES project_task_status(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_task_project_task_category FOREIGN KEY (project_task_category_internal_id) REFERENCES project_task_category(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_task_parent_task FOREIGN KEY (parent_task_internal_id) REFERENCES task(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_task_project FOREIGN KEY (project_internal_id) REFERENCES project(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_task_user_completed FOREIGN KEY (user_completed_internal_id) REFERENCES users(internal_id) ON DELETE SET NULL,
    CONSTRAINT fk_task_user_creator FOREIGN KEY (user_creator_internal_id) REFERENCES users(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_task_user_editor FOREIGN KEY (user_editor_internal_id) REFERENCES users(internal_id) ON DELETE SET NULL
);

CREATE TABLE sub_task (
    internal_id UUID NOT NULL PRIMARY KEY,
    public_id VARCHAR(510) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    completed_at BIGINT,
    task_internal_id UUID NOT NULL,

    CONSTRAINT fk_sub_task_task FOREIGN KEY (task_internal_id) REFERENCES task(internal_id) ON DELETE CASCADE
);

CREATE TABLE task_user (
    task_internal_id UUID NOT NULL,
    user_internal_id UUID NOT NULL,

    PRIMARY KEY (task_internal_id, user_internal_id),

    CONSTRAINT fk_task_user_task FOREIGN KEY (task_internal_id) REFERENCES task(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_task_user_user FOREIGN KEY (user_internal_id) REFERENCES users(internal_id) ON DELETE CASCADE
);

CREATE TABLE task_comment (
    internal_id UUID NOT NULL PRIMARY KEY,
    public_id VARCHAR(510) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    task_internal_id UUID NOT NULL,
    user_author_internal_id UUID NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT,

    CONSTRAINT fk_task_comment_task FOREIGN KEY (task_internal_id) REFERENCES task(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_task_comment_user_author FOREIGN KEY (user_author_internal_id) REFERENCES users(internal_id) ON DELETE CASCADE
);

CREATE TABLE task_comment_file (
    internal_id UUID NOT NULL PRIMARY KEY,
    task_comment_internal_id UUID NOT NULL,
    file_internal_id UUID NOT NULL,

    CONSTRAINT fk_task_comment_file_task_comment FOREIGN KEY (task_comment_internal_id) REFERENCES task_comment(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_task_comment_file_file FOREIGN KEY (file_internal_id) REFERENCES uploaded_file(internal_id) ON DELETE CASCADE
);

CREATE TABLE task_action (
    internal_id UUID NOT NULL PRIMARY KEY,
    public_id VARCHAR(510) UNIQUE NOT NULL,
    type VARCHAR(100) NOT NULL,
    task_internal_id UUID NOT NULL,
    user_internal_id UUID NOT NULL,
    created_at BIGINT NOT NULL,

    CONSTRAINT fk_task_action_task FOREIGN KEY (task_internal_id) REFERENCES task(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_task_action_user FOREIGN KEY (user_internal_id) REFERENCES users(internal_id) ON DELETE CASCADE
);