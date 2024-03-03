CREATE TABLE users
(
    "id"          INTEGER PRIMARY KEY AUTOINCREMENT,
    "telegram_id" INTEGER   NOT NULL UNIQUE,
    "created_at"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE crypto_wallets
(
    "id"         INTEGER PRIMARY KEY AUTOINCREMENT,
    "user_id"    INTEGER      NOT NULL,
    "name"       VARCHAR(128) NOT NULL,
    "amount"     DOUBLE       NOT NULL DEFAULT 0,
    "created_at" TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, name),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE transactions
(
    "id"         INTEGER PRIMARY KEY AUTOINCREMENT,
    "wallet_id"  INTEGER   NOT NULL,
    "amount"     DOUBLE    NOT NULL,
    "price"      DOUBLE    NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (wallet_id) REFERENCES crypto_wallets (id)
);

CREATE TABLE subscriptions
(
    "id"               INTEGER PRIMARY KEY AUTOINCREMENT,
    "type"             VARCHAR(128) NOT NULL,
    "user_id"          INTEGER      NOT NULL,
    "data"             VARCHAR      NOT NULL,
    "notify_interval"  VARCHAR      NOT NULL,
    "last_notified_at" TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
