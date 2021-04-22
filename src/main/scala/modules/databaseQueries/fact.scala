package modules.databaseQueries

import io.getquill._
import java.util.UUID
// val q = quote {
//   query[Book].filter(p => p.pages.contains(25)).allowFiltering
// }
// ctx.run(q)
    // record_uuid uuid,
    // account_uuid uuid,
    // tags set<text>,
    // words set<text>,
    // record text,
    // importance int,

case class Records(
    record_uuid: UUID,
    account_uuid: UUID,
    tags: Set[String],
    words: Set[String],
    record: String,
    importance: Int
)


class FactData(ctx: CassandraAsyncContext[SnakeCase.type]) {
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
    * ON CONFLICT (name) 
    * DO UPDATE 
    *   SET fact_data = EXCLUDED.fact_data, 
    *       common_words = EXCLUDED.common_words, 
    *       related_fact_ids = EXCLUDED.related_fact_ids, 
    *       related_facts = EXCLUDED.related_facts, 
    *       importance = EXCLUDED.importance 
    * RETURNING id, name, fact_data, common_words, related_fact_ids, related_facts, importance
    *
    * @param fact
    * @return
    */
  def upsertFact(fact: Fact): Fact = {
    run {
      query[Fact]
        .insert(
          _.name -> lift(fact.name),
          _.fact_data -> lift(fact.fact_data)
        )
        .onConflictUpdate(_.name)(
          (oldVal, newVal) => oldVal.fact_data -> newVal.fact_data,
          (oldVal, newVal) => oldVal.common_words -> newVal.common_words,
          (oldVal, newVal) => oldVal.related_fact_ids -> newVal.related_fact_ids,
          (oldVal, newVal) => oldVal.related_facts -> newVal.related_facts,
          (oldVal, newVal) => oldVal.importance -> newVal.importance
        )
        .returning(r =>
          (new Fact(
            r.id,
            r.name,
            r.fact_data,
            r.common_words,
            r.related_fact_ids,
            r.related_facts,
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

}
