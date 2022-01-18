import pandas as pd
import sys
import os
import pymysql
import sqlalchemy
import argparse
import logging
import datetime
import time
import yaml
import re



'''
命令行获取扫描结果文件
'''
def getpath():
    parser=argparse.ArgumentParser(description="获取命令")
    parser.add_argument("path",help="获取扫描结果文件")
    args=parser.parse_args()
    path=args.path
    return path

'''
检查文件名是否存在
'''
def checkpath(path):
    if os.path.exists(path):
        return path
    else:
        sys.exit("{}文件不存在".format(path))



'''
通过 pandas 读取解析
'''
def read_frompath(path):
    df=pd.read_table(path,usecols=[0,1,2,3,4,5,6],dtype={"host":str,"tag":str,"ip":str,"port/protocal":str,"status":str,"service":str})
    df=df.astype(str)
    return df



'''
读取某一列
host	tag	ip	port/protocal	status	service	remark/version
'''
def read_fromcolumn(dataframe,column):
    col=dataframe.columns[column]
    return dataframe[col]


'''
读取配置文件
'''
def readfrom_yaml(path):
    try:
        with open(path,"r",encoding='utf-8') as f:
            content=f.read()
            config=yaml.load(content,Loader=yaml.CLoader)
            return config
    except:
        sys.exit("未找到配置文件")


'''
读取白名单
'''
def readwhite(path):
    try: 
        with open(path,"r",encoding='utf-8' ) as f:
            try: 
                content=f.read()
                whitelist=yaml.load(content,Loader=yaml.CLoader)
                try:
                    len(whitelist)
                    return whitelist
                except:
                    print("白名单配置文件格式有误，无法读取")
                    return ""
            except:
                print("白名单配置文件格式有误，无法读取")
                return ""
    except:
        print("未找到白名单配置文件")
        return ""


'''
定义主函数
'''
def main_write_sql():
    path=checkpath(getpath())
    df=read_frompath(path)
    statuslist=read_fromcolumn(df,4)
    name=time.strftime("%Y-%m-%d",time.localtime())+"_Scan_result"
    # 读取配置文件信息
    configname=["config","config","cnf","cnf","con","con"]
    whitelistname=["wl","WL","whitelist"]
    suffixname=["yml","yaml"]
    configpath=""
    whitelistpath=""
    for eachconf in configname:
        for eachsuffix in suffixname:
            x=0
            tmp=eachconf+"."+eachsuffix
            if os.path.exists(tmp):
                x=1
                configpath=tmp
                break
        if x==1:
            break
    for eachwlname in whitelistname:
        for eachsuffix in suffixname:
            x=0
            tmp2=eachwlname+"."+eachsuffix
            if os.path.exists(tmp2):
                x=1
                whitelistpath=tmp2
                break
        if x==1:
            break
    config=readfrom_yaml(configpath)
    whitelist=readwhite(whitelistpath)
    # 整理 mysql 的链接信息
    user=config.get("mysql").get("username")
    password=config.get("mysql").get("password")
    host=config.get("mysql").get("host")
    port=config.get("mysql").get("port")
    database=config.get("mysql").get("database")
    enginestart="mysql+pymysql://"+user+":"+password+"@"+host+":"+str(port)+"/"+database
    engine=sqlalchemy.create_engine(enginestart)
    # 整理白名单配置
    if whitelist!="":
        whitelist=readwhite(whitelistpath).get("whitelist")
    # 查询语句
    # 建立写入数据库的 dataframe
    # 字段分别为主机，项目，ip ，端口，服务，备注，是否白名单，发现时间，处理方式，处理时间
    to_database=pd.DataFrame(columns=('host','tag','ip','port','service','remark','whitelist','scantime','processmode','processfinish'))
    with open("deal_result.csv","w+")as write:
        write.write('host\ttag\tip\tport\tservice\tremark\n')
        for index in df.index:
            if statuslist[index]=="open":
                data=df.loc[index].values
                scantimestamp=time.strftime("%Y-%m-%d %H:%m",time.localtime())
                tmpport=re.findall("\d+",data[3])[0]
                for element in whitelist:
                    if data[2]==element.get("ip") and tmpport==str(element.get("port")):
                        ifwhitelist=True
                        break
                    else:
                        ifwhitelist=False
                try:
                # 写入文件
                    write.write("\t".join(df.loc[index].values[0:-1]))
                    write.write("\n")
                except:
                    for each in df.loc[index]:
                        print(each,type(each))
                try:
                    tmpdict={'host':data[0],'tag':data[1],'ip':data[2],'port':tmpport,'service':data[5],'remark':data[6],'whitelist':ifwhitelist,'scantime':scantimestamp,'processmode':"","processfinish":""}
                    to_database=to_database.append(tmpdict,ignore_index=True)
                except:
                    print(index,"行无法追加")
                    print(df.loc[index].values)
    print("Total records {}".format(len(to_database)))
    try:
        to_database.to_sql(name,engine)
        print("成功写入数据库")
        print("表名为",name)
    except:
        logging.exception(Exception)
        print("无法写入数据库。")





'''
调用主函数
'''
if __name__=="__main__":
    main_write_sql()

