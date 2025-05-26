CREATE SCHEMA IF NOT EXISTS postgis;

-- States --
CREATE TABLE IF NOT EXISTS postgis.states (
    id varchar(2) PRIMARY KEY,
    name varchar(50) NOT NULL,
    fips varchar(2) UNIQUE NOT NULL,
    is_offshore boolean NOT NULL
);

-- Office --
CREATE TABLE IF NOT EXISTS postgis.offices (
    id char(3) PRIMARY KEY,
    icao varchar(4) UNIQUE NOT NULL,
    name varchar(50) NOT NULL,
    state char(2) NOT NULL REFERENCES postgis.states(id),
    location geometry(Point, 4326)
);

-- County Warning Area --
CREATE TABLE IF NOT EXISTS postgis.cwas (
    id char(3) PRIMARY KEY,
    name varchar(50) NOT NULL,
    area real NOT NULL,
    geom geometry(MultiPolygon, 4326) NOT NULL,
    wfo char(3) NOT NULL REFERENCES postgis.offices(id),
    region char(2) NOT NULL,
    valid_from timestamptz NOT NULL
);

-- UGC (Universal Geographic Code) --
CREATE TABLE IF NOT EXISTS postgis.ugcs (
    id serial UNIQUE PRIMARY KEY,
    ugc char(6) NOT NULL,
    name varchar(256) NOT NULL,
    state char(2) NOT NULL REFERENCES postgis.states(id),
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
CREATE INDEX IF NOT EXISTS ugc_ugc ON postgis.ugcs(ugc);
CREATE INDEX IF NOT EXISTS ugc_geom ON postgis.ugcs USING GIST(geom);