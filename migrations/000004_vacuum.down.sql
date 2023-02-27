BEGIN;

ALTER TABLE sample_item SET (autovacuum_enabled = false);
ALTER TABLE sample_batch SET (autovacuum_enabled = false);
ALTER TABLE lock SET (autovacuum_enabled = false);

END;
