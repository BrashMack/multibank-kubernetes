CREATE TABLE clients (
                         id SERIAL PRIMARY KEY,
                         login TEXT NOT NULL UNIQUE,
                         password TEXT NOT NULL, -- хранится уже захешированный
                         created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO clients (login, password)
VALUES
    ('team237-1','$2a$10$./Ao/./m.eVplDUKA85saOAuyGeQlwNJfCShxXGX0bf7Jv9l6gWbS'),
    ('team237-2','$2a$10$iHCSOGMioujK764P3n81WOETVXNtDlX0tB/zK2IrDxwYVBMkD1Q7y'),
    ('team237-3','$2a$10$8Nj2XIxcTGCNFDKyFKD6eelxj3ht.WDOua/u6DDm9a47gv94fuwaq'),
    ('team237-4','$2a$10$ZZmpKFrTgjTmsVVMwLsa8u5zAS5Qs9ik54bwvBC/xyHfKuBE2rJOq'),
    ('team237-5','$2a$10$yWMhhAsWRJx6UUcPgypqBOy3SVsUbxeXGafumYnC6ONc6tsEvaKmK')
