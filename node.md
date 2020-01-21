# 1、在Centos上利用源文件编译的方式安装python3.7和pip：
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


# 5.启动Eth全节点并创建账户
## 1、使用fast模式快速同步ETH Goerli测试网全节点，并在后台运行：
```
nohup ~/geth --goerli >>gelori.txt 2>&1 &
```

## 2、查看以太节点同步状态
```
~/geth attach ~/.ethereum/goerli/geth.ipc
```

![eth节点同步情况](/.gitbook/assets/nucypher/同步情况.png)

当输入eth.syncing时，结果为false表示同步完成，同步过程需要30～60分钟，如下图

```
Welcome to the Geth JavaScript console!

instance: Geth/v1.9.10-stable-58cf5686/linux-amd64/go1.13.6
at block: 2037433 (Tue, 21 Jan 2020 07:21:11 UTC)
 datadir: /home/centos/.ethereum/goerli
 modules: admin:1.0 clique:1.0 debug:1.0 eth:1.0 miner:1.0 net:1.0 personal:1.0 rpc:1.0 txpool:1.0 web3:1.0
rpc:1.0 txpool:1.0 web3:1.0

>
> eth.syncing
false
>
```

# 6.生成NU官方的ETH地址格式

```
> eth.accounts
[]
> personal.newAccount()
Password:
Repeat password:
"0xe160672ef1afdc798f869f79d40e0aa963bfac15"
> eth.accounts
["0xe160672ef1afdc798f869f79d40e0aa963bfac15"]
>web3.toChecksumAddress(eth.accounts[0])
> exit
```

# 7.向官方机器人要测试币

注意账号需要提交上面web3.toChecksumAddress(eth.accounts[0])输出的地址

![NU测试币](/.gitbook/assets/nucypher/nucypher-测试币.png)

# 8.启动stakeholder
nucypher stake init-stakeholder --provider ~/.ethereum/goerli/geth.ipc --network cassandra

```
 ____    __            __
/\  _`\ /\ \__        /\ \
\ \,\L\_\ \ ,_\    __ \ \ \/'\      __   _ __
 \/_\__ \\ \ \/  /'__`\\ \ , <    /'__`\/\`'__\
   /\ \L\ \ \ \_/\ \L\.\\ \ \\`\ /\  __/\ \ \/
   \ `\____\ \__\ \__/.\_\ \_\ \_\ \____\\ \_\
    \/_____/\/__/\/__/\/_/\/_/\/_/\/____/ \/_/

The Holder of Stakes.

Wrote new stakeholder configuration to /home/centos/.local/share/nucypher/stakeholder.json
```

# 9.启动新的stake
```
 ____    __            __
/\  _`\ /\ \__        /\ \
\ \,\L\_\ \ ,_\    __ \ \ \/'\      __   _ __
 \/_\__ \\ \ \/  /'__`\\ \ , <    /'__`\/\`'__\
   /\ \L\ \ \ \_/\ \L\.\\ \ \\`\ /\  __/\ \ \/
   \ `\____\ \__\ \__/.\_\ \_\ \_\ \____\\ \_\
    \/_____/\/__/\/__/\/_/\/_/\/_/\/____/ \/_/

The Holder of Stakes.

| # | Account  ---------------------------------- | Balance -----
=================================================================
| 0 | 0xe160672ef1afDc798F869F79d40E0AA963BfaC15  | 0 NU
Select staking account [0]:
Selected 0:0xe160672ef1afDc798F869F79d40E0AA963BfaC15
Enter stake value in NU (15000 NU - 0 NU) [0]: 15000
Enter stake duration (30 - 365) [365]: 30

============================== STAGED STAKE ==============================

Staking address: 0xe160672ef1afDc798F869F79d40E0AA963BfaC15
~ Chain      -> ID # 5 | Goerli
~ Value      -> 15000 NU (15000000000000000000000 NuNits)
~ Duration   -> 30 Days (30 Periods)
~ Enactment  -> Jan 22 00:00 UTC (period #18283)
~ Expiration -> Feb 21 00:00 UTC (period #18313)

=========================================================================

* Ursula Node Operator Notice *
-------------------------------

By agreeing to stake 15000 NU (15000000000000000000000 NuNits):

- Staked tokens will be locked for the stake duration.

- You are obligated to maintain a networked and available Ursula-Worker node
  bonded to the staker address 0xe160672ef1afDc798F869F79d40E0AA963BfaC15 for the duration
  of the stake(s) (30 periods).

- Agree to allow NuCypher network users to carry out uninterrupted re-encryption
  work orders at-will without interference.

Failure to keep your node online, or violation of re-encryption work orders
will result in the loss of staked tokens as described in the NuCypher slashing protocol.

Keeping your Ursula node online during the staking period and successfully
producing correct re-encryption work orders will result in rewards
paid out in ethers retro-actively and on-demand.

Accept ursula node operator obligation? [y/N]:y
Publish staged stake to the blockchain? [y/N]: y
Enter password to unlock account 0xe160672ef1afDc798F869F79d40E0AA963BfaC15:
Broadcasting stake...

Stake initialization transaction was successful.

Transaction details:
OK | deposit stake | 0xf54f71bd5819d2e164d0c97a3afad2a413afd29298f663f538eeaaf37276d3c6 (214203 gas)
Block #2037643 | 0x8ce450a56320c5d57e3ae74ffc4bc053d6b8ce7c1e8e3b7148d0192847a1b718
 See https://goerli.etherscan.io/tx/0xf54f71bd5819d2e164d0c97a3afad2a413afd29298f663f538eeaaf37276d3c6


StakingEscrow address: 0xdC098916291e1ef683A4f469fa32025c872194df

View your stakes by running 'nucypher stake list'
or set your Ursula worker node address by running 'nucypher stake set-worker'.

See https://docs.nucypher.com/en/latest/guides/staking_guide.html
```
## 如果没有提前要币会出现下列报错
```
nucypher.blockchain.eth.actors.InsufficientTokens: Insufficient token balance (NucypherTokenAgent(registry=InMemoryContractRegistry(id=c2e47a), contract=NuCypherToken)) for new stake initialization of 15000 NU
Sentry is attempting to send 0 pending error messages
Waiting up to 2.0 seconds
Press Ctrl-C to quit
```

# 10.查看现有stakes

```
 ____    __            __
/\  _`\ /\ \__        /\ \
\ \,\L\_\ \ ,_\    __ \ \ \/'\      __   _ __
 \/_\__ \\ \ \/  /'__`\\ \ , <    /'__`\/\`'__\
   /\ \L\ \ \ \_/\ \L\.\\ \ \\`\ /\  __/\ \ \/
   \ `\____\ \__\ \__/.\_\ \_\ \_\ \____\\ \_\
    \/_____/\/__/\/__/\/_/\/_/\/_/\/____/ \/_/

The Holder of Stakes.

======================================= Active Stakes =========================================

| ~ | Staker | Worker | # | Value    | Duration     | Enactment
|   | ------ | ------ | - | -------- | ------------ | -----------------------------------------
| 0 | 0xe160 | 0x0000 | 0 | 15000 NU | 30 periods . | Jan 22 00:00 UTC - Feb 21 00:00 UTC
```