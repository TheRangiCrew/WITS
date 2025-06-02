CREATE SCHEMA IF NOT EXISTS awips;

-- Product --
CREATE TABLE IF NOT EXISTS awips.products (
    id serial,
    product_id varchar(38) NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    received_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    issued timestamptz NOT NULL,
    source char(4) NOT NULL,
    data text NOT NULL,
    wmo char(6) NOT NULL,
    awips char(6) NOT NULL,
    bbb varchar(3),
	PRIMARY KEY (id, issued),
    UNIQUE ( issued, wmo, awips, bbb, id)
) PARTITION BY RANGE (issued);

-- Create all yearly table
CREATE OR REPLACE FUNCTION awips.CREATE_YEARLY_PARTITIONS (YEAR INTEGER) RETURNS VOID AS $$
BEGIN
	-- Products
    PERFORM create_yearly_range_partition('awips.products', year);
	EXECUTE format('
        	CREATE INDEX product_%s_product_id ON awips.products_%s(product_id);',
        	year, year);
END
$$ LANGUAGE PLPGSQL;