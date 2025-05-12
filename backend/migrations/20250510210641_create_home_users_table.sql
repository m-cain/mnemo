-- +goose Up
CREATE TABLE home_users (
    home_id UUID REFERENCES homes(id),
    user_id UUID REFERENCES users(id),
    role VARCHAR(50) NOT NULL,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (home_id, user_id)
);

-- +goose Down
DROP TABLE home_users;
