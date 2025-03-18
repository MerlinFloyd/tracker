-- Create the balance_records table if it doesn't exist
CREATE TABLE IF NOT EXISTS balance_records (
    id SERIAL PRIMARY KEY,
    address VARCHAR(42) NOT NULL,
    balance TEXT NOT NULL,
    balance_eth TEXT NOT NULL,
    fetched_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add constraints only if they don't exist
DO $$
BEGIN
    -- Check if address_format constraint exists before adding it
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'address_format'
    ) THEN
        ALTER TABLE balance_records 
        ADD CONSTRAINT address_format CHECK (address ~ '^0x[a-fA-F0-9]{40}$');
    END IF;
END
$$;

-- Create indexes only if they don't exist
DO $$
BEGIN
    -- Check if index on address exists before creating it
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes WHERE indexname = 'idx_balance_records_address'
    ) THEN
        CREATE INDEX idx_balance_records_address ON balance_records(address);
    END IF;
    
    -- Check if index on fetched_at exists before creating it
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes WHERE indexname = 'idx_balance_records_fetched_at'
    ) THEN
        CREATE INDEX idx_balance_records_fetched_at ON balance_records(fetched_at);
    END IF;
END
$$;
