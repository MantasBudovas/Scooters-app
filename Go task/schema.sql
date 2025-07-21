CREATE TABLE IF NOT EXISTS scooters (
    id UUID PRIMARY KEY,
    status TEXT NOT NULL,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    last_update TIMESTAMP
);

INSERT INTO scooters (id, status, latitude, longitude, last_update) VALUES
('11111111-c073-4e7e-814e-5d6737f4696e', 'free', 45.4215, -75.6972, NOW()),
('22222222-1e32-4427-89b5-d6844a0b80ed', 'occupied', 45.4240, -75.6950, NOW()),
('33333333-c073-4e7e-814e-5d6737f4696e', 'free', 45.4215, -75.6972, NOW()),
('44444444-b048-5e2e-79c5-3e673f4487e0', 'free', 45.4240, -75.6950, NOW()),
('55555555-a159-4e6d-81e4-3e6123f4479a', 'free', 45.4270, -75.6900, NOW());
