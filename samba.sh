#!/bin/bash

if [[ $# != 1 ]]; then
    echo "you should input shared username"
    exit -1 
fi

user_name=$1

yum -y install samba
echo "[$user_name]" >> /etc/samba/smb.conf
echo "	comment = $user_name" >> /etc/samba/smb.conf
echo "	path = /home/$user_name" >> /etc/samba/smb.conf
echo "	writable = yes" >> /etc/samba/smb.conf

smbpasswd -a -n $user_name
service smb restart

chkconfig smb on
chkconfig iptables off

service iptables stop
setenforce 0
sed -i s!SELINUX=enforcing!SELINUX=disabled!g /etc/selinux/config
