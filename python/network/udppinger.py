from socket import *
import time

serverName = '191.101.232.165' # 服务器地址，本例中使用一台远程主机
serverPort = 12000
clientSocket = socket(AF_INET, SOCK_STREAM)
clientSocket.settimeout(1)

for i in range(0, 10):
    sendTime = time.time()
    message = ('ping %d %s' % (i+1, sendTime)).encode()
    try:
        clientSocket.sendto(message, (serverName, serverPort))
        modifiedMessage, serverAddress = clientSocket.recvfrom(1024)
        rtt = time.time() - sendTime
        print('Sequence %d: Ready from %s  RTT = %.3fs' % (i+1, serverName, rtt))
    except Exception as e:
        print('Sequence %d: Request timed out' % (i+1))

clientSocket.close()