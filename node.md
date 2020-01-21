# 1、利用源文件编译的方式安装python3.7和pip：
sudo yum update

sudo yum install -y zlib-devel bzip2-devel openssl-devel ncurses-devel sqlite-devel readline-devel tk-devel gcc make libffi-devel epel-release wget

wget https://www.python.org/ftp/python/3.7.0/Python-3.7.0.tgz

tar -zxvf Python-3.7.0.tgz

cd Python-3.7.0

yum groupinstall "Development Tools”  

./configure prefix=/usr/local/python3 --with-ssl

sudo make & make install

sudo ln -sf /usr/local/python3/bin/python3.7 /usr/bin/python3

sudo ln -sf /usr/local/python3/bin/pip3.7 /usr/bin/pip3

sudo pip3 install --upgrade pip

# 2、安装virtualenv

pip3 install --user virtualen

将~/.local/bin增加到~/.bash_profile文件中

source ~/.bash_profile

# 3、安装nucypher
sudo yum install -y libffi-dev python3-dev python3-virtualenv

cd ~ | virtualenv nucypher-venv

source ~/nucypher-venv/bin/activate

pip3 install -U nucypher

nucypher –help

deactivate

# 4、安装以太客户端Geth
登陆到网页 https://geth.ethereum.org/downloads/ ，查找下载链接进行下载，当前是Geth 1.9.10版本：

wget https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.9.10-58cf5686.tar.gz

tar xf geth-linux-amd64-1.9.10-58cf5686.tar.gz

mv geth-linux-amd64-1.9.10-58cf5686/geth ./

rm geth-linux-amd64-1.9.10-58cf5686.tar.gz


# 启动Eth全节点并创建账户
## 1、使用fast模式快速同步ETH Goerli测试网全节点，并在后台运行：
```
nohup ~/geth --goerli >>gelori.txt 2>&1 &
```

## 2、查看以太节点同步状态
```
~/geth attach ~/.ethereum/goerli/geth.ipc
```

![eth节点同步情况](/.gitbook/assets/nucypher/同步情况.png)

当输入eth.syncing时，结果为false表示同步完成，同步过程需要30～60分钟。