from golang:1.6
RUN mkdir -p /go/src 
COPY src /go/src/
ENV GOPATH /go
ENV GOBIN /go/bin
RUN go install currencyservice 
ENV PATH /go/bin
EXPOSE 9000
ENTRYPOINT currencyservice
