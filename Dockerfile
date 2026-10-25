# to run kubeasz in a container
# @author:  gjmzj
# @repo:     https://github.com/easzlab/kubeasz

FROM easzlab/ansible:2.16.19-lite

ENV TZ="Asia/Shanghai"

COPY . /etc/kubeasz

RUN set -x \
    && ln -s -f /etc/kubeasz/ezctl /usr/bin/ezctl \
    && ln -s -f /etc/kubeasz/ezdown /usr/bin/ezdown \
    && mkdir -p /usr/libexec \
    && ln -s /usr/bin/python3 /usr/libexec/platform-python
