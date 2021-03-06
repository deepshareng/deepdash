FROM r.fds.so:5000/golang1.5.3


ADD . /go/src/github.com/MISingularity/deepdash
WORKDIR /go/src/github.com/MISingularity/deepdash
RUN godep go build -o deepdash /go/src/github.com/MISingularity/deepdash/cmd/deepdash/main.go
CMD /go/src/github.com/MISingularity/deepdash/deepdash
# -mongo-addr=127.0.0.1:27017
# -deepstats-addr=http://127.0.0.1:16759
# -appinfo-addr=http://127.0.0.1:8080
# -auth-addr specify the Oauth authentication service address
# -client-id  Specify the application client id for oauth (default "deepdash-local")
#-client-secret Specify the application client secret for oauth (default "deepdash-deepshare")
# -redirecturi Specify the oauth redirect uri (default "http://localhost:10033/auth/code")

EXPOSE 10033