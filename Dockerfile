FROM golang:1.8.5-jessie
# install glide
RUN go get github.com/Masterminds/glide
# create a working directory
WORKDIR /go/src/app
# add glide.yaml and glide.lock
ADD glide.yaml glide.yaml
ADD glide.lock glide.lock
# install packages
RUN glide install
# add source code
ADD src src
# run main.go
RUN go build -o mongo-cluster-stats src/*.go
ENTRYPOINT ["./mongo-cluster-stats"]
