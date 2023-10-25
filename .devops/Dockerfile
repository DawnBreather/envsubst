FROM golang:alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
# RUN GOOS=linux GOARCH=amd64 go test ./...
RUN GOOS=linux GOARCH=amd64 go build -o /out/envmake envmake/main.go
RUN GOOS=linux GOARCH=amd64 go build -o /out/envsubst envsubst/main.go
RUN GOOS=linux GOARCH=amd64 go build -o /out/set2secret set2secret/main.go
RUN GOOS=linux GOARCH=amd64 go build -o /out/setsubst setsubst/main.go

FROM debian

COPY --from=build /out/envmake /usr/bin
COPY --from=build /out/envsubst /usr/bin
COPY --from=build /out/set2secret /usr/bin
COPY --from=build /out/setsubst /usr/bin
