FROM ubuntu:16.04

RUN apt-get update && apt-get install -y openssh-server curl
RUN mkdir -p ~/.ssh
RUN echo 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDFCV9M5sAhOpVQftVS1VDJlXiR7RSaSvwDIDrSoB1Kgz7C60YmRVA+9+bg0dAHLZdF2IFBHjRiVb222a3w1/Lenrl354vHlHGSoyOFi9R1IKdE3zKDhgx02LAhD+N54Q5qE91ENqM8h/ZbWi2MZI2SB09v7Xnx0gYdsrR7r56B/iWXvVbQM78xcnRahEKrBo7z/reil+DVW/OqtvoWLv14xmqKhboPqgAPTTJHqclQcLxK+l4mG5nGVeNM0OQphiFAlkupXpQdvEweMnn+O1sNuFPOWBSNcs1yqdhZMN0HUHGSescpbxAv3ysXPqh/SqxoMAfLW5XEV9SBBO2Cf9OT minghe@oldmac.local' >> ~/.ssh/authorized_keys
RUN mkdir /var/run/sshd
RUN echo 'root:THEPASSWORDYOUCREATED' | chpasswd
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# SSH login fix. Otherwise user is kicked off after login
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd

ENV NOTVISIBLE "in users profile"
RUN echo "export VISIBLE=now" >> /etc/profile

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
