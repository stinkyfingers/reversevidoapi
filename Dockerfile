FROM golang:stretch as build
COPY . /app
WORKDIR /app
RUN go build -o reversevideoapi main.go

FROM heroku/heroku:16
RUN apt-get update -yq && apt-get install ffmpeg -yq
COPY --from=build /app/reversevideoapi /app/reversevideoapi
CMD ["/app/reversevideoapi"]
