# Series Tracker — Backend

API REST para gestionar series de televisión, construida con Go y SQLite.

**Frontend:** https://github.com/Sofilayerdi/Proyecto1-Web-Frontend.git

**App en producción:** https://seriestracker-231929.onrender.com

**Servidor backend:** https://series-tracker-rk1z.onrender.com

## Tecnologías
- Go con net/http
- Chi (router)
- SQLite (base de datos)

## Correr localmente

1. Clona el repositorio
2. Asegúrate de tener Go instalado o instalar en https://go.dev/doc/install
3. Instala dependencias:
```bash
   go mod tidy
```
4. Corre el servidor:
```bash
   go run .
```
5. El servidor corre en `http://localhost:8000`

## Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | /series | Listar series (soporta ?page=, ?limit=, ?q=) |
| GET | /series/{id} | Obtener una serie |
| POST | /series | Crear una serie |
| PUT | /series/{id} | Editar una serie |
| DELETE | /series/{id} | Eliminar una serie |

## CORS
CORS está configurado para permitir cualquier origen (`*`), necesario porque el cliente y el servidor corren en dominios distintos. El servidor responde a peticiones OPTIONS para el preflight del navegador.

## Challenges implementados
- Códigos HTTP correctos — 201 al crear, 204 al eliminar, 404, 400
- Validación server-side 
- Paginación ?page= y ?limit=
- Búsqueda ?q=
- Ordenamiento ?sort= y ?order=asc|desc
- Exportar la lista de series a Excel (.xlsx) 

## Reflexión
Go me pareció bastante sencillo de usar y entender para construir el servidor HTTP. Me gustó que levantar una API REST fuera tan directo, especialmente usando net/http y Chi, que facilitan mucho el manejo de rutas. Al inicio me costó un poco entender bien cómo usarlos en la separación cliente-servidor, pero cuando logré comprender la estructura del código, todo empezó a fluir de forma más limpia. También aprendí sobre CORS, al principio fue confuso y no tenía claro dónde configurarlo, pero luego entendí que todo se resolvía ajustando correctamente los headers para permitir que el cliente consumiera la API. Si usaría Go de nuevo para backends donde la performance sea importante.
