# Dockerfile for Rotating the indices in elastic of the EFK deployment
#
# @author:  gjmzj
# @repo:    https://github.com/kubeasz/mirrorepo/es-index-rotator
# @ref:     https://github.com/easzlab/kubeasz/tree/master/dockerfiles/es-index-rotator

FROM alpine:3.8

COPY rotate.sh /bin/rotate.sh

RUN echo "===> Installing essential tools..."   && \
    apk --update add bash curl coreutils        && \
    echo "===> Cleaning up cache..."            && \
    rm -rf /var/cache/apk/*                     && \
    chmod +x /bin/rotate.sh

CMD ["/bin/rotate.sh"]
