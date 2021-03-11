package modules.databaseQueries

import io.getquill._

case class Fact(
    id: Option[Int],
    name: String,
    fact_data: String,
    related_fact_ids: Option[String],
    related_facts: Option[String],
    importance: Int
)

case class FactsToWords(
    id: Option[Int],
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

  /** getFactsByUsedWords
    *
    * SELECT x1.id, x1.fact_id, x1.word_id, x1.importance, a.id, a.word
    * FROM facts_to_words x1
    *   INNER JOIN (
    *     SELECT a.id, a.word
    *     FROM word a
    *     WHERE a.word IN (?)
    *   ) AS a
    *   ON x1.word_id = a.id
    *   WHERE a.word = ?
    *
    * @param words
    * @return
    */
  def getFactsByUsedWords(words: List[String]): List[(FactsToWords, Word)] = {
    run {
      query[FactsToWords]
        .join(query[Word].filter(a => liftQuery(words).contains(a.word)))
        .on(_.word_id == _.id)
        .filter({ case (ftw, w) => w.word == lift(words(1)) })
    }
  }

  /** incrementFactToWordImportance
    *
    * UPDATE facts_to_words
    * SET importance = (importance + 1)
    * WHERE id = ?
    *
    * @param id
    */
  def incrementFactToWordImportance(id: Option[Int]) = {
    run {
      query[FactsToWords]
        .filter(_.id == lift(id))
        .update(p => p.importance -> (p.importance + 1))
    }
  }

  /** getFactsByIds
    *
    * SELECT a.id, a.name, a.related_fact_ids, a.related_facts, a.fact_data
    * FROM fact a
    * WHERE a.id IN (?)
    *
    * @param factIds
    * @return
    */
  def getFactsByIds(factIds: List[Int]): List[Fact] = {
    run(query[Fact].filter(a => liftQuery(factIds).contains(a.id)))
  }

  /** checkFactName
    *
    * SELECT x4.id, x4.name, x4.related_fact_ids, x4.related_facts, x4.fact_data, x4.importance
    * FROM fact x4
    * WHERE x4.name = ?
    *
    * @param factName
    * @return
    */
  def checkFactName(factName: String): List[Fact] = {
    run(query[Fact].filter(_.name == lift(factName)))
  }

  /** upsertFact
    *
    * INSERT INTO fact AS t (name,fact_data)
    * VALUES (?, ?)
    * ON CONFLICT DO NOTHING
    * RETURNING id, name, related_fact_ids, related_facts, fact_data, importance
    *
    * @param fact
    * @return
    */
  def upsertFact(fact: Fact): Fact = {
    // FIXME: Update this to update rather than ignore if it already exists "onConflictUpdate"
    run {
      query[Fact]
        .insert(
          _.name -> lift(fact.name),
          _.fact_data -> lift(fact.fact_data)
        )
        .onConflictIgnore
        .returning(r =>
          (new Fact(
            r.id,
            r.name,
            r.related_fact_ids.getOrElse(""),
            r.related_facts,
            Some(r.fact_data),
            r.importance
          ))
        )
    }
  }

  /** batchInsertWordsToFact
    *
    * INSERT INTO facts_to_words (fact_id,word_id)
    * VALUES (?, ?)
    *
    * @param insertValues
    */
  def batchInsertWordsToFact(insertValues: List[FactsToWords]) = {
    run {
      quote {
        liftQuery(insertValues).foreach(e =>
          query[FactsToWords]
            .insert(
              _.fact_id -> e.fact_id,
              _.word_id -> e.word_id
            )
        )
      }
    }
  }

  /** getIdsForWords
    * 
    * SELECT a.id, a.word 
    * FROM word a 
    * WHERE a.word IN (?)
    *
    * @param words
    */
  def getIdsForWords(words: List[String]) = {
    run {
      query[Word]
       .filter(a => liftQuery(words).contains(a.word))
    }
  }

  def updateFactData(fact: Fact) = {
    run {
      query[Fact]
        .filter(_.id == lift(fact.id))
        .update(_.fact_data -> lift(fact.fact_data))
    }
  }

  def getFact(factId: Option[Int]): List[Fact] = {
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

  // def checkFactName(factName: String): List[Fact] = {
  //   run {
  //     query[Fact]
  //       .join(query[FactsToWords])
  //       .on(_.id == _.fact_id)
  //       .join(query[Word])
  //       .on({ case ((a, b), c) => b.word_id == c.id })
  //       .filter(_._2.word == lift(factName))
  //       .map(_._1._1)
  //   }
  // }

  // def getRelatedFactIds(factName: String): List[String] = {
  //   run {
  //     query[Fact]
  //       .filter(_.related_facts like lift(factName))
  //       .map(a => a.related_fact_ids)
  //   }
  // }

  // getRelatedFactsByIds gets the facts back by the provided ids
  def getRelatedFactsByIds(factName: List[Int]): List[Fact] = {
    run(query[Fact].filter(a => liftQuery(factName).contains(a.id)))
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

  def deleteFactData(fact: Fact) = {
    run {
      query[Fact]
        .filter(_.id == lift(fact.id))
        .delete
    }
  }
}
