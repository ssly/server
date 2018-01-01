## task api
1. 接口均符合：Content-Type=application/json
2. 特殊说明：返回字段中的`code`为，0 成功，其他为失败的类型


### 获取所有任务列表
#### URL
GET /task/manager
#### 请求格式
```
/task/manager?finish=true&type=1&difficult=1&minHours=1&maxHours=8
```
#### 请求参数（URL param）
| 参数 | 可选值 | 说明 |
| :- | :- | :- |
| finish | true/false | 是否完成 |
| type | 1/2/3/4 | 类型，紧急重要/紧急不重要/不急重要/不急不重要 |
| difficult | 1/2/3 | 难度，困难/中等/容易 |
| minHours | 0,1,... | 需要时间的最小值 |
| maxHours | 0,1,... | 需要时间的最大值 |
#### 返回格式
```json
{
    "success": true,
    "code": 0,
    "data": [
        {
            "id": "599af5182e02e989cfe19048",
            "name": "测试信息1",
            "type": 2,
            "difficult": 1,
            "deadline": "2017-09-10",
            "hours": 4,
            "finish": false,
            "memo": "测试数据1"
        },
        {
            "id": "5a37da1bf860a92924b3b3a9",
            "name": "测试信息2",
            "type": 1,
            "difficult": 3,
            "deadline": "2017-09-10",
            "hours": 8,
            "finish": true,
            "memo": "测试数据2"
        }
    ]
}
```

### 获取单个任务
#### URL
GET /task/manager/:id
#### 请求格式
```
/task/manager/599af5182e02e989cfe19048
```
#### 返回格式
```json
{
    "success": true,
    "code": 0,
    "data": {
        "id": "599af5182e02e989cfe19048",
        "name": "测试信息1",
        "type": 2,
        "difficult": 1,
        "deadline": "2017-09-10",
        "hours": 4,
        "finish": false,
        "memo": "测试数据1"
    }
}
```

### 增加单个任务
#### URL
POST /task/manager
#### 请求参数
| 参数 | 类型 | 可选值 | 说明 |
| :- | :- | :- | :- |
| name | string | - | 必填 |
| type | number | 1/2/3/4 | 选填，默认为1，类型，紧急重要/紧急不重要/不急重要/不急不重要 |
| difficult | number | 1/2/3 | 选填，默认为2，难度，困难/中等/容易 |
| deadline | number | 0,1,... | 选填，默认为当天，截止至 |
| finish | boolean | true/false | 选填，默认为false，是否完成 |
| hours | number | 0,1,... | 选填，默认为8，需要时限 |
| memo | string | - | 选填，备注 |
#### 请求格式（Body体）
```json
{
    "name": "待添加数据1",
    "type": 1,
    "difficult": 1,
    "deadline": "2017-09-10",
    "hours": 8,
    "finish": false,
    "memo": "备注xinxi"
}
```
#### 返回格式
```json
{
    "success": true,
    "code": 0,
    "data": null
}
```

### 修改单个任务（暂未实现）
#### URL
PUT /task/manager/:id 或<br>
PUT /task/manager
#### 请求参数
同**增加单个任务**
#### 请求格式（Body体）
> URL中带了id则body体可不带id字段
```json
{
    "id": "599af5182e02e989cfe19048",
    "name": "待添加数据1",
    "type": 1,
    "difficult": 1,
    "deadline": "2017-09-10",
    "hours": 8,
    "finish": false,
    "memo": "备注xinxi"
}
```
### 返回格式
```json
{
    "success": true,
    "code": 0,
    "data": null
}
```

### 删除多个任务
#### URL
DELETE /task/manager
#### 请求格式
```json
["599af5182e02e989cfe19048", "5a37da1bf860a92924b3b3a9"]
```
### 返回格式
```json
{
    "success": true,
    "code": 0,
    "data": null
}
```

### 删除单个任务
#### URL
DELETE /task/manager/:id
#### 请求格式
```
/task/manager/599af5182e02e989cfe19048
```
#### 返回格式
```json
{
    "success": true,
    "code": 0,
    "data": null
}
```

