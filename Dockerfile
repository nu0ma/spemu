# Multi-stage build for spemu
FROM scratch

# Copy the binary built by GoReleaser
COPY spemu /usr/local/bin/spemu

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/spemu"]