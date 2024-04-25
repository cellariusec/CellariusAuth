# Para crear una cuenta:
```
curl -X POST -H "Content-Type: application/json" -d '{"Email":"<email>","Password":"<password>"}' https://independent-sparkle-production.up.railway.app/signup
```

# Para ingresar: 

```
curl -X POST -H "Content-Type: application/json" -d '{"Email":"example@example.com","Password":"password123"}' https://independent-sparkle-production.up.railway.app/login
```

# URL de prueba para validar token en MDM

```
http://localhost:5000/ValidateWebJWT?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2Rldi1hZjdhc3NvdmJvbTZhMjhvLnVzLmF1dGgwLmNvbS8iLCJzdWIiOiIyZVNhOFpOUVc2all5enhOb2hsR0RBQTJNZWlrVk90NUBjbGllbnRzIiwiYXVkIjoiaHR0cDovL2xvY2FsaG9zdDo1MDAwL2p3dCIsImlhdCI6MTcxMjYwNzA1NSwiZXhwIjoxNzEyNjkzNDU1LCJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMiLCJhenAiOiIyZVNhOFpOUVc2all5enhOb2hsR0RBQTJNZWlrVk90NSJ9.vl1SbNTblYTp7rvMy68UaOk_cEuLWY7pElXs-zGQGHY&refreshToken=71cf9eff-80c9-40a0-aebf-780a15a94610
```


