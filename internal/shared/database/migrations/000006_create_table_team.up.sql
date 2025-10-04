CREATE TABLE IF NOT EXISTS team (
    internal_id UUID NOT NULL PRIMARY KEY,
    public_id VARCHAR(510) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(510),
    status VARCHAR(100) NOT NULL,
    organization_internal_id UUID NOT NULL,
    user_creator_internal_id UUID NOT NULL,
    user_editor_internal_id UUID,
    created_at BIGINT NOT NULL,
    updated_at BIGINT,

    CONSTRAINT fk_team_organization FOREIGN KEY (organization_internal_id) REFERENCES organization(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_team_user_creator FOREIGN KEY (user_creator_internal_id) REFERENCES users(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_team_user_editor FOREIGN KEY (user_editor_internal_id) REFERENCES users(internal_id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS team_user (
    team_internal_id UUID NOT NULL,
    user_internal_id UUID NOT NULL,

    PRIMARY KEY (team_internal_id, user_internal_id),

    CONSTRAINT fk_team_user_team FOREIGN KEY (team_internal_id) REFERENCES team(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_team_user_user FOREIGN KEY (user_internal_id) REFERENCES users(internal_id) ON DELETE CASCADE
);