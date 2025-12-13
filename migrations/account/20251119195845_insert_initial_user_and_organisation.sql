-- +goose Up
INSERT INTO organisation (
    uuid_organisation,
    name,
    oib,
    address,
    contact_email,
    contact_phone
)
VALUES (
           '02f2c243-6c29-4f21-a98c-955372bc6297',
           'Mshop d.o.o.',
           '12345678901',
           'Ilica 1, Zagreb',
           'kontakt@mshop.hr',
           '+385911234567'
       );
INSERT INTO user_account (
    uuid_user,
    first_name,
    last_name,
    username,
    email,
    phone_number,
    date_of_birth,
    address,
    password_hash,
    role,
    uuid_organisation
)
VALUES (
           '1f859256-33af-45b4-87e4-06c4b6b2ebf0',
           'Ivan',
           'Ivanic',
           'iivanic7',
           'iivanic@gmail.com',
           '+385911234568',
           '1990-01-01',
           'Ilica 1, Zagreb',
           '$2a$12$dJcfMG7QNfeZ0.JnDBwABuD0kuyXwJViLIaUPftrJUpMCimWgTsDa',
           'owner',
           '02f2c243-6c29-4f21-a98c-955372bc6297'
       );
INSERT INTO user_account (
    uuid_user,
    first_name,
    last_name,
    username,
    email,
    phone_number,
    date_of_birth,
    address,
    password_hash,
    role,
    uuid_organisation
)
VALUES (
           '209c27af-1034-4061-8920-d9924bfa858f',
           'Marko',
           'Markic',
           'mmarkic7',
           'mmarkic@gmail.com',
           '+385932234566',
           '1995-03-22',
           'Istarska Ulica 5, Zagreb',
           '$2a$12$A.5I9VWyQXX9jHRQIVAVEeyGSRhiAZn3oXMh8Vws6V7mHVnMHNAai',
           'cashier',
           '02f2c243-6c29-4f21-a98c-955372bc6297'
       );

-- +goose Down
DELETE FROM user_account
WHERE uuid_user = '1f859256-33af-45b4-87e4-06c4b6b2ebf0';

DELETE FROM organisation
WHERE uuid_organisation = '02f2c243-6c29-4f21-a98c-955372bc6297';

