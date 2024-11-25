-- VTEC EVENT --
DEFINE TABLE IF NOT EXISTS vtec_event SCHEMAFULL TYPE NORMAL;

DEFINE FIELD IF NOT EXISTS id ON TABLE vtec_event TYPE 
{ event_number: int, phenomena: string, wfo: string, significance: string, issued: datetime};
DEFINE FIELD IF NOT EXISTS created_at ON TABLE vtec_event TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS updated_at ON TABLE vtec_event TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS start ON TABLE vtec_event TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS expires ON TABLE vtec_event TYPE datetime;
DEFINE FIELD IF NOT EXISTS end ON TABLE vtec_event TYPE datetime;
DEFINE FIELD IF NOT EXISTS end_initial ON TABLE vtec_event TYPE datetime;
DEFINE FIELD IF NOT EXISTS phenomena ON TABLE vtec_event TYPE record<vtec_phenomena>;
DEFINE FIELD IF NOT EXISTS wfo ON TABLE vtec_event TYPE record<office>;
DEFINE FIELD IF NOT EXISTS significance ON TABLE vtec_event TYPE record<vtec_significance>
DEFINE FIELD IF NOT EXISTS title ON TABLE vtec_event TYPE string;

-- VTEC Event History --
DEFINE TABLE IF NOT EXISTS vtec_history SCHEMAFULL TYPE NORMAL;

DEFINE FIELD IF NOT EXISTS id ON TABLE vtec_history TYPE
{ update_number: int, event_number: int, action: string, phenomena: string, wfo: string, significance: string, issued: datetime};
DEFINE FIELD IF NOT EXISTS created_at ON TABLE vtec_history TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS start ON TABLE vtec_history TYPE datetime DEFAULT time::now();
DEFINE FIELD IF NOT EXISTS expires ON TABLE vtec_history TYPE datetime;
DEFINE FIELD IF NOT EXISTS end ON TABLE vtec_history TYPE datetime;
DEFINE FIELD IF NOT EXISTS original ON TABLE vtec_history TYPE string;
DEFINE FIELD IF NOT EXISTS title ON TABLE vtec_history TYPE string;
DEFINE FIELD IF NOT EXISTS action ON TABLE vtec_history TYPE record<vtec_action>;
DEFINE FIELD IF NOT EXISTS phenomena ON TABLE vtec_history TYPE record<vtec_phenomena>;
DEFINE FIELD IF NOT EXISTS wfo ON TABLE vtec_history TYPE record<office>;
DEFINE FIELD IF NOT EXISTS significance ON TABLE vtec_history TYPE record<vtec_significance>
DEFINE FIELD IF NOT EXISTS vtec ON TABLE vtec_history TYPE {
    class: "O" | "T" | "E" | "X",
    action: string,
    wfo: string,
    phenomena: string,
    significance: string,
    event_number: int,
    start: datetime,
    end: datetime
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
    time: datetime,
    original: string
}>;
DEFINE FIELD IF NOT EXISTS product ON TABLE vtec_history TYPE record<product>;