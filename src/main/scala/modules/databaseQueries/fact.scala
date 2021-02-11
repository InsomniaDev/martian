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
    run {
      query[Fact]
        .filter(a => a.id == lift(factId))
    }
  }

  def getFacts(factId: Int): List[Fact] = {
    run {
      query[Fact]
    }
  }

  def getRelatedFactIds(factName: List[String]): List[Option[String]] = {
    run {
      query[Fact]
        .filter(a => liftQuery(factName).contains(a.related_facts))
        .map(a => a.related_fact_ids)
    }
  }

  def insertFact(fact: Fact) = {
    run {
      query[Fact]
        .insert(lift(fact))
    }
  }

  def updateRelatedFacts(fact: Fact) = {
    run {
      query[Fact]
        .filter(_.id == lift(fact.id))
        .update(
          _.related_facts -> lift(fact.related_facts),
          _.related_fact_ids -> lift(fact.related_fact_ids)
        )
    }
  }

  def updateFactData(fact: Fact) = {
    run {
      query[Fact]
        .filter(_.id == lift(fact.id))
        .update(_.fact_data -> lift(fact.fact_data))
    }
  }

  def deleteFactData(fact: Fact) = {
    run {
      query[Fact]
        .filter(_.id == lift(fact.id))
        .delete
    }
  }
}
