-- Creación de la base de datos (Opcional)
-- CREATE DATABASE gestion_futbol;
-- \c gestion_futbol;

-- Eliminación previa si existe para evitar conflictos en scripts de inicialización
DROP TABLE IF EXISTS equipos CASCADE;

-- Definición de la entidad principal: Equipos
CREATE TABLE equipos (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE,
    formacion_actual VARCHAR(10) NOT NULL DEFAULT '4-3-3',
    escudo_url VARCHAR(255),
    
    CONSTRAINT chk_nombre_no_vacio CHECK (LENGTH(TRIM(nombre)) > 0)
);

-- Comentarios explicativos en las columnas
COMMENT ON TABLE equipos IS 'Tabla principal que almacena los equipos del sistema';
COMMENT ON COLUMN equipos.id IS 'Identificador único autoincremental del equipo';
COMMENT ON COLUMN equipos.nombre IS 'Nombre único del club o entidad deportiva';
COMMENT ON COLUMN equipos.formacion_actual IS 'Esquema táctico asignado (ej. 4-3-3, 4-4-2)';