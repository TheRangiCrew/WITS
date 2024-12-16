DEFINE DATABASE IF NOT EXISTS AWIPS;

USE DB AWIPS;

--------------------------------
---- Geograhpic Identifiers ----
--------------------------------
-- States --
DEFINE TABLE IF NOT EXISTS state SCHEMAFULL TYPE NORMAL;
-- ID: state:abbreviation
DEFINE FIELD IF NOT EXISTS abbreviation ON TABLE state TYPE string;
DEFINE FIELD IF NOT EXISTS name ON TABLE state TYPE string;
DEFINE FIELD IF NOT EXISTS fips on TABLE state TYPE string;
DEFINE FIELD IF NOT EXISTS ns ON TABLE state TYPE string | NULL;

-- Offshore Code Listings --
-- https://www.weather.gov/media/directives/010_pdfs/pd01003002curr.pdf --
DEFINE TABLE IF NOT EXISTS offshore SCHEMAFULL TYPE NORMAL;
-- ID offshore:abbreviation
DEFINE FIELD IF NOT EXISTS abbreviation ON TABLE offshore TYPE string;
DEFINE FIELD IF NOT EXISTS name ON TABLE offshore TYPE string;
DEFINE FIELD IF NOT EXISTS fips on TABLE offshore TYPE string;

-- Office (WFO/WSO/Centre) --
DEFINE TABLE IF NOT EXISTS office SCHEMAFULL TYPE NORMAL;
-- ID: office:code
DEFINE FIELD IF NOT EXISTS code ON TABLE office TYPE string;
DEFINE FIELD IF NOT EXISTS icao ON TABLE office TYPE string;
DEFINE FIELD IF NOT EXISTS name ON TABLE office TYPE string;
DEFINE FIELD IF NOT EXISTS state ON TABLE office TYPE record<state>;
DEFINE FIELD IF NOT EXISTS location ON TABLE office TYPE geometry<Point>;

-- CWA (County Warning Area)
DEFINE TABLE IF NOT EXISTS cwa SCHEMAFULL TYPE NORMAL;
-- ID: cwa:code
DEFINE FIELD IF NOT EXISTS code ON TABLE cwa TYPE string;
DEFINE FIELD IF NOT EXISTS name ON TABLE cwa TYPE string;
DEFINE FIELD IF NOT EXISTS centre ON TABLE cwa TYPE geometry<Point>;
DEFINE FIELD IF NOT EXISTS geometry ON TABLE cwa TYPE geometry<MultiPolygon>;
DEFINE FIELD IF NOT EXISTS area ON TABLE cwa TYPE float DEFAULT 0.0;
DEFINE FIELD IF NOT EXISTS wfo ON TABLE cwa TYPE record<office>;
DEFINE FIELD IF NOT EXISTS region ON TABLE cwa TYPE string;
DEFINE FIELD IF NOT EXISTS valid_from ON TABLE cwa TYPE datetime;

-- UGC (Universal Geographic Code) --
DEFINE TABLE IF NOT EXISTS ugc SCHEMAFULL TYPE NORMAL;
-- ID: ugc:AABCCC
-- AA : state
-- B  : "C" | "Z" | "F" where C = county, Z = zones (forecast or marine), F = fire zone
-- CCC: id
DEFINE FIELD IF NOT EXISTS name ON TABLE ugc TYPE string;
DEFINE FIELD IF NOT EXISTS state ON TABLE ugc TYPE record<state | offshore>;
DEFINE FIELD IF NOT EXISTS type ON TABLE ugc TYPE "C" | "Z" | "F";
DEFINE FIELD IF NOT EXISTS number ON TABLE ugc TYPE string;
DEFINE FIELD IF NOT EXISTS area ON TABLE ugc TYPE float DEFAULT 0.0;
DEFINE FIELD IF NOT EXISTS centre ON TABLE ugc TYPE geometry<Point>;
DEFINE FIELD IF NOT EXISTS geometry ON TABLE ugc TYPE geometry<MultiPolygon>;
DEFINE FIELD IF NOT EXISTS cwa ON TABLE ugc TYPE array<record<cwa>>;
DEFINE FIELD IF NOT EXISTS is_marine ON TABLE ugc TYPE bool;
DEFINE FIELD IF NOT EXISTS is_fire ON TABLE ugc TYPE bool;
DEFINE FIELD IF NOT EXISTS valid_from ON TABLE ugc TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS valid_to ON TABLE ugc TYPE datetime | null DEFAULT null;


-----------------------
---- Text Products ----
-----------------------
-- Text Products --
DEFINE TABLE IF NOT EXISTS product SCHEMAFULL TYPE NORMAL;

DEFINE FIELD IF NOT EXISTS id ON TABLE product TYPE 
[string, datetime, int]; // [AWIPS ID, Issued Time, Sequence]
DEFINE FIELD IF NOT EXISTS created_at ON TABLE product TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS product ON TABLE product TYPE string;
DEFINE FIELD IF NOT EXISTS issuer ON TABLE product TYPE record<office>;
DEFINE FIELD IF NOT EXISTS text ON TABLE product TYPE string;
DEFINE FIELD IF NOT EXISTS received_at ON TABLE product TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS wmo ON TABLE product TYPE {
    datatype: string,
    issued: datetime,
    original: string,
    wfo: string,
    bbb: string
};


--------------------------------------------------
---- VTEC Products, Segments, and Static Data ----
--------------------------------------------------
-- VTEC Phenomena --
DEFINE TABLE IF NOT EXISTS vtec_phenomena SCHEMAFULL TYPE NORMAL;

DEFINE FIELD IF NOT EXISTS code ON TABLE vtec_phenomena TYPE string;
DEFINE FIELD IF NOT EXISTS name ON TABLE vtec_phenomena TYPE string;

-- VTEC Actions --
DEFINE TABLE IF NOT EXISTS vtec_action SCHEMAFULL TYPE NORMAL;

DEFINE FIELD IF NOT EXISTS code ON TABLE vtec_action TYPE string;
DEFINE FIELD IF NOT EXISTS name ON TABLE vtec_action TYPE string;

-- VTEC Signifcance --
DEFINE TABLE IF NOT EXISTS vtec_significance SCHEMAFULL TYPE NORMAL;

DEFINE FIELD IF NOT EXISTS code ON TABLE vtec_significance TYPE string;
DEFINE FIELD IF NOT EXISTS name ON TABLE vtec_significance TYPE string;

-- VTEC Event --
DEFINE TABLE IF NOT EXISTS vtec_event SCHEMAFULL TYPE NORMAL;

DEFINE FIELD IF NOT EXISTS id ON TABLE vtec_event TYPE 
{ event_number: int, phenomena: string, office: string, significance: string, year: int};
DEFINE FIELD IF NOT EXISTS created_at ON TABLE vtec_event TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS updated_at ON TABLE vtec_event TYPE datetime VALUE time::now();
DEFINE FIELD IF NOT EXISTS updates ON TABLE vtec_event TYPE int DEFAULT 0;
DEFINE FIELD IF NOT EXISTS issued ON TABLE vtec_event TYPE datetime;
DEFINE FIELD IF NOT EXISTS start ON TABLE vtec_event TYPE datetime | null DEFAULT null;
DEFINE FIELD IF NOT EXISTS expires ON TABLE vtec_event TYPE datetime;
DEFINE FIELD IF NOT EXISTS end ON TABLE vtec_event TYPE datetime | null DEFAULT null;
DEFINE FIELD IF NOT EXISTS end_initial ON TABLE vtec_event TYPE datetime | null DEFAULT null;
DEFINE FIELD IF NOT EXISTS phenomena ON TABLE vtec_event TYPE record<vtec_phenomena>;
DEFINE FIELD IF NOT EXISTS office ON TABLE vtec_event TYPE record<office>;
DEFINE FIELD IF NOT EXISTS significance ON TABLE vtec_event TYPE record<vtec_significance>;
DEFINE FIELD IF NOT EXISTS event_number ON TABLE vtec_event TYPE int;
DEFINE FIELD IF NOT EXISTS title ON TABLE vtec_event TYPE string;
DEFINE FIELD IF NOT EXISTS is_emergency ON TABLE vtec_event TYPE bool;
DEFINE FIELD IF NOT EXISTS is_pds ON TABLE vtec_event TYPE bool;

-- VTEC UGC --
DEFINE TABLE IF NOT EXISTS vtec_ugc SCHEMAFULL TYPE RELATION FROM vtec_event TO ugc ENFORCED;

DEFINE FIELD IF NOT EXISTS id ON TABLE vtec_ugc TYPE
{ event_number: int, phenomena: string, office: string, significance: string, year: int, ugc: string};
DEFINE FIELD IF NOT EXISTS created_at ON TABLE vtec_ugc TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS updated_at ON TABLE vtec_ugc TYPE datetime VALUE time::now();
DEFINE FIELD IF NOT EXISTS issued ON TABLE vtec_ugc TYPE datetime;
DEFINE FIELD IF NOT EXISTS start ON TABLE vtec_ugc TYPE datetime | null DEFAULT null;
DEFINE FIELD IF NOT EXISTS expires ON TABLE vtec_ugc TYPE datetime;
DEFINE FIELD IF NOT EXISTS end ON TABLE vtec_ugc TYPE datetime | null DEFAULT null;
DEFINE FIELD IF NOT EXISTS end_initial ON TABLE vtec_ugc TYPE datetime | null DEFAULT null;
DEFINE FIELD IF NOT EXISTS action ON TABLE vtec_ugc TYPE record<vtec_action>;
DEFINE FIELD IF NOT EXISTS latest ON TABLE vtec_ugc TYPE record<vtec_history>;

-- VTEC Event History --
DEFINE TABLE IF NOT EXISTS vtec_history SCHEMAFULL TYPE NORMAL;

DEFINE FIELD IF NOT EXISTS id ON TABLE vtec_history TYPE
{ event_number: int, phenomena: string, office: string, significance: string, year: int, sequence: int};
DEFINE FIELD IF NOT EXISTS created_at ON TABLE vtec_history TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS issued ON TABLE vtec_history TYPE datetime;
DEFINE FIELD IF NOT EXISTS start ON TABLE vtec_history TYPE datetime | null DEFAULT null;
DEFINE FIELD IF NOT EXISTS expires ON TABLE vtec_history TYPE datetime;
DEFINE FIELD IF NOT EXISTS end ON TABLE vtec_history TYPE datetime | null DEFAULT null;
DEFINE FIELD IF NOT EXISTS original ON TABLE vtec_history TYPE string;
DEFINE FIELD IF NOT EXISTS title ON TABLE vtec_history TYPE string;
DEFINE FIELD IF NOT EXISTS action ON TABLE vtec_history TYPE record<vtec_action>;
DEFINE FIELD IF NOT EXISTS phenomena ON TABLE vtec_history TYPE record<vtec_phenomena>;
DEFINE FIELD IF NOT EXISTS office ON TABLE vtec_history TYPE record<office>;
DEFINE FIELD IF NOT EXISTS significance ON TABLE vtec_history TYPE record<vtec_significance>;
DEFINE FIELD IF NOT EXISTS event_number ON TABLE vtec_history TYPE int;
DEFINE FIELD IF NOT EXISTS vtec ON TABLE vtec_history TYPE {
    class: "O" | "T" | "E" | "X",
    action: string,
    wfo: string,
    phenomena: string,
    significance: string,
    event_number: int,
    start: string,
    end: string
};
---- TODO: Implement H-VTEC type
DEFINE FIELD IF NOT EXISTS h_vtec ON TABLE vtec_history FLEXIBLE TYPE option<object>;
DEFINE FIELD IF NOT EXISTS is_emergency ON TABLE vtec_history TYPE bool;
DEFINE FIELD IF NOT EXISTS is_pds ON TABLE vtec_history TYPE bool;
DEFINE FIELD IF NOT EXISTS lat_lon ON TABLE vtec_history TYPE option<{
    original: string,
    points: geometry<Polygon>
}>;
DEFINE FIELD IF NOT EXISTS polygon ON TABLE vtec_history TYPE option<geometry<Polygon>>;
DEFINE FIELD IF NOT EXISTS tags ON TABLE vtec_history FLEXIBLE TYPE option<object>;
DEFINE FIELD IF NOT EXISTS tml ON TABLE vtec_history TYPE option<{
    direction: int,
    location: geometry<Point>,
    speed: int,
    speedString: string,
    time: datetime,
    original: string
}>;
DEFINE FIELD IF NOT EXISTS product ON TABLE vtec_history TYPE record<product>;
DEFINE FIELD IF NOT EXISTS ugc ON TABLE vtec_history TYPE array<record<ugc>>;

-- RELATE VTEC Events to VTEC History --
DEFINE TABLE IF NOT EXISTS vtec_event_history SCHEMAFULL TYPE RELATION FROM vtec_event TO vtec_history ENFORCED;


------------------
---- Warnings ----
------------------
-- Warnings --
DEFINE TABLE IF NOT EXISTS warning SCHEMAFULL TYPE NORMAL;

DEFINE FIELD IF NOT EXISTS id ON TABLE warning TYPE 
{ event_number: int, phenomena: string, office: string, significance: string, year: int };
DEFINE FIELD IF NOT EXISTS created_at ON TABLE warning TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS updated_at ON TABLE warning TYPE datetime VALUE time::now();
DEFINE FIELD IF NOT EXISTS updates ON TABLE warning TYPE int DEFAULT 0;
DEFINE FIELD IF NOT EXISTS start ON TABLE warning TYPE datetime | null DEFAULT null;
DEFINE FIELD IF NOT EXISTS expires ON TABLE warning TYPE datetime;
DEFINE FIELD IF NOT EXISTS end ON TABLE warning TYPE datetime | null DEFAULT null;
DEFINE FIELD IF NOT EXISTS end_initial ON TABLE warning TYPE datetime DEFAULT time::now() + 24h;
DEFINE FIELD IF NOT EXISTS phenomena ON TABLE warning TYPE record<vtec_phenomena>;
DEFINE FIELD IF NOT EXISTS office ON TABLE warning TYPE record<office>;
DEFINE FIELD IF NOT EXISTS significance ON TABLE warning TYPE record<vtec_significance>;
DEFINE FIELD IF NOT EXISTS event_number ON TABLE warning TYPE int;
DEFINE FIELD IF NOT EXISTS title ON TABLE warning TYPE string;
DEFINE FIELD IF NOT EXISTS is_emergency ON TABLE warning TYPE bool;
DEFINE FIELD IF NOT EXISTS is_pds ON TABLE warning TYPE bool;


-- Warning UGC --
DEFINE TABLE IF NOT EXISTS warning_ugc SCHEMAFULL TYPE RELATION FROM warning TO ugc ENFORCED;

DEFINE FIELD IF NOT EXISTS id ON TABLE warning_ugc TYPE
{ event_number: int, phenomena: string, office: string, significance: string, year: int, ugc: string};
DEFINE FIELD IF NOT EXISTS created_at ON TABLE warning_ugc TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS updated_at ON TABLE warning_ugc TYPE datetime VALUE time::now();
DEFINE FIELD IF NOT EXISTS issued ON TABLE warning_ugc TYPE datetime;
DEFINE FIELD IF NOT EXISTS start ON TABLE warning_ugc TYPE datetime | null DEFAULT null;
DEFINE FIELD IF NOT EXISTS expires ON TABLE warning_ugc TYPE datetime;
DEFINE FIELD IF NOT EXISTS end ON TABLE warning_ugc TYPE datetime | null DEFAULT null;
DEFINE FIELD IF NOT EXISTS end_initial ON TABLE warning_ugc TYPE datetime | null DEFAULT null;
DEFINE FIELD IF NOT EXISTS action ON TABLE warning_ugc TYPE record<vtec_action>;
DEFINE FIELD IF NOT EXISTS latest ON TABLE warning_ugc TYPE record<warning_history>;

-- Warning History --
DEFINE TABLE IF NOT EXISTS warning_history SCHEMAFULL TYPE NORMAL;

DEFINE FIELD IF NOT EXISTS id ON TABLE warning_history TYPE
{ event_number: int, phenomena: string, office: string, significance: string, year: int, sequence: int};
DEFINE FIELD IF NOT EXISTS created_at ON TABLE warning_history TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS issued ON TABLE warning_history TYPE datetime;
DEFINE FIELD IF NOT EXISTS start ON TABLE warning_history TYPE datetime | null DEFAULT null;
DEFINE FIELD IF NOT EXISTS expires ON TABLE warning_history TYPE datetime;
DEFINE FIELD IF NOT EXISTS end ON TABLE warning_history TYPE datetime | null DEFAULT null;
DEFINE FIELD IF NOT EXISTS original ON TABLE warning_history TYPE string;
DEFINE FIELD IF NOT EXISTS title ON TABLE warning_history TYPE string;
DEFINE FIELD IF NOT EXISTS action ON TABLE warning_history TYPE record<vtec_action>;
DEFINE FIELD IF NOT EXISTS phenomena ON TABLE warning_history TYPE record<vtec_phenomena>;
DEFINE FIELD IF NOT EXISTS office ON TABLE warning_history TYPE record<office>;
DEFINE FIELD IF NOT EXISTS significance ON TABLE warning_history TYPE record<vtec_significance>;
DEFINE FIELD IF NOT EXISTS event_number ON TABLE warning_history TYPE int;
DEFINE FIELD IF NOT EXISTS vtec ON TABLE warning_history TYPE {
    class: "O" | "T" | "E" | "X",
    action: string,
    wfo: string,
    phenomena: string,
    significance: string,
    event_number: int,
    start: string,
    end: string
};
---- TODO: Implement H-VTEC type
DEFINE FIELD IF NOT EXISTS h_vtec ON TABLE warning_history FLEXIBLE TYPE option<object>;
DEFINE FIELD IF NOT EXISTS is_emergency ON TABLE warning_history TYPE bool;
DEFINE FIELD IF NOT EXISTS is_pds ON TABLE warning_history TYPE bool;
DEFINE FIELD IF NOT EXISTS polygon ON TABLE warning_history TYPE option<geometry<Polygon>>;
DEFINE FIELD IF NOT EXISTS tags ON TABLE warning_history FLEXIBLE TYPE option<object>;
DEFINE FIELD IF NOT EXISTS tml ON TABLE warning_history TYPE option<{
    direction: int,
    location: geometry<Point>,
    speed: int,
    speedString: string,
    time: datetime,
    original: string
}>;
DEFINE FIELD IF NOT EXISTS product ON TABLE warning_history TYPE record<product>;
DEFINE FIELD IF NOT EXISTS ugc ON TABLE warning_history TYPE array<record<ugc>>;

-- RELATE Warning to Warning History --
DEFINE TABLE IF NOT EXISTS warning_event_history SCHEMAFULL TYPE RELATION FROM warning TO warning_history ENFORCED;


----------------------------------------
---- Severe Weather Watches/WWs ----
----------------------------------------
-- Severe Watches --
DEFINE TABLE IF NOT EXISTS severe_watch SCHEMAFULL TYPE NORMAL;

DEFINE FIELD IF NOT EXISTS id ON TABLE severe_watch TYPE
{ number: int, phenomena: "SV" | "TO", issued: datetime };
DEFINE FIELD IF NOT EXISTS expires ON TABLE severe_watch TYPE datetime;
DEFINE FIELD IF NOT EXISTS sel ON TABLE severe_watch TYPE record<product>;
DEFINE FIELD IF NOT EXISTS wou ON TABLE severe_watch TYPE record<product>;
DEFINE FIELD IF NOT EXISTS wwp ON TABLE severe_watch TYPE {
    text: record<product>,
    degrees: int,
    speed: int,
    max_tops: number,
    max_hail: number,
    one_or_more_hail: string,
    ten_or_more_sev_hail: string,
    max_wind: number,
    one_or_more_wind: string,
    ten_or_more_sev_wind: string,
    two_or_more_tor: string,
    strong_tor: string,
    six_or_more_combo: string
};
DEFINE FIELD IF NOT EXISTS is_pds ON TABLE severe_watch TYPE bool DEFAULT false;
DEFINE FIELD IF NOT EXISTS polygon ON TABLE severe_watch TYPE geometry<Polygon>;


-------------------------------------------------
---- Mesoscale Convective Discussions (MCDs) ----
-------------------------------------------------
-- MCD --
DEFINE TABLE IF NOT EXISTS mcd SCHEMAFULL TYPE NORMAL;

DEFINE FIELD IF NOT EXISTS id ON TABLE mcd TYPE
{ number: int, year: int };
DEFINE FIELD IF NOT EXISTS created_at ON TABLE mcd TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS updated_at ON TABLE mcd TYPE datetime VALUE time::now();
DEFINE FIELD IF NOT EXISTS issued ON TABLE mcd TYPE datetime;
DEFINE FIELD IF NOT EXISTS expires ON TABLE mcd TYPE datetime;
DEFINE FIELD IF NOT EXISTS original ON TABLE mcd TYPE record<product>;
DEFINE FIELD IF NOT EXISTS concerning ON TABLE mcd TYPE string;
DEFINE FIELD IF NOT EXISTS polygon ON TABLE mcd TYPE geometry<Polygon>;
DEFINE FIELD IF NOT EXISTS watch_probability ON TABLE mcd TYPE int;

-- RELATE MCD to WW --
DEFINE TABLE IF NOT EXISTS mcd_watch SCHEMAFULL TYPE RELATION FROM mcd TO severe_watch ENFORCED;


--------------
---- Logs ----
--------------
-- Logs --
DEFINE TABLE IF NOT EXISTS log TYPE NORMAL;

DEFINE FIELD IF NOT EXISTS id ON TABLE log DEFAULT rand::uuid();
DEFINE FIELD IF NOT EXISTS created_at ON TABLE log TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS time ON TABLE log TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS level ON TABLE log TYPE string;
DEFINE FIELD IF NOT EXISTS product ON TABLE log TYPE option<record<product>>;
DEFINE FIELD IF NOT EXISTS awips ON TABLE log TYPE option<string>;
DEFINE FIELD IF NOT EXISTS wmo ON TABLE log TYPE option<string>;
DEFINE FIELD IF NOT EXISTS text ON TABLE log TYPE option<string>;
DEFINE FIELD IF NOT EXISTS message ON TABLE log TYPE string;