# 商业化产品中台

商业化产品中台是一系列基础服务和基础库，满足商业化产品的公共需求。

## 目标

* 提升开发速度
* 可复用
* 可扩展

## 策略

* 垂直方向：减少非必要的组件间依赖。大部分组件可以独立使用，也可以互相组合。
* 水平方向：分层架构，核心通用需求与个性化需求分离。

## 约定

为了便利客户端对接，作出以下约定。

### 响应结构：

```json
{
  "code": 0,
  "message": "", 
  "data": {}
}
```

```proto
message GenericReply {
  int32 code = 1;
  string message = 2;
  google.protobuf.Any data = 3;
}
```

#### 状态码 `code`

正确响应时返回 0。非0即为错误响应。0-99为预留通用状态码。100以上为非通用状态码，由各业务线自行定义。

当前预留状态码如下：

##### OK(0)：成功

操作成功完成

##### CANCELLED(1)：被取消

操作被取消（通常是被调用者取消）

##### UNKNOWN(2)：未知

未知错误。这个错误可能被返回的一个例子是，如果从其他地址空间接收到的状态值属于在当前地址空间不知道的错误空间（注：看不懂。。。）。此外，API发起的没有返回足够信息的错误也可能被转换到这个错误。

##### INVALID_ARGUMENT(3)：无效参数

客户端给出了一个无效参数。注意，这和 FAILED_PRECONDITION 不同。INVALID_ARGUMENT 指明是参数有问题，和系统的状态无关。

##### DEADLINE_EXCEEDED(4)：超过最后期限

在操作完成前超过最后期限。对于修改系统状态的操作，甚至在操作被成功完成时也可能返回这个错误。例如，从服务器返回的成功的应答可能被延迟足够长时间以至于超过最后期限。

##### NOT_FOUND(5)：无法找到

某些请求实体(例如文件或者目录)无法找到

##### ALREADY_EXISTS(6)：已经存在

某些我们试图创建的实体(例如文件或者目录)已经存在

##### PERMISSION_DENIED(7)：权限不足

调用者没有权限来执行指定操作。PERMISSION_DENIED 不可以用于因为某些资源被耗尽而导致的拒绝（对于这些错误请使用 RESOURCE_EXHAUSTED）。当调用者无法识别身份时不要使用 PERMISSION_DENIED （对于这些错误请使用 UNAUTHENTICATED）

##### RESOURCE_EXHAUSTED(8)：资源耗尽

某些资源已经被耗尽，可能是用户配额，或者可能是整个文件系统没有空间。

##### FAILED_PRECONDITION(9): 前置条件失败

操作被拒绝，因为系统不在这个操作执行所要求的状态下。例如，要被删除的目录是非空的，rmdir操作用于非目录等。

> 下面很容易见分晓的测试可以帮助服务实现者来决定使用 FAILED_PRECONDITION, ABORTED 和 UNAVAILABLE:
> * 如果客户端可以重试刚刚这个失败的调用，使用 UNAVAILABLE。
> * 如果客户端应该在更高级别做重试（例如，重新开始一个 读-修改-写 序列操作），使用 ABORTED。
> * 如果客户端不应该重试，直到系统状态被明确修复，使用 FAILED_PRECONDITION 。例如，如果 "rmdir" 因为目录非空而失败，应该返回 FAILED_PRECONDITION ，因为客户端不应该重试，除非先通过删除文件来修复目录。

##### ABORTED(10): 中途失败

操作中途失败，通常是因为并发问题如时序器检查失败，事务失败等。

##### OUT_OF_RANGE(11)：超出范围

操作试图超出有效范围，例如，搜索或者读取超过文件结尾。

> 和 INVALID_ARGUMENT 不同，这个错误指出的问题可能被修复，如果系统状态修改。例如，32位文件系统如果被要求读取不在范围[0,2^32-1]之内的offset将生成 INVALID_ARGUMENT，但是如果被要求读取超过当前文件大小的offset时将生成 OUT_OF_RANGE 。

> 在 FAILED_PRECONDITION 和 OUT_OF_RANGE 之间有一点重叠。当OUT_OF_RANGE适用时我们推荐使用 OUT_OF_RANGE （更具体的错误）.

##### UNIMPLEMENTED(12): 未实现

操作没有实现，或者在当前服务中没有支持/开启。

##### INTERNAL(13)：内部错误

内部错误。意味着某些底层系统期待的不变性被打破。如果看到这些错误，说明某些东西被严重破坏。

##### UNAVAILABLE(14)：不可用

服务当前不可用。这大多数可能是一个临时情况，可能通过稍后的延迟重试而被纠正。

##### DATA_LOSS(15)：数据丢失

无法恢复的数据丢失或者损坏。

##### UNAUTHENTICATED(16)：未经认证

请求没有操作要求的有效的认证凭证。

#### 消息 `Message`

当状态码是0时，Message可以被忽略。否则，Message可以向用户进行展示，解释出现的错误。Message中不应包括系统细节，比如错误栈等。

#### 内容 `Data`

由每个接口自行定义。

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

#### 数据同步

在用户中心对用户进行了任何新增和修改后，都对在 kafka 里写入一条用户完整信息。业务线如果有需要可以通过消费 kafka 来实时同步用户数据。

kafka 数据结构为如下 protobuf ：

```proto
message UserInfo {
  enum Gender {
    GENDER_UNKNOWN = 0;
    GENDER_MALE = 1;
    GENDER_FEMALE = 2;
  }
  message TaobaoExtra {
    string userid = 1;
    string open_sid = 2;
    string top_access_token = 3;
    string avatar_url = 4;
    string havana_sso_token = 5;
    string nick = 6;
    string open_id = 7;
    string top_auth_code = 8;
    string top_expire_time = 9;
  }

  message WechatExtra {
    string access_token = 1;
    int64 expiresIn = 2;
    string refresh_token = 3;
    string open_id = 4;
    string scope = 5;
    string nick_name = 6;
    int32 sex = 7;
    string province = 8;
    string city = 9;
    string country = 10;
    string headimgurl = 11;
    repeated string privilege = 12;
    string unionid = 13;
  }
  uint64 id = 1;
  string user_name = 2;
  string wechat = 3;
  // 头像地址
  string head_img = 4;
  Gender gender = 5;
  string birthday = 6;
  string token = 7;
  // 第三方ID
  string third_party_id = 8;
  bool is_new = 9 [(gogoproto.jsontag) = "is_new"];
  WechatExtra wechat_extra = 10;
  TaobaoExtra taobao_extra = 11;
  string mobile = 12;
}
```

用户结构会随着后续支持的登陆方式增多（如QQ登陆等）而逐渐扩展，但是会保持向后兼容。

Topic名称：

* common_user_info_channel_test (测试)
* common_user_info_channel (生产)

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

上述配置内容会在接口中输出如下内容：

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

配置`DSL`可以依照需求去不断扩展，实现灰度下发、地域定向、机型定向等高级功能。当前版本配置平台DSL只支持真值判断(`==`, `!=`)和逻辑判断(`&&`, `||`, `()`)。

#### 客户端消费配置

假设配置平台地址localhost:8080, 消费配置`/foo`（配置平台yaml对应路径）

GET http://localhost:8080/rule/v1/calculate/foo?package_name=xxx&channel=yyy&&version_code=zzz

> 配置平台接口不需要`jwt`,因此可以在启动阶段调用。

`Querystring` 中的内容主要用于高级配置中的`DSL`进行判断。

完整的 `Querystring` 字段包括：

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

Go语言可以直接使用 `glab.tagtic.cn/ad_gains/kitty/pkg/rule/client`

```go
package main

import (
	"context"
    "fmt"

    "glab.tagtic.cn/ad_gains/kitty/pkg/rule/client"
    "go.etcd.io/etcd/clientv3"
)

func main() {
	etcd, _ := clientv3.New(clientv3.Config{
		Endpoints: []string{"etcd-1:2379", "etcd-2:2379", "etcd-3:2379"},
		Context:   context.Background(),
	})
	engine, _ := client.NewRuleEngine(
        client.WithClient(etcd), 
        client.Rule("kitty-testing"),
        client.Rule("whatever"), // 使用任何rule必须先在这里注册, 可以注册多个
    )
    go engine.Watch(context.Background())
	reader, _ := engine.Of("kitty-testing").Payload(&rule.Payload{}) // 配合DSL高级配置
	fmt.Println(reader.String("foo")) //bar
}
```

