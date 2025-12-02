-- +goose Up
-- +goose StatementBegin


INSERT INTO transaction (
    uuid_transaction,
    total_amount,
    currency,
    payment_method,
    is_successful,
    transaction_type,
    description,
    created_at,
    uuid_organisation,
    uuid_user
)
VALUES (
    '3d11383f-c664-4b2e-96c1-5cf037a68391', 
    9.99,
    'EUR',
    'credit_card',
    true,
    'Purchase',
    'Kupnja USB-C kabela',
    NOW(),
    '02f2c243-6c29-4f21-a98c-955372bc6297', 
    '1f859256-33af-45b4-87e4-06c4b6b2ebf0' 
);


INSERT INTO transaction_item (
    uuid_transaction_item,
    item_name,
    item_price,
    quantity,
    subtotal,
    uuid_item,
    uuid_transaction
)
VALUES (
    'c981e99a-6be8-44d9-aa85-a823a20c2dc0',
    'Premium USB-C Kabel 1m',
    9.99,
    1,
    9.99,
    '7c0d9f7a-b48e-4bb9-a8f3-4946b278c321', 
    '3d11383f-c664-4b2e-96c1-5cf037a68391' 
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM transaction_item
WHERE uuid_transaction_item = 'c981e99a-6be8-44d9-aa85-a823a20c2dc0';

DELETE FROM transaction
WHERE uuid_transaction = '3d11383f-c664-4b2e-96c1-5cf037a68391';

-- +goose StatementEnd
