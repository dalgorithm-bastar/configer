# 1.基本语法

## 1.1 填充语法

在文件中使用双重花括号{{}}，表示其中的内容需要按照go template的规则替换

## 1.2 取值和作用域

使用"."取结构体当前作用域中的值，使用with、range等关键字改变当前作用域，配合使用end关键字退出当前作用域

示例： 源结构体：

```json
{
  "replicator_number": "2",
  "deployment_info": [
    {
      "NODE_ID": "1",
      "hostname": "hostforMCtest1",
      "BIND_IP": "{ip}"
    },
    {
      "NODE_ID": "2",
      "hostname": "hostforMCtest2",
      "BIND_IP": "{ip}"
    }
  ]
}
```

模板：

```json
{{
  .replicator_number
}}
{{range .deployment_info}}
{{.NODE_ID}
}
{
{
end
}
}
```

填充结果：

```json
2

1

2

```

## 1.3 去除空白

在{{}}中添加"-"来去除填充部分前方或后方的空白。需要注意的是关键字也会产生换行符，而且循环关键字range在每次迭代时都会产生换行。

示例：

模板：

```json
{{
  .replicator_number
}}
{{range .deployment_info -}}
{{.NODE_ID
}
}
{
{
end
}
}
```

填充结果：

```json
2
1
2

```

## 1.4 定义变量

使用如下语句定义变量：

```shell
$i:=1
$j:=(function)
```

其中值可以为常量，也可以是函数的返回值，变量的作用域跟随所处语句的作用域。使用变量时，在变量名前加$即可。

# 2.条件判断和循环

## 2.1 条件语句

使用if、else、end关键字实现条件判断

示例： 模板：

```json
{{
  $i: =1
  -
}}
{{range .deployment_info -}}
{{if eq i .NODE_ID}}         #注意"."的作用域变化
{{
.NODE_ID
}
}
{
{
end
}
}
```

填充结果：

```json
1

```

## 2.2 循环语句

循环语句由range关键字和end关键字组成，有多种使用方法，这里介绍两种：

- 仅取值，见上述各例

- 同时取索引和值，用法如下：

```shell
{{range $index,$val:= .deployment_info}}
```

根据数据类型的不同，index分别为数组索引值或map的键值，val为值

# 4.函数调用

## 4.1 内置函数

内置函数类型均返回bool值，常配合条件语句使用，常用函数如下：

```shell
index
    对可索引对象进行索引取值。第一个参数是索引对象，后面的参数是索引位。
    "index x 1 2 3"代表的是x[1][2][3]
not
    布尔取反。只能一个参数
len
    返回参数的length
eq arg1 arg2：
    arg1 == arg2时为true
ne arg1 arg2：
    arg1 != arg2时为true
lt arg1 arg2：
    arg1 < arg2时为true
le arg1 arg2：
    arg1 <= arg2时为true
gt arg1 arg2：
    arg1 > arg2时为true
ge arg1 arg2：
    arg1 >= arg2时为true
```

## 4.2 配置中心预设函数

见配置中心使用文档

## 4.3 函数嵌套

使用一个空格将函数的多个参数分隔，使用()实现函数嵌套。

示例： 模板：

```json
{{
  $i: =1
  -
}}
{{if ne (index .deployment_info 1) i}}
{
{
end
}
}
```
