package modules.databaseQueries

import io.getquill._

case class Fact(
    id: Int,
    name: String,
    related_fact_ids: String,
    related_facts: String,
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

  def checkFactName(factName: String): List[Fact] = {
    run {
      query[Fact]
        .filter(_.name like lift(factName))
    }
  }

  // checkFactNames gets the fact back by the provided name
  def checkFactNames(factName: List[String]): List[Fact] = {
    run {
      query[Fact]
        .filter(a => liftQuery(factName).contains(a.name))
    }
  }

  def getRelatedFactIds(factName: String): List[String] = {
    run {
      query[Fact]
        .filter(_.related_facts like lift(factName))
        .map(a => a.related_fact_ids)
    }
  }

  // getRelatedFactsByIds gets the facts back by the provided ids
  def getRelatedFactsByIds(factName: List[Int]): List[Fact] = {
    run {
      query[Fact]
        .filter(a => liftQuery(factName).contains(a.id))
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
