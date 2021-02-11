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
    run {
      query[Config]
        .filter(_.key == lift(configKey))
    }
  }

  def getConfig(): List[Config] = {
    run {
      query[Config]
    }
  }

  def insertConfigItem(config: Config) = {
    run {
      query[Config]
        .insert(
          _.key -> lift(config.key),
          _.value -> lift(config.value)
        )
    }
  }

  def updateConfig(config: Config) = {
    run {
      query[Config]
        .filter(_.id == lift(config.id))
        .update(
          _.key -> lift(config.key),
          _.value -> lift(config.value)
        )
    }
  }

  def deleteConfig(config: Config) = {
    run {
      query[Config]
        .filter(_.id == lift(config.id))
        .delete
    }
  }
}
