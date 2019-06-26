FROM scratch
COPY kfilt /
ENTRYPOINT ["/kfilt"]
