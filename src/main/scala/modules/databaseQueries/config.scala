package modules.databaseQueries
import io.getquill._
import java.util.UUID
import scala.concurrent.ExecutionContext
import scala.concurrent.Future

case class Config(
    config_uuid: UUID,
    name: String,
    record: String
)

class ConfigData(ctx: CassandraAsyncContext[SnakeCase.type])(implicit ec: ExecutionContext) {
  import ctx._

  def getConfigByName(configKey: String): Future[List[Config]] = {
    run {
      query[Config]
        .filter(_.name == lift(configKey))
    }
  }

  def getConfig(): Future[List[Config]] = {
    run {
      query[Config]
    }
  }

  def insertConfigItem(config: Config) = {
    run {
      query[Config]
        .insert(lift(config))
    }
    true
  }

  def updateConfig(config: Config) = {
    run {
      query[Config]
        .filter(_.config_uuid == lift(config.config_uuid))
        .update(
          _.config_uuid -> lift(config.config_uuid),
          _.name -> lift(config.name),
          _.record -> lift(config.record)
        )
    }
    true
  }

  def deleteConfig(config: Config) = {
    run {
      query[Config]
        .filter(_.config_uuid == lift(config.config_uuid))
        .delete
    }
    true
  }
}
