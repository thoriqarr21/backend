CREATE TABLE barang (
    id BIGSERIAL PRIMARY KEY,
    nama_barang VARCHAR(255) NOT NULL,
    stok VARCHAR(255) NOT NULL,
    image_url VARCHAR(500) DEFAULT NULL,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);

-- Buat function untuk update timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Buat trigger
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON barang
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();