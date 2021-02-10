package modules.knowledge

import modules.databaseQueries._
import io.getquill._

// TODO: Could we possibly export all of this knowledge data into markdown files and then display them through something like Hugo?
// https://gohugo.io/documentation/

class FactParser(ctx: PostgresJdbcContext[SnakeCase.type])
    extends FactData(ctx) {
        // TODO: Need to check for fact existence here
        // TODO: Need to use config values in here and have a class with those values
    }
