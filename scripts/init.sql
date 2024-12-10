CREATE TABLE users (
    id INT PRIMARY KEY,
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0.00
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    service_id INT DEFAULT 0,
    order_id INT DEFAULT 0,
    amount DECIMAL(15, 2) NOT NULL,
    type VARCHAR(255) NOT NULL, -- 'deposit', 'withdrawal', 'reserve', 'revenue_recognition', 'transfer_from', 'transfer_to'
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);