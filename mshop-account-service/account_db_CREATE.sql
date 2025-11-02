
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE user_role AS ENUM ('owner', 'admin', 'cashier');

CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS trigger AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION normalize_user_account_email_username()
RETURNS trigger AS $$
BEGIN
  IF NEW.email IS NOT NULL THEN
    NEW.email := lower(NEW.email);
  END IF;
  IF NEW.username IS NOT NULL THEN
    NEW.username := lower(NEW.username);
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION normalize_organisation_contact_email()
RETURNS trigger AS $$
BEGIN
  IF NEW.contact_email IS NOT NULL THEN
    NEW.contact_email := lower(NEW.contact_email);
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE organisation (
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


CREATE TABLE user_account (
  uuid_user uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  first_name text NOT NULL,
  last_name text NOT NULL,
  username text NOT NULL,
  email text NOT NULL,
  is_email_verified boolean NOT NULL DEFAULT false,
  phone_number text NOT NULL,
  date_of_birth date NOT NULL,
  address text NOT NULL,
  password_hash text NOT NULL,
  last_login_at timestamptz,
  recovery_token_hash text,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  is_active boolean NOT NULL DEFAULT true,
  email_verification_token_hash text,
  email_verification_expires_at timestamptz,
  role user_role NOT NULL DEFAULT 'cashier',
  uuid_organisation uuid,
  CONSTRAINT fk_user_org FOREIGN KEY (uuid_organisation) REFERENCES organisation (uuid_organisation) ON DELETE SET NULL
);

CREATE UNIQUE INDEX user_account_email_uq ON user_account(email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX user_account_username_uq ON user_account(username) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX user_account_phonenumber_uq ON user_account(phone_number) WHERE deleted_at IS NULL;

CREATE TRIGGER user_account_set_updated_at
BEFORE UPDATE ON user_account
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column();

CREATE TRIGGER user_account_normalize_before_insert_update
BEFORE INSERT OR UPDATE ON user_account
FOR EACH ROW
EXECUTE FUNCTION normalize_user_account_email_username();