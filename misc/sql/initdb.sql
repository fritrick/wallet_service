CREATE TABLE wallet (
    wid SERIAL PRIMARY KEY,
    name varchar(40) NOT NULL CHECK (name <> '') UNIQUE,
    balance INTEGER DEFAULT 0,
    user_id INTEGER
);

CREATE TABLE transactions (
    tid SERIAL PRIMARY KEY,
    wid INTEGER,
    amount INTEGER,
    create_date timestamptz NOT NULL DEFAULT now(),
    client_operation_hash VARCHAR(40) NOT NULL UNIQUE,
    CONSTRAINT fk_wallet
      FOREIGN KEY(wid) 
	  REFERENCES wallet(wid)
);