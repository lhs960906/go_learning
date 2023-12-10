# 一、访问关系型数据库

该教程为 Go 访问数据库的一个概览。

使用 Go，您可以将各种（a variety of）数据库和数据访问方法合并（incorporate）到您的应用程序中。这篇章节的主题就是如何使用标准库 `database/sql` 去访问关系型数据库。

一个介绍性的用 Go 访问数据的教程，请参阅：[Tutorial: Accessing a relational database](https://go.dev/doc/tutorial/database-access)。

Go 还支持其他数据访问的技术，包括更高级别的关系型数据库的访问 ORM 库，以及非关系 NoSQL 数据存储。 

* 对象关系映射（ORM）库。标准库中的 `database/sql` 包提供了低级的数据访问逻辑，您还可以用 Go 在更高的抽象等级去访问数据存储。关于 Go 的两个受欢迎的对象关系映射（ORM）库，请参阅：[GORM](https://gorm.io/index.html) ([package reference](https://pkg.go.dev/gorm.io/gorm)) 和 [ent](https://entgo.io/) ([package reference](https://pkg.go.dev/entgo.io/ent))
* 非关系行数据存储。Go 社区有已经开发好的主流 NoSQL 数据存储驱动，包括[MongoDB](https://docs.mongodb.com/drivers/go/) and [Couchbase](https://docs.couchbase.com/go-sdk/current/hello-world/overview.html)。更多信息请搜索：[pkg.go.dev](https://pkg.go.dev/)

## 1.1、支持的数据库管理系统

Go 支持绝大多数的通用关系型数据库管理系统，包括：MySQL，Oracle，Postgres，SQL Server，SQLite 以及更多。

您可以在 [SQLDrivers](https://github.com/golang/go/wiki/SQLDrivers) 找到完整的驱动列表。

## 1.2、执行查询和数据库更改的方法

`database/sql` 包包括一系列设计用来执行数据库操作的特定函数。例如，你可以使用 Query 或 QueryRow 去执行查询，QueryRow 被设计用于您只需要单行数据的场景，省略（omitting）返回仅包含单行开销（overhead）的 sql.Rows的开销。您可以使用 Exec 函数执行例如 INSERT，UPDATE 或 DELETE 语句以改变数据库。

更多信息请参阅：

* [Executing SQL statements that don’t return data](https://go.dev/doc/database/change-data)
* [Querying for data](https://go.dev/doc/database/querying)

## 1.3、事务

通过 sql.Tx，你可以编写在一个事务中执行数据库操作的代码。在一个事务中，多个操作能够一起被执行（performed）并且以一个最终的 commit 结束，在一个原子步骤中要么应用所有更改，要么丢弃它们。

更多事务相关信息，参阅：[Executing transactions](https://go.dev/doc/database/execute-transactions)

## 1.4、查询取消

当您希望可以取消一个数据库操作，可以使用 context.Context，例如当客户端的连接或一个操作比您预计执行的时间更长。

对于任何数据库操作，您可以使用将 Context 作为参数的 `database/sql` 包函数。 使用 Context，您可以为操作指定超时或截止日期。 您还可以使用 Context 通过应用程序将取消请求传播到执行 SQL 语句的函数，确保资源在不再需要时得到释放。

更多信息，请参阅：[Canceling in-progress operations](https://go.dev/doc/database/cancel-operations)

## 1.5、管理连接池

当您使用 `sql.DB` 数据库句柄，您正在连接一个内置的连接池，该连接池会根据您的代码需要创建和处理连接。通过 sql.DB 的句柄是使用 Go 进行数据库访问的最常见方式。更多信息请参阅：[Opening a database handle](https://go.dev/doc/database/open-handle)。

`database.sql` 包为您管理了连接池，但是如果有更高级的需求，你可以按 [Setting connection pool properties](https://go.dev/doc/database/manage-connections#connection_pool_properties) 中的说明设置数据库的连接池参数。

对于那些需要单个保留连接的操作，database/sql 包提供了 sql.Conn。 当使用 sql.Tx 的事务不是一个好的选择时，Conn 尤其有用。

例如，您的代码可能需要：

* 通过 DDL 进行模式更改，包括包含其自身事务语义的逻辑。 将 sql 包事务函数与 SQL 事务语句混合是一种不良做法，如 [Executing transactions](https://go.dev/doc/database/execute-transactions) 中所述
* 执行创建临时表的查询锁定操作

更多信息，请参阅：[Using dedicated connections](https://go.dev/doc/database/manage-connections#dedicated_connections)。

# 二、教程：访问一个关系型数据库

这篇文章介绍用 Go 访问访问数据库的基本方法以及标准库中的 `database/sql` 包。

如果您对 Go 及其工具有基本的了解，您将充分利用本教程。如果这是您第一次接触 Go，请参阅教程：[Tutorial: Get started with Go](https://go.dev/doc/tutorial/getting-started) 以获取快速介绍。

您将使用到 `database/sql` 包中连接到数据库的类型和函数、执行事务、取消正在进行的操作等等。更多使用该包的信息见 [Accessing databases](https://go.dev/doc/database/index)。

在这边教程中，你将创建一个数据库，然后去编写代码访问这个数据库。您的示例项目将是有关 ''复古爵士乐" 唱片的数据存储库。

在本教程中，您将逐步（progress through）完成以下部分：

* 为你的代码创建一个目录
* 建立（set up）一个数据库
* 导入数据库驱动
* 获取一个数据库句柄和连接
* 查询多行
* 查询单行
* 添加数据

## 2.1、前置条件

* 安装 MySQL 关系型数据库管理系统
* 安装 Go。安装建议见 [Installing Go](https://go.dev/doc/install)

* 编写代码的工具。任何的文本编辑器都可以
* 命令行终端。Go 在 Linux 和 Mac 终端上都能够工作得很好，以及 Windows 的 cmd 和 PowerShell 终端

## 2.2、为你的代码创建一个目录

首先（to begin），为你要写的代码创建一个目录。

1. 打开命令提示符并到您的家目录中

   On Linux or Mac：

   ```
   $ cd
   ```

   On Windows：

   ```cmd
   C:\> cd %HOMEPATH%
   ```

2. 从命令提示符为你的代码创建一个名为 data-access 目录 

   ```shell
   $ mkdir data-access
   $ cd data-access
   ```

3. 创建一个模块，您可以在其中管理将在本教程中添加的依赖项

   ```shell
   $ go mod init github.com/lhs960906/data-access
   go: creating new go.mod: module github.com/lhs960906/data-access
   ```

   这个命令创建一个 go.mod 文件，其中（in which）将列出您添加的依赖项以供跟踪。更多信息参阅：[Managing dependencies](https://go.dev/doc/modules/managing-dependencies)

   > Note：在实际开发中，您会指定一个更符合您自己需求的模块路径。有关更多信息，请参阅：Managing dependencies

接下来，您将创建一个数据库。

## 2.3、创建一个数据库

在这一步中，您将创建您将要使用（working with）的数据库。您将使用 DBMS 自带的 CLI 去创建一个数据库和表，以及添加一些数据。

您将创建一个数据库，其中包括有关黑胶唱片（vinyl）上的老式爵士乐（vintage jazz）唱片的数据。

这里我们使用 [MySQL CLI](https://dev.mysql.com/doc/refman/8.0/en/mysql.html)，但是大部分 DBMSes 有他们自己的类似特性的 CLI。

1. 打开一个新的命令提示符

2. 在这个命令行中，登录到你的 DBMS，就像下面这个 MySQL 的示例一样

   ```shell
   $ mysql -ulhs -plhs960906 -h192.168.229.160
   mysql> 
   ```

3. 在 mysql 的命令提示符中，创建数据库

   ```mysql
   mysql> create database recordings;
   ```

4. 使用你刚刚创建的数据库然后添加表

   ```mysql
   mysql> use recordings;
   Database changed
   ```

5. 使用你的文本编辑器，在 data-access 目录，创建一个叫 create-tables.sql 的 SQL 脚本来添加表

6. 在这个文件中，粘贴如下 SQL 代码，然后保存该文件

   ```shell
   $ vi create-tables.sql
   DROP TABLE IF EXISTS album;
   CREATE TABLE album (
     id         INT AUTO_INCREMENT NOT NULL,
     title      VARCHAR(128) NOT NULL,
     artist     VARCHAR(255) NOT NULL,
     price      DECIMAL(5,2) NOT NULL,
     PRIMARY KEY (`id`)
   );
   
   INSERT INTO album
     (title, artist, price)
   VALUES
     ('Blue Train', 'John Coltrane', 56.99),
     ('Giant Steps', 'John Coltrane', 63.99),
     ('Jeru', 'Gerry Mulligan', 17.99),
     ('Sarah Vaughan', 'Sarah Vaughan', 34.98);
   ```

7. 执行你刚才创建的脚本

   ```mysql
   mysql> source /path/to/create-tables.sql
   ```

8. 在你的 DBMS 命令行提示符中，使用 SELECT 语句确认您已经成功导入了表数据

   ```mysql
   mysql> select * from album;
   +----+---------------+----------------+-------+
   | id | title         | artist         | price |
   +----+---------------+----------------+-------+
   |  1 | Blue Train    | John Coltrane  | 56.99 |
   |  2 | Giant Steps   | John Coltrane  | 63.99 |
   |  3 | Jeru          | Gerry Mulligan | 17.99 |
   |  4 | Sarah Vaughan | Sarah Vaughan  | 34.98 |
   +----+---------------+----------------+-------+
   4 rows in set (0.00 sec)
   ```

## 2.4、导入一个数据库驱动

你已经得到了一个有一些数据的数据库，现在我们开始写 Go 代码。

定位和导入一个数据库驱动，会将 `database/sql` 包中的函数请求翻译为数据库可以理解的请求。

1. 在你的浏览器中，访问 [SQLDrivers](https://github.com/golang/go/wiki/SQLDrivers) 维基页面以确认一个你可以使用的驱动。在本篇的示例中，我们将使用 [Go-MySQL-Driver](https://github.com/go-sql-driver/mysql/)

2. 注意驱动的包名为：`github.com/go-sql-driver/mysql`

3. 使用你的文本编辑器，创建一个文件，其中写了您的 Go 代码，然后在 data-access 目录中将其存储为 main.go 文件

4. 在 main.go 中，赋值如下代码导入驱动包

   ```go
   package main
   
   import "github.com/go-sql-driver/mysql"
   ```

随着驱动被导入，你可以开始写代码去访问数据库了。

## 2.5、获取数据库句柄和连接

### 2.5.1、编写代码

现在编写一些 Go 代码，让您可以使用数据库句柄访问数据库。

您将使用一个 `sql.DB` 结构体的指针，它表示对特定数据库的访问。

1. 在 main.go 中，import 语句下面，粘贴如下的 Go 代码创建一个数据库句柄

   ```go
   var db *sql.DB
   
   func main() {
       // Capture connection properties.
       cfg := mysql.Config{
           User:   os.Getenv("DBUSER"),
           Passwd: os.Getenv("DBPASS"),
           Net:    "tcp",
           Addr:   "192.168.229.160:3306",
           DBName: "recordings",
           AllowNativePasswords: true,
       }
       // Get a database handle.
       var err error
       db, err = sql.Open("mysql", cfg.FormatDSN())
       if err != nil {
           log.Fatal(err)
       }
   
       pingErr := db.Ping()
       if pingErr != nil {
           log.Fatal(pingErr)
       }
       fmt.Println("Connected!")
   }
   ```

   在这块的代码中，您将：

   * 声明一个 `*sql.DB` 类型的变量。这是你的数据库句柄

     > Note：在这让 db 成为一个全局变量能够简化（simplifies）这个例子。再生产中，你要避免这种全局变量，例如（such as）通过将变量传递给需要它的函数或将其包装在结构中。

   * 使用 MySQL 驱动的 Config，以及 Config 类型的 FormatDSN 方法，去收集连接相关的属性并将它们格式化为 DSN 的连接字符串。

     > Note：这个 Config 结构体让代码相较于连接字符串更易阅读

   * 调用 `sql.Open` 去初始化一个 db 变量，通过 FormatDSN

   * 检查 `sql.Open` 的 error。它可能会失败，例如，你的数据库连接格式细节（specifics）不正确（well-formed）

     > Note：为了简化代码，您将会调用 `log.Fatal` 去结束执行并打印错误到控制台。在生产代码中，您要以更加优雅的（graceful）方式去处理错误。

   * 调用 `DB.Ping` 以确认连接到数据库是否工作。在运行时，`sql.Open` 可能不会立即连接，当然这取决于你的驱动。您在这使用 `Ping` 来确认 `database/sql` 包在您需要的时候可以进行连接

   * 检查 `Ping` 返回的 error，以防（in case）数据库连接失败

   * 如果 `Ping` 成功则打印一个信息

2. 在 main.go 文件顶部附件，保声明的下方，导入你代码中需要用到的一些包

   ```go
   package main
   
   import (
       "database/sql"
       "fmt"
       "log"
       "os"
   
       "github.com/go-sql-driver/mysql"
   )
   ```

3. 保存 main.go 文件

### 2.5.2、运行代码

1. 开始将 MySQL 驱动模块作为依赖项进行跟踪（tracking）

   使用 `go get` 去添加 `github.com/go-sql-driver/mysql`模块作为你自己模块的依赖。使用句号 `.` 意味着获取当前目录中代码的依赖项

   ```shell
   $ go get .
   go: added github.com/go-sql-driver/mysql v1.7.1
   ```

   Go 将会下载这个依赖，因为你在之前的 import 声明中添加了它。更多依赖跟踪的相关信息，请参阅：[Adding a dependency](https://go.dev/doc/modules/managing-dependencies#adding_dependency)

2. 在命令提示符中，设置 DBUSER 和 DBPASS 环境变量供 Go 程序使用

   On Linux or Mac：

   ```shell
   $ export DBUSER=lhs
   $ export DBPASS=lhs960906
   ```

   On Windows:

   ```shell
   C:\Users\you\data-access> set DBUSER=lhs
   C:\Users\you\data-access> set DBPASS=lhs960906
   ```

3. 在包含 main.go 的目录中，通过 `go run .` 命令运行代码，`.` 意味着运行当前目录中的包。

   ```shell
   $ go run .
   Connected!
   ```

你可以连接到数据库了！接下来，您将查询一些数据库中的数据。

## 2.6、查询多行

在这个部分，我们将用 Go 执行 SQL 查询返回多行。

 使用 `database/sql` 包中的 Query 方法可以返回多行，然后我们可以循环遍历返回的行（稍后您将在 [Query for a single row ](https://go.dev/doc/tutorial/database-access#single_row)部分中学习如何查询单行）。

### 2.6.1、编写代码

1. 在 main.go 中，main 函数的上方，粘贴如下 Ablum 结构体的定义，您将使用它来承接查询返回的行数据。

   ```go
   type Album struct {
       ID     int64
       Title  string
       Artist string
       Price  float32
   }
   ```

2. 在 main 函数的下方，粘贴如下 albumsByArtist 方法以查询数据库

   ```go
   // albumsByArtist queries for albums that have the specified artist name.
   func albumsByArtist(name string) ([]Album, error) {
       // An albums slice to hold data from returned rows.
       var albums []Album
   
       rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
       if err != nil {
           return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
       }
       defer rows.Close()
       // Loop through rows, using Scan to assign column data to struct fields.
       for rows.Next() {
           var alb Album
           if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
               return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
           }
           albums = append(albums, alb)
       }
       if err := rows.Err(); err != nil {
           return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
       }
       return albums, nil
   }
   ```

   在如上代码中，您做了如下几件事情：

   * 定义了一个 `Album`  类型的切片 albums。它将用来承载返回的行数据。结构体的属性名和类型对应（correspond）数据库的字段名和类型。

   * 使用 `DB.Query` 执行 SELECT 语句去查询具有指定艺术家（artist）姓名的专辑

     Query 方法的第一个参数是 SQL 语句，随后你可以传一到N个任意类型的参数。这些为您提供了在 SQL 语句中指定参数值的位置。通过将 SQL 语句与参数值分开（而不是使用 `fmt.Sprintf` 拼接它们），您启用 `database/sql` 包来发送与 SQL 文本分开的值，从而消除任何 SQL 注入风险

   * 延迟关闭 rows，这样他承载的资源将在方法退出时被释放

   * 循环迭代返回的 rows，使用 `rows.Scan` 方法可以将每一行中的字段值传递给 Album 结构的属性中

     Scan 方法接受（take）指向 Go 值的指针列表，字段值将会写入这些地方。在这里，您将指针指向 alb 变量的属性（通过使用 & 操作符）。Scan 通过指针写入以更新结构体的属性

   * 在循环中，检查扫描字段值到结构属性时产生的 error

   * 在循环中，将新的 alb 变量添加到 albums 切片中

   * 循环结束后，使用 `rows.Err` 检查整个查询中的错误。请注意，如果查询本身失败，则检查此处的错误是发现结果不完整的唯一方法

3. 更新你的 main 方法去调用 albumsByArtist

   在 main 方法的末尾，添加如下方法：

   ```go
   albums, err := albumsByArtist("John Coltrane")
   if err != nil {
       log.Fatal(err)
   }
   fmt.Printf("Albums found: %v\n", albums)
   ```

   在新的代码中，您正在做：

   * 调用新添加的 albumsByArtist 方法，用 albums 变量接受它的返回值
   * 打印结果

### 2.6.2、运行代码

在包含 main.go 目录的命令行下，执行如下代码：

```go
$ go run .
Connected!
Albums found: [{1 Blue Train John Coltrane 56.99} {2 Giant Steps John Coltrane 63.99}]
```

接下来，您将查询单行。

## 2.7、查询单行

在这部分，您将使用 Go 在数据库中查询单行。

对于您知道最多返回一行的 SQL 语句，您可以使用 QueryRow，它比使用 Query 循环更简单。

### 2.7.1、编写代码

1. 在 albumsByArtist 之下，粘贴如下 albumByID 方法

   ```go
   // albumByID queries for the album with the specified ID.
   func albumByID(id int64) (Album, error) {
       // An album to hold data from the returned row.
       var alb Album
   
       row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
       if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
           if err == sql.ErrNoRows {
               return alb, fmt.Errorf("albumsById %d: no such album", id)
           }
           return alb, fmt.Errorf("albumsById %d: %v", id, err)
       }
       return alb, nil
   }
   ```

   在这段代码中，您：

   * 使用 `DB.QueryRow` 去执行了一个 SELECT 语句，用指定的 ID 去查询专辑

     他返回了一个 `sql.Row`。为了简化调用代码，QueryRow 不会返回 error。它被安排稍后从 `Rows.Scan` 返回任何查询错误（例如 `sql.ErrNoRows`）

   * 使用 `Row.Scan` 去拷贝字段值到结构体属性中

   * 检查 Scan 产生的错误

     这个指定的 `sql.ErrNoRows` 的错误表明这个查询未返回行。通常，该错误值得用更具体的文本替换，例如此处的 "no such album"

2. 更新 main 函数以调用 albumByID

   在 main 函数的末尾，添加如下代码

   ```go
   // Hard-code ID 2 here to test the query.
   alb, err := albumByID(2)
   if err != nil {
       log.Fatal(err)
   }
   fmt.Printf("Album found: %v\n", alb)
   ```

   在这段新代码中，您：

   * 调用了您添加的 albumByID 方法
   * 打印了返回的专辑 ID

### 2.7.2、运行代码

在包含 main.go 目录的命令行下，执行如下代码：

```shell
$ go run .
Connected!
Albums found: [{1 Blue Train John Coltrane 56.99} {2 Giant Steps John Coltrane 63.99}]
Album found: {2 Giant Steps John Coltrane 63.99}
```

## 2.8、添加数据

在这一部分，你将会使用 Go 去执行一个 INSERT 的 SQL 语句添加一个新行到数据库中。

你已经知道了如何使用 Query 和 QueryRow 执行返回数据的 SQL。要执行不返回数据的 SQL，您需要使用 Exec 方法。

### 2.8.1、编写代码

1. 在 albumByID 的下方，粘贴如下 addAlbum 方法往数据库中插入一个新专辑，然后保存 main.go 文件

   ```go
   // addAlbum adds the specified album to the database,
   // returning the album ID of the new entry
   func addAlbum(alb Album) (int64, error) {
       result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
       if err != nil {
           return 0, fmt.Errorf("addAlbum: %v", err)
       }
       id, err := result.LastInsertId()
       if err != nil {
           return 0, fmt.Errorf("addAlbum: %v", err)
       }
       return id, nil
   }
   ```

   在这段代码中，您：

   * 使用 `DB.Exec` 执行了一个 INSERT 语句

     像 Query 一样，Exec 携带了一个 SQL 语句，后面跟着一些列的参数值

   * 检查 INSERT 操作的 error

   * 使用 `Result.LastInsertId` 获取插入数据后的 ID

   * 检查获取 ID 时的 error

2. 更新 main 函数去调用 addAlbum 方法

   在 main 函数的底部，添加如下代码

   ```go
   albID, err := addAlbum(Album{
       Title:  "The Modern Sound of Betty Carter",
       Artist: "Betty Carter",
       Price:  49.99,
   })
   if err != nil {
       log.Fatal(err)
   }
   fmt.Printf("ID of added album: %v\n", albID)
   ```

   在新代码中，您：

   * 用一个新 album 调用 addAlbum 方法，将分配的 ID 赋值给 albID 变量

### 2.8.2、运行代码

在包含 main.go 目录的命令行下，执行如下代码：

```shell
$ go run .
Connected!
Albums found: [{1 Blue Train John Coltrane 56.99} {2 Giant Steps John Coltrane 63.99}]
Album found: {2 Giant Steps John Coltrane 63.99}
ID of added album: 5
```

# 三、打开数据库句柄

database/sql 包通过减少管理连接的需要来简化数据库访问。与许多数据访问 API 不同，使用 database/sql，您不需要显式打开连接、执行工作，然后关闭连接。相反，您的代码会打开一个表示连接池的数据库句柄，然后使用该句柄执行数据访问操作，仅在需要释放资源（例如由检索到的 rows 或 prepared statement 持有的资源）时调用 Close 方法。

换句话说，它是由 sql.DB 表示的数据库句柄，代表您的代码处理连接、打开和关闭它们。当您的代码使用句柄执行数据库操作时，这些操作可以并发访问数据库。 有关更多信息，请参阅：[Managing connections](https://go.dev/doc/database/manage-connections)。

> 注意：您还可以保留数据库连接。更多信息，请参阅：[Using dedicated connections](https://go.dev/doc/database/manage-connections#dedicated_connections)。

除了 database/sql 包中可用的API之外，Go社区还为所有最常见（和许多不常见）数据库管理系统（DBMS）开发了驱动程序。

打开数据库句柄时，您需要遵循以下高级步骤：

1. 定位和导入数据库 driver

   驱动程序在您的 Go 代码和数据库之间转换请求和响应。更多信息，请参阅：[Locating and importing a database driver](https://go.dev/doc/database/open-handle#database_driver)

2. 打开一个数据库句柄

   导入驱动程序后，您可以打开特定数据库的句柄。

3. 确认连接

   打开数据库句柄后，您的代码可以检查连接是否可用。

您的代码通常不会显式打开或关闭数据库连接 — 这是由数据库句柄完成的。但是，您的代码应该释放它一路获得的资源，例如包含查询结果的 sql.Rows。 有关更多信息，请参阅：[释放资源](#Freeing resources)

## 3.1、定位和导入数据库驱动

您需要一个支持您正在使用的 DBMS 的数据库驱动程序。找到适合您的数据库的驱动程序，请参阅： see [SQLDrivers](https://github.com/golang/go/wiki/SQLDrivers)。

为了使驱动程序可用于您的代码，您可以像导入另一个 Go 包一样导入它。 这是一个例子：

```go
import "github.com/go-sql-driver/mysql"
```

请注意，如果您不直接从驱动程序包调用任何函数（例如当 sql 包隐式使用它时），您将需要使用空白导入，它在导入路径前添加下划线前缀：

```go
import _ "github.com/go-sql-driver/mysql"
```

> 注意：最佳实践是，避免使用数据库驱动程序自己的 API 进行数据库操作。相反，请使用 database/sql 包中的函数。 这将有助于保持代码与 DBMS 松散耦合，从而在需要时更轻松地切换到不同的 DBMS。

## 3.2、打开一个数据库句柄

数据库句柄 sql.DB 提供了单独或在事务中读取和写入数据库的能力。

您可以通过调用 sql.Open（采用连接字符串）或 sql.OpenDB（采用 driver.Connector）来获取数据库句柄。两者都返回一个指向 sql.DB 的指针。

> 注意：请务必将您的数据库凭据保存在您的 Go 源代码之外。

### 3.2.1、使用 connection string 打开

当您想要使用 connection string 进行连接时，请使用 sql.Open 函数。字符串的格式会因您使用的驱动程序而异。

这里是一个 MySQL 驱动的示例：

```go
db, err = sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/jazzrecords")
if err != nil {
    log.Fatal(err)
}
```

但是，您可能会发现以更结构化的方式捕获连接属性可以使您的代码更具可读性。 细节因 driver 而异。

例如，您可以将前面的示例替换为以下示例，该示例使用 MySQL 驱动程序的 Config 来指定属性，并使用其 FormatDSN 方法来构建连接字符串。

```go
// Specify connection properties.
cfg := mysql.Config{
    User:   username,
    Passwd: password,
    Net:    "tcp",
    Addr:   "127.0.0.1:3306",
    DBName: "jazzrecords",
}

// Get a database handle.
db, err = sql.Open("mysql", cfg.FormatDSN())
if err != nil {
    log.Fatal(err)
}
```

### 3.2.2、使用 Connector 打开【待补充】

```go
// Specify connection properties.
cfg := mysql.Config{
    User:   username,
    Passwd: password,
    Net:    "tcp",
    Addr:   "127.0.0.1:3306",
    DBName: "jazzrecords",
}

// Get a driver-specific connector.
connector, err := mysql.NewConnector(&cfg)
if err != nil {
    log.Fatal(err)
}

// Get a database handle.
db = sql.OpenDB(connector)
```

### 3.2.3、错误处理

您的代码应该检查尝试创建句柄（例如使用 sql.Open）时是否出现错误。这不会是连接错误。相反，如果 sql.Open 无法初始化句柄，您将收到错误消息。例如，如果它无法解析您指定的 DSN，则可能会发生这种情况。

## 3.3、确认连接

当您打开一个数据库句柄时，database/sql 包本身可能不会立即创建一个新的数据库连接。相反，它可能会在您的代码需要时创建连接。如果您不会立即使用数据库并想确认是否可以建立连接，请调用 Ping 或 PingContext。

```go
db, err = sql.Open("mysql", connString)

// Confirm a successful connection.
if err := db.Ping(); err != nil {
    log.Fatal(err)
}
```

## <a id="Freeing resources">3.4、释放资源</a>

尽管您不需要用 database/sql 包去显式管理或关闭的连接，但您的代码应该在不再需要时释放连接所获得的资源。这些可以包括表示从查询返回的数据的 sql.Rows 或表示准备好的语句的 sql.Stmt 所持有的资源。

通常，您可以通过推迟对 Close 函数的调用来关闭资源，以便在封闭函数退出之前释放资源。

以下示例中的代码推迟 Close 以释放 sql.Rows 持有的资源。

```go
rows, err := db.Query("SELECT * FROM album WHERE artist = ?", artist)
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

// 循环返回的rows.
...
```



# 四、执行无数据返回的 SQL 语句

当你要执行不返回数据的数据库动作，使用 `database/sql` 中的 Exec 或 ExecContext 方法。这种 SQL 语句包括  `INSERT`, `DELETE` 和 `UPDATE`。

当你的查询可能返回 rows，使用 Query 或 QueryContext 代替，更多信息，请参阅：[Querying a database](https://go.dev/doc/database/querying)。

ExecContext 方法和 Exec 方法一样，但是有一个额外的 `context.Context` 参数，如  [Canceling in-progress operations](https://go.dev/doc/database/cancel-operations) 中描述。

下面的示例中使用 DB.Exec 去执行语句往 album 表中添加新的记录。

```go
func AddAlbum(alb Album) (int64, error) {
    result, err := db.Exec("INSERT INTO album (title, artist) VALUES (?, ?)", alb.Title, alb.Artist)
    if err != nil {
        return 0, fmt.Errorf("AddAlbum: %v", err)
    }

    // Get the new album's generated ID for the client.
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("AddAlbum: %v", err)
    }
    // Return the new album's ID.
    return id, nil
}
```

DB.Exec 返回一个一个 sql.Result 和一个 error。当 error 为 nil 的时候，你可以使用 Result 去获取最后一次插入条目的 ID 或者获取操作影响的行数。

> Note：在 prepared 语句使用的参数占位符（placeholders）因您使用的 DBMS 和驱动各不相同（vary）。例如，Postgres 的 [pq driver](https://pkg.go.dev/github.com/lib/pq) 使用像 `$1` 这样的占位符来代替 `?`

如果您的代码将重复执行相同的 SQL 语句，请考虑使用 sql.Stmt 从 SQL 语句创建可重用的准备语句。更多信息请参阅：[Using prepared statements](https://go.dev/doc/database/prepared-statements)

> Caution：不要使用 fmt.Sprintf 之类的字符串格式化函数来组装（assemble） SQL 语句！ 您可能会引入（introduce） SQL 注入风险。更多信息请参阅：[Avoiding SQL injection risk](https://go.dev/doc/database/sql-injection)

执行 SQL 语句不返回 rows 的函数：

<table>
    <tr>
        <th>Function</th>
        <th>Description</th>
    </tr>
    <tr>
        <td>
            DB.Exec</br>
			DB.ExecContext
        </td>
        <td>执行单个SQL语句</td>
    </tr>
	<tr>
        <td>
            Tx.Exec</br>
			Tx.ExecContext
        </td>
        <td>在一个更大的事务中执行SQL语句</td>
    </tr>
	<tr>
        <td>
            Stmt.Exec</br>
			Stmt.ExecContext
        </td>
        <td>执行一个已准备好的SQL语句</td>
    </tr>
	<tr>
        <td>Conn.ExecContext</td>
        <td>使用保留连接</td>
    </tr>
</table>



# 五、查询数据【待补充】

当执行返回数据的 SQL 语句时，请使用 database/sql 包中提供的 Query 方法之一（例如 DB.Query 或 DB.QueryRow）。其中每一个都会返回一个或多个行，您可以使用 Scan 方法将其数据复制到变量中。database/sql 包提供了两种执行结果查询的方法。

* 查询单行 – QueryRow 最多返回从数据库返回单行数据。更多信息，请参考：[Querying for a single row](https://go.dev/doc/database/querying#single_row)
* 查询多行 – Query 用一个 Rows 结构返回数据库中所有匹配的行，你的程序可以循环这个 Rows。更多信息，请参考：[Querying for multiple rows](https://go.dev/doc/database/querying#multiple_rows)

执行不返回数据的语句时，可以改用 Exec 或 ExecContext 方法。有关更多信息，请参阅：[Executing statements that don’t return data](https://go.dev/doc/database/change-data)。

如果您的代码将重复执行相同的 SQL 语句，请考虑使用准备好的语句。有关更多信息，请参阅：[Using prepared statements](https://go.dev/doc/database/prepared-statements)。

> 警告：不要使用字符串格式化函数（例如 fmt.Sprintf）来组装 SQL 语句！您可能会引入 SQL 注入风险。有关更多信息，请参阅：[Avoiding SQL injection risk](https://go.dev/doc/database/sql-injection)。

## 5.1、查询单行

QueryRow 最多检索单个数据库行，例如当您想通过唯一 ID 查找数据时。如果查询返回多行，则 Scan 方法将丢弃除第一行之外的所有行。

QueryRowContext 的工作方式与 QueryRow 类似，但带有 context.Context 参数。有关更多信息，请参阅：取消正在进行的操作。

以下示例使用查询来查明是否有足够的库存来支持购买。如果足够则 SQL 语句返回 true，否则返回 false。 Row.Scan 通过指针将布尔返回值复制到 enough 变量中。

**返回单行的方法**：

<table>
    <tr>
        <th>Function</th>
        <th>Description</th>
	</tr>
    <tr>
        <td>
            <a href="https://pkg.go.dev/database/sql#DB.QueryRow">DB.QueryRow</a><br/>
        	<a href="https://pkg.go.dev/database/sql#DB.QueryRowContext">DB.QueryRowContext</a>
        </td>
        <td>单独运行单行查询</td>
	</tr>
    <tr>
    	<td>
            <a href="https://pkg.go.dev/database/sql#Tx.QueryRow">Tx.QueryRow</a><br/>
        	<a href="https://pkg.go.dev/database/sql#Tx.QueryRowContext">Tx.QueryRowContext</a>
        </td>
        <td>在较大的事务中运行单行查询。有关更多信息，请参阅：</td>
    </tr>
    <tr>
    	<td>
            <a href="https://pkg.go.dev/database/sql#Tx.QueryRow">Stmt.QueryRow</a><br/>
        	<a href="https://pkg.go.dev/database/sql#Tx.QueryRowContext">Stmt.QueryRowContext</a>
        </td>
        <td>使用已准备好的语句运行单行查询。有关更多信息，请参阅：</td>
    </tr>
    <tr>
    	<td>
            <a href="https://pkg.go.dev/database/sql#Conn.QueryRowContext">Conn.QueryRowContext</a>
        </td>
        <td>与保留连接一起使用。有关更多信息，请参阅：</td>
    </tr>
</table>

## 5.2、查询多行

## 5.3、处理可为空的列值

## 5.4、处理多个结果集



# 六、使用 prepared statement

您可以是定义一个 prepared statement 以重复使用，这样可以免去您每次进行数据库操作时创建 statement 的开销，以致于可以让您的代码运行地更快一些。

> Note：prepared statements 中的参数占位符各种各样，取决于您使用的数据库驱动。例如，Postgres 的 [pq driver](https://pkg.go.dev/github.com/lib/pq) 使用像 $1 这样的占位符代替 ?。

## 6.1、什么是 prepared statement

prepared statement 是由 DBMS 解析和保存的 SQL，通常包含占位符但没有实际参数值。稍后，可以使用一组参数值来执行该语句。

## 6.2、如何使用 prepared statement

当您希望重复执行相同的 SQL 时，可以使用 sql.Stmt 提前准备 SQL 语句，然后根据需要执行它。

以下示例创建一个 prepared statement，用于从数据库中选择特定的专辑。DB.Prepare 返回一个 sql.Stmt，表示给定 SQL 文本的 prepared statement。然后您可以将 SQL 语句的参数传递给 Stmt.Exec、Stmt.QueryRow 或 Stmt.Query 来运行该语句。

```go
// AlbumByID 获取指定的专辑
func AlbumByID(id int) (Album, error) {
    // 定义一个 prepared statement. You'd typically define the statement
    // elsewhere and save it for use in functions such as this one.
    stmt, err := db.Prepare("SELECT * FROM album WHERE id = ?")
    if err != nil {
        log.Fatal(err)
    }

    var album Album

	// 执行 prepared statement，将 id 的值传入占位符 ?
    err := stmt.QueryRow(id).Scan(&album.ID, &album.Title, &album.Artist, &album.Price, &album.Quantity)
    if err != nil {
        if err == sql.ErrNoRows {
            // Handle the case of no rows returned.
        }
        return album, err
    }
    return album, nil
}
```

## 6.3、Prepared statement 的行为

一个预准备的 sql.Stmt 提供常用的 Exec、QueryRow 和 Query 方法来调用语句。有关使用这些方法的更多信息，请参阅：[Querying for data](https://go.dev/doc/database/querying) 和 [Executing SQL statements that don’t return data](https://go.dev/doc/database/change-data)。

但是，由于 sql.Stmt 已经代表了预设的 SQL 语句，因此它的 Exec、QueryRow 和 Query 方法仅采集占位符对应的 SQL 参数值，而省略了 SQL 文本。

您可以通过不同的方式定义新的 sql.Stmt，具体取决于您将如何使用它。

* DB.Prepare 和 DB.PrepareContext 创建一个准备好的语句，该语句可以在事务之外单独执行，就像 DB.Exec 和 DB.Query 一样。
* Tx.Prepare、Tx.PrepareContext、Tx.Stmt 和 Tx.StmtContext 创建用于特定事务的 prepared statement。Prepare 和 PrepareContext 使用 SQL 文本来定义语句。Stmt 和 StmtContext 使用 DB.Prepare 或 DB.PrepareContext 的结果。也就是说，它们将非用于事务的 sql.Stmt 转换为用于此事务的 sql.Stmt。
* Conn.PrepareContext 从 sql.Conn 创建一条 prepared statement，它表示一个保留的连接。

确保当代码完成语句时调用 stmt.Close。这将释放可能与其关联的任何数据库资源（例如底层连接）。对于函数中仅是局部变量的语句，推迟 stmt.Close() 就足够了。

用于创建准备语句的方法：

<table>
    <tr>
        <th>方法</th>
        <th>描述</th>
    </tr>
    <tr>
        <td>
        	DB.Prepare<br/>
            DB.PrepareContext
        </td>
        <td>
        	准备一条单独执行的语句，或者使用 Tx.Stmt 将其转换为事务内的 prepared statement。
        </td>
    </tr>
    <tr>
        <td>
        	Tx.Prepare<br/>
            Tx.PrepareContext<br/>
            Tx.Stmt<br/>
            Tx.StmtContext
        </td>
        <td>
            在特定事务中使用的 prepared statement。更多信息，请参阅：<a href="Executing transactions.">Executing transactions</a>
        </td>
    </tr>
    <tr>
        <td>
        	Conn.PrepareContext
        </td>
        <td>
            与保留连接一起使用。 有关更多信息，请参阅：<a href="https://go.dev/doc/database/manage-connections">Managing connections
        </td>
    </tr>
</table>



# 七、执行事务

您可以使用代表事务的 sql.Tx 执行数据库事务。 除了表示事务特定语义的 Commit 和 Rollback 方法之外，sql.Tx 还具有您用来执行常见数据库操作的所有方法。要获取 sql.Tx，您可以调用 DB.Begin 或 DB.BeginTx。

数据库事务将多个操作分组为更大目标的一部分。所有操作都必须成功，否则都不能成功，在任何一种情况下都保留数据的完整性。通常，交易工作流程包括：

* 1）开始交易。
* 2）执行一组数据库操作。
* 3）如果没有错误发生，提交事务以进行数据库更改。
* 4）如果发生错误，回滚事务以保持数据库不变。

sql 包提供了开始和结束事务的方法，以及执行中间数据库操作的方法。这些方法对应于上述工作流程中的四个步骤。

* 1）开始一个事务

  > [DB.Begin](https://pkg.go.dev/database/sql#DB.Begin) 或 [DB.BeginTx](https://pkg.go.dev/database/sql#DB.BeginTx) 开启一个新的数据库事务， 返回的 sql.Tx 就代表这个事务

* 2）执行数据库操作

  使用 sql.Tx，您可以在使用单个连接的一系列操作中查询或更新数据库。为了支持这一点，Tx 导出以下方法：

  * [Exec](https://pkg.go.dev/database/sql#Tx.Exec) 和 [ExecContext](https://pkg.go.dev/database/sql#Tx.ExecContext) 用于通过 SQL 语句（如 INSERT、UPDATE 和 DELETE）更改数据库。

    > 有关详细信息，请参阅： [Executing SQL statements that don’t return data](https://go.dev/doc/database/change-data)

  * [Query](https://pkg.go.dev/database/sql#Tx.Query)，[QueryContext](https://pkg.go.dev/database/sql#Tx.QueryContext)，[QueryRow](https://pkg.go.dev/database/sql#Tx.QueryRow) 和 [QueryRowContext](https://pkg.go.dev/database/sql#Tx.QueryRowContext) 用于不返回 rows 的操作

    > 有关详细信息，请参阅： [Querying for data](https://go.dev/doc/database/querying)

  * [Prepare](https://pkg.go.dev/database/sql#Tx.Prepare), [PrepareContext](https://pkg.go.dev/database/sql#Tx.PrepareContext), [Stmt](https://pkg.go.dev/database/sql#Tx.Stmt), and [StmtContext](https://pkg.go.dev/database/sql#Tx.StmtContext) 用于预定义 prepared statements

    > 有关详细信息，请参阅： [Using prepared statements](https://go.dev/doc/database/prepared-statements)

* 3）使用如下方式的其中一种结束事务

  * 使用 Tx.Commit 提交事务

    > 如果 Commit 成功（返回 nil 错误），则所有查询结果都被确认为有效，并且所有已执行的更新都作为单个原子更改应用于数据库。 如果 Commit 失败，则 Tx 上 Query 和 Exec 的所有结果都应被视为无效而丢弃。

  * 使用 Tx.Rollback 回滚事务

    >  即使 Tx.Rollback 失败，事务也不再有效，也不会提交到数据库。

## 7.1、最佳实践

遵循以下最佳实践，以更好地应对事务有时需要的复杂语义和连接管理。

* 使用本节中描述的 API 来管理事务。 不要直接使用与事务相关的 SQL 语句（例如 BEGIN 和 COMMIT），否则会使数据库处于不可预测的状态，尤其是在并发程序中。
* 使用事务时，请注意不要直接调用非事务 sql.DB 方法，因为这些方法将在事务之外执行，从而使您的代码对数据库状态的视图不一致，甚至导致死锁。

## 7.2、示例

以下示例中的代码使用事务为相册创建新的客户订单。 在此过程中，代码将：

1. 开始交易
2. 推迟事务的回滚。如果事务成功，它将在函数退出之前提交，从而使延迟回滚调用成为无操作。如果事务失败，则不会提交，这意味着将在函数退出时调用回滚
3. 确认客户订购的专辑有足够的库存
4. 如果有足够的库存，请更新库存计数，将其减少所订购的专辑数量
5. 创建一个新订单并检索新订单为客户端生成的 ID
6. 提交交易并返回 ID

此示例使用采用 context.Context 参数的 Tx 方法。这使得函数的执行（包括数据库操作）在运行时间过长或客户端连接关闭时被取消。有关更多信息，请参阅 [Canceling in-progress operations](https://go.dev/doc/database/cancel-operations)

```go
// CreateOrder创建一个专辑的订单，并返回一个订单的id
func CreateOrder(ctx context.Context, albumID, quantity, custID int) (orderID int64, err error) {

    // 创建用于准备失败结果的辅助函数。
    fail := func(err error) (int64, error) {
        return 0, fmt.Errorf("CreateOrder: %v", err)
    }

    // 发起一个事务请求并获取一个代表事务的 Tx
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return fail(err)
    }
    // 延迟rollback以防出现任何失败
    defer tx.Rollback()

    // 确认专辑库存(inventory)足够下订单。
    var enough bool
    if err = tx.QueryRowContext(ctx, "SELECT (quantity >= ?) from album where id = ?",
        quantity, albumID).Scan(&enough); err != nil {
        if err == sql.ErrNoRows {
            return fail(fmt.Errorf("no such album"))
        }
        return fail(err)
    }
    if !enough {
        return fail(fmt.Errorf("not enough inventory"))
    }

    // 减去订单中专辑的数量以更新专辑库存
    _, err = tx.ExecContext(ctx, "UPDATE album SET quantity = quantity - ? WHERE id = ?",
        quantity, albumID)
    if err != nil {
        return fail(err)
    }

    // 在album_order表中插入一个新row
    result, err := tx.ExecContext(ctx, "INSERT INTO album_order (album_id, cust_id, quantity, date) VALUES (?, ?, ?, ?)",
        albumID, custID, quantity, time.Now())
    if err != nil {
        return fail(err)
    }
    // 获取刚刚创建的订单条目的id
    orderID, err = result.LastInsertId()
    if err != nil {
        return fail(err)
    }

    // 提交事务
    if err = tx.Commit(); err != nil {
        return fail(err)
    }

    // 返回订单ID
    return orderID, nil
}
```

# 八、取消正在进行的操作

你可以通过使用 Go 的 `context.Context` 管理正在进行（in-progress）的操作。Context 是一个标准的 Go 数据值，可以报告它所代表（represents）的整体操作是否已经（has been）被取消并且不再需要。通过在应用程序中跨函数调用和服务传递 `context.Context`，它们可以提前停止工作并在不再需要处理时返回错误。更多 Context 相关信息，请参阅：[Go Concurrency Patterns: Context](https://blog.golang.org/context)。

例如，你可能希望：

* 终止长时间运行的（long-running）操作，包括需要花费很长时间完成的数据库操作
* 从其他地方（elsewhere）传播（propagete）取消请求，例如当客户端关闭连接时。

许多提供给 Go 开发者的 APIs 包含携带 Context 参数的方法，使您可以更轻松地在整个（throughout）应用程序中使用 Context。

## 8.1、在超时后取消数据库操作

您可以使用 Context 设置超时或截止时间，在该时间后（after which）操作将被取消。要派生（derive）具有超时或截止日期的 Context，请调用 `context.WithTimeout` 或 `context.WithDeadline` 方法。

如下的超时示例代码派生了一个 Context 并且将其传递给了 `sql.DB.QueryContext` 的方法。

```go
func QueryWithTimeout(ctx context.Context) {
    // Create a Context with a timeout.
    queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    // Pass the timeout Context with a query.
    rows, err := db.QueryContext(queryCtx, "SELECT * FROM album")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    // Handle returned rows.
    ...
}

func InsertWithTimeout(ctx context.Context) {
    // Create a Context with a timeout.
    queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    // Pass the timeout Context with a insert.
	result, err := db.ExecContext(ctx,
		"INSERT INTO User (FirstName, LastName, Email, Password) VALUES (?, ?, ?)",
		user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		log.Fatal(err)
	}
}
```

当一个 Context 派生自外部 Context 时，如本例中的 queryCtx 派生自 ctx，如果取消外部 Context，则派生的 Context 也会自动取消。例如，在 HTTP 服务器中，`http.Request.Context` 方法返回与请求关联的 Context。 如果 HTTP 客户端断开连接或取消 HTTP 请求（可能在 HTTP/2 中），则该 Context 将被取消。传递一个 HTTP 请求的 Context 给上面的 QueryWithTimeout 会导致在整个 HTTP 请求被取消或查询花费超过 5s 时，数据的查询操作会提前停止。

> Note：始终延迟对 cancel 函数的调用，该函数在您创建具有超时或截止日期的新 Context 时返回。在包含函数退出时，会释放调新 Context 承载的资源。它也会取消 queryCtx，但是当函数返回时，应该不再（anymore）使用 queryCtx。

# 九、管理连接

对于绝大多数（vast majority）程序，您不需要调整 sql.DB 连接池默认值。但对于某些高级程序，您可能需要调整连接池参数或显式使用连接。本主题解释了如何进行。

sql.DB 数据库句柄对于多个 goroutine 并发使用是安全的（这意味着该句柄是其他语言可能所说的 "线程安全"）。其他的某些 database access library 基于同一时间只能用于一次操作的连接。为了弥补这一差距，每个 sql.DB 管理一个到底层数据库的活动连接池，根据 Go 程序中并行性的需要创建新的连接。

连接池适合大多数数据访问需求。 当您调用 sql.DB Query 或 Exec 方法时，sql.DB 实现会从池中检索可用连接，或者根据需要创建一个连接。当不再需要连接时，sql 包将其返回到池中，这中特性使得 Go 支持高级的数据库访问并行性。

## 9.1、设置连接池参数

您可以通过设置参数以指导 sql 包如何管理连接池的属性。要获取有关这些属性影响的统计信息，请使用 [`DB.Stats`](https://pkg.go.dev/database/sql#DB.Stats)。

### 9.1.1、设置最大连接数

[DB.SetMaxOpenConns](https://pkg.go.dev/database/sql#DB.SetMaxOpenConns) 对打开的连接数施加（imposes）限制。超过此限制，新的数据库操作将等待现有操作完成，sql.DB 才创建另一个连接。默认情况下，当所有现有连接都在使用且需要连接时，sql.DB 就会创建一个新连接。

请记住，设置限制会使数据库使用类似于获取锁或信号量，结果您的应用程序可能会死锁等待新的数据库连接。

### 9.1.2、设置连接的最大生命周期

使用 [DB.SetConnMaxLifetime](https://pkg.go.dev/database/sql#DB.SetConnMaxLifetime) 设置连接在关闭之前可以保持打开状态的最长时间。

默认情况下，连接可以使用并重复使用任意长的时间，但须遵守上述限制。在某些系统中，例如那些使用负载平衡数据库服务器的系统，确保应用程序永远不会在不重新连接的情况下使用特定连接的时间过长会很有帮助。

### 9.1.3、使用专用连接

database/sql 包含当数据库可能为在特定连接上执行的一系列操作分配隐式含义时可以使用的函数。

最常见的示例是事务，它通常以 BEGIN 命令开始，以 COMMIT 或 ROLLBACK 命令结束，并且包括在整个事务中这些命令之间的连接上发出的所有命令。对于此用例，请使用 sql 包的事务支持。 请参阅：[Executing transactions](https://go.dev/doc/database/execute-transactions)。

对于一系列单独操作必须全部在同一连接上执行的其他用例，sql 包提供了专用连接。DB.Conn 获得一个专用连接，即 sql.Conn。 sql.Conn 具有方法 BeginTx、ExecContext、PingContext、PrepareContext、QueryContext 和 QueryRowContext，它们的行为类似于 DB 上的等效方法，但仅使用专用连接。 完成专用连接后，您的代码必须使用 Conn.Close 释放它。



# 十、避免 SQL 注入

您可以通过提供 SQL 参数值作为 sql 包函数参数来避免 SQL 注入风险。 sql 包中的许多函数为 SQL 语句以及该语句的参数中使用的值提供参数（其他函数为准备好的语句和参数提供参数）。
