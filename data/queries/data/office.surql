OPTION IMPORT;

-- ------------------------------
-- TRANSACTION
-- ------------------------------

BEGIN TRANSACTION;

-- ------------------------------
-- TABLE DATA: office
-- ------------------------------

UPSERT office:ABQ CONTENT { code: 'ABQ', icao: 'KABQ', id: office:ABQ, location: (-106.016611696, 35.1125387845), name: 'Albuquerque', state: state:NM };
UPSERT office:ABR CONTENT { code: 'ABR', icao: 'KABR', id: office:ABR, location: (-99.1578687708, 45.0559400154), name: 'Aberdeen', state: state:SD };
UPSERT office:AFG CONTENT { code: 'AFG', icao: 'PAFG', id: office:AFG, location: (-153.4724743932, 66.38456499521), name: 'Fairbanks', state: state:AK };
UPSERT office:ALY CONTENT { code: 'ALY', icao: 'KALY', id: office:ALY, location: (-73.9211191031, 42.8140232043), name: 'Albany', state: state:NY };
UPSERT office:AMA CONTENT { code: 'AMA', icao: 'KAMA', id: office:AMA, location: (-101.5164304, 35.869648463), name: 'Amarillo', state: state:TX };
UPSERT office:APX CONTENT { code: 'APX', icao: 'KAPX', id: office:APX, location: (-84.6640103835, 45.012398868), name: 'Gaylord', state: state:MI };
UPSERT office:ARX CONTENT { code: 'ARX', icao: 'KARX', id: office:ARX, location: (-91.3508377932, 43.7775964754), name: 'La Crosse', state: state:WI };
UPSERT office:BGM CONTENT { code: 'BGM', icao: 'KBGM', id: office:BGM, location: (-75.8994230946, 42.2331900531), name: 'Binghamton', state: state:NY };
UPSERT office:BIS CONTENT { code: 'BIS', icao: 'KBIS', id: office:BIS, location: (-101.417583385, 47.3646553978), name: 'Bismarck', state: state:ND };
UPSERT office:BMX CONTENT { code: 'BMX', icao: 'KBMX', id: office:BMX, location: (-86.6758909443, 33.0269285434), name: 'Birmingham', state: state:AL };
UPSERT office:BOI CONTENT { code: 'BOI', icao: 'KBOI', id: office:BOI, location: (-116.93582353, 43.4103077552), name: 'Boise', state: state:ID };
UPSERT office:BOU CONTENT { code: 'BOU', icao: 'KBOU', id: office:BOU, location: (-104.538294118, 40.0123848727), name: 'Denver', state: state:CO };
UPSERT office:BOX CONTENT { code: 'BOX', icao: 'KBOX', id: office:BOX, location: (-71.7528617974, 42.1104014769), name: 'Norton', state: state:MA };
UPSERT office:BRO CONTENT { code: 'BRO', icao: 'KBRO', id: office:BRO, location: (-98.2654323918, 26.7021935443), name: 'Brownsville', state: state:TX };
UPSERT office:BTV CONTENT { code: 'BTV', icao: 'KBTV', id: office:BTV, location: (-73.4550903218, 44.3571291119), name: 'Burlington', state: state:VT };
UPSERT office:BUF CONTENT { code: 'BUF', icao: 'KBUF', id: office:BUF, location: (-77.5485364015, 42.999565196), name: 'Buffalo', state: state:NY };
UPSERT office:BYZ CONTENT { code: 'BYZ', icao: 'KBYZ', id: office:BYZ, location: (-107.367814548, 45.7909260146), name: 'Billings', state: state:MT };
UPSERT office:CAE CONTENT { code: 'CAE', icao: 'KCAE', id: office:CAE, location: (-81.1819345186, 33.8479583811), name: 'Columbia', state: state:SC };
UPSERT office:CAR CONTENT { code: 'CAR', icao: 'KCAR', id: office:CAR, location: (-68.6686131947, 45.8574115262), name: 'Caribou', state: state:ME };
UPSERT office:CHS CONTENT { code: 'CHS', icao: 'KCHS', id: office:CHS, location: (-81.0307561962, 32.5323319035), name: 'Charleston', state: state:SC };
UPSERT office:CLE CONTENT { code: 'CLE', icao: 'KCLE', id: office:CLE, location: (-81.9404833155, 41.1867098763), name: 'Cleveland', state: state:OH };
UPSERT office:CRP CONTENT { code: 'CRP', icao: 'KCRP', id: office:CRP, location: (-98.2275211475, 28.0715945832), name: 'Corpus Christi', state: state:TX };
UPSERT office:CTP CONTENT { code: 'CTP', icao: 'KCTP', id: office:CTP, location: (-77.6739500657, 40.8077376001), name: 'State College', state: state:PA };
UPSERT office:CYS CONTENT { code: 'CYS', icao: 'KCYS', id: office:CYS, location: (-105.0343672, 42.0421926748), name: 'Cheyenne', state: state:WY };
UPSERT office:DDC CONTENT { code: 'DDC', icao: 'KDDC', id: office:DDC, location: (-100.162302332, 37.8615170994), name: 'Dodge City', state: state:KS };
UPSERT office:DLH CONTENT { code: 'DLH', icao: 'KDLH', id: office:DLH, location: (-92.4575991436, 47.0511957473), name: 'Duluth', state: state:MN };
UPSERT office:DMX CONTENT { code: 'DMX', icao: 'KDMX', id: office:DMX, location: (-93.7409295336, 41.9757484115), name: 'Des Moines', state: state:IA };
UPSERT office:DTX CONTENT { code: 'DTX', icao: 'KDTX', id: office:DTX, location: (-83.5116679433, 42.9460461437), name: 'Detroit/Pontiac', state: state:MI };
UPSERT office:DVN CONTENT { code: 'DVN', icao: 'KDVN', id: office:DVN, location: (-91.0153363019, 41.5122136108), name: 'Quad Cities', state: state:IL };
UPSERT office:EAX CONTENT { code: 'EAX', icao: 'KEAX', id: office:EAX, location: (-93.9301897675, 39.489005053), name: 'Kansas City/Pleasant Hill', state: state:MO };
UPSERT office:EKA CONTENT { code: 'EKA', icao: 'KEKA', id: office:EKA, location: (-123.433348578, 40.2614966956), name: 'Eureka', state: state:CA };
UPSERT office:EPZ CONTENT { code: 'EPZ', icao: 'KEPZ', id: office:EPZ, location: (-106.913536781, 32.3506149352), name: 'El Paso Tx/Santa Teresa', state: state:NM };
UPSERT office:EWX CONTENT { code: 'EWX', icao: 'KEWX', id: office:EWX, location: (-98.9440801587, 29.650191409600005), name: 'Austin/San Antonio', state: state:TX };
UPSERT office:FFC CONTENT { code: 'FFC', icao: 'KFFC', id: office:FFC, location: (-83.9246225026, 33.3266385387), name: 'Peachtree City', state: state:GA };
UPSERT office:FGF CONTENT { code: 'FGF', icao: 'KFGF', id: office:FGF, location: (-96.8022847193, 47.6559391988), name: 'Grand Forks', state: state:ND };
UPSERT office:FGZ CONTENT { code: 'FGZ', icao: 'KFGZ', id: office:FGZ, location: (-111.060570581, 35.3670427976), name: 'Flagstaff', state: state:AZ };
UPSERT office:FSD CONTENT { code: 'FSD', icao: 'KFSD', id: office:FSD, location: (-96.9500288663, 43.5162792024), name: 'Sioux Falls', state: state:SD };
UPSERT office:FWD CONTENT { code: 'FWD', icao: 'KFWD', id: office:FWD, location: (-97.1392386202, 32.3664908102), name: 'Fort Worth', state: state:TX };
UPSERT office:GGW CONTENT { code: 'GGW', icao: 'KGGW', id: office:GGW, location: (-106.23605129, 47.857447188), name: 'Glasgow', state: state:MT };
UPSERT office:GID CONTENT { code: 'GID', icao: 'KGID', id: office:GID, location: (-98.7020370769, 40.4305168423), name: 'Hastings', state: state:NE };
UPSERT office:GJT CONTENT { code: 'GJT', icao: 'KGJT', id: office:GJT, location: (-108.387497494, 38.9821003933), name: 'Grand Junction', state: state:CO };
UPSERT office:GLD CONTENT { code: 'GLD', icao: 'KGLD', id: office:GLD, location: (-101.46830328, 39.4262964096), name: 'Goodland', state: state:KS };
UPSERT office:GRB CONTENT { code: 'GRB', icao: 'KGRB', id: office:GRB, location: (-88.8922230206, 44.9880692719), name: 'Green Bay', state: state:WI };
UPSERT office:GRR CONTENT { code: 'GRR', icao: 'KGRR', id: office:GRR, location: (-85.3513970145, 43.0913071789), name: 'Grand Rapids', state: state:MI };
UPSERT office:GSP CONTENT { code: 'GSP', icao: 'KGSP', id: office:GSP, location: (-82.0799887399, 35.1556677834), name: 'Greenville-Spartanburg', state: state:SC };
UPSERT office:GUM CONTENT { code: 'GUM', icao: 'PGUM', id: office:GUM, location: (145.158812026, 14.510822985300003), name: 'Tiyan', state: state:GU };
UPSERT office:GYX CONTENT { code: 'GYX', icao: 'KGYX', id: office:GYX, location: (-70.7651431493, 44.2035123412), name: 'Gray', state: state:ME };
UPSERT office:HFO CONTENT { code: 'HFO', icao: 'PHFO', id: office:HFO, location: (-156.347431421, 20.2492312725), name: 'Honolulu', state: state:HI };
UPSERT office:HGX CONTENT { code: 'HGX', icao: 'KHGX', id: office:HGX, location: (-95.6747570264, 30.0460458533), name: 'Houston/Galveston', state: state:TX };
UPSERT office:HNX CONTENT { code: 'HNX', icao: 'KHNX', id: office:HNX, location: (-119.327932783, 36.3440025462), name: 'Hanford', state: state:CA };
UPSERT office:HUN CONTENT { code: 'HUN', icao: 'KHUN', id: office:HUN, location: (-86.7676101113, 34.6661265472), name: 'Huntsville', state: state:AL };
UPSERT office:ICT CONTENT { code: 'ICT', icao: 'KICT', id: office:ICT, location: (-97.1760285131, 37.9143149228), name: 'Wichita', state: state:KS };
UPSERT office:ILM CONTENT { code: 'ILM', icao: 'KILM', id: office:ILM, location: (-79.0207894179, 34.1675975871), name: 'Wilmington', state: state:NC };
UPSERT office:ILN CONTENT { code: 'ILN', icao: 'KILN', id: office:ILN, location: (-83.8855711516, 39.5419976043), name: 'Wilmington', state: state:OH };
UPSERT office:ILX CONTENT { code: 'ILX', icao: 'KILX', id: office:ILX, location: (-89.0048723254, 39.944180538), name: 'Lincoln', state: state:IL };
UPSERT office:IND CONTENT { code: 'IND', icao: 'KIND', id: office:IND, location: (-86.4451128122, 39.6212899731), name: 'Indianapolis', state: state:IN };
UPSERT office:IWX CONTENT { code: 'IWX', icao: 'KIWX', id: office:IWX, location: (-85.4551485287, 41.2653426932), name: 'Northern Indiana', state: state:IN };
UPSERT office:JAN CONTENT { code: 'JAN', icao: 'KJAN', id: office:JAN, location: (-90.1995519875, 32.5381934483), name: 'Jackson', state: state:MS };
UPSERT office:JAX CONTENT { code: 'JAX', icao: 'KJAX', id: office:JAX, location: (-82.2078258322, 30.462115851600004), name: 'Jacksonville', state: state:FL };
UPSERT office:JKL CONTENT { code: 'JKL', icao: 'KJKL', id: office:JKL, location: (-83.5749843307, 37.3902895602), name: 'Jackson', state: state:KY };
UPSERT office:KEY CONTENT { code: 'KEY', icao: 'KKEY', id: office:KEY, location: (-81.0866832409, 24.8422289042), name: 'Key West', state: state:FL };
UPSERT office:LBF CONTENT { code: 'LBF', icao: 'KLBF', id: office:LBF, location: (-100.712043682, 41.8755271684), name: 'North Platte', state: state:NE };
UPSERT office:LCH CONTENT { code: 'LCH', icao: 'KLCH', id: office:LCH, location: (-92.9684126126, 30.4523804621), name: 'Lake Charles', state: state:LA };
UPSERT office:LKN CONTENT { code: 'LKN', icao: 'KLKN', id: office:LKN, location: (-116.208390723, 40.2484934479), name: 'Elko', state: state:NV };
UPSERT office:LMK CONTENT { code: 'LMK', icao: 'KLMK', id: office:LMK, location: (-85.6681042847, 37.7254114055), name: 'Louisville', state: state:KY };
UPSERT office:LOT CONTENT { code: 'LOT', icao: 'KLOT', id: office:LOT, location: (-88.2582240597, 41.4841083442), name: 'Chicago', state: state:IL };
UPSERT office:LOX CONTENT { code: 'LOX', icao: 'KLOX', id: office:LOX, location: (-119.372865381, 34.7187246855), name: 'Los Angeles/Oxnard', state: state:CA };
UPSERT office:LSX CONTENT { code: 'LSX', icao: 'KLSX', id: office:LSX, location: (-90.8305486493, 38.8166723605), name: 'St Louis', state: state:MO };
UPSERT office:LUB CONTENT { code: 'LUB', icao: 'KLUB', id: office:LUB, location: (-101.528683937, 33.8520368548), name: 'Lubbock', state: state:TX };
UPSERT office:LWX CONTENT { code: 'LWX', icao: 'KLWX', id: office:LWX, location: (-77.997722212, 38.8774555788), name: 'Baltimore MD/Washington', state: state:DC };
UPSERT office:LZK CONTENT { code: 'LZK', icao: 'KLZK', id: office:LZK, location: (-92.4254667234, 34.9573071673), name: 'Little Rock', state: state:AR };
UPSERT office:MAF CONTENT { code: 'MAF', icao: 'KMAF', id: office:MAF, location: (-103.101250695, 31.402868353799995), name: 'Midland/Odessa', state: state:TX };
UPSERT office:MEG CONTENT { code: 'MEG', icao: 'KMEG', id: office:MEG, location: (-89.4991175198, 35.162902217), name: 'Memphis', state: state:TN };
UPSERT office:MFL CONTENT { code: 'MFL', icao: 'KMFL', id: office:MFL, location: (-80.8496653484, 26.1967539601), name: 'Miami', state: state:FL };
UPSERT office:MFR CONTENT { code: 'MFR', icao: 'KMFR', id: office:MFR, location: (-122.020351801, 42.478496575), name: 'Medford', state: state:OR };
UPSERT office:MKX CONTENT { code: 'MKX', icao: 'KMKX', id: office:MKX, location: (-89.0291764316, 43.1760047836), name: 'Milwaukee/Sullivan', state: state:WI };
UPSERT office:MLB CONTENT { code: 'MLB', icao: 'KMLB', id: office:MLB, location: (-81.0595919602, 28.2150343643), name: 'Melbourne', state: state:FL };
UPSERT office:MOB CONTENT { code: 'MOB', icao: 'KMOB', id: office:MOB, location: (-87.6230789026, 31.2617708012), name: 'Mobile', state: state:AL };
UPSERT office:MPX CONTENT { code: 'MPX', icao: 'KMPX', id: office:MPX, location: (-93.8279129996, 44.9941981513), name: 'Twin Cities/Chanhassen', state: state:MN };
UPSERT office:MQT CONTENT { code: 'MQT', icao: 'KMQT', id: office:MQT, location: (-87.812806266, 46.3821282321), name: 'Marquette', state: state:MI };
UPSERT office:MRX CONTENT { code: 'MRX', icao: 'KMRX', id: office:MRX, location: (-83.7137292339, 36.0098107779), name: 'Morristown', state: state:TN };
UPSERT office:MSO CONTENT { code: 'MSO', icao: 'KMSO', id: office:MSO, location: (-114.522911024, 46.7929701534), name: 'Missoula', state: state:MT };
UPSERT office:MTR CONTENT { code: 'MTR', icao: 'KMTR', id: office:MTR, location: (-121.81958439, 37.2528822756), name: 'San Francisco', state: state:CA };
UPSERT office:OAX CONTENT { code: 'OAX', icao: 'KOAX', id: office:OAX, location: (-96.6469058038, 41.3771817325), name: 'Omaha/Valley', state: state:NE };
UPSERT office:OHX CONTENT { code: 'OHX', icao: 'KOHX', id: office:OHX, location: (-86.565647905, 35.9382613694), name: 'Nashville', state: state:TN };
UPSERT office:OKX CONTENT { code: 'OKX', icao: 'KOKX', id: office:OKX, location: (-73.3881080473, 41.1774564599), name: 'Upton', state: state:NY };
UPSERT office:OTX CONTENT { code: 'OTX', icao: 'KOTX', id: office:OTX, location: (-118.183426918, 47.7432407919), name: 'Spokane', state: state:WA };
UPSERT office:OUN CONTENT { code: 'OUN', icao: 'KOUN', id: office:OUN, location: (-98.2194096537, 35.1995667465), name: 'Norman', state: state:OK };
UPSERT office:PAH CONTENT { code: 'PAH', icao: 'KPAH', id: office:PAH, location: (-88.7169444426, 37.412209524), name: 'Paducah', state: state:KY };
UPSERT office:PBZ CONTENT { code: 'PBZ', icao: 'KPBZ', id: office:PBZ, location: (-80.2596909574, 40.3936579819), name: 'Pittsburgh', state: state:PA };
UPSERT office:PDT CONTENT { code: 'PDT', icao: 'KPDT', id: office:PDT, location: (-119.72734266100001, 45.4082541141), name: 'Pendleton', state: state:OR };
UPSERT office:PHI CONTENT { code: 'PHI', icao: 'KPHI', id: office:PHI, location: (-75.1959242024, 39.9817188158), name: 'Mount Holly', state: state:NJ };
UPSERT office:PIH CONTENT { code: 'PIH', icao: 'KPIH', id: office:PIH, location: (-112.868432894, 43.3522622112), name: 'Pocatello', state: state:ID };
UPSERT office:PPG CONTENT { code: 'PPG', icao: 'NSTU', id: office:PPG, location: (-170.372204502, -14.264897808), name: 'Pago Pago', state: state:AS };
UPSERT office:PQE CONTENT { code: 'PQE', icao: 'PPQE', id: office:PQE, location: (161.459811475, 6.70681047545), name: 'Micronesia Domain East', state: state:PQ };
UPSERT office:PQR CONTENT { code: 'PQR', icao: 'KPQR', id: office:PQR, location: (-122.847432544, 45.0966532396), name: 'Portland', state: state:OR };
UPSERT office:PQW CONTENT { code: 'PQW', icao: 'PPQW', id: office:PQW, location: (138.080232588, 7.74523961438), name: 'Micronesia Domain West', state: state:PQ };
UPSERT office:PSR CONTENT { code: 'PSR', icao: 'KPSR', id: office:PSR, location: (-113.515617572, 33.3110899834), name: 'Phoenix', state: state:AZ };
UPSERT office:PUB CONTENT { code: 'PUB', icao: 'KPUB', id: office:PUB, location: (-104.559332039, 37.9223077543), name: 'Pueblo', state: state:CO };
UPSERT office:RAH CONTENT { code: 'RAH', icao: 'KRAH', id: office:RAH, location: (-78.9605476302, 35.6623665311), name: 'Raleigh', state: state:NC };
UPSERT office:REV CONTENT { code: 'REV', icao: 'KREV', id: office:REV, location: (-119.20944460400001, 39.7992447696), name: 'Reno', state: state:NV };
UPSERT office:RIW CONTENT { code: 'RIW', icao: 'KRIW', id: office:RIW, location: (-108.791194398, 43.1754373327), name: 'Riverton', state: state:WY };
UPSERT office:RLX CONTENT { code: 'RLX', icao: 'KRLX', id: office:RLX, location: (-81.4497271215, 38.5741087143), name: 'Charleston', state: state:WV };
UPSERT office:RNK CONTENT { code: 'RNK', icao: 'KRNK', id: office:RNK, location: (-80.1517015339, 37.067549617), name: 'Blacksburg', state: state:VA };
UPSERT office:SEW CONTENT { code: 'SEW', icao: 'KSEW', id: office:SEW, location: (-122.53979531, 47.6667294474), name: 'Seattle', state: state:WA };
UPSERT office:SGF CONTENT { code: 'SGF', icao: 'KSGF', id: office:SGF, location: (-93.0520111076, 37.4040182415), name: 'Springfield', state: state:MO };
UPSERT office:SGX CONTENT { code: 'SGX', icao: 'KSGX', id: office:SGX, location: (-116.929645155, 33.6871449946), name: 'San Diego', state: state:CA };
UPSERT office:SHV CONTENT { code: 'SHV', icao: 'KSHV', id: office:SHV, location: (-93.8745803088, 32.6254771794), name: 'Shreveport', state: state:LA };
UPSERT office:SJT CONTENT { code: 'SJT', icao: 'KSJT', id: office:SJT, location: (-100.07201423, 31.5840987903), name: 'San Angelo', state: state:TX };
UPSERT office:SJU CONTENT { code: 'SJU', icao: 'TJSJ', id: office:SJU, location: (-66.4065281471, 18.2127063237), name: 'San Juan', state: state:PR };
UPSERT office:SLC CONTENT { code: 'SLC', icao: 'KSLC', id: office:SLC, location: (-112.104704031, 39.4941478912), name: 'Salt Lake City', state: state:UT };
UPSERT office:STO CONTENT { code: 'STO', icao: 'KSTO', id: office:STO, location: (-121.436252029, 39.2400356582), name: 'Sacramento', state: state:CA };
UPSERT office:TAE CONTENT { code: 'TAE', icao: 'KTAE', id: office:TAE, location: (-84.5415963743, 30.834201186900003), name: 'Tallahassee', state: state:FL };
UPSERT office:TFX CONTENT { code: 'TFX', icao: 'KTFX', id: office:TFX, location: (-111.223817297, 47.1820554746), name: 'Great Falls', state: state:MT };
UPSERT office:TOP CONTENT { code: 'TOP', icao: 'KTOP', id: office:TOP, location: (-96.3939234895, 39.1496112526), name: 'Topeka', state: state:KS };
UPSERT office:TSA CONTENT { code: 'TSA', icao: 'KTSA', id: office:TSA, location: (-95.2473348521, 35.7542046697), name: 'Tulsa', state: state:OK };
UPSERT office:TWC CONTENT { code: 'TWC', icao: 'KTWC', id: office:TWC, location: (-110.669143573, 32.3278995352), name: 'Tucson', state: state:AZ };
UPSERT office:UNR CONTENT { code: 'UNR', icao: 'KUNR', id: office:UNR, location: (-103.021544961, 44.2578770656), name: 'Rapid City', state: state:SD };
UPSERT office:VEF CONTENT { code: 'VEF', icao: 'KVEF', id: office:VEF, location: (-115.620760185, 36.2869821386), name: 'Las Vegas', state: state:NV };

-- ------------------------------
-- TRANSACTION
-- ------------------------------

COMMIT TRANSACTION;