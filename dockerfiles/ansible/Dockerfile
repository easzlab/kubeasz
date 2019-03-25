# Dockerfile for building Ansible image for Alpine 3
# Origin from https://github.com/William-Yeh/docker-ansible

FROM alpine:3.8

RUN \ 
    echo "===> Adding Python runtime..."  && \
    apk --update add python py-pip openssl ca-certificates    && \
    apk --update add --virtual build-dependencies \
        python-dev libffi-dev openssl-dev build-base          && \
    pip install --upgrade pip cffi                            && \
    \
    \
    echo "===> Installing Ansible..."  && \
    pip install ansible==2.6.12        && \
    \
    \
    echo "===> Installing handy tools..."          && \
    pip install --upgrade pycrypto                 && \
    apk --update add bash openssh-client rsync     && \
    \
    \
    echo "===> Removing package list..."  && \
    apk del build-dependencies            && \
    rm -rf /var/cache/apk/*               && \
    rm -rf /root/.cache

# default command: display Ansible version
CMD [ "ansible", "--version" ]
