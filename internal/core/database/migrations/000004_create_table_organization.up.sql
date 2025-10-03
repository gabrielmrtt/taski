CREATE TABLE IF NOT EXISTS organization (
    internal_id UUID NOT NULl PRIMARY KEY,
    public_id VARCHAR(510) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(100) NOT NULL,
    user_creator_internal_id UUID NOT NULL,
    user_editor_internal_id UUID,
    created_at BIGINT NOT NULL,
    updated_at BIGINT,
    deleted_at BIGINT,

    CONSTRAINT fk_organizations_user_creator FOREIGN KEY (user_creator_internal_id) REFERENCES users(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_organizations_user_editor FOREIGN KEY (user_editor_internal_id) REFERENCES users(internal_id) ON DELETE SET NULL
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
    CONSTRAINT fk_roles_user_editor FOREIGN KEY (user_editor_internal_id) REFERENCES users(internal_id) ON DELETE SET NULL,
    CONSTRAINT fk_roles_organization FOREIGN KEY (organization_internal_id) REFERENCES organization(internal_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS role_permission (
    role_internal_id UUID NOT NULL,
    permission_internal_id UUID NOT NULL,

    CONSTRAINT fk_role_permission_role FOREIGN KEY (role_internal_id) REFERENCES roles(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permission_permission FOREIGN KEY (permission_internal_id) REFERENCES permissions(internal_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization_user (
    organization_internal_id UUID NOT NULL,
    user_internal_id UUID NOT NULL,
    role_internal_id UUID NOT NULL,
    status VARCHAR(100) NOT NULL,

    PRIMARY KEY (organization_internal_id, user_internal_id),

    CONSTRAINT fk_organization_user_organization FOREIGN KEY (organization_internal_id) REFERENCES organization(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_organization_user_user FOREIGN KEY (user_internal_id) REFERENCES users(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_organization_user_role FOREIGN KEY (role_internal_id) REFERENCES roles(internal_id) ON DELETE CASCADE
);