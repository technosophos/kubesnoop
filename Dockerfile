FROM scratch
EXPOSE 5000
ENV PORT 5000
COPY kubesnoop kubesnoop
ENTRYPOINT ["/kubesnoop"]
