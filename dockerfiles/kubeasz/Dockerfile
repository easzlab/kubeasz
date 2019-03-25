# NOTE: Prepare following Requirements and Build the image "kubeasz:$TAG"
# Requirement 1: git clone --depth=1 https://github.com/gjmzj/kubeasz.git 
# Requirement 2: download binaries at https://pan.baidu.com/s/1c4RFaA, and put into dir 'kubeasz/bin'
# Build: docker build -t kubeasz:$TAG .

FROM jmgao1983/ansible:v2.6 

COPY kubeasz/ /etc/ansible

RUN ln -s /etc/ansible/tools/easzctl /usr/bin/easzctl

CMD [ "sleep", "360000000" ]
