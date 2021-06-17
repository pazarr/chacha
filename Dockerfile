FROM busybox
ADD chacha /bin/chacha
RUN chmod +x /bin/chacha
ADD html /var/html
EXPOSE 80/tcp
ENTRYPOINT ["/bin/chacha"]
