BEGIN;

ALTER TABLE sample_item SET (autovacuum_enabled = true);
ALTER TABLE sample_batch SET (autovacuum_enabled = true);
ALTER TABLE lock SET (autovacuum_enabled = true);

END;