# Usa una imagen base oficial de Go
FROM golang:1.19-alpine

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia el archivo go.mod y go.sum al directorio de trabajo
COPY go.mod go.sum ./

# Descarga las dependencias del módulo
RUN go mod download

# Copia el resto del código fuente de la aplicación
COPY . .

# Compila la aplicación
RUN go build -o app main.go

# Expone el puerto en el que la aplicación escucha
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./app"]
