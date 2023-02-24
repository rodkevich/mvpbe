BEGIN;

-- sample items to store in db
CREATE TYPE Sample_Item_Status AS ENUM ('CREATED', 'PENDING', 'COMPLETE', 'DELETED');
CREATE TABLE Sample_Item (
    item_id SERIAL PRIMARY KEY,
    start_timestamp TIMESTAMP NOT NULL,
    end_timestamp TIMESTAMP NOT NULL,
    status Sample_Item_Status NOT NULL DEFAULT 'CREATED'
);

END;