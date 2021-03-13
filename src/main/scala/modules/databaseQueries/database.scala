package modules

import io.getquill._
import org.mongodb.scala._
import scala.collection.JavaConverters._

class PostgresDatabase {

  lazy val ctx = new PostgresJdbcContext(SnakeCase, "ctx")

  import ctx._
}

class MongoMan {
  val mongoClient = MongoClient(
    "mongodb://rusty:rusty@192.168.1.19:30933,192.168.1.19:30934,192.168.1.19:30935"
  )

  def buildHost() = {
    val user: String = "rusty" // the user name
    val source: String = "admin" // the source where the user is defined
    val password: Array[Char] =
      "rusty".toCharArray // the password as a character array
// ...
    val credential = MongoCredential.createCredential(user, source, password)
    val settings: MongoClientSettings = MongoClientSettings
      .builder()
      .applyToClusterSettings(b =>
        b.hosts(
          List(
            new ServerAddress("192.168.1.19", 30933),
            new ServerAddress("192.168.1.19", 30934),
            new ServerAddress("192.168.1.19", 30935)
          ).asJava
        )
      )
      .credential(credential)
      .build()
    val mongoClient: MongoClient = MongoClient(settings)
  }

  def getFactConfiguration(): MongoCollection[Document] = {
    val mdb = mongoClient.getDatabase("config")
    mdb.getCollection("factConfiguration")
  }

  def getFactDatabase(userUuid: String): MongoCollection[Document] = {
    val mdb = mongoClient.getDatabase("facts")
    mdb.getCollection(userUuid)
  }

  def insertFactForUser(userUuid: String, fact: String) = {
    val db = getFactDatabase(userUuid)
    val document = Document("name" -> fact)
    db.insertOne(document).subscribe(println(_))
  }
}
