package modules

import io.getquill._
import org.bson.conversions.Bson
import org.bson.types.ObjectId
import org.bson.codecs.configuration.CodecRegistries.{
  fromRegistries,
  fromProviders,
  fromCodecs
}
import org.mongodb.scala.MongoClient.DEFAULT_CODEC_REGISTRY

import org.mongodb.scala._
import org.mongodb.scala.model.Filters._
import org.mongodb.scala.bson.codecs.Macros._
import scala.collection.JavaConverters._
import fansi.ErrorMode
import scala.util.Success
import scala.util.Failure
import java.{util => ju}
import org.bson.BsonDocument

// Using case classes with MongoDB
// http://mongodb.github.io/mongo-java-driver/4.2/driver-scala/getting-started/quick-start-case-class/

//http://mongodb.github.io/mongo-scala-driver/2.2/getting-started/quick-tour/

class PostgresDatabase {

  lazy val ctx = new PostgresJdbcContext(SnakeCase, "ctx")

  import ctx._
}

case class User(
    _id: ObjectId,
    name: String
)
object User {
  def apply(name: String): User =
    User(new ObjectId(), name)
}

class MongoMan {
  val url = MongoClient(
    "mongodb://rusty:rusty@192.168.1.19:30933"
  )

  lazy val mongoClient: MongoClient = MongoClient(buildHost())

  def buildHost(): MongoClientSettings = {
    val user: String = "rusty" // the user name
    val source: String = "admin" // the source where the user is defined
    val password: Array[Char] =
      "rusty".toCharArray // the password as a character array
// ...
    val credential = MongoCredential.createCredential(user, source, password)
    MongoClientSettings
      .builder()
      .applyToClusterSettings(b =>
        b.hosts(
          List(
            new ServerAddress("192.168.1.19", 30933)
          ).asJava
        )
      )
      .credential(credential)
      .build()
  }

  def getFactConfiguration(): MongoCollection[Document] = {
    val mdb = mongoClient.getDatabase("config")
    mdb.getCollection("factConfiguration")
  }

  def getFactDatabase(userUuid: String): MongoCollection[User] = {
    val codecRegistry = fromRegistries(
      fromProviders(classOf[User]),
      DEFAULT_CODEC_REGISTRY
    )
    val mdb = mongoClient.getDatabase("facts").withCodecRegistry(codecRegistry)
    mdb.getCollection(userUuid)
  }

  /** Retrieves the fact that matches the provided query for the user
    * executes the passed in function
    *
    * @param userUuid that we are matching the fact for
    * @param queryString that we are searching on
    */
  def getFactForUser(
      userUuid: String,
      queryString: String
  ): Observable[String] = {
    val db = getFactDatabase(userUuid)
    db.find(BsonDocument.parse(queryString))
      .map(_.name)
  }

  /** Insert the facts for the user
    *
    * @param userUuid collection to insert into
    * @param user object to insert as a document
    */
  def insertFactForUser(userUuid: String, user: User) = {
    val db = getFactDatabase(userUuid)

    // val document = Document("name" -> fact)
    val newUser = User("test")
    db
      .insertOne(newUser)
      .subscribe(println(_))
  }
}
