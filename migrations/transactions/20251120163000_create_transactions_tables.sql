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

CREATE TYPE transaction_type_enum AS ENUM ('Purchase', 'Refund');


CREATE TABLE transaction (
  uuid_transaction uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  total_amount numeric(10,2) NOT NULL CHECK (total_amount >= 0),
  currency text NOT NULL,
  payment_method text NOT NULL,
  is_successful boolean NOT NULL DEFAULT false,
  transaction_type transaction_type_enum NOT NULL,
  description text,

  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  completed_at timestamptz,
  refunded_at timestamptz,

  uuid_organisation uuid NOT NULL,
  uuid_user uuid NOT NULL,

  uuid_refund_to_transaction uuid,

  CONSTRAINT fk_refund_to_transaction
    FOREIGN KEY (uuid_refund_to_transaction)
    REFERENCES transaction(uuid_transaction)
);

CREATE TRIGGER transaction_set_updated_at
  BEFORE UPDATE ON transaction
  FOR EACH ROW
  EXECUTE FUNCTION set_updated_at_column();


CREATE TABLE transaction_item (
  uuid_transaction_item uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  item_name text NOT NULL,
  item_price numeric(10,2) NOT NULL CHECK (item_price >= 0),
  quantity integer NOT NULL CHECK (quantity > 0),
  subtotal numeric(10,2) NOT NULL CHECK (subtotal >= 0),

  uuid_item uuid NOT NULL,
  uuid_transaction uuid NOT NULL,

  CONSTRAINT fk_transaction_item_tx
    FOREIGN KEY (uuid_transaction)
    REFERENCES transaction(uuid_transaction)
    ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS transaction_item;

DROP TRIGGER IF EXISTS transaction_set_updated_at ON transaction;
DROP TABLE IF EXISTS transaction;

DROP TYPE IF EXISTS transaction_type_enum;

DROP FUNCTION IF EXISTS set_updated_at_column;

-- +goose StatementEnd
