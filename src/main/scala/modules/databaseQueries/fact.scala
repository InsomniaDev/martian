package modules.databaseQueries

import io.getquill._

case class Fact(
    id: Int,
    name: String,
    related_fact_ids: String,
    related_facts: String,
    fact_data: Option[String]
)

case class FactsToWords(
    id: Int,
    fact_id: Int,
    word_id: Int,
    importance: Int
)

case class Word(
    id: Int,
    word: String
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
        .join(query[FactsToWords])
        .on(_.id == _.fact_id)
        .join(query[Word])
        .on({ case ((a, b), c) => b.word_id == c.id })
        .filter(_._2.word == lift(factName))
        .map(_._1._1)
    }
  }

  // checkFactNames gets the fact back by the provided name
  def checkFactNames(factName: List[String]): List[Fact] = {
    run(query[Fact].filter(a => liftQuery(factName).contains(a.name)))
  }

  // getFactsByUsedWords returns the list of relationships and the words used
  def getFactsByUsedWords(words: List[String]): List[(FactsToWords, Word)] = {
    run {
      query[FactsToWords]
        .join(query[Word].filter(a => liftQuery(words).contains(a.word)))
        .on(_.word_id == _.id)
        .filter({ case (ftw, w) => w.word == lift(words(1)) })
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
    run(query[Fact].filter(a => liftQuery(factName).contains(a.id)))
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
