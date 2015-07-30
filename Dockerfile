FROM scratch
EXPOSE 5000
ENV PORT 5000
COPY server-linux server-linux
ENTRYPOINT ["/server-linux"]
