from socket import *

tcpSerPort = 8899
tcpSerSock = socket(AF_INET, SOCK_STREAM)

tcpSerSock.bind(('', tcpSerPort))
tcpSerSock.listen(5)

while True:
    print('Ready to serve...')
    tcpCliSock, addr = tcpSerSock.accept()
    print('Received a connection from: ', addr)
    message = tcpCliSock.recv(4096).decode()

    filename = message.split()[1].partition("//")[2].replace('/', '_')
    fileExist = "false"
    try:
        f = open(filename, "r")
        outputdata = f.readlines()
        fileExist = "true"
        print('File Exists')

        for i in range(0, len(outputdata)):
            tcpCliSock.send(outputdata[i].encode())
        print('Read from cache')

    except IOError:
        print('File Exist: ', fileExist)
        if fileExist == "false":
            print('Creating socket on proxyserver')
            c = socket(AF_INET, SOCK_STREAM)

            hostn = message.split()[1].partition("//")[2].partition("/")[0]
            print('Host name: ', hostn)
            try:
                c.connect((hostn, 80))
                print('Socket connected to port 80 of the host')

                c.sendall(message.encode())
                buff = c.recv(4096)

                tcpCliSock.sendall(buff)
                tmpFile = open("./" + filename + "w")
                tmpFile.writelines(buff.decode().replace('\r\n', '\n'))
                tmpFile.close()

            except:
                print("Illegal request")

        else:
            print('File Not Found...Stupid Andy')

    tcpCliSock.close()
tcpSerSock.close()