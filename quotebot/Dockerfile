FROM golang
RUN git clone https://github.com/bwmarrin/discordgo $GOPATH/src/github.com/bwmarrin/discordgo
RUN git clone https://github.com/gorilla/websocket $GOPATH/src/github.com/gorilla/websocket
RUN go get golang.org/x/crypto/nacl/secretbox
COPY . $GOPATH/src/github.com/kflau/omybot/quotebot/
WORKDIR $GOPATH/src/github.com/kflau/omybot/quotebot/
RUN go install -a
RUN chmod a+x $GOPATH/bin/quotebot
RUN rm -fr $GOPATH/src/github.com/kflau/omybot/
RUN chsh -s /bin/false
RUN useradd -ms /usr/sbin/nologin app
USER app
CMD ["/go/bin/quotebot"]
