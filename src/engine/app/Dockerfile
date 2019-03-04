FROM fedora:latest


LABEL ios.k8s.display-name="fossil-app" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

RUN dnf install -y golang

RUN mkdir /app

WORKDIR /app

ENV GOPATH=/app
ENV GOBIN=/app

RUN go get -u github.com/gorilla/mux

COPY . /app

RUN cd /app && go build engine/util && go install engine/app

RUN chown -R 1001:0 /app && \
    chmod -R ug+rwX /app && \
    chmod 755 /app/app

RUN chmod -R 777 /tmp

RUN echo "1.0" > /etc/imageversion

USER 1001

CMD /app/app