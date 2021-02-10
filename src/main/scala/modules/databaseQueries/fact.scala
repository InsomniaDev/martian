package modules.databaseQueries
import io.getquill._

case class Fact(
    id: Int,
    related_fact_ids: Option[String],
    related_facts: Option[String],
    fact_data: Option[String]
)

class FactData(ctx: PostgresJdbcContext[SnakeCase.type]) {
  import ctx._

  def getFact(factId: Int): List[Fact] = {
    val q = quote {
      query[Fact]
        .filter(a => a.id == lift(factId))
    }
    val resp = ctx.run(q)
    resp
  }

  def getFacts(factId: Int): List[Fact] = {
    val q = quote {
      query[Fact]
    }
    val resp = ctx.run(q)
    resp
  }

  def getRelatedFactIds(factName: List[String]): List[Option[String]] = {
    val q = quote {
      query[Fact]
        .filter(a => liftQuery(factName).contains(a.related_facts))
        .map(a => a.related_fact_ids)
    }
    val resp = ctx.run(q)
    resp
  }

  def insertFact(fact: Fact) = {
    val q = quote {
      query[Fact]
        .insert(lift(fact))
    }
    val resp = ctx.run(q)
    resp
  }

  def updateRelatedFacts(fact: Fact) = {
    val q = quote {
      query[Fact]
        .filter(_.id == lift(fact.id))
        .update(
          _.related_facts -> lift(fact.related_facts),
          _.related_fact_ids -> lift(fact.related_fact_ids)
        )
    }
    val resp = ctx.run(q)
    resp
  }

  def updateFactData(fact: Fact) = {
    val q = quote {
      query[Fact]
        .filter(_.id == lift(fact.id))
        .update(_.fact_data -> lift(fact.fact_data))
    }
    val resp = ctx.run(q)
    resp
  }

  def deleteFactData(fact: Fact) = {
    val q = quote {
      query[Fact]
        .filter(_.id == lift(fact.id))
        .delete
    }
    val resp = ctx.run(q)
    resp
  }
}
