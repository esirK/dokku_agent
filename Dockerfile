FROM golang:1.17 as go


FROM dokku/dokku:latest
# copy go binary from go image as dokku user
COPY --from=go --chown=dokku /usr/local/go /usr/local/go

ENV PATH="{$PATH}:/usr/local/go/bin"

COPY ./contrib/cmd.sh /tmp

RUN chmod +x /tmp/cmd.sh

CMD ["/tmp/cmd.sh"]
