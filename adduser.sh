#!/bin/bash

if [[ $# != 2 ]]; then
    echo "should input username userpwd"
    exit -1 
fi

user_name=$1
user_pwd=$2

if id -u $user_name >/dev/null 2>&1; then
    echo "$user_name has existed"
    exit -1
else
    echo "create user $user_name..."
    useradd $user_name
    echo "$user_name:$user_pwd" | chpasswd

    echo "" >> /home/$user_name/.bash_profile
    echo "ulimit -c unlimited" >> /home/$user_name/.bash_profile
    
    if ! egrep $user_name /etc/security/limits.conf > /dev/null;then
        echo "$user_name          soft    nofile          10240" >> /etc/security/limits.conf
        echo "$user_name          hard    nofile          10240" >> /etc/security/limits.conf
    fi
	
    yum -y install vim

    echo "set encoding=utf-8" >> /home/$user_name/.vimrc
    echo "set termencoding=&encoding" >> /home/$user_name/.vimrc
    echo "set fileencodings=utf-8,chinese,gb2312" >> /home/$user_name/.vimrc
    echo "set fileencoding=utf-8" >> /home/$user_name/.vimrc
fi
