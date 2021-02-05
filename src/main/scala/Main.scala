// import java.nio.file.Path;
// import com.lambdazen.bitsy.BitsyGraph;
import io.getquill._

case class Accounts (user_id: Int, username: String)

object Main extends App {
  println("Hello, World!")

  lazy val ctx = new PostgresJdbcContext(SnakeCase, "ctx")

  import ctx._

  val q = quote {
    query[Accounts].map(a => a.username)
  }
  
  val resp = ctx.run(q)
  resp.map(println(_))
}
