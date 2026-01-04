-- +goose Up
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO categories (id, name) VALUES
    (1, 'pizza'),
    (2, 'drinks'),
    (3, 'snacks'),
    (4, 'desserts'),
    (5, 'sauces');

CREATE TABLE IF NOT EXISTS menu_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price BIGINT NOT NULL,
    category_id INT NOT NULL REFERENCES categories(id),
    image_url VARCHAR(500),
    weight INT,
    is_available BOOLEAN DEFAULT true,
    ingredients TEXT[],
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pizza_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id INT NOT NULL REFERENCES categories(id) DEFAULT 1,
    image_url VARCHAR(500),
    is_available BOOLEAN DEFAULT true,
    ingredients TEXT[],
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pizza_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pizza_id UUID NOT NULL REFERENCES pizza_items(id) ON DELETE CASCADE,
    size VARCHAR(20) NOT NULL,
    price BIGINT NOT NULL,
    weight INT,
    diameter INT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_menu_items_category ON menu_items(category_id);
CREATE INDEX idx_menu_items_available ON menu_items(is_available);
CREATE INDEX idx_pizza_variants_pizza_id ON pizza_variants(pizza_id);

-- +goose Down
DROP TABLE IF EXISTS pizza_variants;
DROP TABLE IF EXISTS pizza_items;
DROP TABLE IF EXISTS menu_items;
DROP TABLE IF EXISTS categories;
