-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE FUNCTION set_updated_at_column()
RETURNS trigger AS $$
BEGIN
  NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION normalize_organisation_contact_email()
RETURNS trigger AS $$
BEGIN
  IF NEW.contact_email IS NOT NULL THEN
    NEW.contact_email := lower(NEW.contact_email);
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE organisation(
    uuid_organisation uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    oib text NOT NULL,
    address text NOT NULL,
    contact_email text NOT NULL,
    contact_phone text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,
    is_active boolean NOT NULL DEFAULT true
);

ALTER TABLE organisation ADD CONSTRAINT organisation_oib_uq UNIQUE (oib);
ALTER TABLE organisation ADD CONSTRAINT organisation_email_uq UNIQUE (contact_email);
ALTER TABLE organisation ADD CONSTRAINT organisation_phone_uq UNIQUE (contact_phone);

CREATE TRIGGER organisation_set_updated_at
    BEFORE UPDATE ON organisation
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_column();

CREATE TRIGGER organisation_normalize_contact_email
    BEFORE INSERT OR UPDATE ON organisation
    FOR EACH ROW
    EXECUTE FUNCTION normalize_organisation_contact_email();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS organisation_set_updated_at ON organisation;
DROP TRIGGER IF EXISTS organisation_normalize_contact_email ON organisation;

DROP TABLE IF EXISTS organisation;

DROP FUNCTION IF EXISTS set_updated_at_column;
DROP FUNCTION IF EXISTS normalize_organisation_contact_email;
-- +goose StatementEnd
