INSERT INTO "equipos" ("id", "nombre", "formacion_actual", "escudo_url", "valoracion", "posicion_tabla")
OVERRIDING SYSTEM VALUE VALUES
(1, 'Real Madrid', '4-3-3', 'https://imgs.search.brave.com/dwqyu2PcGJpnnqt0I5oa_ypMh74XwzjJW8lxxauKNHo/rs:fit:860:0:0:0/g:ce/aHR0cHM6Ly9hc3Nl/dHMuZm9vdHlsb2dv/cy5jb20vbG9nb3Mv/cmVhbC1tYWRyaWQv/cmVhbC1tYWRyaWQt/bG9nby1mb290eWxv/Z29zLnN2Zw', 88, 1),
(2, 'FC Barcelona', '4-3-3', 'https://imgs.search.brave.com/q8dZl_cSX5SD6bkSiEs7JjWjtpmmDGRr369DnKUjJoA/rs:fit:860:0:0:0/g:ce/aHR0cHM6Ly9hc3Nl/dHMuZm9vdHlsb2dv/cy5jb20vbG9nb3Mv/ZmMtYmFyY2Vsb25h/L2ZjLWJhcmNlbG9u/YS1sb2dvLWZvb3R5/bG9nb3Muc3Zn', 86, 2)
ON CONFLICT (id) DO NOTHING;

INSERT INTO "jugadores" ("id", "equipo_id", "nombre", "posicion", "fecha_nacimiento", "media_general", "altura", "ritmo", "tiro", "pase", "regate", "defensa", "fisico", "goles", "asistencias", "foto_url")
OVERRIDING SYSTEM VALUE VALUES
(1, 1, 'Vinicius Jr', 'DEL', '2000-07-12', 89, '1.76', 95, 82, 78, 90, 29, 68, 12, 7, 'https://imgs.search.brave.com/v6g0WREoDji_rzapqMoymMoa8qmQqMtEP4ziKTuEDSg/rs:fit:860:0:0:0/g:ce/aHR0cHM6Ly9tZWRp/YS5nZXR0eWltYWdl/cy5jb20vaWQvMjI4/MDc4Mzg5Ni9waG90/by92aW5pY2l1cy1q/dW5pb3ItZ2VzdHVy/ZXMtZHVyaW5nLWdy/b3VwLWMtZmlmYS13/b3JsZC1jdXAtMjAy/Ni1icmF6aWwtYW5k/LW1vcm9jY28tbWV0/LWxpZmUuanBnP3M9/NjEyeDYxMiZ3PTAm/az0yMCZjPV81VUpo/aXFxRHZJMnIyelZz/OUZXSkJOWjJVTlBv/bUlCeTYzSWkwdUJF/YlU9'),
(2, 2, 'Pedri', 'MED', '2002-11-25', 86, '1.74', 78, 68, 85, 88, 68, 70, 4, 9, 'https://imgs.search.brave.com/fRKO4u5_8P5DCvq_zJwObhfhRsP_dzP0G8oni1RUXK0/rs:fit:860:0:0:0/g:ce/aHR0cHM6Ly9tZWRp/YS5nZXR0eWltYWdl/cy5jb20vaWQvMjIx/ODUxMzQzMS9waG90/by9tYWRyaWQtc3Bh/aW4tcGVkcmktb2Yt/c3BhaW4tcG9zZXMt/Zm9yLWEtcG9ydHJh/aXQtYWhlYWQtb2Yt/dGhlLXVlZmEtbmF0/aW9ucy1sZWFndWUt/MjAyNS1maW5hbHMu/anBnP3M9NjEyeDYx/MiZ3PTAmaz0yMCZj/PU9RVzkxUVpvZ2tC/Y1k4a2xVRVVFMXJT/Q3I4Qjdfb1F6U0JV/RkVFOERoOEE9')
ON CONFLICT (id) DO NOTHING;

UPDATE "equipos" SET "mejor_jugador" = 1 WHERE "id" = 1;
UPDATE "equipos" SET "mejor_jugador" = 2 WHERE "id" = 2;

INSERT INTO "partidos" ("id", "equipo_local_id", "equipo_visitante_id", "goles_local", "goles_visitante", "fecha", "estado")
OVERRIDING SYSTEM VALUE VALUES
(1, 1, 2, 2, 1, '2026-09-10', 'FINALIZADO'),
(2, 2, 1, 0, 0, '2026-09-20', 'PENDIENTE')
ON CONFLICT (id) DO NOTHING;

SELECT setval(pg_get_serial_sequence('equipos', 'id'), coalesce(max(id), 1)) FROM "equipos";
SELECT setval(pg_get_serial_sequence('jugadores', 'id'), coalesce(max(id), 1)) FROM "jugadores";
SELECT setval(pg_get_serial_sequence('partidos', 'id'), coalesce(max(id), 1)) FROM "partidos";