PARA EJECUTAR EL SERVIDOR, UTILIZAR:
-clonar el repositorio
-cd Tp-ProgWeb
-ejecutar el comando: docker compose up -d --build (en este paso, pueden comentar esta linea en docker-compose.yml en caso de que no quieran cargar los datos de prueba de equipos:"./db/test-data.sql:/docker-entrypoint-initdb.d/2_test-data.sql")
-luego acceder desde el navegador a: http://localhost:8080

DETALLES DE ESTA SEGUNDA PARTE:
-Se incorporó docker, tanto un contenedor para la aplicacion como un contenedor para la base de datos. Se pueden ver detalles acerca de como están administrados en docker-compose.yml
-Se incorporó postgres
-Se agrego la entidad básica "equipos", junto con los datos de prueba en test-data.sql y se los muestra en una tabla muy sencilla que hace petición a la api: "/api/equipos"
-Se estructuró el proyecto en diferentes secciones: 
    +Por un lado esta la carpeta "db" donde se guarda todo lo relacionado a la base de datos
    +Por otro lado, se agregó la carpeta "internal" que contiene datos propios del funcionamiento interno del servidor (como los handlers de la API y el servidor de paginas HTML según la petición realizada) y los datos generados por "sqlc" para la peticiones 
-Se modifico la direccion del archivo "main.go" que inicializa el servidor a la carpeta "server", además de que se mantuvo su estructura simple para que solo se encargue de iniciar la base de datos, el servidor y delegue el manejo de rutas a los Handlers

PARA FUTURAS EDICIONES:
-Agregaremos el resto de entidades de la aplicación, junto con sus peticiones requeridas 
-Si se nos solicita agregaremos .env y .gitignore para contraseñas y archivos que no sean necesarios 

DETALLES DE NUESTRA PÁGINA:
Nuestra página consta de un tablero que enseña equipos cargados por el administrador, donde los usuarios que la utilicen pueden ver información acerca de los mismos. Por ejemplo, plantel, valoración media del equipo, mejores jugadores, estadísticas individuales de cada jugador, últimos resultados de ese equipo, posición en la tabla, etc.
Para ello, cada elemento tendrá:
+ Tarjeta del equipo
    Es la información resumida que se muestra en el tablero:

    -Nombre
    -Valoración media
    -Cantidad de jugadores
    -Posición en la tabla

+ Detalle del equipo
    Se muestra al hacer clic en la tarjeta del equipo:

    -Nombre
    -Formación
    -Plantel de jugadores
    -Mejor jugador
    -Valoración media
    -Últimos resultados
    -Estadísticas generales del equipo

+ Tarjeta de jugador
    Información resumida de cada jugador dentro del plantel:

    -Nombre
    -Media del jugador
    -Posición

+ Detalle del jugador
    Informacion detallada que se muestra al hacer click sobre la carta del jugador:

    -Nombre
    -Posición
    -Edad
    -Media
    -Partidos jugados
    -Goles
    -Asistencias
    -Otras estadísticas individuales (velocidad, regate, tiro, etc)