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

CREATE TABLE item (
  uuid_item uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  uuid_organisation uuid NOT NULL,
  name text NOT NULL,
  description text,
  price numeric(10,2) NOT NULL CHECK (price >= 0),
  currency text NOT NULL DEFAULT 'EUR',
  sku text,
  stock_quantity integer NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
  image_url text,
  is_active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE UNIQUE INDEX item_sku_uq
  ON item(sku, uuid_organisation)
  WHERE deleted_at IS NULL;

CREATE TRIGGER item_set_updated_at
  BEFORE UPDATE ON item
  FOR EACH ROW
  EXECUTE FUNCTION set_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS item_set_updated_at ON item;
DROP INDEX IF EXISTS item_sku_uq;
DROP TABLE IF EXISTS item;

-- +goose StatementEnd