## **Cellarius Auth** 

## **Servidor de autenticación Cellarius** 

###__Para correr localmente:__
docker pull joseignnacio/cellariusauth
docker run -p 8080:80 joseignnacio/cellariusauth:latest

###__Para crear una cuenta nueva:__
curl -X POST -H "Content-Type: application/json" -d '{"Email":<email>,"Password":<password>}' http://localhost:8080/signup


###__Para iniciar sesión y generar JWTs:__
curl -X POST -H "Content-Type: application/json" -d '{"Email":"example@example.com","Password":"password123"}' http://localhost:8080/login



#__URL de prueba para validar token en MDM:__

http://localhost:5000/ValidateWebJWT?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2Rldi1hZjdhc3NvdmJvbTZhMjhvLnVzLmF1dGgwLmNvbS8iLCJzdWIiOiIyZVNhOFpOUVc2all5enhOb2hsR0RBQTJNZWlrVk90NUBjbGllbnRzIiwiYXVkIjoiaHR0cDovL2xvY2FsaG9zdDo1MDAwL2p3dCIsImlhdCI6MTcxMjYwNzA1NSwiZXhwIjoxNzEyNjkzNDU1LCJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMiLCJhenAiOiIyZVNhOFpOUVc2all5enhOb2hsR0RBQTJNZWlrVk90NSJ9.vl1SbNTblYTp7rvMy68UaOk_cEuLWY7pElXs-zGQGHY&refreshToken=71cf9eff-80c9-40a0-aebf-780a15a94610
