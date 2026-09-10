INSERT INTO "equipos" ("id", "nombre", "formacion_actual", "escudo_url", "valoracion", "posicion_tabla")
OVERRIDING SYSTEM VALUE VALUES
(1, 'Real Madrid', '4-3-3', 'https://example.com/real.png', 88, 1),
(2, 'FC Barcelona', '4-3-3', 'https://example.com/barca.png', 86, 2)
ON CONFLICT (id) DO NOTHING;

INSERT INTO "jugadores" ("id", "equipo_id", "nombre", "posicion", "fecha_nacimiento", "media_general", "altura", "ritmo", "tiro", "pase", "regate", "defensa", "fisico", "goles", "asistencias")
OVERRIDING SYSTEM VALUE VALUES
(1, 1, 'Vinicius Jr', 'DEL', '2000-07-12', 89, 1.76, 95, 82, 78, 90, 29, 68, 12, 7),
(2, 2, 'Pedri', 'MED', '2002-11-25', 86, 1.74, 78, 68, 85, 88, 68, 70, 4, 9)
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