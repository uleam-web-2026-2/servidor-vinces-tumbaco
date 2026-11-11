# Etapa 1: compilar, con todas las herramientas de Go
FROM golang:1.27 AS compilacion
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /servidor .

# Etapa 2: la imagen de producción, solo con el binario y usuario sin privilegios
FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=compilacion /servidor /servidor
ENTRYPOINT ["/servidor"]
