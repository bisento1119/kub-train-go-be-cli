# syntax=docker/dockerfile:1

#FROM alpine:3.18
FROM amd64/ubuntu:22.04


WORKDIR /.

# Download Go modules
COPY ./go-be-cli ./
COPY ./environments/dockerConf.yml ./environments/localConf.yml

# Set destination for COPY


# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 4882

# Run
CMD ["./go-be-cli"]