// import java.nio.file.Path;
// import com.lambdazen.bitsy.BitsyGraph;
import io.getquill._

import modules.databaseQueries._

object Main extends App {
  println("Hello, World!")

  lazy val ctx = new PostgresJdbcContext(SnakeCase, "ctx")

  // new FactData(ctx).(getFactData)
  val config = new ConfigData(ctx)

  if (config.getConfig().map(println(_)).length == 0) {
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
