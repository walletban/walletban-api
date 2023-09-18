FROM golang:1.20

ARG dbuser
ARG dburl
ARG dbpass
ARG jwtpass

ENV DB_SERVER $dburl
ENV DB_USER $dbuser
ENV DB_PASSWORD $dbpass
ENV JWT_SECRET $jwtpass
WORKDIR /go/src/app
COPY . .

RUN make build

EXPOSE 8000
CMD ["./bin/walletban-api"]