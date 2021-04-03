// import java.nio.file.Path;
// import com.lambdazen.bitsy.BitsyGraph;
import io.getquill._

import modules.databaseQueries._
import routes._
import akka.actor.typed.ActorSystem
import akka.actor.typed.scaladsl.Behaviors
import akka.http.scaladsl.Http
import akka.http.scaladsl.model._
import akka.http.scaladsl.server.Directives._
import scala.io.StdIn
import modules.MongoMan
import org.mongodb.scala.bson.collection.immutable.Document
import modules.User
import org.mongodb.scala.model.Filters._
import akka.actor.Status
import scala.util.Success
import scala.util.Failure

object Martian {

  private def loadConfig() = {
    lazy val ctx = new PostgresJdbcContext(SnakeCase, "ctx")

    // new FactData(ctx).(getFactData)
    val config = new ConfigData(ctx)

    if (config.getConfig().map(println(_)).size == 0) {
      // Let's have common words inserted into the database for the time being
      val commonWordsConfig = new Config(
        0,
        Some("commonWords"),
        Some(
          "the,be,to,of,and,a,in,that,have,I,it,for,not,on,with,he,as,you,do,at,this,but,his,by,from,they,we,say,her,she,or,an,will,my,one,all,would,there,their,what,so,up,out,if,about,who,get,which,go,me,when,make,can,like,time,no,just,him,know,take,people,into,year,your,good,some,could,them,see,other,than,then,now,look,only,come,its,over,think,also,back,after,use,two,how,our,work,first,well,way,even,new,want,because,any,these,give,day,most,us"
        )
      )
      config.insertConfigItem(commonWordsConfig)
      println("Config updated")
    } else println("Config already set")

  }

  // def processGet(value: String): Route = {
  //   onComplete(new MongoMan().getFactForUser("test", null)) {

  //   }
  // }

  def main(args: Array[String]): Unit = {

    implicit val system = ActorSystem(Behaviors.empty, "my-system")
    // needed for the future flatMap/onComplete in the end
    implicit val executionContext = system.executionContext

    val route = {
      get {
        path("hello") {
          complete("Hello To You")
        }
        path("user" / Segment) { query =>
          val values = query
            .split("=")
            .grouped(2)
            .collect { case Array(k, v) => k -> v }
            .toMap
          val (userUuid, factName) = values.head
          val getValue = new MongoMan()
            .getFactForUser(userUuid, s"{'name':'${factName}'}")
          onSuccess(getValue) { value =>
            complete(value.name)
          }
        // resp match {
        //   case User(id, name) => complete(HttpResponse(entity = e.name))
        // }
        // complete(HttpResponse(entity = resp.name))
        }
      } ~ post {
        path("new" / Segment) { newData =>
          val mc = new MongoMan().insertFactForUser("test", User(newData))
          complete("done")
        }
      }
    }

    val bindingFuture = Http().newServerAt("localhost", 8080).bind(route)

    println(s"Server online at http://localhost:8080/\nPress RETURN to stop...")
    StdIn.readLine() // let it run until user presses return
    bindingFuture
      .flatMap(_.unbind()) // trigger unbinding from the port
      .onComplete(_ => system.terminate()) // and shutdown when done
  }
}

// TODO: See if we want to do Akka Streams here
// - https://doc.akka.io/docs/alpakka/current/mongodb.html
