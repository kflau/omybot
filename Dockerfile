FROM golang
RUN git clone https://github.com/bwmarrin/discordgo $GOPATH/src/github.com/bwmarrin/discordgo
RUN git clone https://github.com/gorilla/websocket $GOPATH/src/github.com/gorilla/websocket
RUN go get golang.org/x/crypto/nacl/secretbox
COPY src/ $GOPATH/src/
WORKDIR $GOPATH/bin
RUN go build -a $GOPATH/src/github.com/lkfai/omybot/*.go
RUN chmod a+x $GOPATH/bin/bot
CMD ["/go/bin/bot"]
