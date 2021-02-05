// import java.nio.file.Path;
// import com.lambdazen.bitsy.BitsyGraph;
import io.getquill._

import modules.databaseQueries._

object Main extends App {
  println("Hello, World!")

  lazy val ctx = new PostgresJdbcContext(SnakeCase, "ctx")
  
  new Select(ctx).getAccounts().map(println(_))
}
