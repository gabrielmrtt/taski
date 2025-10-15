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

ALTER TABLE roles ADD CONSTRAINT fk_roles_organization FOREIGN KEY (organization_internal_id) REFERENCES organization(internal_id) ON DELETE CASCADE;

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