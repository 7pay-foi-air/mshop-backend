-- +goose Up
INSERT INTO item (
    uuid_item,
    uuid_organisation,
    name,
    description,
    price,
    currency,
    sku,
    stock_quantity,
    image_url
)
VALUES 
(
    '7c0d9f7a-b48e-4bb9-a8f3-4946b278c321',
    '02f2c243-6c29-4f21-a98c-955372bc6297',
    'Premium USB-C Kabel 1m',
    'Visokokvalitetni USB-C kabel za brzo punjenje i prijenos podataka.',
    9.99,
    'EUR',
    'USB-C-001',
    150,
    'https://example.com/images/usb-c.jpg'
),
(
    'f1b9d6c7-2fc5-4e9b-b6dd-59f0be1ebb88',
    '02f2c243-6c29-4f21-a98c-955372bc6297',
    'Bluetooth Slušalice X100',
    'Bežične slušalice s redukcijom buke i 20h trajanja baterije.',
    49.90,
    'EUR',
    'HEAD-X100',
    80,
    'https://example.com/images/headphones.jpg'
);

-- +goose Down
DELETE FROM item
WHERE uuid_item IN (
    '7c0d9f7a-b48e-4bb9-a8f3-4946b278c321',
    'f1b9d6c7-2fc5-4e9b-b6dd-59f0be1ebb88'
);
