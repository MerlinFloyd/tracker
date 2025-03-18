-- Drop indexes if they exist
DROP INDEX IF EXISTS idx_balance_records_address;
DROP INDEX IF EXISTS idx_balance_records_fetched_at;

-- Remove constraints if they exist
DO $$
BEGIN
    -- Check if address_format constraint exists before removing it
    IF EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'address_format'
    ) THEN
        ALTER TABLE balance_records DROP CONSTRAINT address_format;
    END IF;
END
$$;

-- Drop the balance_records table if it exists
DROP TABLE IF EXISTS balance_records;
