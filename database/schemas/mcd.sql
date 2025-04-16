CREATE SCHEMA IF NOT EXISTS mcd;

-- MCD --
CREATE TABLE IF NOT EXISTS mcd.mcd (
    id int,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    product varchar(38) NOT NULL,
    issued timestamptz NOT NULL,
    expires timestamptz NOT NULL,
    year int NOT NULL,
    concerning varchar(255) NOT NULL,
    geom geometry(Polygon, 4326) NOT NULL,
    watch_probability int,
    most_prob_tornado text,
    most_prob_gust text,
    most_prob_hail text,
	PRIMARY KEY (id, year)
) PARTITION BY LIST (year);

CREATE OR REPLACE FUNCTION mcd.CREATE_YEARLY_PARTITIONS (YEAR INTEGER) RETURNS VOID AS $$
BEGIN
    -- MCD
	PERFORM create_yearly_list_partition('mcd.mcd', year);
    EXECUTE format('
        	CREATE INDEX mcd_%s_geom ON mcd.mcd_%s USING GIST (geom);',
        	year, year);
END
$$ LANGUAGE PLPGSQL;