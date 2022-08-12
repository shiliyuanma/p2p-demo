## 传统p2p穿透(hole-punch)的原型:

> 通过中继节点(relay)协调打洞，主要是交换要打洞的两个对等方(peer)的public addr
> 如果两个对等方的nat类型都是对称nat，几乎无法打通，其他情况均可.
> ./p2p-demo nat //检测本机nat类型

## 启动:

1.启动中继节点(relay):   
    `./p2p-demo hole relay` //默认监听6002端口

2.启动节点A(provider):   
    `./p2p-demo hole peer --relay "103.44.247.16:6002" --role p2pp`  

3.启动节点B(visitor):  
    `./p2p-demo hole peer --relay "103.44.247.16:6002" --role p2pv`  
