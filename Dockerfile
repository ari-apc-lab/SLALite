###

FROM golang:1.12.7-alpine

RUN apk add --no-cache git curl

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh


WORKDIR /go/src/SLALite

COPY . .
RUN rm -rf vendor
#RUN dep init PAOLO
RUN dep ensure

ARG VERSION
ARG DATE

RUN CGO_ENABLED=0 GOOS=linux go build -a -o SLALite -ldflags="-X main.version=${VERSION} -X main.date=${DATE}" .

#FROM golang:alpine
RUN mkdir /opt/slalite

RUN cp -r /go/src/SLALite /opt/slalite


WORKDIR /opt/slalite
#COPY --from=builder /go/src/SLALite/SLALite .
#COPY /go/src/SLALite .
#
RUN mkdir /etc/slalite
COPY /docker/slalite.yml /etc/slalite/slalite.yml

#ADD /euxdat-accounting /accounting

#COPY /euxdat-accounting /accounting

COPY euxdat-accounting /opt/slalite/euxdat-accounting

RUN apk update && apk add python python-dev py-pip python3 python3-dev libffi-dev openssl-dev gcc libc-dev make
RUN pip3 install pipenv
RUN pip3 install requests  
RUN virtualenv -p /usr/bin/python3.7 py37
ENV VIRTUAL_ENV /py37
ENV PATH /py37/bin:$PATH

 
RUN source py37/bin/activate
RUN pip install requests

RUN addgroup -S slalite && adduser -D -G slalite slalite
RUN chown -R slalite:slalite /etc/slalite && chmod 700 /etc/slalite
RUN chown -R slalite:slalite /opt/slalite && chmod 700 /opt/slalite


CMD . source py37/bin/activate
#EXPOSE 8090
EXPOSE 8094
#ENTRYPOINT ["./run_slalite.sh"]
USER slalite
#ENTRYPOINT ["./run_slalite.sh"]


ENTRYPOINT ["/opt/slalite/SLALite/SLALite"]

#ENTRYPOINT ["/go/src/SLALite"]


