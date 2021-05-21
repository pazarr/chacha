FROM busybox
ADD chacha /bin/chacha
RUN chmod +x /bin/chacha
ENTRYPOINT ["/bin/chacha"]
