-- +goose Up
-- +goose StatementBegin

-- Сначала создаём функцию (если её ещё нет в других миграциях)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Затем таблицу
CREATE TABLE IF NOT EXISTS cards (
    id VARCHAR(36) PRIMARY KEY,
    card_number VARCHAR(19) NOT NULL,
    date_expire DATE NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'RUB',
    owner_id VARCHAR(36) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_cards_owner_id ON cards(owner_id);

-- И только потом триггер
CREATE TRIGGER update_cards_updated_at
    BEFORE UPDATE ON cards
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_cards_updated_at ON cards;
DROP TABLE IF EXISTS cards;
-- Опционально: DROP FUNCTION IF EXISTS update_updated_at_column(); 
-- (но лучше не дропать, если функция используется в других таблицах)
-- +goose StatementEnd