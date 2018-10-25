FROM golang

#specify the working directory
WORKDIR /go/src/aws_api

#copy files to container
COPY . .

#fetch dependencies
RUN go get -d -v ./...

#install the application
RUN go install -v ./...

EXPOSE 8000
#
CMD ["aws_api"]
