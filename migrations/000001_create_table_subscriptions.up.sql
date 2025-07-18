CREATE TABLE IF NOT EXISTS subscriptions (
    user_id UUID NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    PRIMARY KEY (user_id, service_name)
);