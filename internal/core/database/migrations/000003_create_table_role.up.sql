CREATE TABLE IF NOT EXISTS permissions (
    internal_id UUID NOT NULL PRIMARY KEY,
    slug VARCHAR(510) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(510)
);

