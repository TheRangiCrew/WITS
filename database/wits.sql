-- States --
CREATE TABLE state (
    id varchar(2) PRIMARY KEY,
    name varchar(50) NOT NULL,
    fips varchar(2) UNIQUE NOT NULL,
    is_offshore boolean NOT NULL
);

-- Office --
CREATE TABLE office (
    id char(3) PRIMARY KEY,
    icao varchar(4) UNIQUE NOT NULL,
    name varchar(50) NOT NULL,
    state char(2) NOT NULL REFERENCES state(id),
    location geometry(Point, 4326)
);

-- County Warning Area --
CREATE TABLE cwa (
    id char(3) PRIMARY KEY,
    name varchar(50) NOT NULL,
    area real NOT NULL,
    geom geometry(MultiPolygon, 4326) NOT NULL,
    wfo char(3) NOT NULL REFERENCES office(id),
    region char(2) NOT NULL,
    valid_from timestamptz NOT NULL
);

-- UGC (Universal Geographic Code) --
CREATE TABLE ugc (
    id serial UNIQUE PRIMARY KEY,
    ugc char(6) NOT NULL,
    name varchar(256) NOT NULL,
    state char(2) NOT NULL REFERENCES state(id),
    type char(1) NOT NULL,
    number smallint NOT NULL,
    area real NOT NULL,
    geom geometry(MultiPolygon, 4326) NOT NULL,
    cwa char(3)[] NOT NULL,
    is_marine boolean,
    is_fire boolean,
    valid_from timestamptz DEFAULT CURRENT_TIMESTAMP,
    valid_to timestamptz
);
CREATE INDEX ugc_ugc ON ugc(ugc);
CREATE INDEX ugc_geom ON ugc USING GIST(geom);

-----------------------
---- Text Products ----
-----------------------
-- Product --
CREATE TABLE product (
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

--------------------------------------------------
---- VTEC Products, Segments, and Static Data ----
--------------------------------------------------
-- Phenomena types
CREATE TABLE phenomena (
    id char(2) PRIMARY KEY,
    name varchar(64) NOT NULL,
    description varchar(64)
);

-- Significance levels
CREATE TABLE significance (
    id char(1) PRIMARY KEY,
    name varchar(64) NOT NULL,
    description varchar(64)
);

-- Action types
CREATE TABLE action (
    id char(3) PRIMARY KEY,
    name varchar(64) NOT NULL,
    description varchar(64)
);

-- VTEC Event --
CREATE TABLE vtec_event (
    id serial,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    issued timestamptz NOT NULL,
    starts timestamptz,
    expires timestamptz NOT NULL,
    ends timestamptz DEFAULT NULL,
    end_initial timestamptz DEFAULT NULL,
    class char(1) NOT NULL,
    phenomena char(2) NOT NULL REFERENCES phenomena(id),
    wfo char(4) NOT NULL REFERENCES office(icao),
    significance char(1) NOT NULL REFERENCES significance(id),
    event_number smallint NOT NULL,
    year smallint NOT NULL,
    title varchar(128) NOT NULL,
    is_emergency boolean DEFAULT false,
    is_pds boolean DEFAULT false,
    polygon_start geometry(Polygon, 4326),
	PRIMARY KEY (wfo, phenomena, significance, event_number, year)
) PARTITION BY LIST (year);

-- VTEC UGC --
CREATE TABLE vtec_ugc (
    id serial,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    wfo char(4) NOT NULL,
    phenomena char(2) NOT NULL,
    significance char(1) NOT NULL,
    event_number smallint NOT NULL,
    ugc integer NOT NULL REFERENCES ugc(id),
    issued timestamptz NOT NULL,
    starts timestamptz DEFAULT NULL,
    expires timestamptz NOT NULL,
    ends timestamptz DEFAULT NULL,
    end_initial timestamptz DEFAULT NULL,
    action char(3) NOT NULL REFERENCES action(id),
    year smallint NOT NULL,
	FOREIGN KEY (wfo, phenomena, significance, event_number, year) 
        REFERENCES vtec_event(wfo, phenomena, significance, event_number, year) ON DELETE CASCADE,
    PRIMARY KEY (wfo, phenomena, significance, event_number, year, ugc)
) PARTITION BY LIST (year);

-- VTEC Event Updates --
CREATE TABLE vtec_update (
    id serial,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    issued timestamptz NOT NULL,
    starts timestamptz DEFAULT NULL,
    expires timestamptz NOT NULL,
    ends timestamptz DEFAULT NULL,
    text text NOT NULL,
    product varchar(38) NOT NULL,
    wfo char(4) NOT NULL,
    action char(3) NOT NULL,
    class char(1) NOT NULL,
    phenomena char(2) NOT NULL,
    significance char(1) NOT NULL,
    event_number smallint NOT NULL,
    year smallint NOT NULL,
    title varchar(128) NOT NULL,
    is_emergency boolean DEFAULT false,
    is_pds boolean DEFAULT false,
    polygon geometry(Polygon, 4326),
    direction int,
    location geometry(Point, 4326),
    speed int,
    speed_text varchar(30),
    tml_time timestamptz,
    ugc char(6)[],
    tornado varchar(64),
    damage varchar(64),
    hail_threat varchar(64),
    hail_tag varchar(64),
    wind_threat varchar(64),
    wind_tag varchar(64),
    flash_flood varchar(64),
    rainfall_tag varchar(64),
    flood_tag_dam varchar(64),
    spout_tag varchar(64),
    snow_squall varchar(64),
    snow_squall_tag varchar(64),
	PRIMARY KEY (wfo, phenomena, significance, event_number, year, id),
    CONSTRAINT fk_vtec_event
    FOREIGN KEY (wfo, phenomena, significance, event_number, year)
    REFERENCES vtec_event(wfo, phenomena, significance, event_number, year) ON DELETE CASCADE
) PARTITION BY LIST (year);


CREATE TABLE warning (
    id serial,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    issued timestamptz NOT NULL,
    starts timestamptz DEFAULT NULL,
    expires timestamptz NOT NULL,
    ends timestamptz DEFAULT NULL,
    end_initial timestamptz DEFAULT NULL,
    text text NOT NULL,
    wfo char(4) NOT NULL,
    action char(3) NOT NULL,
    class char(1) NOT NULL,
    phenomena char(2) NOT NULL,
    significance char(1) NOT NULL,
    event_number smallint NOT NULL,
    year smallint NOT NULL,
    title varchar(128) NOT NULL,
    is_emergency boolean DEFAULT false,
    is_pds boolean DEFAULT false,
    polygon geometry(Polygon, 4326),
    direction int,
    location geometry(Point, 4326),
    speed int,
    speed_text varchar(30),
    tml_time timestamptz,
    ugc char(6)[],
    tornado varchar(64),
    damage varchar(64),
    hail_threat varchar(64),
    hail_tag varchar(64),
    wind_threat varchar(64),
    wind_tag varchar(64),
    flash_flood varchar(64),
    rainfall_tag varchar(64),
    flood_tag_dam varchar(64),
    spout_tag varchar(64),
    snow_squall varchar(64),
    snow_squall_tag varchar(64),
	PRIMARY KEY (wfo, phenomena, significance, event_number, year, id)
);
CREATE INDEX warnings_issued ON warning(issued);
CREATE INDEX warnings_starts ON warning(starts);
CREATE INDEX warnings_expires ON warning(expires);
CREATE INDEX warnings_ends ON warning(ends);
CREATE INDEX warnings_phenomena_significance ON warning(phenomena, significance);

-------------
---- MCD ----
-------------
-- MCD --
CREATE TABLE mcd (
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

-- Create a single range partition of a table
CREATE OR REPLACE FUNCTION CREATE_YEARLY_RANGE_PARTITION (TABLE_NAME TEXT, YEAR INTEGER) RETURNS VOID AS $$
DECLARE
    start_date DATE := make_date(year, 1, 1);
    end_date DATE := make_date(year + 1, 1, 1);
    partition_name TEXT := format('%I_%s', table_name, year);
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = partition_name) THEN
    	EXECUTE format('
        	CREATE TABLE IF NOT EXISTS %I PARTITION OF %I
        	FOR VALUES FROM (%L) TO (%L);',
        	partition_name, table_name, start_date, end_date);
			RAISE NOTICE 'Created partition: %', partition_name;
    ELSE
        RAISE NOTICE 'Partition already exists: %', partition_name;
    END IF;
END
$$ LANGUAGE PLPGSQL;

-- Create a single list partition of a table
CREATE OR REPLACE FUNCTION CREATE_YEARLY_LIST_PARTITION (TABLE_NAME TEXT, YEAR INTEGER) RETURNS VOID AS $$
DECLARE
    partition_name TEXT := format('%I_%s', table_name, year);
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = partition_name) THEN
    	EXECUTE format('
        	CREATE TABLE IF NOT EXISTS %I PARTITION OF %I
        	FOR VALUES IN (%L);',
        	partition_name, table_name, year);
		RAISE NOTICE 'Created partition: %', partition_name;
    ELSE
        RAISE NOTICE 'Partition already exists: %', partition_name;
    END IF;
END
$$ LANGUAGE PLPGSQL;

-- Create all yearly table
CREATE OR REPLACE FUNCTION CREATE_YEARLY_PARTITIONS (YEAR INTEGER) RETURNS VOID AS $$
BEGIN
	-- Products
    PERFORM create_yearly_range_partition('product', year);
	EXECUTE format('
        	CREATE INDEX product_%s_product_id ON product_%s(product_id);',
        	year, year);

	-- VTEC Tables
	PERFORM create_yearly_list_partition('vtec_event', year);
    EXECUTE format('
        	CREATE INDEX vtec_event_%s_issued ON vtec_event_%s(issued);',
        	year, year);
    EXECUTE format('
        	CREATE INDEX vtec_event_%s_starts ON vtec_event_%s(starts);',
        	year, year);
    EXECUTE format('
            CREATE INDEX vtec_event_%s_expires ON vtec_event_%s(expires);',
            year, year);
    EXECUTE format('
        	CREATE INDEX vtec_event_%s_ends ON vtec_event_%s(ends);',
        	year, year);
    EXECUTE format('
        	CREATE INDEX vtec_event_%s_phenomena_significance ON vtec_event_%s(phenomena, significance);',
        	year, year);
    EXECUTE format('
        	CREATE INDEX vtec_event_%s_is_emergency ON vtec_event_%s(is_emergency) WHERE is_emergency = true;',
        	year, year);
    EXECUTE format('
        	CREATE INDEX vtec_event_%s_is_pds ON vtec_event_%s(is_pds) WHERE is_pds = true;',
        	year, year);
    EXECUTE format('
        	CREATE INDEX vtec_event_%s_polygon_start ON vtec_event_%s USING GIST (polygon_start);',
        	year, year);

	PERFORM create_yearly_list_partition('vtec_ugc', year);
    EXECUTE format('
        	CREATE INDEX vtec_ugc_%s_ugc ON vtec_ugc_%s(ugc);',
        	year, year);
    EXECUTE format('
        	CREATE INDEX vtec_ugc_%s_action ON vtec_ugc_%s(action);',
        	year, year);

	PERFORM create_yearly_list_partition('vtec_update', year);
    
	-- MCD
	PERFORM create_yearly_list_partition('mcd', year);
    EXECUTE format('
        	CREATE INDEX mcd_%s_geom ON mcd_%s USING GIST (geom);',
        	year, year);
END
$$ LANGUAGE PLPGSQL;
