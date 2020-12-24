#!/bin/bash

printf "%-5s %-10s %-4s\n" no name mark

# 查看某个应用的进程号pid
pgrep docker
12384

# 查看某个进程的环境变量,并格式化一下
cat /proc/$PID/environ | tr '\0' '\n'

# 添加环境变量
export PATH="$PATH:/home/user/bin"

# 获取变量长度，和Lua的#不谋而合
length=${#var}

# 判断是否是超级用户
if [$UID -ne 0]; then 
    echo Non root user
else
    echo Root user
fi

# 运算
no1=1
no2=2
let res1=no1+no2
echo $res
let no2++
echo $no2
let no1+=2
echo $no1
res2=$[no1+no2]

# bc 数学运算
echo "3*4" | bc
# scale 控制精度，小数点后2位
echo "scale=2;33/7" | bc

# alias
alias rm='cp $@ ~/backup && rm $@'
# \转义将不再执行别名命令，转而执行原本命令，如果原本没有就no command
\rm

# 禁止回显输入密码
#!/bin/bash
echo -e "enter password"
stty -echo #禁止回显	
read password
stty echo  #允许回显
echo $password

# 日期格式
# 工作日（weekday）
# %a （例如：Sat）
# %A （例如：Saturday）

# 月
# %b （例如：Nov）
# %B （例如：November）

# 日
# %d （例如：31）

# 特定格式日期（mm/dd/yy）
# %D （例如：10/18/10）

# 年
# %y （例如：10）
# %Y （例如：2010）

# 小时
# %I 或%H （例如：08）

# 分钟
# %M （例如：33）

# 秒
# %S （例如：10）

# 纳秒
# %N （例如：695208515）

# Unix纪元时（以秒为单位）
# %s （例如：1290049486）


#!/bin/bash
#文件名: sleep.sh
echo Count:
tput sc     # 保存光标位置

# 循环40秒
for count in `seq 0 4`
do
  tput rc   # 恢复光标位置
  tput ed   # 清除到光标的所有内容
  echo -n $count
  sleep 1
done

function DEBUG()
{
    ["$_DEBUG"=="on"] && $@ || :
}
for i in {1..10}
do
    DEBUG echo "I is $i"
done

fname()
{
    echo $1, $2;
    echo "$@"
    echo "$*"
    return 0
}

# /etc/security/limits.conf中的nproc 来限制可生成的最大进程数

# while 内置避免重复新建进程
repeat() { while :; do $@ && return; sleep 30; done }
repeat wget -c http://www.baidu.com

#!/bin/bash
#用途: 演示IFS的用法
line="root:x:0:0:root:/root:/bin/bash"
oldIFS=$IFS;
IFS=":"
count=0
for item in $line;
do

     [ $count -eq 0 ]  && user=$item;
     [ $count -eq 6 ]  && shell=$item;
    let count++
done;
IFS=$oldIFS
echo $user shell is $shell;

# [ -f $file_var ] ：如果给定的变量包含正常的文件路径或文件名，则返回真。
# [ -x $var ] ：如果给定的变量包含的文件可执行，则返回真。
# [ -d $var ] ：如果给定的变量包含的是目录，则返回真。
# [ -e $var ] ：如果给定的变量包含的文件存在，则返回真。
# [ -c $var ] ：如果给定的变量包含的是一个字符设备文件的路径，则返回真。
# [ -b $var ] ：如果给定的变量包含的是一个块设备文件的路径，则返回真。
# [ -w $var ] ：如果给定的变量包含的文件可写，则返回真。
# [ -r $var ] ：如果给定的变量包含的文件可读，则返回真。
# [ -L $var ] ：如果给定的变量包含的是一个符号链接，则返回真。

# 文件类型	类型参数
# 普通文件	f
# 符号链接	l
# 目录	d
# 字符设备	c
# 块设备	b
# 套接字	s
# FIFO	p
find . -type f

# 改变所有者    注意该命令结尾的\; 。必须对分号进行转义，否则shell会将其视为find 命令的结束，而非chown 命令的结束。
find . -type f -user root -exec chown slynux {} \;

# 执行命令，+会让find搜索一个完整列表后一次性作为参数执行 {} 表示一个命令，都需要加上
fine . -type f -name '*.c' -exec cat {} >all_c_files.txt +

# xargs 接受一个stdin参数
find /smbMount -iname '*.docx' -print0 | xargs -0 grep -L image

# 匹配删除
find . -type f -name "*.txt" -print0 | xargs -0 rm -f

# 匹配删除
echo "Hello 123 world 456" | tr -d '0-9'

# 删除0-9之外的，从输入文本中删除不在补集中的所有字符
echo hello 1 char 2 next 4 | tr -d -c '0-9 \n' 

# 替换
echo hello 1 char 2 next 4 | tr -c '0-9' ' '

# salt
openssl passwd -1 -salt SALT_STRING PASSWORD

count=1
for img in `find . iname '*.png' -or -iname '*.jpg' -type f -maxdepth 1`
do
    new=image-$count.${img##*.}
    echo renaming 
    mv "$img" "$new"
    let count++
done
