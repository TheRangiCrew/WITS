-- For vtec.warnings
ALTER TABLE vtec.warnings
ALTER COLUMN location TYPE geometry(MultiPoint, 4326)
USING
    CASE
        WHEN location IS NULL THEN NULL
        ELSE ST_Multi(location)
    END;

-- For vtec.updates
ALTER TABLE vtec.updates
ALTER COLUMN location TYPE geometry(MultiPoint, 4326)
USING
    CASE
        WHEN location IS NULL THEN NULL
        ELSE ST_Multi(location)
    END;