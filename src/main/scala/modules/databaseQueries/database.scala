package modules

import io.getquill._

// Using case classes with MongoDB
// http://mongodb.github.io/mongo-java-driver/4.2/driver-scala/getting-started/quick-start-case-class/

//http://mongodb.github.io/mongo-scala-driver/2.2/getting-started/quick-tour/

class CassandraDatabase {

  lazy val ctx = new CassandraAsyncContext(SnakeCase, "db")

  import ctx._
}
