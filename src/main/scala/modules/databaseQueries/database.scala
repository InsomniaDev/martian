package modules

import io.getquill._

class PostgresDatabase {

  lazy val ctx = new PostgresJdbcContext(SnakeCase, "ctx")

  import ctx._
}
