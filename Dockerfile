# Stage 1 - build frontend
FROM node:16-alpine3.15 AS frontend

WORKDIR /app

COPY frontend .

RUN npm install

RUN npm run build


# Stage 2 - build backend
FROM golang:1.19.2-bullseye AS backend

WORKDIR /app

COPY backend .

RUN go mod download

RUN go build -o main . 


# Stage 3 - build final image
FROM debian:bullseye-slim

WORKDIR /app

SHELL ["/bin/bash", "-c"] 

RUN apt-get update && apt-get install -y python3 ffmpeg wget unzip

RUN wget https://www.bok.net/Bento4/binaries/Bento4-SDK-1-6-0-639.x86_64-unknown-linux.zip -O /tmp/Bento4.zip \
  && unzip /tmp/Bento4.zip -d /home/ \
  && rm -f /tmp/Bento4.zip

ENV PATH "$PATH:/home/Bento4-SDK-1-6-0-639.x86_64-unknown-linux/bin"

ENV LISTEN_ADDR ":8000"

COPY --from=frontend /app/dist ./www

COPY --from=backend /app/main .

EXPOSE 8000

VOLUME ["/app/atus_data"]

CMD ["./main"]