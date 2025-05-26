-- Logs
CREATE TABLE IF NOT EXISTS public.logs (
    id serial,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    level varchar(5),
	source varchar(64),
    message text,
    PRIMARY KEY (id),
)

-- Create a single range partition of a table
CREATE OR REPLACE FUNCTION public.CREATE_YEARLY_RANGE_PARTITION (TABLE_NAME TEXT, YEAR INTEGER) RETURNS VOID AS $$
DECLARE
    start_date DATE := make_date(year, 1, 1);
    end_date DATE := make_date(year + 1, 1, 1);
    partition_name TEXT := format('%s_%s', table_name, year);
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = partition_name) THEN
    	EXECUTE format('
        	CREATE TABLE IF NOT EXISTS %s PARTITION OF %s
        	FOR VALUES FROM (%L) TO (%L);',
        	partition_name, table_name, start_date, end_date);
			RAISE NOTICE 'Created partition: %', partition_name;
    ELSE
        RAISE NOTICE 'Partition already exists: %', partition_name;
    END IF;
END
$$ LANGUAGE PLPGSQL;

-- Create a single list partition of a table
CREATE OR REPLACE FUNCTION public.CREATE_YEARLY_LIST_PARTITION (TABLE_NAME TEXT, YEAR INTEGER) RETURNS VOID AS $$
DECLARE
    partition_name TEXT := format('%s_%s', table_name, year);
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = partition_name) THEN
    	EXECUTE format('
        	CREATE TABLE IF NOT EXISTS %s PARTITION OF %s
        	FOR VALUES IN (%L);',
        	partition_name, table_name, year);
		RAISE NOTICE 'Created partition: %', partition_name;
    ELSE
        RAISE NOTICE 'Partition already exists: %', partition_name;
    END IF;
END
$$ LANGUAGE PLPGSQL;