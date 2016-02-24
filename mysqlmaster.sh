#!/bin/bash

#配置 安装目录 root密码 
install_dir=""
root_pwd="123456"
#读写用户名字 密码 只读用户名字 密码
master_name="master"
master_pwd="123456"
readonly_name="readonly"
readonly_pwd="123456"

echo "install mysql"
if yum list installed mysql >/dev/null 2>&1; then
    echo "has installed mysql"
else
    echo "yum install mysql"
    echo ""
    yum -y install mysql
fi

if yum list installed mysql-server >/dev/null 2>&1; then
    echo "has installed mysql-server"
else
    echo "yum install mysql-server"
    echo ""
    yum -y install mysql-server
fi

if yum list installed mysql-devel >/dev/null 2>&1; then
    echo "has installed mysql-devel"    
else
    echo "yum install mysql-devel"
    echo ""
    yum -y install mysql-devel
fi

mysql_dir="$install_dir/mysql"
mysql_data_dir="$mysql_dir/data"
mysql_log_dir="$mysql_dir/log"
mkdir -p $mysql_data_dir
mkdir -p $mysql_log_dir
chown -R mysql:mysql $mysql_dir

echo "conf my.cnf"
echo "" > /etc/my.cnf
echo "[mysqld]" >> /etc/my.cnf
echo "datadir = $mysql_data_dir" >> /etc/my.cnf
echo "socket = /var/lib/mysql/mysql.sock" >> /etc/my.cnf
echo "user = mysql" >> /etc/my.cnf
echo "# Disabling symbolic-links is recommended to prevent assorted security risks" >> /etc/my.cnf
echo "symbolic-links = 0" >> /etc/my.cnf
echo "" >> /etc/my.cnf
echo "# master conf" >> /etc/my.cnf
echo "server-id = 1" >> /etc/my.cnf
echo "log_bin = $mysql_log_dir/mysql-bin.log" >> /etc/my.cnf
echo "read-only = 0" >> /etc/my.cnf
echo "binlog-ignore-db = mysql" >> /etc/my.cnf
echo "" >> /etc/my.cnf
echo "[mysqld_safe]" >> /etc/my.cnf
echo "log-error = $mysql_log_dir/mysqld.log" >> /etc/my.cnf
echo "pid-file = /var/run/mysqld/mysqld.pid" >> /etc/my.cnf
echo "" >> /etc/my.cnf

echo "start mysql"
chkconfig mysqld on
service mysqld restart
sleep 3
/usr/bin/mysqladmin -u root password "$root_pwd"

echo "add mysql user"
mysql -uroot -p"$root_pwd" -e "drop user ''@'localhost';"
host_name=`hostname`
host_user=`echo $host_name | sed 's!_!\\_!g'`
mysql -uroot -p"$root_pwd" -e "drop user ''@'$host_user';"

if [[ $master_name != "" ]]; then
    mysql -uroot -p"$root_pwd" -e "CREATE USER '$master_name'@'%' IDENTIFIED BY '$master_pwd';"
	mysql -uroot -p"$root_pwd" -e "GRANT all privileges on *.* to '$master_name'@'%' identified by '$master_pwd' WITH GRANT OPTION;"
	mysql -uroot -p"$root_pwd" -e "flush privileges;"
fi

if [[ $readonly_name != "" ]]; then
    mysql -uroot -p"$root_pwd" -e "CREATE USER '$readonly_name'@'%' IDENTIFIED BY '$readonly_pwd';"
	mysql -uroot -p"$root_pwd" -e "GRANT SELECT on *.* to '$readonly_name'@'%' identified by '$readonly_pwd';"
	mysql -uroot -p"$root_pwd" -e "flush privileges;"
fi

echo "mysql status"
service mysqld status
