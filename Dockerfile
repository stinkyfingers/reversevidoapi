FROM golang:stretch as build
COPY . /app
WORKDIR /app
RUN go build -o reversevideoapi main.go
RUN apt-get update -yq && apt-get install ffmpeg -yq

FROM heroku/heroku:16
COPY --from=build /app/reversevideoapi /app/reversevideoapi
CMD ["/app/reversevideoapi"]
