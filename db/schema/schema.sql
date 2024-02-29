CREATE TABLE "users"
(
    "id"          INTEGER PRIMARY KEY AUTOINCREMENT,
    "telegram_id" INTEGER   NOT NULL UNIQUE,
    "created_at"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "crypto_wallets"
(
    "id"         INTEGER PRIMARY KEY AUTOINCREMENT,
    "user_id"    INTEGER      NOT NULL,
    "name"       VARCHAR(128) NOT NULL,
    "created_at" TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE "transactions"
(
    "id"         INTEGER PRIMARY KEY AUTOINCREMENT,
    "wallet_id"  INTEGER   NOT NULL,
    "amount"     INTEGER   NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (wallet_id) REFERENCES crypto_wallets (id)
);