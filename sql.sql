DROP TABLE "files" CASCADE;
DROP TABLE "history";

CREATE TABLE "files" (
	"song_id" SERIAL PRIMARY KEY, --айди песни
	"name" TEXT NOT NULL UNIQUE, --имя
    "url" TEXT NOT NULL UNIQUE, --ссылка на файл
    "active" boolean NOT NULL DEFAULT TRUE --показывать песню в голосованиях
);

CREATE TABLE "history" (
	"voting_id" SERIAL PRIMARY KEY, --айди голосования
	"date" DATE NOT NULL UNIQUE, --дата окончания
	"songs" integer ARRAY NOT NULL, --массив учавствоваших песен
	"winner" integer DEFAULT NULL, --победитель
    "votes" integer ARRAY NOT NULL --голоса
    FOREIGN KEY ("winner") REFERENCES "files" ("song_id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE "votes" (
	"voting_id" integer PRIMARY KEY, --айди голосования
    FOREIGN KEY ("voting_id") REFERENCES "history" ("voting_id") ON DELETE CASCADE ON UPDATE CASCADE
);








INSERT INTO "files" ("name", "url", "active") VALUES
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

INSERT INTO "history" ("date", "songs", "winner") VALUES
('2021-01-03', '{1,2,3,4,5}', '{12345,1233,333332,131313,23131313,131313,313131,3131313,5675675,452231,1231313}', 4),
('2021-01-10', '{1,2,3,4,5}', '{5345345,54445,444,76767}', 1),
('2021-01-17', '{5,4,3}', '{12345,123123,756767,4345345,12313,242423}', 3),
('2021-01-24', '{1,2,7,8}', '{131345,43535345,42424,1231313,745445,23424,675675,224,6756,234234,657657,234234,67567}', 8),
('2021-01-31', '{1,2,5,6,7,8,9,10}', '{12345,768345,348756,32456,234346,1285,34653,256,587856,75678}', 2);


-------------------------


-- SELECT array(SELECT id FROM files);
-- SELECT unnest(array[1, 2, 3]), '2015-01-01', 3;








