FROM alpine:3.1
MAINTAINER Technosophos <technosophos@gmail.com>
EXPOSE 5000
ENV PORT 5000
COPY kubesnoop kubesnoop
ENTRYPOINT ["/kubesnoop"]
