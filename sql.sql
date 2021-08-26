DROP TABLE "songs" CASCADE;
DROP TABLE "history" CASCADE;
DROP TABLE "votes" CASCADE;

CREATE TABLE "songs" (
	"song_id" SERIAL PRIMARY KEY, --айди песни
	"name" TEXT NOT NULL UNIQUE, --имя
    "url" TEXT NOT NULL UNIQUE, --ссылка на файл
    "active" boolean NOT NULL DEFAULT TRUE --добавлять песню в новые голосованиях
);

CREATE TABLE "history" (
	"voting_id" SERIAL PRIMARY KEY, --айди голосования
	"date" DATE NOT NULL UNIQUE, --дата окончания
	"winner" integer NOT NULL DEFAULT 0, --айди победителя
    FOREIGN KEY ("winner") REFERENCES "songs" ("song_id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE "votes" (
	"voting_id" integer NOT NULL, --айди голосования
    "song_id" integer NOT NULL, --айди песни в голосовании
    "votes" integer ARRAY, --голоса за песню
    "points" integer NOT NULL, --количество очков за песню
    FOREIGN KEY ("voting_id") REFERENCES "history" ("voting_id") ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY ("song_id") REFERENCES "songs" ("song_id") ON DELETE CASCADE ON UPDATE CASCADE
);


INSERT INTO "songs" VALUES (0,'','',FALSE);
INSERT INTO "songs" ("name", "url", "active") VALUES
    ('Чайковский - Вальс Цветов', 'https://ru.drivemusic.me/dl/XFPds6Z9SPPNskSR-CzlCA/1629406386/download_music/2015/06/chajjkovskijj-vals-cvetov.mp3', TRUE),
    ('Чайковский - Танец Феи Драже', 'https://ru.drivemusic.me/dl/KMEzHd_7SHZBtuxE-cINcQ/1629406387/download_music/2015/06/chajjkovskijj-tanec-fei-drazhe.mp3', TRUE),
    ('Чайковский - Танец маленьких лебедей', 'https://ru.drivemusic.me/dl/tDmiaqLckQlm5sO_Blwxxg/1629406387/download_music/2013/02/chajjkovskijj-tanec-malenkikh-lebedejj.mp3', TRUE),
    ('Чайковский - Лебединое Озеро', 'https://ru.drivemusic.me/dl/IdArgLlyLR_x0ji2Bi9NRw/1629406389/download_music/2015/06/chajjkovskijj-lebedinoe-ozero-scena.mp3', TRUE),
    ('Чайковский - Щелкунчик', 'https://ru.drivemusic.me/dl/XOWoZrgw58xPv-U8BnvbtQ/1629406390/download_music/2015/06/chajjkovskijj-shhelkunchik-marsh.mp3', TRUE),
    ('Чайковский - Времена Года', 'https://ru.drivemusic.me/dl/x5Q_L4ErLKv6jmVOF-sLSA/1629406391/download_music/2015/06/chajjkovskijj-vremena-goda-osennjaja-pesnja.mp3', TRUE),
    ('Чайковский - Allegro Non Troppo', 'https://ru.drivemusic.me/dl/QDWu8wBviRhR4WKkYFiotw/1629406391/download_music/2015/06/chajjkovskijj-allegro-non-troppo.mp3', TRUE),
    ('Чайковский - Barcarole', 'https://ru.drivemusic.me/dl/HI_YkGWLZeBbdyejVVIeaA/1629406392/download_music/2015/06/chajjkovskijj-barcarole.mp3', TRUE),
    ('Чайковский - Lake In Moonlight', 'https://ru.drivemusic.me/dl/qLd6sgceSBQxWchu2_iPlA/1629406393/download_music/2015/06/chajjkovskijj-lake-in-moonlight.mp3', FALSE),
    ('Чайковский - Испанский Танец', 'https://ru.drivemusic.me/dl/Q5KVx37LdMLEe_SPmuAH4A/1629406393/download_music/2015/06/chajjkovskijj-ispanskijj-tanec.mp3', FALSE);

INSERT INTO "history" ("date", "winner") VALUES
    ('2021-01-03', 4),
    ('2021-01-10', 1),
    ('2021-01-17', 3);

INSERT INTO "votes" ("voting_id", "song_id", "votes", "points") VALUES
    ('1', '1', '{1870,1136,1094,1333,1358,1310,1023,1119,1668}', '9'),
    ('1', '2', '{1170,1186,1504,1945,1694}', '5'),
    ('1', '3', '{12345,1040,1421,1008,1623,1269,1629}', '17'),
    ('1', '4', '{1920,1485,1705,1264,1787,1438,1246,1270,1784,1874}', '10'),
    ('1', '5', '{1323,1402,1288,1996,1752,1670,1971}', '16'),
    ('2', '1', '{1831,1298,1692,1335,1788,1153,1488,1519,1770,1595}', '10'),
    ('2', '7', '{1969,1765,1421,1648,1997,1400,1103,1854,1182}', '9'),
    ('2', '8', '{1317,1341,1421,1665,1179,1730,1930,1117}', '8'),
    ('2', '9', '{1311,1607,1942,1232}', '4'),
    ('2', '10','{1778}', '10'),
    ('3', '1', '{1189,1534,1949,1770,1718,1749}', '6'),
    ('3', '5', '{1656,1849,1834,1449,1141,1235,1147,1569,1705,1529}', '10'),
    ('3', '3', '{12345,1893,1501,1889,1694,1521,1317,1019,1302,1158}', '19'),
    ('3', '2', '{1632,1325,1510,1286,1220,1503,1332,1071,1029}', '9'),
    ('3', '4', '{1753,1793,1414}', '14');





-------------------------


-- SELECT array(SELECT id FROM songs);
-- SELECT unnest(array[1, 2, 3]), '2015-01-01', 3;








