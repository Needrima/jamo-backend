FROM golang:latest

WORKDIR /app # declare /app as root directory "."

ARG mongo_db_port=27017

ARG mongo_db_username=amirdeen

ARG mongo_db_password=27017

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /jamo-backend

CMD [ "/jamo-backend" ]