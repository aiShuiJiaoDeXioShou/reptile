go build # 先把应用打包一下

rm -rf build

mkdir build

# 移动编译好的文件
cp -R reptile.exe static templates build

# 更新git版本
git add .
git commit -m "shell自动更新"
git push -u origin