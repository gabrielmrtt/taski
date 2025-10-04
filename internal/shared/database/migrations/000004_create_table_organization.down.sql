ALTER TABLE role DROP CONSTRAINT fk_roles_organization;

DROP TABLE IF EXISTS organization_user;

DROP TABLE IF EXISTS organization;