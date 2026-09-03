-- 2_test-data.sql

-- 1. Equipos
INSERT INTO "equipos" ("id", "nombre", "formacion_actual", "escudo_url", "valoracion", "posicion_tabla") 
VALUES
(1, 'Real Madrid', '4-3-3', 'https://via.placeholder.com/100', 88, 1),
(2, 'FC Barcelona', '4-3-3', 'https://via.placeholder.com/100', 86, 2)
ON CONFLICT DO NOTHING;

-- Resetear secuencia de IDs de equipos
SELECT setval(pg_get_serial_sequence('equipos', 'id'), COALESCE((SELECT MAX(id) FROM equipos), 1));

-- 2. Jugadores
INSERT INTO "jugadores" (
    "id", "equipo_id", "nombre", "posicion", "fecha_nacimiento", 
    "media_general", "altura", "ritmo", "tiro", "pase", 
    "regate", "defensa", "fisico", "goles", "asistencias", "foto_url"
) VALUES
(1, 1, 'Kylian Mbappé', 'DEL', '1998-12-20', 91, 1.78, 97, 90, 80, 92, 36, 78, 12, 4, 'https://via.placeholder.com/100'),
(2, 2, 'Lamine Yamal', 'DEL', '2007-07-13', 81, 1.80, 88, 75, 78, 82, 23, 50, 5, 7, 'https://via.placeholder.com/100')
ON CONFLICT DO NOTHING;

-- Resetear secuencia de IDs de jugadores
SELECT setval(pg_get_serial_sequence('jugadores', 'id'), COALESCE((SELECT MAX(id) FROM jugadores), 1));

-- 3. Partidos
INSERT INTO "partidos" ("id", "equipo_local_id", "equipo_visitante_id", "goles_local", "goles_visitante", "fecha", "estado") 
VALUES
(1, 1, 2, 2, 1, '2026-03-01', 'FINALIZADO')
ON CONFLICT DO NOTHING;

-- Resetear secuencia de IDs de partidos
SELECT setval(pg_get_serial_sequence('partidos', 'id'), COALESCE((SELECT MAX(id) FROM partidos), 1));