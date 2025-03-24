INSERT INTO phenomena (id, name) VALUES
('AF', 'Ashfall (land)'),
('AS', 'Air Stagnation'),
('BH', 'Beach Hazard'),
('BW', 'Brisk Wind'),
('BZ', 'Blizzard'),
('CF', 'Coastal Flood'),
('CW', 'Cold Weather'),
('DF', 'Debris Flow'),
('DS', 'Dust Storm'),
('DU', 'Blowing Dust'),
('EC', 'Extreme Cold'),
('EH', 'Excessive Heat'),
('EW', 'Extreme Wind'),
('FA', 'Flood'),
('FF', 'Flash Flood'),
('FG', 'Dense Fog (land)'),
('FL', 'Flood (forecast point)'),
('FR', 'Frost'),
('FW', 'Fire Weather'),
('FZ', 'Freeze'),
('GL', 'Gale'),
('HF', 'Hurricane Force Wind'),
('HT', 'Heat'),
('HU', 'Hurricane'),
('HW', 'High Wind'),
('HY', 'Hydrologic'),
('HZ', 'Hard Freeze'),
('IS', 'Ice Storm'),
('LE', 'Lake Effect Snow'),
('LO', 'Low Water'),
('LS', 'Lakeshore Flood'),
('LW', 'Lake Wind'),
('MA', 'Marine'),
('MF', 'Dense Fog (marine)'),
('MH', 'Ashfall (marine)'),
('MS', 'Dense Smoke (marine)'),
('RB', 'Small Craft for Rough Bar'),
('RP', 'Rip Current Risk'),
('SC', 'Small Craft'),
('SE', 'Hazardous Seas'),
('SI', 'Small Craft for Winds'),
('SM', 'Dense Smoke (land)'),
('SQ', 'Snow Squall'),
('SR', 'Storm'),
('SS', 'Storm Surge'),
('SU', 'High Surf'),
('SV', 'Severe Thunderstorm'),
('SW', 'Small Craft for Hazardous Seas'),
('TO', 'Tornado'),
('TR', 'Tropical Storm'),
('TS', 'Tsunami'),
('TY', 'Typhoon'),
('UP', '(Heavy) Freezing Spray'),
('WC', 'Wind Chill'),
('WI', 'Wind'),
('WS', 'Winter Storm'),
('WW', 'Winter Weather'),
('XH', 'Extreme Heat'),
('ZF', 'Freezing Fog'),
('ZR', 'Freezing Rain');


INSERT INTO significance (id, name) VALUES
('A', 'Watch'),
('S', 'Statement'),
('W', 'Warning'),
('Y', 'Advisory');


INSERT INTO action (id, name) VALUES
('CAN', 'Event Cancelled'),
('CON', 'Event Continued'),
('COR', 'Correction'),
('EXA', 'Event Extended (Area)'),
('EXB', 'Event Area Extended and Time Changed'),
('EXP', 'Event Expired'),
('EXT', 'Event Extended (Time)'),
('NEW', 'New Event'),
('ROU', 'Routine'),
('UPG', 'Event Upgraded');
