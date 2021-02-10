package modules.knowledge

import modules.databaseQueries._
import io.getquill._

class FactParser(ctx: PostgresJdbcContext[SnakeCase.type])
    extends FactData(ctx) {
        // TODO: Need to check for fact existence here
        // TODO: Need to use config values in here and have a class with those values
    }
