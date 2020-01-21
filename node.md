# 相关计划

1.阶段1任务 https://blog.nucypher.com/casi-phase-1-tasks/


# 1. 在Centos上利用源文件编译的方式安装python3.7和pip：
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

# 11.添加worker
首先按照[生成NU官方的ETH地址格式](#6.生成NU官方的ETH地址格式)，生成一个新的eth账号

```
nucypher stake set-worker
```
输出：
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
| 0 | 0xe160672ef1afDc798F869F79d40E0AA963BfaC15  | 30000 NU
| 1 | 0x7B6BCC437e0D4B7f857E48364Ec098B2A53001f2  | 0 NU
Select staking account [0]:
Selected 0:0xe160672ef1afDc798F869F79d40E0AA963BfaC15
Enter worker address: 0x7B6BCC437e0D4B7f857E48364Ec098B2A53001f2
Enter password to unlock account 0xe160672ef1afDc798F869F79d40E0AA963BfaC15:
Commit to bonding worker 0x7B6BCC437e0D4B7f857E48364Ec098B2A53001f2 to staker 0xe160672ef1afDc798F869F79d40E0AA963BfaC15 for a minimum of 2 periods? [y/N]: y

Worker 0x7B6BCC437e0D4B7f857E48364Ec098B2A53001f2 successfully bonded to staker 0xe160672ef1afDc798F869F79d40E0AA963BfaC15
OK | set_worker | 0x9d78b8fd64d514601b3ad618ced47488580964d6236ed08ab7e3331a127ddbb6 (75287 gas)
Block #2037750 | 0xfa7f8db748fa4a52898a41db0ae037635272e02f99bde9eb2021192cc1847bf1
 See https://goerli.etherscan.io/tx/0x9d78b8fd64d514601b3ad618ced47488580964d6236ed08ab7e3331a127ddbb6

Bonded at period #18282 (2020-01-21 08:40:07.008259+00:00)
This worker can be replaced or detached after period #18284 (2020-01-23 00:00:00+00:00)
```
此时查看stakelist，发现worker不再是0000了
```
nucypher stake list
```
输出：
```
The Holder of Stakes.

======================================= Active Stakes =========================================

| ~ | Staker | Worker | # | Value    | Duration     | Enactment
|   | ------ | ------ | - | -------- | ------------ | -----------------------------------------
| 0 | 0xe160 | 0x7B6B | 0 | 15000 NU | 30 periods . | Jan 22 00:00 UTC - Feb 21 00:00 UTC
```

# 12.配置Ursha

## 申请eth的Gorli测试网币
```
geth attach ~/.ethereum/goerli/geth.ipc
> personal.newAccount();
> eth.accounts[1]
["0xc080708026a3a280894365efd51bb64521c45147"]
The new account is 0xc080708026a3a280894365efd51bb64521c45147 in this case.
```
然后转到 https://goerli-faucet.slock.it/
去申请,不申请的话后面一步会报错：
```
    txhash = self.client.send_raw_transaction(signed_raw_transaction)
  File "/home/centos/nucypher-venv/lib/python3.7/site-packages/nucypher/blockchain/eth/clients.py", line 235, in send_raw_transaction
    return self.w3.eth.sendRawTransaction(raw_transaction=transaction)
  File "/home/centos/nucypher-venv/lib/python3.7/site-packages/web3/eth.py", line 388, in sendRawTransaction
    [raw_transaction],
  File "/home/centos/nucypher-venv/lib/python3.7/site-packages/web3/manager.py", line 152, in request_blocking
    raise ValueError(response["error"])
builtins.ValueError: {'code': -32000, 'message': 'insufficient funds for gas * price + value'}
```

## 使用CLI运行Ursula

```
nucypher ursula init --provider ~/.ethereum/goerli/geth.ipc --poa --staker-address 0xe160672ef1afDc798F869F79d40E0AA963BfaC15 --network cassandra
```
输出:
```
 ,ggg,         gg
dP""Y8a        88                                   ,dPYb,
Yb, `88        88                                   IP'`Yb
 `"  88        88                                   I8  8I
     88        88                                   I8  8'
     88        88   ,gggggg,    ,g,     gg      gg  I8 dP    ,gggg,gg
     88        88   dP""""8I   ,8'8,    I8      8I  I8dP    dP"  "Y8I
     88        88  ,8'    8I  ,8'  Yb   I8,    ,8I  I8P    i8'    ,8I
     Y8b,____,d88,,dP     Y8,,8'_   8) ,d8b,  ,d8b,,d8b,_ ,d8,   ,d8b,
      "Y888888P"Y88P      `Y8P' "YY8P8P8P'"Y88P"`Y88P'"Y88P"Y8888P"`Y8


the Untrusted Re-Encryption Proxy.


| # | Account  ---------------------------------- | Balance -----
=================================================================
| 0 | 0xe160672ef1afDc798F869F79d40E0AA963BfaC15
| 1 | 0x7B6BCC437e0D4B7f857E48364Ec098B2A53001f2
Select worker account [0]: 1
Selected 1:0x7B6BCC437e0D4B7f857E48364Ec098B2A53001f2
Is this the public-facing IPv4 address (13.229.53.90) you want to use for Ursula? [y/N]: y
Enter NuCypher keyring password (16 character minimum):
Repeat for confirmation:
Generated keyring /home/centos/.local/share/nucypher/keyring
Saved configuration file /home/centos/.local/share/nucypher/ursula.json

If you haven't done it already, initialize a NU stake with 'nucypher stake' or

To run an Ursula node from the default configuration filepath run:

'nucypher ursula run'
```

### 运行ursula

```
nucypher ursula run --interactive
```
成功后显示：
```

Enter password to unlock account 0x7B6BCC437e0D4B7f857E48364Ec098B2A53001f2:
Enter NuCypher keyring password:
Decrypting NuCypher keyring...
Connecting to preferred teacher nodes...
Starting Ursula on 13.229.53.90:9151
Connecting to cassandra
Working ~ Keep Ursula Online!
Attached 0xe160672ef1afDc798F869F79d40E0AA963BfaC15@13.229.53.90:9151
♄ ⛇ | SlateGray Saturn Cyan Snowman


Type 'help' or '?' for help
Ursula(0xe160672) >>>
```

### 查看Ursula节点状态

```
Ursula(0xe160672) >>> status

⇀URSULA ♄ ⛇↽
(Ursula)⇀SlateGray Saturn Cyan Snowman↽ (0xe160672ef1afDc798F869F79d40E0AA963BfaC15)
Uptime .............. 0:01:12
Start Time .......... 1 minute ago
Fleet State.......... c941129 ⇀Peru Crossbones↽ ☠
Learning Status ..... Learning at 5s Intervals
Learning Round ...... Round #16
Operating Mode ...... Decentralized
Rest Interface ...... 13.229.53.90:9151
Node Storage Type ... Local
Known Nodes ......... 15
Work Orders ......... 0
Current Teacher ..... (Ursula)⇀LightSlateGray Dharma PaleTurquoise Star↽ (0x7820aDA8554197e41Ff0FA54aF30BbB98c716765)
Current Period ...... 18282
Worker Address ...... 0x7B6BCC437e0D4B7f857E48364Ec098B2A53001f2
```

### 查看已知的Ursula节点地址

```
Ursula(0xe160672) >>> known_nodes

Known Nodes (connected 15 / seen 17)
Fleet State c941129 ⇀Peru Crossbones↽ ☠
54.152.254.4:9151    | (Ursula)⇀DarkMagenta Diamond Red Mountain↽ (0xcC0678E51a8237b762c09d6548d2d07285609e98)
54.169.230.89:9151   | (Ursula)⇀DarkViolet Hermes DodgerBlue Sharp↽ (0x250D36cbE375d104b4668a97303350A861A1491a)
54.254.231.79:9151   | (Ursula)⇀Wheat Gear WhiteSmoke Taurus↽ (0xA37FBC68ab044F1D6ECB18ec58BF0855d8b0d750)
95.217.4.43:9151     | (Ursula)⇀Khaki Virgo SlateBlue Leo↽ (0xA47f8D1Df610DC56DD523ec1Ac335392E0891B2c)
198.199.95.24:9151   | (Ursula)⇀DarkGray Heart Cyan Fleur-de-lis↽ (0x8b5CE9324069cEe2FCbA3382C18506cf55d4dD82)
104.248.133.94:9151  | (Ursula)⇀CadetBlue Key CadetBlue Rain↽ (0xd6C288d7494E425C6E436e907B0343e15440983C)
198.199.94.71:9151   | (Ursula)⇀Brown Airplane MediumOrchid Flower↽ (0x9259109490b1EE912e453cF321E0e01cD06F47a7)
128.0.51.144:9151    | (Ursula)⇀LightSlateGray Dharma PaleTurquoise Star↽ (0x7820aDA8554197e41Ff0FA54aF30BbB98c716765)
15.164.214.68:9151   | (Ursula)⇀Aqua Jupiter Beige Key↽ (0x6dda3D258e21Abc37668E82aE543e8C77567F220)
64.227.4.109:9151    | (Ursula)⇀AntiqueWhite Juno LawnGreen Bishop↽ (0xD9b6B55b005f1B23b45a9a4aC9669deFac6dAd67)
206.81.17.135:9151   | (Ursula)⇀Crimson King Ivory Alembic↽ (0x2160DCf3EAE12e21DCC1b7294859d80Ba2065589)
78.47.245.141:9151   | (Ursula)⇀Tan Flag LemonChiffon Virgo↽ (0x3B686e73F3c7D5E9e5F1fDFa4eF41Ca0e5E1f60A)
129.211.56.4:9151    | (Ursula)⇀YellowGreen Taurus Turquoise Pawn↽ (0x9705dF28eA06a38ad6042e37C546a9F01c1BE056)
128.0.51.143:9151    | (Ursula)⇀LightSkyBlue Queen OldLace Circle↽ (0xE78e5351857744e83fa85A17EE29657C574eEF36)
51.38.237.243:9151   | (Ursula)⇀LightSeaGreen Crossbones MediumBlue Shield↽ (0x1Ef837917c26a2f9aa31652A36aa7EA5aF7582c8)
```

还可以在线访问(把ip替换为你自己的):
https://13.229.53.90:9151/status

![在线节点](/.gitbook/assets/nucypher/在线节点.png)

### 停止节点

```
Ursula(0xe160672) >>> stop
```

# 13.提交表单，获取测试网激励

https://www.nucypher.com/incentivized-testnet/

其中的ETH地址是上面生成的两个eth地址(采用标准的ETH格式，不是转换后的！)