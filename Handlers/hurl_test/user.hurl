

# 1. Crear nuevo usuario
POST http://localhost:8080/users
Content-Type: application/json
{
    "user_name":    "Usuario de prueba",
    "amigos":       0,
    "photo":        1
}
HTTP 201
[Captures]
user_id: jsonpath "$.id_user"
[Asserts]
jsonpath "$.user_name" == "Usuario de prueba"
jsonpath "$.amigos" == 0
jsonpath "$.photo" == 1

# 2. Obtener el usuario recién creado
GET http://localhost:8080/users/{{user_id}}
HTTP 200
[Asserts]
jsonpath "$.id_user" == {{user_id}}
jsonpath "$.user_name" == "Usuario de prueba"
jsonpath "$.amigos" == 0
jsonpath "$.photo" == 1

# 3. Listar usuarios
GET http://localhost:8080/users
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"

# 4. Actualizar usuario
PUT http://localhost:8080/users/{{user_id}}
Content-Type: application/json
{
    "user_name":    "Usuario",
    "amigos":       1,
    "photo":        2
}
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.user_name" == "Usuario"
jsonpath "$.amigos" == 1
jsonpath "$.photo" == 2

# 5. Verificar que se actualizó el usuario
GET http://localhost:8080/users/{{user_id}}
HTTP 200
[Asserts]
jsonpath "$.id_user" == {{user_id}}
jsonpath "$.user_name" == "Usuario"
jsonpath "$.amigos" == 1
jsonpath "$.photo" == 2

# 6. Eliminar usuario
DELETE http://localhost:8080/users/{{user_id}}
HTTP 204

# 7. Verificar que se eliminó 
GET http://localhost:8080/users/{{user_id}}
HTTP 404

# ===============
# Tests de error

# 8. Intentar obtener usuario no existente
GET http://localhost:8080/users/999999
HTTP 404

# 9. Intento de creación con campos vacíos o inválidos
POST http://localhost:8080/users
Content-Type: application/json
{
    "user_name":    "       ",
    "amigos":       -3,
    "photo":        0
}
HTTP 400

# 10. Búsqueda con ID no numérico
GET http://localhost:8080/users/abc
HTTP 400

# 11. Intentar actualizar un usuario que no existe
PUT http://localhost:8080/users/999999
Content-Type: application/json
{
    "id_user":      999999,
    "user_name":    "Usuario Inexistente",
    "amigos":       0,
    "photo":        1
}
HTTP 404

# 12. Intentar eliminar un usuario que no existe
DELETE http://localhost:8080/users/999999
HTTP 404