go build # 先把应用打包一下
cp reptile.exe build

#! /bin/bash
function Copy_H_File_To_PrjInclude() {
	for file in `ls $1`
	do
		if [ -d $1"/"$file ];then # 判断是否是目录，是目录则递归
			Copy_H_File_To_PrjInclude $1"/"$file

		elif [ -f $1"/"$file ];then # 判断是否是文件		
			str=1"/"$file
			cp $1"/"$file build #这里是目标路径
			echo $1"/"$file
		fi
	done
}

SourcePath = 'static' #源路径
Copy_H_File_To_PrjInclude $SourcePath  

SourcePath2 = 'templates'
Copy_H_File_To_PrjInclude $SourcePath2