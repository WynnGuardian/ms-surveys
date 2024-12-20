FROM golang:1.23.1

WORKDIR /go/src
COPY ./backend /go/src

ENV PATH="/go/bin:${PATH}"

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

EXPOSE 8086

CMD ["tail", "-f", "/dev/null"]