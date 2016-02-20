#!/bin/bash

#配置 安装目录 mysqlroot密码
install_dir=""
root_pwd="123456"
slave_port=3307
#master配置
master_ip="192.168.0.100"
master_port=3306
master_mysql_name="master"
master_mysql_pwd="123456"
#slave配置
slave_mysql_name="backup"
slave_mysql_pwd="123456"
slave_mysql_ip="192.168.0.101"

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

mysql_dir="$install_dir/mysql/$slave_port"
mysql_data_dir="$mysql_dir/data"
mysql_log_dir="$mysql_dir/log"
mkdir -p $mysql_data_dir
mkdir -p $mysql_log_dir
chown -R mysql:mysql $mysql_dir

mysqlcnf=`cat /etc/my.cnf | grep "^\[mysqld_multi\]"`
if [[ $mysqlcnf == "" ]]; then
    echo "" > /etc/my.cnf
    echo "[mysqld_multi]" >> /etc/my.cnf
    echo "mysqld = /usr/bin/mysqld_safe" >> /etc/my.cnf
    echo "mysqladmin = /usr/bin/mysqladmin" >> /etc/my.cnf
    echo "" >> /etc/my.cnf
fi

mysqlcnf=`cat /etc/my.cnf | grep "^\[mysqld$slave_port\]"`
if [[ $mysqlcnf != "" ]]; then
    echo "slave has existed"
    exit -1
fi

echo "" >> /etc/my.cnf
echo "[mysqld$slave_port]" >> /etc/my.cnf
echo "port = $slave_port" >> /etc/my.cnf
echo "datadir = $mysql_data_dir" >> /etc/my.cnf
echo "socket = $mysql_dir/mysql.sock" >> /etc/my.cnf
echo "pid-file = $mysql_dir/mysqld.pid" >> /etc/my.cnf
echo "log-error = $mysql_log_dir/mysqld.log" >> /etc/my.cnf
echo "user = mysql" >> /etc/my.cnf
echo "server-id = $slave_port" >> /etc/my.cnf
echo "read-only = 1" >> /etc/my.cnf
echo "replicate-ignore-db = mysql" >> /etc/my.cnf
echo "" >> /etc/my.cnf

/usr/bin/mysql_install_db --datadir=$mysql_data_dir --user=mysql
/usr/bin/mysqld_multi start $slave_port
sleep 3 
/usr/bin/mysqladmin -uroot password "$root_pwd" -S $mysql_dir/mysql.sock

echo "mysql status"
/usr/bin/mysqld_multi report $slave_port

if [[ $master_ip != "" && $slave_mysql_name != "" ]]; then
    echo "start slave"
    mysql -h$master_ip -P$master_port -u$master_mysql_name -p$master_mysql_pwd -e "GRANT REPLICATION SLAVE ON *.* TO '$slave_mysql_name'@'$slave_mysql_ip' IDENTIFIED BY '$slave_mysql_pwd';flush privileges;"

    masterinfo=`mysql --skip-column-names -h$master_ip -P$master_port -u$master_mysql_name -p$master_mysql_pwd -e "show master status;"`
    if [[ $masterinfo == "" ]]; then
        echo "get master info failed"
        exit -1
    fi
    masterfile=`echo $masterinfo | awk '{print $1}'`
    masterpostion=`echo $masterinfo | awk '{print $2}'`

    mysql -uroot -p$root_pwd -S $mysql_dir/mysql.sock -e "slave stop;"
    mysql -uroot -p$root_pwd -S $mysql_dir/mysql.sock -e "change master to master_host='$master_ip', master_user='$slave_mysql_name', master_password='$slave_mysql_pwd', master_log_file='$masterfile', master_log_pos=$masterpostion;"
    mysql -uroot -p$root_pwd -S $mysql_dir/mysql.sock -e "slave start;"
    mysql -uroot -p$root_pwd -S $mysql_dir/mysql.sock -e "SHOW SLAVE STATUS\G;"
fi
