package modules.databaseQueries
import io.getquill._

case class Config(
    id: Int,
    key: Option[String],
    value: Option[String]
)

// FIXME: Need to figure out what to do with ctx => this class would be better as a trait
class ConfigData(ctx: PostgresJdbcContext[SnakeCase.type]) {
  import ctx._

  def getConfigByKey(configKey: Option[String]): List[Config] = {
    val q = quote {
      query[Config]
        .filter(_.key == lift(configKey))
    }
    val resp = ctx.run(q)
    resp
  }

  def getConfig(): List[Config] = {
    val q = quote {
      query[Config]
    }
    val resp = ctx.run(q)
    resp
  }

  def insertConfigItem(config: Config) = {
    val q = quote {
      query[Config]
        .insert(
          _.key -> lift(config.key),
          _.value -> lift(config.value)
        )
    }
    val resp = ctx.run(q)
    resp
  }

  def updateConfig(config: Config) = {
    val q = quote {
      query[Config]
        .filter(_.id == lift(config.id))
        .update(
          _.key -> lift(config.key),
          _.value -> lift(config.value)
        )
    }
    val resp = ctx.run(q)
    resp
  }

  def deleteConfig(config: Config) = {
    val q = quote {
      query[Config]
        .filter(_.id == lift(config.id))
        .delete
    }
    val resp = ctx.run(q)
    resp
  }
}
