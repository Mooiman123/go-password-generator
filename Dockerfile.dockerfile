# Gebruik de officiële Go-image
FROM golang:1.20

# Zet de werkdirectory in de container
WORKDIR /app

# Kopieer go-modules en download dependencies
COPY go.mod ./
RUN go mod tidy

# Kopieer de rest van
