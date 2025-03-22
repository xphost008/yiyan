# 一个很不错的一言网站源码

### 须知

1. 本网站仅用于我的学校【武汉工程职业技术学院】举办活动所需要。
2. 为了使得下一届乃至下下一届的学弟学妹能也能使用到本仓库，特此开源~
3. 本网站如何安装，请看下面：

### 安装方式

1. 这里推荐使用GoLand进行开发、MySQL进行数据库连接。（如果你使用别的技术，请参阅官方文档）
2. 首先安装MySQL语言，如何安装MySQL数据库，请自行百度。并导入一个数据库，名为【yiyandata】，创建表的命令如下：

```sql
CREATE DATABASE IF NOT EXISTS yiyandata CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci
```

3. 创建好数据库后，使用`use yiyandata`进入到数据库内。
4. 紧接着，再创建4个表。创建4个表的代码如下：

```sql
create table yiyan
(
    id          int auto_increment comment '提交的id值'
        primary key,
    submitter   varchar(11)                                                                                                          not null comment '提交者',
    content     varchar(255)                                                                                                         not null comment '名言',
    source      varchar(255)                                                                                                         not null comment '来源',
    author      varchar(255)                                                                                                         not null comment '作者',
    classifiers enum ('anime', 'comic', 'game', 'literature', 'myself', 'internet', 'other', 'video', 'poem', 'ncm', 'philosophy', 'funny') not null comment '分类'
)
    comment '名人名言';
create table users
(
    id       varchar(11)  not null comment '学号'
        primary key,
    username varchar(255) not null comment '名字',
    password varchar(255) not null comment '密码'
)
    comment '用户表';
create table like_record
(
    yiyan_id int         not null,
    user_id  varchar(11) not null,
    primary key (user_id, yiyan_id),
    constraint like_record_ibfk_1
        foreign key (user_id) references users (id),
    constraint like_record_ibfk_2
        foreign key (yiyan_id) references yiyan (id)
);
create table admin(
    id int not null comment '管理员id'
    primary key,
)
create index yiyan_id
    on like_record (yiyan_id);

```

5. 其中，第一个创建的表是存储用户的一言表，第二个创建的表是存储用户的表，第三个表是存储点赞的表。第4个表是管理员表，用来存储所有管理员账号的id。
6. 由于本网站使用`邀请制`，并没有提供`注册按钮`键。因此表设计得挺简单得。
7. 使用`GoLand`导入本工程，然后自动下载`go.mod`下的内容。
8. 之后你可以在本网站的隐私条款目录下，自行写你自己的内容。并且直接运行你的`main.go`即可。所有网站该有的功能均在你用浏览器访问`localhost:8080`即可看到，这里不再赘述。

### 使用须知

1. 本软件已经开源，如果你想要使用本软件做点事情，你不能做违法的事！本软件虽然说并不仅限【武汉工程职业技术学院】使用，但也不代表你可以随意使用源码。
2. 本软件仅适合使用自定义frp开设公网让同学进行访问，不建议直接使用域名。【因为该死的bing会自动收录你的域名，这就很💩】
3. 本软件使用`No License`发布，因为我不希望有任何人复制我的源代码到处传播，也不希望有任何人对我的代码进行二次修改后再次分发。
4. 本软件其实做出来的初衷仅是对【武汉工程职业技术学院】单独使用，而不是发布到github，但是我们【计算机协会】有成员希望将其开源，因此我才开源至github。

### 致谢

衷心感谢[@Red-Feng](https://github.com/Red-Feng)在项目全过程中提供的全方位支持。在技术攻关阶段，您在算法优化和数据验证环节给予的关键性建议，使研究得以突破瓶颈。同时其在在需求分析和测试验证中的协同配合，你的多维度视角为方案完善提供了重要启发。