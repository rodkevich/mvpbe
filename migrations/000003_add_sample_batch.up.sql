BEGIN;

-- batch of sample items
CREATE TYPE Sample_Batch_Status AS ENUM ('OPEN', 'PENDING', 'COMPLETE', 'DELETED');
CREATE TABLE Sample_Batch (
    batch_id SERIAL PRIMARY KEY,
    item_id INT NOT NULL REFERENCES Sample_Item(item_id),
    start_timestamp TIMESTAMP NOT NULL,
    end_timestamp TIMESTAMP NOT NULL,
    status Sample_Batch_Status NOT NULL DEFAULT 'OPEN'
);

END;
