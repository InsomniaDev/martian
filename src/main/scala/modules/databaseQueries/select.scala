package modules.databaseQueries
import io.getquill._

case class Accounts(user_id: Int, username: String)

class Select(ctx: PostgresJdbcContext[SnakeCase.type]) {
  import ctx._
  
  def getAccounts(): List[String] = {
    val q = quote {
      query[Accounts].map(a => a.username)
    }
    val resp = ctx.run(q)
    resp
  }
}
