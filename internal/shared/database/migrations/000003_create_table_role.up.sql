CREATE TABLE IF NOT EXISTS permissions (
    internal_id UUID NOT NULL PRIMARY KEY,
    slug VARCHAR(510) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(510)
);

CREATE TABLE IF NOT EXISTS roles (
    internal_id UUID NOT NULL PRIMARY KEY,
    public_id VARCHAR(510) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(510),
    organization_internal_id UUID,
    user_creator_internal_id UUID,
    user_editor_internal_id UUID,
    is_system_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT,
    deleted_at BIGINT,

    CONSTRAINT fk_roles_user_creator FOREIGN KEY (user_creator_internal_id) REFERENCES users(internal_id) ON DELETE SET NULL,
    CONSTRAINT fk_roles_user_editor FOREIGN KEY (user_editor_internal_id) REFERENCES users(internal_id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS role_permission (
    role_internal_id UUID NOT NULL,
    permission_internal_id UUID NOT NULL,

    CONSTRAINT fk_role_permission_role FOREIGN KEY (role_internal_id) REFERENCES roles(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permission_permission FOREIGN KEY (permission_internal_id) REFERENCES permissions(internal_id) ON DELETE CASCADE
);