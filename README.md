# 商业化产品中台

商业化产品中台是一系列基础服务和基础库，满足商业化产品的公共需求。

## 目标

* 提升开发速度
* 可复用
* 可扩展

## 策略

* 垂直方向：减少非必要的组件间依赖。大部分组件可以独立使用，也可以互相组合。
* 水平方向：分层架构，核心通用需求与个性化需求分离。

## 组件

### 用户中心

用户中心实现了基于 `jwt` 的无状态登陆。

登录成功后接口返回 token 字段即为 [jwt](jwt.io)。`jwt` 同时具有两种功能。其一是为验证用户身份，其二是携带部分用户信息。
客户端登陆后，请求业务接口时，添加http header `Authorization`， 内容为 `bearer xxxxx` （将 `xxxxx` 替换为 `jwt`）。
服务端验证用户身份时，首先判断 header 中的 `jwt` 可以正确的被密钥解开，然后判断过期时间未超期，则认为 `jwt` 真实有效。解开后获取的 json 内容，即为 `jwt` 携带的用户信息。
这部分用户信息可以商业化平台组件间的数据追踪，数据关联。例如，`jwt` 中包含的 `wechat` 字段，即为微信 openid，提现服务可以向此 openid 返现。

`jwt` 基本数据结构：

```go
type Claim struct {
	stdjwt.StandardClaims
	PackageName string `json:"PackageName,omitempty"`
	UserId      uint64 `json:"UserId,omitempty"`
	Suuid       string `json:"Suuid,omitempty"`
	Channel     string `json:"Channel,omitempty"`
	VersionCode string `json:"VersionCode,omitempty"`
	Wechat      string `json:"Wechat,omitempty"`
	Mobile      string `json:"Mobile,omitempty"`
}
```

json 形式如下：

```json
{"exp":1605232490,"iat":1605146090,"iss":"signCmd","PackageName":"com.donews.www","UserId":1,"Suuid":"","Channel":"","VersionCode":"","Wechat":"","Mobile":""}
```

`jwt` 签名算法：
* HS256

`jwt` 密钥：
* 参见配置文件

### 配置中心

配置中心提供了简单的配置下发以及灰度下发等功能。

[配置中心后台](http://monetization-config.xg.tagtic.cn/) 提供了简单快速的GUI配置页面和配置下发接口，并且实现了 "`自举`"， 即通过配置平台来配置配置平台。

#### 创建一项新的配置

点击右上角的齿轮图标。进入配置中心的配置页。编辑如下`yaml`结构并保存。

```yaml
style: basic
rule:
  list:
    - name: 商业化平台
      icon: home2
      children:
        - name: 用户体系
          path: /kitty
          id: user
        - name: 积分体系
          path: /score
          id: score
    - name: 活动
      icon: home2
      children:
        - name: 砸金蛋
          path: /egg
          id: egg
```

保存后刷新页面，配置中心会如 list 字段所示，左侧菜单结构变为：

```
+-- 商业化平台
|   +-- 用户体系
|   +-- 积分体系
+-- 活动
|   +-- 砸金蛋
```

创建项目后，每个项目可以在配置平台对应页面进行编辑。

#### 创建基本配置

基本配置：yaml内容里的rule字段即为下发的json结构。

```yaml
style: basic
rule:
  foo: bar
```

上述配置内容会在接口中输出如下内容

```json
{"foo": "bar"}
```

#### 创建高级配置

高级配置：yaml中的内容为一种配置 [DSL]( https://zh.wikipedia.org/wiki/%E9%A2%86%E5%9F%9F%E7%89%B9%E5%AE%9A%E8%AF%AD%E8%A8%80 )，计算后得出下发json

```yaml
style: advanced
rule:
  - if: PackageName == "com.infinities.reward.shopping"
    then:
      foo: bar
  - if: true
    then:
      foo: baz
```

当包名为"com.infinities.reward.shopping"时，上述配置内容会在接口中输出如下内容

```json
{"foo": "bar"}
```

其他包名，上述内容将会输出

```json
{"foo": "baz"}
```

配置`DSL`可以依照需求去不断扩展，实现灰度下发、地域定向、机型定向等高级功能。目前配置平台DSL只支持真值判断和逻辑判断。

#### 客户端消费配置

假设配置平台地址localhost:8080, 消费配置`/foo`（配置平台yaml对应路径）

GET localhost:8080/rule/v1/calculate/foo?package_name=xxx&channel=yyy

`Querystring` 中的内容主要用于高级配置中的`DSL`进行判断。

完整的可用query配置包括：

```
	string "version_code" 版本号
	string "channel" 渠道
	uint8  "os" 1为ios，2为安卓
	uint64 "user_id" 用户id
	string "imei" IMEI
	string "idfa" IDFA
	string "oaid" OAID
	string "suuid" SUUID地址
	string "mac" MAC地址
	string "android_id" 安卓ID
	string "package_name" 包名
```

#### 服务端消费配置

异构语言可以使用HTTP方式消费配置，同客户端。




