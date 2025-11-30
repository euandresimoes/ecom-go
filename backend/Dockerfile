# ---- Build ----
FROM golang:1.25-alpine AS build

WORKDIR /app
    
COPY go.mod go.sum ./
RUN go mod download
    
COPY . .
    
RUN go build -o app ./cmd
    
# ---- Final ----
FROM alpine
    
WORKDIR /app
    
COPY --from=build /app/app .
    
EXPOSE 7020
    
CMD ["./app"]