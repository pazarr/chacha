FROM busybox
ADD chacha /bin/chacha
RUN chmod +x /bin/chacha
EXPOSE 80/tcp
ENTRYPOINT ["/bin/chacha"]
