FROM ubuntu:16.04

RUN apt update
RUN DEBIAN_FRONTEND=noninteractive apt -y install file curl man vim git gcc


WORKDIR /tmp
RUN curl -o go1.16.tar.gz -L https://golang.org/dl/go1.16.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.16.tar.gz
RUN rm go1.16.tar.gz

RUN mkdir /root/go
WORKDIR /root/go

RUN git clone https://github.com/magefile/mage
WORKDIR /root/go/mage
RUN PATH=$PATH:/usr/local/go/bin go run bootstrap.go

ENV PATH=$PATH:/usr/local/go/bin:/root/go/bin

WORKDIR /root/go

EXPOSE 8000

CMD [ "/bin/bash" ]